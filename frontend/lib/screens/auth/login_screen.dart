import 'package:flutter/material.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/widgets/text_field_input.dart';
import 'package:get/get.dart';
import 'package:frontend/screens/auth/register_screen.dart';

class LoginScreen extends StatelessWidget {
  const LoginScreen({super.key});

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
                    child: Column(
                      children: [
                        const SizedBox(height: 12),
                        const Text(
                          "Welcome",
                          style: TextStyle(
                            color: Color(0xFF0B0870),
                            fontWeight: FontWeight.bold,
                            fontSize: 24,
                          ),
                        ),
                        const SizedBox(height: 10),
                        const Text(
                          "Please log in to continue",
                          style: TextStyle(
                            color: Color(0xFF0B0870),
                            fontWeight: FontWeight.bold,
                            fontSize: 13,
                          ),
                        ),
                        const SizedBox(height: 30),
                        const TextFieldInput(
                          labelname: "ID card number",
                          preIcon: Icon(
                            Icons.people_outlined,
                            color: Colors.black,
                          ),
                          sufIcon: Icon(
                            Icons.visibility_off_outlined,
                            color: Colors.black,
                          ),
                        ),
                        const SizedBox(height: 25),
                        const TextFieldInput(
                          labelname: "Password",
                          preIcon: Icon(
                            Icons.lock_open_outlined,
                            color: Colors.black,
                          ),
                          sufIcon: Icon(
                            Icons.visibility_off_outlined,
                            color: Colors.black,
                          ),
                          isEye: true,
                          hideText: true,
                        ),
                        const Spacer(),
                        FilledButton(
                          onPressed: () {},
                          style: FilledButton.styleFrom(
                            minimumSize: const Size(288, 49),
                            backgroundColor: color.AppColors.buttonColor,
                          ),
                          child: const Text(
                            "SIGN IN",
                            style: TextStyle(color: Colors.white, fontSize: 24),
                          ),
                        ),
                        const SizedBox(height: 15),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Text(
                              "Don't have an account?",
                              style: TextStyle(
                                color: Colors.grey,
                                fontSize: 15,
                              ),
                            ),
                            TextButton(
                              style: TextButton.styleFrom(
                                padding: const EdgeInsets.all(5),
                                minimumSize: const Size(0, 0),
                              ),
                              onPressed: () {
                                 Get.to(()=> const RegisScreen());
                              },
                              child: Text(
                                "Sign Up",
                                style: TextStyle(
                                  fontWeight: FontWeight.bold,
                                  fontSize: 15,
                                  color: color.AppColors.buttonColor,
                                ),
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 20),
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
