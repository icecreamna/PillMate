import 'dart:convert';
import 'package:http/http.dart' as http;

class AuthService {
  static const baseURL = "http://10.0.2.2:8080";

  Future<Map<String, dynamic>> register({
    required String email,
    required String password,
  }) async {
    final uri = Uri.parse("$baseURL/register");

    final resp = await http.post(
      uri,
      headers: {"Content-type": "application/json"},
      body: jsonEncode({"email": email, "password": password}),
    );

    final data = jsonDecode(resp.body);

    if (resp.statusCode != 200 && resp.statusCode != 201) {
      throw Exception(data["message"] ?? "Error: ${resp.statusCode}");
    }
    return data;
  }
}
