import 'package:flutter/widgets.dart';
import 'package:frontend/app/routes/app_pages.dart';
import 'package:get/get.dart';

enum OTPType { register, forgot }

class OtpController extends GetxController {
  //TODO: Implement OtpController
  final otpController = TextEditingController();
  late final OTPType otpType;
  RxString errorOtp = "".obs;

  @override
  void onInit() {
    super.onInit();
    otpType = Get.arguments['otpPage'];
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {
    otpController.dispose();
    super.onClose();
  }

  void validateOtp() {
    final otp = otpController.text.trim();
    
    if(otp.isEmpty){
      errorOtp.value = "กรุณากรอกค่า";
    }else if(otp != "123456"){
      errorOtp.value = "รหัส OTP ไม่ถูกต้อง";
    }
    else {
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
}
