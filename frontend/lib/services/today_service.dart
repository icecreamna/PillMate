import 'dart:convert';

import 'package:frontend/providers/today_provider.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:http/http.dart' as http;

class TodayService {
  static const baseUrl = "http://10.0.2.2:8080";

  Future<List<DoseGroup>> fetchTodayNoti(DateTime date) async {
    final String? token = AuthService.jwtToken;
    if (token == null) throw Exception("Token missing");

    final formatted = date.toIso8601String().split('T').first;
    print("ส่ง formatted ไปแบบนี้ ${formatted}");
    final url = Uri.parse("$baseUrl/api/noti-items?");
    final res = await http.get(
      url,
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );
    if (res.statusCode != 200) {
      throw Exception("โหลดข้อมูลไม่สำเร็จ (${res.statusCode})");
    }
    final body = jsonDecode(res.body);
    print("โหลด item $body");
    final List data = body["data"] ?? [];
    final List groups = body["group_cards"] ?? [];

    final List<DoseGroup> all = [];

    for (var d in data) {
      final id = d["id"];
      final name = d["med_name"] ?? "-";
      final unit = d["unit_name"] ?? "";
      final instruction = d["instruction_name"] ?? "";
      final amountPerday = d["amount_per_time"] ?? "";
      final notifyDate = DateTime.parse(
        "${d["notify_date"]} ${d["notify_time"]}",
      );
      final group = DoseGroup(
        id: id,
        nameGroup: "-",
        key: "${d["notify_time"]}-${instruction}",
        at: notifyDate,
        instruction: instruction,
        doses: [
          DoseSingle(name: name, amountPerTime: amountPerday, unit: unit),
        ],
      );
      all.add(group);
    }

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
              amountPerTime: item["amount_per_time"] ?? ""
            ),
          )
          .toList();

      final group = DoseGroup(
        id: groupId,
        nameGroup: nameGroup,
        key: "${g["notify_time"]}-${instruction}",
        at: notifyDate,
        instruction: instruction,
        doses: doses,
      );
      all.add(group);
    }
    return all;
  }
}
