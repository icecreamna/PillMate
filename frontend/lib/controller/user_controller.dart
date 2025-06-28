import 'package:get/get.dart';

class UserController extends GetxController {
  RxBool obsecurePassword = true.obs;
  RxBool obsecureConfirmPassword = true.obs;

  void toggleObsecurePassword() {
    obsecurePassword.value = !obsecurePassword.value;
  }

  void toggleObsecureConfirmPassword() {
    obsecureConfirmPassword.value = !obsecureConfirmPassword.value;
  }
}
