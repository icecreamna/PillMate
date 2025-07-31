import 'package:get/get.dart';

class RegisterController extends GetxController {
  //TODO: Implement RegisterController

  final RxBool _obsecurePassword = true.obs;
  final RxBool _obsecureConfirmPassword = true.obs;
  bool get obsecurePassword => _obsecurePassword.value;
  bool get obsecureConfirmPassword => _obsecureConfirmPassword.value;
  @override
  void onInit() {
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {
    super.onClose();
  }

  void toggleObsecurePassword() =>_obsecurePassword.value = !_obsecurePassword.value;
  void toggleObsecureConfirmPassword() => _obsecureConfirmPassword.value = !_obsecureConfirmPassword.value;
}
