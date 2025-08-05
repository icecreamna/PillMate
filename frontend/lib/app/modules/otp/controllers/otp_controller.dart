import 'package:frontend/app/routes/app_pages.dart';
import 'package:get/get.dart';

enum OTPType { register, forgot }

class OtpController extends GetxController {
  //TODO: Implement OtpController 
  late final OTPType otpType ;

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
    super.onClose();
  }
  void submitOtp() {
    if(otpType == OTPType.forgot){
      Get.offNamed(Routes.NEW_PASSWORD);
    }
    else {
      Get.offNamed(Routes.PROFILE_SETUP);
    }
  }
}
