import 'dart:convert';

import 'package:frontend/providers/today_provider.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:http/http.dart' as http;
import 'package:intl/intl.dart';

class TodayService {
  static const baseUrl = "http://10.0.2.2:8080";

  Future<List<DoseGroup>> fetchTodayNoti(DateTime date) async {
    final String? token = AuthService.jwtToken;
    if (token == null) throw Exception("Token missing");

    final url = Uri.parse("$baseUrl/api/noti-items?");
    final res = await http.get(
      url,
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );
    if (res.statusCode != 200) {
      throw Exception("โหลดข้อมูลไม่สำเร็จ (${res.statusCode})");
    }
    final body = jsonDecode(res.body);
    final List data = body["data"] ?? [];
    final List groups = body["group_cards"] ?? [];

    final List<DoseGroup> all = [];

    // singles
    for (var d in data) {
      final notiSingleId = d["id"];
      final medicineId = d["my_medicine_id"];
      final name = d["med_name"] ?? "-";
      final unit = d["unit_name"] ?? "";
      final instruction = d["instruction_name"] ?? "";
      final amountPerTime = d["amount_per_time"] ?? "";
      final notifyDate = DateTime.parse(
        "${d["notify_date"]} ${d["notify_time"]}",
      );

      all.add(
        DoseGroup(
          notiSingleId: notiSingleId,
          medicineId: medicineId,
          nameGroup: "-",
          key: "${d["notify_time"]}-$instruction",
          at: notifyDate,
          instruction: instruction,
          doses: [
            DoseSingle(name: name, amountPerTime: amountPerTime, unit: unit),
          ],
        ),
      );
    }

    // groups
    for (var g in groups) {
      final groupId = g["group_id"];
      final nameGroup = g["group_name"] ?? "-";
      final instruction = g["items"].isNotEmpty
          ? (g["items"][0]["instruction_name"] ?? "")
          : "";
      final notifyDate = DateTime.parse(
        "${g["notify_date"]} ${g["notify_time"]}",
      );

      final doses = (g["items"] as List)
          .map(
            (item) => DoseSingle(
              name: item["med_name"] ?? "-",
              unit: item["unit_name"] ?? "",
              amountPerTime: item["amount_per_time"] ?? "",
            ),
          )
          .toList();

      final notiGroupIds = (g["items"] as List)
          .map<int>((item) => item["noti_item_id"] as int)
          .toList();

      all.add(
        DoseGroup(
          notiGroupIds: notiGroupIds,
          groupId: groupId,
          nameGroup: nameGroup,
          key: "${g["notify_time"]}-$instruction",
          at: notifyDate,
          instruction: instruction,
          doses: doses,
        ),
      );
    }
    return all;
  }

  Future<List<Map<String, dynamic>>> fetchAllSymptoms() async {
    final String? token = AuthService.jwtToken;
    if (token == null) throw Exception("Token missing");

    final res = await http.get(
      Uri.parse("$baseUrl/api/symptoms"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );

    if (res.statusCode == 200) {
      final body = jsonDecode(res.body);
      final List data = body["data"] ?? [];
      return data.map((e) => Map<String, dynamic>.from(e)).toList();
    } else {
      throw Exception(
        "โหลด symptoms ไม่สำเร็จ: ${res.statusCode} - ${res.body}",
      );
    }
  }

  Future<bool> markTaken({required int notiItemId, required bool taken}) async {
    final String? token = AuthService.jwtToken;
    if (token == null) throw Exception("Token missing");

    final url = Uri.parse("$baseUrl/api/noti-items/$notiItemId/taken");

    final res = await http.patch(
      url,
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
      body: jsonEncode({"taken": taken}),
    );

    if (res.statusCode == 200) {
      print("✅ อัปเดตสถานะกินยาเรียบร้อย: ${res.body}");
      return true;
    } else {
      throw Exception(
        "❌ อัปเดตสถานะกินยาไม่สำเร็จ: ${res.statusCode} - ${res.body}",
      );
    }
  }


  Future<Map<String, dynamic>?> createSymptom({
    required String symptomNote,
    required int notiItemId,
    int? myMedicineId,
    int? groupId,
  }) async {
    final String? token = AuthService.jwtToken;
    if (token == null) throw Exception("Token missing");

    if (myMedicineId == null && groupId == null) {
      throw Exception("ต้องส่ง myMedicineId หรือ groupId อย่างใดอย่างหนึ่ง");
    }

    final Map<String, dynamic> body = {
      "noti_item_id": notiItemId,
      "symptom_note": symptomNote,
      if (myMedicineId != null) "my_medicine_id": myMedicineId,
      if (groupId != null) "group_id": groupId,
    };

    print("ส่งข้อมูล Symptom ไป: $body");

    final res = await http.post(
      Uri.parse("$baseUrl/api/symptom"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
      body: jsonEncode(body),
    );

    if (res.statusCode == 201) {
      final body = jsonDecode(res.body);
      print("✅ เพิ่มอาการสำเร็จ: $body");
      return body["data"];
    } else {
      print("❌ เพิ่มอาการไม่สำเร็จ: ${res.statusCode} - ${res.body}");
      return null;
    }
  }

  Future<bool> updateSymptom({
    required int symptomId,
    required String symptomNote,
  }) async {
    final String? token = AuthService.jwtToken;
    if (token == null) throw Exception("Token missing");

    final res = await http.patch(
      Uri.parse("$baseUrl/api/symptom/$symptomId"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
      body: jsonEncode({"symptom_note": symptomNote}),
    );

    if (res.statusCode == 200) {
      print("✅ อัปเดตอาการสำเร็จ: ${res.body}");
      return true;
    } else {
      print("❌ อัปเดตอาการไม่สำเร็จ: ${res.statusCode} - ${res.body}");
      return false;
    }
  }
}
