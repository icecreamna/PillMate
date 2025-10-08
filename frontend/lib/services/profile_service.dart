import 'dart:convert';

import 'package:http/http.dart' as http;

class ProfileService {
  static const baseUrl = "http://10.0.2.2:8080";

  Future<bool> updateProfile({
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
}
