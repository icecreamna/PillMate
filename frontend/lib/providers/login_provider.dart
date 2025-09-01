import 'package:flutter/material.dart';

class LoginProvider extends ChangeNotifier {
  bool _obsecurePassword = true;
  bool get obsecurePassword => _obsecurePassword;

  void toggleObsecurePassword() {
    _obsecurePassword = !_obsecurePassword;
    notifyListeners();
  }
}
