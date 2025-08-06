import 'package:flutter/material.dart';
import 'package:frontend/app/modules/otp/controllers/otp_controller.dart';
import 'package:frontend/app/routes/app_pages.dart';
import 'package:frontend/app/widgets/filled_button_custom.dart';
import 'package:frontend/app/widgets/text_field_input.dart';

import 'package:get/get.dart';

import '../controllers/forget_password_controller.dart';

import '../../../../utils/colors.dart' as color;

class ForgetPasswordView extends GetView<ForgetPasswordController> {
  const ForgetPasswordView({super.key});
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
                    height: 191,
                    width: 191,
                  ),
                ),
                const SizedBox(height: 15),
                Container(
                  height: 514,
                  width: double.infinity,
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(20),
                    color: Colors.white,
                  ),
                  child: Column(
                    children: [
                      const SizedBox(height: 30),
                      const Text(
                        "ลืมรหัสผ่าน",
                        style: TextStyle(
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                          letterSpacing: 0,
                        ),
                      ),
                      const SizedBox(height: 10),
                      const Text(
                        "กรุณาใส่อีเมลเพื่อยืนยันตัวตน",
                        style: TextStyle(fontSize: 16),
                      ),
                      const SizedBox(height: 60),
                      const TextFieldInput(
                        labelname: "E-mail",
                        preIcon: Icon(
                          Icons.email_outlined,
                          color: Colors.black,
                        ),
                        textInputType: TextInputType.emailAddress,
                      ),
                      const SizedBox(height: 70),
                      const Padding(
                        padding: EdgeInsets.symmetric(horizontal: 29),
                        child: Text(
                          "หากที่อยู่อีเมลที่คุณระบุมีการลงทะเบียนไว้กับเรา ระบบจะส่งรหัส OTP ไปที่อีเมลของคุณ",
                          style: TextStyle(
                            color: Color(0xFFFF0000),
                            fontSize: 16,
                          ),
                        ),
                      ),
                      const Spacer(),
                      FilledButtonCustom(
                        text: "ถัดไป",
                        onPressed: () => Get.offNamed(
                          Routes.OTP,
                          arguments: {"otpPage": OTPType.forgot},
                        ),
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
                                        fontSize: 20,
                                        color: Colors.black,
                                      ),
                                    ),
                                  ),
                                  TextSpan(
                                    text: "กลับ",
                                    style: TextStyle(
                                      fontSize: 20,
                                      color: Colors.black,
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
              ],
            ),
          ),
        ),
      ),
    );
  }
}
