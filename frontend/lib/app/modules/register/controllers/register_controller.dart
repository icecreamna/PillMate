import 'package:flutter/material.dart';
import 'package:get/get.dart';

class RegisterController extends GetxController {
  //TODO: Implement RegisterController

  final RxBool _obsecurePassword = true.obs;
  final RxBool _obsecureConfirmPassword = true.obs;
  bool get obsecurePassword => _obsecurePassword.value;
  bool get obsecureConfirmPassword => _obsecureConfirmPassword.value;
  final emailController = TextEditingController();
  final passwordController = TextEditingController();
  final confirmPasswordController = TextEditingController();

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
    emailController.dispose();
    passwordController.dispose();
    confirmPasswordController.dispose();
    super.onClose();
  }

  void toggleObsecurePassword() =>
      _obsecurePassword.value = !_obsecurePassword.value;
  void toggleObsecureConfirmPassword() =>
      _obsecureConfirmPassword.value = !_obsecureConfirmPassword.value;
}
