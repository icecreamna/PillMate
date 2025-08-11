import 'package:flutter/material.dart';
import 'package:frontend/app/modules/otp/controllers/otp_controller.dart';
import 'package:frontend/app/routes/app_pages.dart';
import 'package:frontend/app/widgets/filled_button_custom.dart';
import 'package:frontend/app/widgets/text_field_input.dart';

import 'package:get/get.dart';
//
import '../controllers/new_password_controller.dart';

import '../../../../../utils/colors.dart' as color;

final _fromKey = GlobalKey<FormState>();

class NewPasswordView extends GetView<NewPasswordController> {
  NewPasswordView({super.key});

  final OtpController otpController = Get.find<OtpController>();

  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();

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
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: Form(
                    key: _fromKey,
                    child: Column(
                      children: [
                        const SizedBox(height: 30),
                        const Text(
                          "สร้างรหัสผ่านใหม่",
                          style: TextStyle(
                            color: Colors.black,
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 30),
                         Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 42),
                          child: Row(
                            children: [
                              const Icon(Icons.email_outlined, color: Colors.black),
                              const SizedBox(width: 15),
                              Text(
                                otpController.emailText.value,
                                style: const TextStyle(
                                  color: Colors.black,
                                  fontSize: 16,
                                ),
                              ),
                            ],
                          ),
                        ),
                        const SizedBox(height: 15),
                        Obx(
                          () => TextFieldInput(
                            labelname: "Password",
                            textInputType: TextInputType.text,
                            hideText: controller.obsecurePassword,
                            controller: _passwordController,
                            validator: (password) {
                              if (password == null || password.trim().isEmpty) {
                                return "กรุณากรอกค่า";
                              } else if (password.length < 6) {
                                return "รหัสควรมีความยาวมากกว่าเท่ากับ 6 ตัว";
                              }
                              return null;
                            },
                            preIcon: const Icon(
                              Icons.lock_outline,
                              color: Colors.black,
                            ),
                            isSuf: true,
                            sufIcon: IconButton(
                              onPressed: () =>
                                  controller.toggleObsecurePassword(),
                              icon: Icon(
                                controller.obsecurePassword
                                    ? Icons.visibility_off_outlined
                                    : Icons.visibility_outlined,
                                color: Colors.black,
                              ),
                            ),
                          ),
                        ),
                        const SizedBox(height: 15),
                        Obx(
                          () => TextFieldInput(
                            labelname: "Confirm Password",
                            textInputType: TextInputType.text,
                            hideText: controller.obsecureConfirmPassword,
                            controller: _confirmPasswordController,
                            validator: (cfmp) {
                              if (cfmp == null || cfmp.trim().isEmpty) {
                                return "กรุณากรอกค่า";
                              } else if (cfmp != _passwordController.text) {
                                return "รหัสไม่ตรงกัน";
                              }
                              return null;
                            },
                            preIcon: const Icon(
                              Icons.lock_outline,
                              color: Colors.black,
                            ),
                            isSuf: true,
                            sufIcon: IconButton(
                              onPressed: () =>
                                  controller.toggleObsecureConfirmPassword(),
                              icon: Icon(
                                controller.obsecureConfirmPassword
                                    ? Icons.visibility_off_outlined
                                    : Icons.visibility_outlined,
                                color: Colors.black,
                              ),
                            ),
                          ),
                        ),
                        const Spacer(),
                        FilledButtonCustom(
                          text: "ตกลง",
                          onPressed: () => _fromKey.currentState!.validate()
                              ? Get.offNamed(Routes.LOGIN)
                              : null,
                        ),
                        const SizedBox(height: 30),
                        Padding(
                          padding: const EdgeInsets.fromLTRB(5, 0, 0, 10),
                          child: Align(
                            alignment: Alignment.centerLeft,
                            child: TextButton(
                              onPressed: () => Get.offNamed(Routes.LOGIN),
                              child: const Text.rich(
                                TextSpan(
                                  children: [
                                    WidgetSpan(
                                      alignment: PlaceholderAlignment.middle,
                                      child: Text(
                                        "<",
                                        style: TextStyle(
                                          color: Colors.black,
                                          fontSize: 20,
                                        ),
                                      ),
                                    ),
                                    TextSpan(
                                      text: "กลับ",
                                      style: TextStyle(
                                        color: Colors.black,
                                        letterSpacing: 0,
                                        fontSize: 20,
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
