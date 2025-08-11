import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:frontend/app/routes/app_pages.dart';
import 'package:frontend/app/widgets/filled_button_custom.dart';
import 'package:frontend/app/widgets/text_field_input.dart';

import 'package:get/get.dart';

import '../controllers/profile_setup_controller.dart';

import '../../../utils/colors.dart' as color;

final _fromKey = GlobalKey<FormState>();

class ProfileSetupView extends GetView<ProfileSetupController> {
  const ProfileSetupView({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor1st,
      body: SafeArea(
        child: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.all(20.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                SizedBox(
                  width: double.infinity,
                  child: Image.asset(
                    "assets/images/drugs.png",
                    width: 191,
                    height: 191,
                  ),
                ),
                const SizedBox(height: 15),
                Container(
                  width: double.infinity,
                  height: 514,
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(20),
                    color: Colors.white,
                  ),
                  child: Form(
                    key: _fromKey,
                    child: Column(
                      children: [
                        const SizedBox(height: 30),
                        const Text(
                          "ข้อมูลผู้ใช้",
                          style: TextStyle(
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        TextFieldInput(
                          labelname: "ชื่อ",
                          textInputType: TextInputType.text,
                          controller: controller.firstnameController,
                          inputFormatters: [
                            FilteringTextInputFormatter.allow(
                              RegExp(r'[a-zA-zก-ฮ]'),
                            ),
                          ],
                        ),
                        Obx(
                          () => Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 27),
                            child: Align(
                              alignment: Alignment.centerLeft,
                              child: Visibility(
                                visible: controller.hasError[0].isNotEmpty,
                                child: Text(
                                  controller.hasError[0],
                                  style: const TextStyle(
                                    fontSize: 12,
                                    color: Color(0xFFFF0000),
                                  ),
                                ),
                              ),
                            ),
                          ),
                        ),
                        const SizedBox(height: 5),
                        TextFieldInput(
                          labelname: "นามสกุล",
                          textInputType: TextInputType.text,
                          controller: controller.lastnameController,
                          inputFormatters: [
                            FilteringTextInputFormatter.allow(
                              RegExp(r'[a-zA-zก-ฮ]'),
                            ),
                          ],
                        ),
                        Obx(
                          () => Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 27),
                            child: Align(
                              alignment: Alignment.centerLeft,
                              child: Visibility(
                                visible: controller.hasError[1].isNotEmpty,
                                child: Text(
                                  controller.hasError[1],
                                  style: const TextStyle(
                                    fontSize: 12,
                                    color: Color(0xFFFF0000),
                                  ),
                                ),
                              ),
                            ),
                          ),
                        ),
                        const SizedBox(height: 5),
                        TextFieldInput(
                          labelname: "เลขบัตรประชาชน",
                          textInputType: TextInputType.number,
                          controller: controller.idcardController,
                          inputFormatters: [
                            FilteringTextInputFormatter.allow(RegExp(r'\d')),
                            LengthLimitingTextInputFormatter(13),
                          ],
                        ),
                        Obx(
                          () => Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 27),
                            child: Align(
                              alignment: Alignment.centerLeft,
                              child: Visibility(
                                visible: controller.hasError[2].isNotEmpty,
                                child: Text(
                                  controller.hasError[2],
                                  style: const TextStyle(
                                    fontSize: 12,
                                    color: Color(0xFFFF0000),
                                  ),
                                ),
                              ),
                            ),
                          ),
                        ),
                        const SizedBox(height: 5),
                        TextFieldInput(
                          labelname: "เบอร์โทรศัพท์",
                          textInputType: TextInputType.number,
                          controller: controller.phoneController,
                          inputFormatters: [
                            FilteringTextInputFormatter.allow(RegExp(r'[0-9]')),
                            LengthLimitingTextInputFormatter(10),
                          ],
                        ),
                        Obx(
                          () => Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 27),
                            child: Align(
                              alignment: Alignment.centerLeft,
                              child: Visibility(
                                visible: controller.hasError[3].isNotEmpty,
                                child: Text(
                                  controller.hasError[3],
                                  style: const TextStyle(
                                    fontSize: 12,
                                    color: Color(0xFFFF0000),
                                  ),
                                ),
                              ),
                            ),
                          ),
                        ),
                        const Spacer(),
                        FilledButtonCustom(
                          text: "ลงทะเบียน",
                          onPressed: () => controller.checkSetUp(),
                        ),
                        const SizedBox(height: 30),
                        Padding(
                          padding: const EdgeInsets.fromLTRB(5, 0, 0, 10),
                          child: Align(
                            alignment: Alignment.centerLeft,
                            child: TextButton(
                              onPressed: () => Get.offNamed(Routes.REGISTER),
                              child: const Text.rich(
                                TextSpan(
                                  children: [
                                    WidgetSpan(
                                      alignment: PlaceholderAlignment.middle,
                                      child: Text(
                                        "<",
                                        style: TextStyle(
                                          fontSize: 20,
                                          color: Colors.black,
                                        ),
                                      ),
                                    ),
                                    TextSpan(
                                      text: "กลับ",
                                      style: TextStyle(
                                        color: Colors.black,
                                        fontSize: 20,
                                        letterSpacing: 0,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
