import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:frontend/app/modules/otp/controllers/otp_controller.dart';
import 'package:frontend/app/modules/register/controllers/register_controller.dart';
import 'package:frontend/app/routes/app_pages.dart';
import 'package:frontend/app/utils/colors.dart' as color;
import 'package:frontend/app/widgets/filled_button_custom.dart';
import 'package:frontend/app/widgets/text_field_input.dart';
import 'package:get/get.dart';

final _fromKey = GlobalKey<FormState>();

class RegisterScreen extends GetView<RegisterController> {
  RegisterScreen({super.key});

  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor1st,
      body: SafeArea(
        child: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.all(20),
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
                          "ลงทะเบียน",
                          style: TextStyle(
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                            letterSpacing: 0,
                          ),
                        ),
                        const SizedBox(height: 30),
                        TextFieldInput(
                          labelname: "E-mail",
                          textInputType: TextInputType.emailAddress,
                          inputFormatters: [
                            FilteringTextInputFormatter.deny(RegExp(r'\s')),
                          ],
                          validator: (email) {
                            if (email == null || email.trim().isEmpty) {
                              return "กรุณากรอกค่า";
                            }
                            return null;
                          },
                          controller: _emailController,
                          preIcon: const Icon(
                            Icons.email_outlined,
                            color: Colors.black,
                          ),
                        ),

                        const SizedBox(height: 15),
                        Obx(
                          () => TextFieldInput(
                            labelname: "Password",
                            controller: _passwordController,
                            textInputType: TextInputType.text,
                            hideText: controller.obsecurePassword,
                            inputFormatters: [
                              FilteringTextInputFormatter.deny(RegExp(r'\s')),
                            ],
                            validator: (pwd) {
                              if (pwd == null || pwd.trim().isEmpty) {
                                return "กรุณากรอกค่า";
                              } else if (pwd.trim().length < 6) {
                                return "รหัสควรมีความยาวมากกว่าเท่ากับ 6 ตัว";
                              }
                              return null;
                            },
                            isSuf: true,
                            preIcon: const Icon(
                              Icons.lock_outline,
                              color: Colors.black,
                            ),
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
                            controller: _confirmPasswordController,
                            textInputType: TextInputType.text,
                            hideText: controller.obsecureConfirmPassword,
                            inputFormatters: [
                              FilteringTextInputFormatter.deny(RegExp(r'\s')),
                            ],
                            validator: (cmfp) {
                              if (cmfp == null || cmfp.trim().isEmpty) {
                                return "กรุณากรอกค่า";
                              } else if (cmfp.trim() !=
                                  _passwordController.text) {
                                return "รหัสไม่ตรงกัน";
                              }
                              return null;
                            },
                            isSuf: true,
                            preIcon: const Icon(
                              Icons.lock_outline,
                              color: Colors.black,
                            ),
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
                          text: "ถัดไป",
                          onPressed: () => _fromKey.currentState!.validate()
                              ? Get.offNamed(
                                  Routes.OTP,
                                  arguments: {"otpPage": OTPType.register},
                                )
                              : null,
                        ),
                        const SizedBox(height: 15),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Text(
                              "มีบัญชีอยู่แล้ว?",
                              style: TextStyle(color: Colors.grey),
                            ),
                            const SizedBox(width: 5),
                            TextButton(
                              style: TextButton.styleFrom(
                                padding: const EdgeInsets.all(0),
                                minimumSize: const Size(0, 0),
                              ),
                              onPressed: () {
                                Get.offAllNamed(Routes.LOGIN);
                              },
                              child: Text(
                                "เข้าสู่ระบบ",
                                style: TextStyle(
                                  color: color.AppColors.buttonColor,
                                  fontWeight: FontWeight.bold,
                                  fontSize: 15,
                                ),
                              ),
                            ),
                          ],
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
