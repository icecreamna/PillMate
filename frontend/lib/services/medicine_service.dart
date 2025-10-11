import 'dart:convert';

import 'package:frontend/models/dose.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:http/http.dart' as http;

class MedicineService {
  static const baseUrl = "http://10.0.2.2:8080";

  Future<List<Dose>?> getMyMedicines() async {
    final String? token = AuthService.jwtToken;
    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return null;
    }
    final res = await http.get(
      Uri.parse("$baseUrl/api/my-medicines"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );

    if (res.statusCode == 200) {
      final body = jsonDecode(res.body);
      final List data = body["data"] ?? [];
      return data.map((e) => Dose.fromJson(e)).toList();
    } else {
      print("❌ Error loading medicines: ${res.body}");
      return [];
    }
  }

  Future<bool> addMedicineInfo({
    required String medName,
    required String genericName,
    required String properties,
    required int formId,
    required int unitId,
    required int instructionId,
    required String amountPerTime,
    required String timePerDay,
  }) async {
    final String? token = AuthService.jwtToken;
    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return false;
    }
    final res = await http.post(
      Uri.parse("$baseUrl/api/my-medicine"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
      body: jsonEncode({
        "med_name": medName,
        "generic_name": genericName,
        "properties": properties,
        "form_id": formId,
        "unit_id": unitId,
        "instruction_id": instructionId,
        "amount_per_time": amountPerTime,
        "times_per_day": timePerDay,
        "med_status": "active",
      }),
    );
    if (res.statusCode == 201) {
      print("✅ Medicine created: ${res.body}");
      return true;
    } else {
      print("❌ Error: ${res.body}");
      return false;
    }
  }

  Future<bool> updatedMedicineInfo({
    required int id,
    required String medName,
    required String genericName,
    required String properties,
    required int formId,
    required int unitId,
    required int instructionId,
    required String amountPerTime,
    required String timePerDay,
  }) async {
    final String? token = AuthService.jwtToken;
    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return false;
    }
    final res = await http.put(
      Uri.parse("$baseUrl/api/my-medicine/$id"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
      body: jsonEncode({
        "med_name": medName,
        "generic_name": genericName,
        "properties": properties,
        "form_id": formId,
        "unit_id": unitId,
        "instruction_id": instructionId,
        "amount_per_time": amountPerTime,
        "times_per_day": timePerDay,
        "med_status": "active",
      }),
    );
    if (res.statusCode == 200) {
      print("✅ Medicine updated: ${res.body}");
      return true;
    } else {
      print("❌ Error: ${res.body}");
      return false;
    }
  }

  Future<bool> deleteMedicineInfo({required int id}) async {
    final String? token = AuthService.jwtToken;
    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return false;
    }
    final res = await http.delete(
      Uri.parse("$baseUrl/api/my-medicine/$id"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );
    if (res.statusCode == 200) {
      print("✅ Medicine updated: ${res.body}");
      return true;
    } else {
      print("❌ Error: ${res.body}");
      return false;
    }
  }
}
