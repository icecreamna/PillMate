import 'dart:convert';
import 'package:frontend/models/notification_info.dart';
import 'package:http/http.dart' as http;
import 'package:frontend/services/auth_service.dart';

class NotiService {
  static const baseUrl = "http://10.0.2.2:8080";

  Future<List<NotiFormatModel>> getFormats() async {
    final token = AuthService.jwtToken;
    if (token == null) throw Exception("Missing JWT token");

    final res = await http.get(
      Uri.parse("$baseUrl/api/noti-formats"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );

    if (res.statusCode == 200) {
      final body = jsonDecode(res.body);
      final List data = body["data"] ?? [];
      return data.map((e) => NotiFormatModel.fromJson(e)).toList();
    } else {
      throw Exception("โหลด noti formats ไม่สำเร็จ: ${res.body}");
    }
  }

  String _toIsoDate(String thaiDate) {
    try {
      // ตัวอย่าง thaiDate: "12 ต.ค. 2568"
      final parts = thaiDate.split(' ');
      if (parts.length < 3) return thaiDate;
      final day = parts[0];
      final monthStr = parts[1].replaceAll('.', '');
      final year = (int.parse(parts[2]) - 543).toString();

      final monthMap = {
        "มค": "01",
        "กพ": "02",
        "มีค": "03",
        "เมย": "04",
        "พค": "05",
        "มิย": "06",
        "กค": "07",
        "สค": "08",
        "กย": "09",
        "ตค": "10",
        "พย": "11",
        "ธค": "12",
      };
      final month = monthMap[monthStr] ?? "01";
      return "$year-$month-${day.padLeft(2, '0')}";
    } catch (_) {
      return thaiDate; // fallback ถ้า parse ไม่ได้
    }
  }

  Future<bool> addNotification({
    int? myMedicineId, // optional
    int? groupId, // optional
    required int notiFormatId,
    required NotificationInfo info,
  }) async {
    final token = AuthService.jwtToken;
    if (token == null) throw Exception("Missing JWT token");

    if (myMedicineId == null && groupId == null) {
      throw Exception("ต้องส่ง myMedicineId หรือ groupId อย่างใดอย่างหนึ่ง");
    }

    String endpoint;
    switch (notiFormatId) {
      case 1:
        endpoint = "fixed-times";
        break;
      case 2:
        endpoint = "interval";
        break;
      case 3:
        endpoint = "every-n-days";
        break;
      case 4:
        endpoint = "cycle";
        break;
      default:
        throw Exception("Unknown noti format");
    }
    final url = Uri.parse("$baseUrl/api/noti/$endpoint");
    final startDate = _toIsoDate(info.startDate);
    final endDate = _toIsoDate(info.endDate);
    final Map<String, dynamic> body = {
      if (myMedicineId != null) "my_medicine_id": myMedicineId,
      if (groupId != null) "group_id": groupId,
      "start_date": startDate,
      "end_date": endDate,
      "noti_format_id": notiFormatId,
    };
    switch (notiFormatId) {
      case 1: // Fixed Times
        if (info.times != null) body["times"] = info.times;
        break;
      case 2: // Interval
        if (info.intervalHours != null) body["interval_hours"] = info.intervalHours;
        if (info.intervalTake != null) body["times_per_day"] = info.intervalTake;
        if (info.times != null) body["times"] = info.times;
        break;
      case 3: // Every N Days
        if (info.daysGap != null) body["interval_day"] = info.daysGap;
        if (info.times != null) body["times"] = info.times;
        break;
      case 4: // Cycle
        if (info.takeDays != null && info.breakDays != null) {
          body["cycle_pattern"] = [info.takeDays, info.breakDays];
        }
        if (info.times != null) body["times"] = info.times;
        break;
    }
    body.removeWhere((key, value) => value == null);
    print("ส่ง${body}");
    final res = await http.post(
      url,
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
      body: jsonEncode(body),
    );
    if (res.statusCode >= 200 && res.statusCode < 300) {
      print("✅ เพิ่มการแจ้งเตือนสำเร็จ (${res.statusCode})");
      return true;
    } else {
      throw Exception("❌ เพิ่มการแจ้งเตือนไม่สำเร็จ: ${res.body}");
    }
  }
}
