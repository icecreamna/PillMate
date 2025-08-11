import 'package:get/get.dart';

class LoginController extends GetxController {
  //TODO: Implement LoginController

  final RxBool _obsecurePassword = true.obs ;
  bool get obsecurePassword => _obsecurePassword.value;
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

  void toggleObsecurePassword() => _obsecurePassword.value =!_obsecurePassword.value;

}
