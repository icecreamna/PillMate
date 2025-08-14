import 'package:flutter/material.dart';

class RegisterProvider extends ChangeNotifier {
  bool _obsecurePassword = true;
  bool _obsecureConfirmPassword = true;
  bool get obsecurePassword => _obsecurePassword;
  bool get obsecureConfirmPassword => _obsecureConfirmPassword;

  @override
  void dispose() {
    // TODO: implement dispose
    super.dispose();
  }

  void toggleObsecurePassword() {
    _obsecurePassword = !_obsecurePassword;
    notifyListeners();
  }

  void toggleObsecureConfirmPassword() {
    _obsecureConfirmPassword = !_obsecureConfirmPassword;
    notifyListeners();
  }
}
