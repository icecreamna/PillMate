import 'dart:convert';
import 'package:http/http.dart' as http;

class OtpService {
  static const baseUrl = "http://10.0.2.2:8080";

  Future<void> requestOtp(int patientId) async {
    final res = await http.post(
      Uri.parse("$baseUrl/patient/$patientId/otp/request"),
    );
    print("requestOtp${patientId}");
    if (res.statusCode != 200) {
      throw Exception("Request OTP failed: ${res.body}");
    }
  }

  // Future<void> resendOtp(int patientId) async {
  //   final res = await http.post(
  //     Uri.parse("$baseUrl/patient/$patientId/otp/resend"),
  //   );
  //   if (res.statusCode != 200) {
  //     throw Exception("Resend OTP failed: ${res.body}");
  //   }
  // }

  Future<Map<String, dynamic>> verifyOtp(int patientId, String otpCode) async {
    final res = await http.post(
      Uri.parse("$baseUrl/patient/$patientId/otp/verify"),
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({"otp_code": otpCode}),
    );
    return {"statusCode": res.statusCode, "body": jsonDecode(res.body)};
  }
}
