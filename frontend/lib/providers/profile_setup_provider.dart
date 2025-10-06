import 'package:flutter/material.dart';
import 'package:frontend/screens/login_screen.dart';

class ProfileSetupProvider extends ChangeNotifier {
  final firstnameController = TextEditingController();
  final lastnameController = TextEditingController();
  final idcardController = TextEditingController();
  final phoneController = TextEditingController();

  late final List<TextEditingController> controllers = [
    firstnameController,
    lastnameController,
    idcardController,
    phoneController,
  ];

  List<String> hasError = List.filled(4, "");

  void checkSetUp(BuildContext context) {
    bool hasAnyError = false;
    for (int i = 0; i < controllers.length; i++) {
      if (controllers[i].text.trim().isEmpty) {
        hasError[i] = "กรุณากรอกค่า";
        hasAnyError = true;
      } else if (i == 2 && controllers[2].text.trim().length != 13) {
        hasError[2] = "กรุณากรอกให้ครบ 13 ตัว";
        hasAnyError = true;
      } else if (i == 3 && controllers[3].text.trim().length != 10) {
        hasError[3] = "กรุณากรอกให้ครบ 10 ตัว";
        hasAnyError = true;
      } else {
        hasError[i] = "";
      }
    }
    notifyListeners();
    if (!hasAnyError) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (context) => const LoginScreen()),
      );
    }
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
