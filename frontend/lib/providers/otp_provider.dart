import 'dart:async';

import 'package:flutter/material.dart';
import 'package:frontend/enums/page_type.dart';
import 'package:frontend/services/otp_service.dart';

class OtpProvider extends ChangeNotifier {
  final OtpService otpService;

  final PageType otpType;
  String emailText = "";
  String errorOtp = "";
  int countdown = 0;
  final int patientId;
  bool _isSending = false;
  bool get isSending => _isSending;

  Timer? _timer;

  OtpProvider({
    required this.otpType,
    required this.emailText,
    required this.patientId,
    required this.otpService,
  });

  // void init() {
  //   sendOtp();
  // }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  Future<bool> validateOtp(String otp) async {
    try {
      final result = await otpService.verifyOtp(patientId, otp.trim());
      final status = result["statusCode"] as int;
      if (status == 200) {
        errorOtp = "";
        notifyListeners();
        return true;
      }
      if (otp.trim().isEmpty) {
        errorOtp = "กรุณากรอกค่า";
      } else if (status == 401) {
        errorOtp = "รหัส OTP ไม่ถูกต้องหรือหมดอายุ";
      } else if (status == 404) {
        errorOtp = "ไม่พบรหัส OTP โปรดขอใหม่อีกครั้ง";
      } else {
        errorOtp = "เกิดข้อผิดพลาด: ${result['body']}";
      }
      notifyListeners();
      return false;
    } catch (e) {
      errorOtp = "เกิดข้อผิดพลาด: ${e.toString()}";
      notifyListeners();
      return false;
    }
  }

  void countDownTime({int second = 60}) {
    countdown = second;
    notifyListeners();

    _timer?.cancel();
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (countdown == 0) {
        timer.cancel();
      } else {
        countdown--;
      }
      notifyListeners();
    });
  }

  Future<void> sendOtp() async {
      _isSending = true;
  notifyListeners();
    try {
      await otpService.requestOtp(patientId);
      countDownTime();
      debugPrint("sendOtp เรียกละ");
    } catch (e) {
      debugPrint("ส่ง OTP ล้มเหลว: ${e.toString()}");
    }finally{
      _isSending = false;
      notifyListeners();
    }
  }
}
