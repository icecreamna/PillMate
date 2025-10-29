import 'dart:convert';

import 'package:frontend/services/auth_service.dart';
import 'package:http/http.dart' as http;

class GroupService {
  static const baseUrl = "http://10.0.2.2:8080";

  Future<List<dynamic>> getGroups() async {
    final String? token = AuthService.jwtToken;
    if (token == null) return [];

    final res = await http.get(
      Uri.parse("$baseUrl/api/groups"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );

    if (res.statusCode == 200) {
      final body = jsonDecode(res.body);
      print("โหลดกลุ่มสำเร็จ${body["data"]}");
      return body["data"] ?? [];
    } else {
      print("❌ โหลดกลุ่มไม่สำเร็จ: ${res.statusCode}");
      return [];
    }
  }

  Future<Map<String, dynamic>?> getGroupWithDetail({
    required int groupId,
  }) async {
    final String? token = AuthService.jwtToken;
    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return null;
    }

    final res = await http.get(
      Uri.parse("$baseUrl/api/groups/$groupId"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );

    if (res.statusCode == 200) {
      final data = jsonDecode(res.body);
      return data["data"];
    } else {
      throw Exception("โหลดกลุ่มไม่สำเร็จ (${res.statusCode})");
    }
  }

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

  Future<bool> updateGroup({
    required int groupId,
    required String newGroupName,
    required List<String> medicineIds,
  }) async {
    final String? token = AuthService.jwtToken;
    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return false;
    }
    final res = await http.put(
      Uri.parse("$baseUrl/api/groups/$groupId"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
      body: jsonEncode({
        "new_group_name": newGroupName,
        "my_medicine_ids": medicineIds.map((e) => int.parse(e)).toList(),
      }),
    );
    if (res.statusCode == 200) {
      print("✅ อัปเดตกลุ่มสำเร็จ: ${res.body}");
      return true;
    } else {
      throw Exception(
        "❌ อัปเดตกลุ่มไม่สำเร็จ (${res.statusCode}): ${res.body}",
      );
    }
  }

  Future<bool> deleteGroup({required String groupId}) async {
    final String? token = AuthService.jwtToken;
    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return false;
    }
    final res = await http.delete(
      Uri.parse("$baseUrl/api/groups/$groupId"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );
    if(res.statusCode == 200){
      print("✅ ลบกลุ่มสำเร็จ: ${res.body}");
      return true;
    } else {
      throw Exception("❌ ลบกลุ่มไม่สำเร็จ (${res.statusCode}): ${res.body}");
    }
  }
}
