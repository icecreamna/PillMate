import 'package:flutter/material.dart';

class NewPasswordProvider extends ChangeNotifier {
  bool _obsecurePassword = true;
  bool get obsecurePassword => _obsecurePassword;
  bool _obsecureConfirmPassword = true;
  bool get obsecureConfirmPassword => _obsecureConfirmPassword;

  void toggleObsecurePassword() {
    _obsecurePassword = !_obsecurePassword;
    notifyListeners();
  }

  void toggleObsecureConfirmPassword() {
    _obsecureConfirmPassword = !_obsecureConfirmPassword;
    notifyListeners();
  }
}
