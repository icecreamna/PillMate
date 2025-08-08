import 'package:flutter/widgets.dart';
import 'package:frontend/app/routes/app_pages.dart';
import 'package:get/get.dart';

class ProfileSetupController extends GetxController {
  //TODO: Implement ProfileSetupController

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

  List<String> hasError = List.filled(4, "").obs;

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

  void checkSetUp() {
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
      }else {
        hasError[i] = "";
      }
    }
    if (!hasAnyError) {
      Get.offNamed(Routes.LOGIN);
    }
  }
}
