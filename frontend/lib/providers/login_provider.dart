import 'package:flutter/material.dart';
import 'package:frontend/services/auth_service.dart';

class LoginProvider extends ChangeNotifier {
  bool _obsecurePassword = true;
  bool _isLoading = false;
  bool get obsecurePassword => _obsecurePassword;
    bool get isLoading => _isLoading;


  final AuthService authService;
  LoginProvider(this.authService);

  void toggleObsecurePassword() {
    _obsecurePassword = !_obsecurePassword;
    notifyListeners();
  }

  Future<Map<String, dynamic>?> login({
    required String email,
    required String password,
  }) async {
    _isLoading = true;
    notifyListeners();

    try {
      final res = await authService.login(email: email, password: password);
      return res;
    } catch (e) {
      debugPrint("❌ LoginProvider error: $e");
      return null;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  
}
