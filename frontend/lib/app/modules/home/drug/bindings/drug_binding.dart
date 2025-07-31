import 'package:get/get.dart';

import '../controllers/drug_controller.dart';

class DrugBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<DrugController>(
      () => DrugController(),
    );
  }
}
