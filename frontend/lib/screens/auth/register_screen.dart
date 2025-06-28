import 'package:flutter/material.dart';
import 'package:frontend/controller/user_controller.dart';
import 'package:frontend/screens/auth/login_screen.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/widgets/text_field_input.dart';
import 'package:get/get.dart';

class RegisScreen extends StatelessWidget {
  const RegisScreen({super.key});

  @override
  Widget build(BuildContext context) {
    UserController uc = Get.put(UserController());
    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor1st,
      body: SafeArea(
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
                      preIcon: Icon(Icons.people_outlined, color: Colors.black),
                      sufIcon: Icon(null),
                    ),
                    const SizedBox(height: 25),
                    Obx(
                      () => TextFieldInput(
                        labelname: "Password",
                        preIcon: const Icon(
                          Icons.lock_outline,
                          color: Colors.black,
                        ),
                        sufIcon: IconButton(
                          onPressed: () => uc.toggleObsecurePassword(),
                          icon: Icon(
                            uc.obsecurePassword.value
                                ? Icons.visibility_off_outlined
                                : Icons.visibility_outlined,
                            color: Colors.black,
                          ),
                        ),
                        hideText: uc.obsecurePassword.value,
                        isSuf: true,
                      ),
                    ),
                    const SizedBox(height: 25),
                    Obx(
                      () => TextFieldInput(
                        labelname: "Confirm Password",
                        preIcon: const Icon(
                          Icons.lock_outline,
                          color: Colors.black,
                        ),
                        sufIcon: IconButton(
                          onPressed: () => uc.toggleObsecureConfirmPassword(),
                          icon: Icon(
                            uc.obsecureConfirmPassword.value
                                ? Icons.visibility_off_outlined
                                : Icons.visibility_outlined,
                            color: Colors.black,
                          ),
                        ),
                        hideText: uc.obsecureConfirmPassword.value,
                        isSuf: true,
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
                        style: TextStyle(color: Colors.white, fontSize: 24),
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
                            Get.off(() => const LoginScreen());
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
    );
  }
}
