import 'dart:async';

import 'package:flutter/material.dart';
import 'package:frontend/enums/page_type.dart';
import 'package:frontend/screens/forget_password_screen.dart';
import 'package:frontend/screens/new_password_screen.dart';
import 'package:frontend/screens/profile_setup_screen.dart';
import 'package:frontend/screens/register_screen.dart';

class OtpProvider extends ChangeNotifier {

  final PageType otpType;
  String emailText = "";
  String errorOtp = "";
  int countdown = 0;

  Timer? _timer;

  OtpProvider({required this.otpType, required this.emailText});

  void init() {
    sendOtp();
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  bool validateOtp(String otp) {

    if (otp.trim().isEmpty) {
      errorOtp = "กรุณากรอกค่า";
      notifyListeners();
      return false;
    } else if (otp.trim() != "123456") {
      errorOtp = "รหัส OTP ไม่ถูกต้อง";
      notifyListeners();
      return false;
    } else {
      errorOtp = "";
      notifyListeners();
      return true;
    }
  }

  void goNextScreen(BuildContext context) {
    switch (otpType) {
      case PageType.register:
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(builder: (context) => const ProfileSetupScreen()),
        );
        break;
      case PageType.forgot:
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(builder: (context) => const NewPasswordScreen()),
        );
        break;
    }
  }

  void goBackScreen(BuildContext context) {
    switch (otpType) {
      case PageType.register:
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(builder: (context) => const RegisterScreen()),
        );
        break;
      case PageType.forgot:
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(builder: (context) => const ForgetPasswordScreen()),
        );
        break;
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

  void sendOtp() {
    // TODO: ใส่ logic ส่ง OTP จริง
    countDownTime();
  }
}
