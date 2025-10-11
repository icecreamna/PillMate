import 'dart:convert';

import 'package:frontend/services/auth_service.dart';
import 'package:http/http.dart' as http;

class GroupService {
  static const baseUrl = "http://10.0.2.2:8080";

  Future<bool> createGroup({
    required String groupName,
    required List<String> medicineIds,
  }) async {
    final String? token = AuthService.jwtToken;
    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return false;
    }
    final res = await http.post(
      Uri.parse("$baseUrl/api/groups"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
      body: jsonEncode({
        "group_name": groupName,
        "my_medicine_ids": medicineIds.map((e) => int.parse(e)).toList(),
      }),
    );
    if (res.statusCode == 201) {
      print("✅ กลุ่มถูกสร้างแล้ว: ${res.body}");
      return true;
    } else {
      print("❌ สร้างกลุ่มไม่สำเร็จ: ${res.statusCode} - ${res.body}");
      return false;
    }
  }

  Future<List<dynamic>> getGroups() async {
    final String? token = AuthService.jwtToken;
    if (token == null) return [];

    final res = await http.get(
      Uri.parse("$baseUrl/api/groups"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );

    if (res.statusCode == 200) {
      final body = jsonDecode(res.body);
      return body["data"] ?? [];
    } else {
      print("❌ โหลดกลุ่มไม่สำเร็จ: ${res.statusCode}");
      return [];
    }
  }
}
