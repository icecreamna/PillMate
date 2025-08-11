import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:frontend/app/widgets/filled_button_custom.dart';
import '../../../routes/app_pages.dart';
import 'package:frontend/app/utils/colors.dart' as color;
import 'package:frontend/app/widgets/text_field_input.dart';
import 'package:get/get.dart';
import '../controllers/login_controller.dart';

final _fromKey = GlobalKey<FormState>();

class LoginScreen extends GetView<LoginController> {
  LoginScreen({super.key});

  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();

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
                    height: 191,
                    width: 191,
                  ),
                ),
                const SizedBox(height: 15),
                Container(
                  height: 514,
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(20),
                    color: Colors.white,
                  ),
                  width: double.infinity,
                  child: SizedBox.expand(
                    child: Form(
                      key: _fromKey,
                      child: Column(
                        children: [
                          const SizedBox(height: 25),
                          const Text(
                            "ยินดีต้อนรับ",
                            style: TextStyle(
                              color: Colors.black,
                              fontWeight: FontWeight.bold,
                              fontSize: 24,
                              letterSpacing: 0,
                            ),
                          ),
                          const SizedBox(height: 10),
                          const Text(
                            "กรุณาเข้าสู่ระบบเพื่อใช้งาน",
                            style: TextStyle(color: Colors.black, fontSize: 16),
                          ),
                          const SizedBox(height: 30),
                          TextFieldInput(
                            labelname: "E-mail",
                            controller: _emailController,
                            inputFormatters: [
                              FilteringTextInputFormatter.deny(RegExp(r'\s')),
                            ],
                            validator: (email) {
                              if (email == null || email.trim().isEmpty) {
                                return "กรุณากรอกค่า";
                              }
                              return null;
                            },
                            preIcon: const Icon(
                              Icons.email_outlined,
                              color: Colors.black,
                            ),
                            textInputType: TextInputType.emailAddress,
                          ),

                          const SizedBox(height: 25),
                          Obx(
                            () => TextFieldInput(
                              labelname: "Password",
                              controller: _passwordController,
                              inputFormatters: [
                                FilteringTextInputFormatter.deny(RegExp(r'\s')),
                              ],
                              validator: (password) {
                                if (password == null ||
                                    password.trim().isEmpty) {
                                  return "กรุณากรอกค่า";
                                }
                                return null;
                              },
                              preIcon: const Icon(
                                Icons.lock_outline,
                                color: Colors.black,
                              ),
                              sufIcon: IconButton(
                                onPressed: () {
                                  controller.toggleObsecurePassword();
                                },
                                icon: Icon(
                                  color: Colors.black,
                                  controller.obsecurePassword
                                      ? Icons.visibility_off_outlined
                                      : Icons.visibility_outlined,
                                ),
                              ),
                              isSuf: true,
                              hideText: controller.obsecurePassword,
                              textInputType: TextInputType.text,
                            ),
                          ),
                          const SizedBox(height: 20),
                          Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 20.0,
                            ),
                            child: Align(
                              alignment: Alignment.centerRight,
                              child: TextButton(
                                onPressed: () =>
                                    Get.offNamed(Routes.FORGET_PASSWORD),
                                child: Text(
                                  "ลืมรหัสผ่าน",
                                  style: TextStyle(
                                    color: color.AppColors.buttonColor,
                                    fontSize: 15,
                                    fontWeight: FontWeight.bold,
                                    decoration: TextDecoration.underline,
                                    decorationColor:
                                        color.AppColors.buttonColor,
                                  ),
                                ),
                              ),
                            ),
                          ),
                          const Spacer(),
                          FilledButtonCustom(
                            text: "เข้าสู่ระบบ",
                            onPressed: () => _fromKey.currentState!.validate()
                                ? Get.offNamed(Routes.HOME)
                                : null,
                          ),
                          const SizedBox(height: 15),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              const Text(
                                "ยังไม่มีบัญชี?",
                                style: TextStyle(
                                  color: Colors.grey,
                                  fontSize: 15,
                                ),
                              ),
                              TextButton(
                                style: TextButton.styleFrom(
                                  padding: const EdgeInsets.all(5),
                                ),
                                onPressed: () {
                                  Get.offAllNamed(Routes.REGISTER);
                                },
                                child: Text(
                                  "ลงทะเบียน",
                                  style: TextStyle(
                                    fontWeight: FontWeight.bold,
                                    fontSize: 15,
                                    color: color.AppColors.buttonColor,
                                    decoration: TextDecoration.underline,
                                    decorationColor:
                                        color.AppColors.buttonColor,
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ],
                      ),
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
