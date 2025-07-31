import 'package:flutter/material.dart';
import 'package:frontend/app/modules/register/controllers/register_controller.dart';
import 'package:frontend/app/routes/app_pages.dart';
import 'package:frontend/app/utils/colors.dart' as color;
import 'package:frontend/app/widgets/text_field_input.dart';
import 'package:get/get.dart';

class RegisterScreen extends GetView<RegisterController> {
  const RegisterScreen({super.key});

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
                  child: Column(
                    children: [
                      const SizedBox(height: 80),

                      const TextFieldInput(
                        labelname: "ID card number",
                        textInputType: TextInputType.number,
                        preIcon: Icon(
                          Icons.person_outline,
                          color: Colors.black,
                        ),
                      ),

                      const SizedBox(height: 25),
                      Obx(
                        () => TextFieldInput(
                          labelname: "Password",
                          textInputType: TextInputType.text,
                          hideText: controller.obsecurePassword,
                          isSuf: true,
                          preIcon: const Icon(
                            Icons.lock_outline,
                            color: Colors.black,
                          ),
                          sufIcon: IconButton(
                            onPressed: () => controller.toggleObsecurePassword(),
                            icon: Icon(
                              controller.obsecurePassword
                                  ? Icons.visibility_off_outlined
                                  : Icons.visibility_outlined,
                              color: Colors.black,
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(height: 25),
                      Obx(
                        () => TextFieldInput(
                          labelname: "Confirm Password",
                          textInputType: TextInputType.text,
                          hideText: controller.obsecureConfirmPassword,
                          isSuf: true,
                          preIcon: const Icon(
                            Icons.lock_outline,
                            color: Colors.black,
                          ),
                          sufIcon: IconButton(
                            onPressed: () => controller.toggleObsecureConfirmPassword(),
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
                      FilledButton(
                        style: FilledButton.styleFrom(
                          minimumSize: const Size(288, 49),
                          backgroundColor: color.AppColors.buttonColor,
                        ),
                        onPressed: () {},
                        child: const Text(
                          "SIGN UP",
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                      const SizedBox(height: 15),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          const Text(
                            "Already member?",
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
                              "Sign In",
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
              ],
            ),
          ),
        ),
      ),
    );
  }
}