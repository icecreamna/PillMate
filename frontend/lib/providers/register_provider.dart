import 'package:flutter/material.dart';
import 'package:frontend/services/auth_service.dart';

class RegisterProvider extends ChangeNotifier {
  final AuthService authService;
  RegisterProvider(this.authService);

  bool _obsecurePassword = true;
  bool _obsecureConfirmPassword = true;
  bool _isLoading = false;
  bool get obsecurePassword => _obsecurePassword;
  bool get obsecureConfirmPassword => _obsecureConfirmPassword;
  bool get isLoading => _isLoading;

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

  Future<Map<String, dynamic>> register({
    required String email,
    required String password,
  }) async {
    _isLoading = true;
    notifyListeners();

    try {
      final res = await authService.register(email: email, password: password);
      return res;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
