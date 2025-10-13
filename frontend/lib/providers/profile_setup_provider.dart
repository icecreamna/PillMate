import 'package:flutter/material.dart';
import 'package:frontend/screens/login_screen.dart';
import 'package:frontend/services/profile_service.dart';

class ProfileSetupProvider extends ChangeNotifier {
  final int patientId;
  final ProfileService profileService;

  ProfileSetupProvider({required this.profileService, required this.patientId});

  final firstnameController = TextEditingController();
  final lastnameController = TextEditingController();
  final idcardController = TextEditingController();
  final phoneController = TextEditingController();

  List<String> hasError = List.filled(4, "");
  bool _isLoading = false;
  bool get isLoading => _isLoading;

  Future<void> checkSetUp(BuildContext context) async {
    bool hasAnyError = false;

    for (int i = 0; i < 4; i++) {
      final text = [
        firstnameController,
        lastnameController,
        idcardController,
        phoneController,
      ][i].text.trim();
      if (text.isEmpty) {
        hasError[i] = "กรุณากรอกค่า";
        hasAnyError = true;
      } else if (i == 2 && text.length != 13) {
        hasError[i] = "กรุณากรอกให้ครบ 13 ตัว";
        hasAnyError = true;
      } else if (i == 3 && text.length != 10) {
        hasError[i] = "กรุณากรอกให้ครบ 10 ตัว";
        hasAnyError = true;
      } else {
        hasError[i] = "";
      }
    }
    notifyListeners();

    if (hasAnyError) return;

    _isLoading = true;
    notifyListeners();
    final success = await profileService.setUpProfile(
      patientId: patientId,
      firstName: firstnameController.text.trim(),
      lastName: lastnameController.text.trim(),
      idCardNumber: idcardController.text.trim(),
      phoneNumber: phoneController.text.trim(),
    );

    if (success) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text(
            "สมัครเสร็จสิ้น",
            style: TextStyle(color: Colors.white),
          ),
          backgroundColor: Colors.green,
          behavior: SnackBarBehavior.floating,
          duration: Duration(seconds: 2),
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text(
            "สมัครไม่ผ่าน",
            style: TextStyle(color: Colors.white),
          ),
          backgroundColor: Colors.red,
          behavior: SnackBarBehavior.floating,
          duration: Duration(seconds: 2),
        ),
      );
    }

    Navigator.pushReplacement(
      context,
      MaterialPageRoute(builder: (context) => const LoginScreen()),
    );
  }

  @override
  void dispose() {
    firstnameController.dispose();
    lastnameController.dispose();
    idcardController.dispose();
    phoneController.dispose();
    super.dispose();
  }
}
