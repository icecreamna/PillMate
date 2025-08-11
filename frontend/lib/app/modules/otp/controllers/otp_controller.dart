import 'dart:async';

import 'package:flutter/widgets.dart';
import 'package:frontend/app/modules/login/forget_password/controllers/forget_password_controller.dart';
import 'package:frontend/app/routes/app_pages.dart';
import 'package:get/get.dart';

import '../../register/controllers/register_controller.dart';

enum OTPType { register, forgot }

class OtpController extends GetxController {
  //TODO: Implement OtpController

  final otpController = TextEditingController();
  late final OTPType otpType;
  RxString emailText = "".obs ; 
  RxString errorOtp = "".obs;
  RxInt countdown = 0.obs;
  // bool isCount = false;
  Timer? _timer;

  @override
  void onInit() {
    super.onInit();
    otpType = Get.arguments['otpPage'];
    sendOtp();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {
    _timer?.cancel();
    super.onClose();
  }

  void validateOtp() {
    final otp = otpController.text.trim();

    if (otp.isEmpty) {
      errorOtp.value = "กรุณากรอกค่า";
    } else if (otp != "123456") {
      errorOtp.value = "รหัส OTP ไม่ถูกต้อง";
    } else {
      errorOtp.value = "";
      goNextScreen();
    }
  }

  void goNextScreen() {
    switch (otpType) {
      case OTPType.register:
        Get.offNamed(Routes.PROFILE_SETUP);
        break;
      case OTPType.forgot:
        Get.offNamed(Routes.NEW_PASSWORD);
        break;
    }
  }

  void goBackScreen() {
    switch (otpType) {
      case OTPType.register:
        Get.offNamed(Routes.REGISTER);
        break;
      case OTPType.forgot:
        Get.offNamed(Routes.FORGET_PASSWORD);
        break;
    }
  }

  void countDownTime({int second = 60}) {
    countdown.value = second;
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (countdown.value == 0) {
        // isCount = false;
        timer.cancel();
      } else {
        // isCount = true;
        countdown.value--;
      }
    });
  }

  void sendOtp() {
    //sendOtplogic
    countDownTime();
  }

  String get emailShow {
    switch (otpType) {
      case OTPType.forgot:
        final ForgetPasswordController forgetPasswordController =
            Get.find<ForgetPasswordController>();
            emailText.value = forgetPasswordController.emailController.text;
        return forgetPasswordController.emailController.text;
      case OTPType.register:
        final RegisterController registerController =
            Get.find<RegisterController>();
        return registerController.emailController.text;
    }
  }
}
