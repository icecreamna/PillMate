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

  Future<List<Dose>?> getMedeicineHospital() async {
    final String? token = AuthService.jwtToken;
    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return null;
    }

    final res = await http.post(
      Uri.parse("$baseUrl/api/my-medicine/sync-from-prescription"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );

    if (res.statusCode == 200) {
      final body = jsonDecode(res.body);
      final data = body['data'] ?? [];
      final List<Dose> newDoses = List<Dose>.from(
        (data as List).map(
          (e) => Dose.fromJson({
            "id": e["mymedicine_id"],
            "med_name": e["med_name"],
            "properties": e["properties"] ?? "-",
            "form_id": e["form_id"],
            "unit_id": e["unit_id"],
            "instruction_id": e["instruction_id"],
            "start_date": e["start_date"],
            "end_date": e["end_date"],
            "note": e["note"],
            "form_name": "-",
            "unit_name": "-",
            "instruction_name": "-",
            "amount_per_time": e["amount_per_time"],
            "times_per_day": e["times_per_day"],
            "source": e["source"] ?? "hospital",
          }),
        ),
      );
      print("✅ ซิงค์จากโรงพยาบาลสำเร็จ: ${data.length} รายการ");
      return newDoses;
    } else {
      print("❌ Error loading medicines from hospital: ${res.body}");
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

  Future<int?> getCountNotification() async {
    final String? token = AuthService.jwtToken;
    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return null;
    }

    final res = await http.get(
      Uri.parse("$baseUrl/api/prescriptions/sync-status"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );
    if (res.statusCode == 200) {
      final body = jsonDecode(res.body);
      print("Number of count notification ${body["count"]}");
      return body["count"] ?? 0;
    } else {
      throw Exception("❌ Cant load number of count: ${res.statusCode} - ${res.body}");
    }
  }
}
