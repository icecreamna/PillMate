import 'dart:convert';

import 'package:frontend/services/auth_service.dart';
import 'package:http/http.dart' as http;

class ProfileService {
  static const baseUrl = "http://10.0.2.2:8080";

  Future<bool> setUpProfile({
    required int patientId,
    required String firstName,
    required String lastName,
    required String idCardNumber,
    required String phoneNumber,
  }) async {
    final res = await http.put(
      Uri.parse("$baseUrl/patient/$patientId/profile"),
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({
        "first_name": firstName,
        "last_name": lastName,
        "id_card_number": idCardNumber,
        "phone_number": phoneNumber,
      }),
    );
    if (res.statusCode == 200) return true;
    print("❌ update failed: ${res.body}");
    return false;
  }

  Future<Map<String, dynamic>?> fetchProfile() async {
    final token = AuthService.jwtToken;

    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return null;
    }

    final res = await http.get(
      Uri.parse("$baseUrl/api/patient/me"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
    );

    if (res.statusCode == 200) {
      final data = jsonDecode(res.body);
      print("✅ Fetch profile success: ${data["data"]}");
      return data["data"];
    } else {
      print("❌ Fetch profile failed: ${res.body}");
      return null;
    }
  }

  Future<Map<String, dynamic>?> fetchAppointment() async {
    final token = AuthService.jwtToken;

    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return null;
    }
    try {
      final res = await http.get(
        Uri.parse("$baseUrl/api/appointments/next"),
        headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
      );

      if (res.statusCode == 200) {
        final data = jsonDecode(res.body);
        print("✅ Fetch Appoint success: ${data}");
        return data;
      } else if (res.statusCode == 404) {
        print("📭 No upcoming appointment");
        return null;
      } else {
        print("❌ Fetch appointment failed: ${res.statusCode} ${res.body}");
        return null;
      }
    } catch (e) {
      print("🚨 Error fetching appointment: $e");
      return null;
    }
  }

  Future<bool> updateInfoProfile({
    required String idCard,
    required String firstName,
    required String lastName,
    required String tel,
  }) async {
    final token = AuthService.jwtToken;
    if (token == null) {
      print("⚠️ Missing token — user not logged in");
      return false;
    }

    final res = await http.put(
      Uri.parse("$baseUrl/api/patient/me"),
      headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
      body: jsonEncode({
        "first_name": firstName,
        "last_name": lastName,
        "id_card_number": idCard,
        "phone_number": tel,
      }),
    );
    if (res.statusCode == 200) {
      print("updated complete ${res.body}");
      return true;
    } else {
      print("cant updated ${res.body}-${res.statusCode} ");
      return false;
    }
  }
}
