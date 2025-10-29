import 'package:flutter/material.dart';
import 'package:frontend/services/auth_service.dart';

class NewPasswordProvider extends ChangeNotifier {
  bool _obsecurePassword = true;
  bool _obsecureConfirmPassword = true;

  final String email;
  final int patientId;
  bool _isLoading = false;

  bool get obsecurePassword => _obsecurePassword;
  bool get obsecureConfirmPassword => _obsecureConfirmPassword;
  bool get isLoading => _isLoading;

  final AuthService authService;
  NewPasswordProvider({
    required this.authService,
    required this.email,
    required this.patientId,
  });

  void toggleObsecurePassword() {
    _obsecurePassword = !_obsecurePassword;
    notifyListeners();
  }

  void toggleObsecureConfirmPassword() {
    _obsecureConfirmPassword = !_obsecureConfirmPassword;
    notifyListeners();
  }

  Future<Map<String, dynamic>> resetPassword(String newPassword) async {
    _isLoading = true;
    notifyListeners();
    
    try {
      final res = await authService.resetPassword(
        patientId: patientId,
        newPassword: newPassword,
      );
      return res;
    } catch (e) {
      throw Exception("Provider resetPassword failed");
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
