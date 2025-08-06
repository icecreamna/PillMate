import 'package:flutter/material.dart';
import 'package:frontend/app/widgets/text_field_input.dart';

import 'package:get/get.dart';

import '../controllers/new_password_controller.dart';

import '../../../../../utils/colors.dart' as color;

class NewPasswordView extends GetView<NewPasswordController> {
  const NewPasswordView({super.key});
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
                  child: const Column(
                    children: [
                      SizedBox(height: 30),
                      Text(
                        "สร้างรหัสผ่านใหม่",
                        style: TextStyle(
                          color: Colors.black,
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      TextFieldInput(
                        labelname: "Password",
                        textInputType: TextInputType.text,
                        preIcon: Icon(Icons.lock_outline, color: Colors.black),
                      ),
                      TextFieldInput(
                        labelname: "Confirm Password",
                        textInputType: TextInputType.text,
                        preIcon: Icon(Icons.lock_outline, color: Colors.black),
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
