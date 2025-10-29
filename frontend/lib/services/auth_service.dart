import 'dart:convert';
import 'package:path_provider/path_provider.dart';
import 'dart:io';
import 'package:http/http.dart' as http;

class AuthService {
  static const baseURL = "http://10.0.2.2:8080";
  static String? jwtToken;

  static Future<File> _tokenFile() async {
    final dir = await getApplicationDocumentsDirectory();
    return File('${dir.path}/jwt_token.txt');
  }

  // ✅ บันทึก token ลงไฟล์
  static Future<void> _saveTokenToFile(String token) async {
    final file = await _tokenFile();
    await file.writeAsString(token);
    print("💾 Token saved to file: ${file.path}");
  }

  // ✅ โหลด token จากไฟล์ (ใช้ตอน background)
  static Future<String?> loadTokenFromFile() async {
    try {
      final file = await _tokenFile();
      if (await file.exists()) {
        final token = await file.readAsString();
        jwtToken = token;
        print("🔑 Loaded token from file");
        return token;
      }
    } catch (e) {
      print("❌ loadTokenFromFile error: $e");
    }
    return null;
  }

  Future<Map<String, dynamic>> register({
    required String email,
    required String password,
  }) async {
    final uri = Uri.parse("$baseURL/register");

    final lowerEmail = email.trim().toLowerCase();

    final resp = await http.post(
      uri,
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({"email": lowerEmail, "password": password}),
    );

    if (resp.statusCode != 200) {
      throw Exception("Register failed: ${resp.body}");
    }

    final data = jsonDecode(resp.body);

    final patientId = data["patient_id"];
    await http.post(Uri.parse("$baseURL/patient/$patientId/otp/request"));

    return data;
  }

  Future<Map<String, dynamic>?> login({
    required String email,
    required String password,
  }) async {
    final lowerEmail = email.trim().toLowerCase();

    final res = await http.post(
      Uri.parse("$baseURL/login"),
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({"email": lowerEmail, "password": password}),
    );

    if (res.statusCode == 200) {
      final data = jsonDecode(res.body);
      final token = data["token"];
      if (token != null) {
        jwtToken = token;
        print("✅ Login success & token saved to file");
        await _saveTokenToFile(token); // ✅ เก็บลงไฟล์แทน SharedPreferences
      }
      return data;
    } else {
      print("❌ Login failed (${res.statusCode}): ${res.body}");
      return null;
    }
  }

  Future<bool> logout() async {
    try {
      jwtToken = null;
      final res = await http.post(
        Uri.parse("$baseURL/logout"),
        headers: {"Content-Type": "application/json"},
      );

      if (res.statusCode == 200) {
        print("✅ Logout success: ${res.body}");
        return true;
      } else {
        print("❌ Logout failed (${res.statusCode}): ${res.body}");
        return false;
      }
    } catch (e) {
      print("Logout exception: $e");
      return false;
    }
  }

  Future<int?> forgotPasswordEmail(String email) async {
    final res = await http.post(
      Uri.parse("$baseURL/patient/password/forgot"),
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({"email": email}),
    );

    if (res.statusCode == 200) {
      final data = jsonDecode(res.body);
      return data["patient_id"];
    } else {
      print("❌ Forgot password error: ${res.body}");
      return null;
    }
  }

  Future<Map<String, dynamic>> resetPassword({
    required int patientId,
    required String newPassword,
  }) async {
    final res = await http.put(
      Uri.parse("$baseURL/patient/$patientId/reset-password"),
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({"new_password": newPassword}),
    );

    final data = jsonDecode(res.body);

    if (res.statusCode == 200) return data;
    print("updated failed");
    throw Exception("Service Reset password failed: ${res.body}");
  }

}
