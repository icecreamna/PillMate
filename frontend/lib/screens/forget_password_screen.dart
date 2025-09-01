import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:frontend/enums/page_type.dart';
import 'package:frontend/providers/new_password_provider.dart';
import 'package:frontend/screens/login_screen.dart';
import 'package:frontend/screens/otp_screens.dart';
import 'package:frontend/widgets/filled_button_custom.dart';
import 'package:frontend/widgets/text_field_input.dart';
import 'package:provider/provider.dart';

import '../../../../utils/colors.dart' as color;

class ForgetPasswordScreen extends StatelessWidget {
  const ForgetPasswordScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (_) => NewPasswordProvider(),
      child: _ForgetPasswordView(),
    );
  }
}

class _ForgetPasswordView extends StatefulWidget {
  @override
  State<_ForgetPasswordView> createState() => _ForgetPasswordScreenState();
}

class _ForgetPasswordScreenState extends State<_ForgetPasswordView> {
  final _formKey = GlobalKey<FormState>();
  final emailController = TextEditingController();

  @override
  void dispose() {
    emailController.dispose();
    super.dispose();
  }

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
                  child: Form(
                    key: _formKey,
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
                        TextFieldInput(
                          labelname: "E-mail",
                          controller: emailController,
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
                          onPressed: () => _formKey.currentState!.validate()
                              ? Navigator.pushReplacement(
                                  context,
                                  MaterialPageRoute(
                                    builder: (context) => const OtpScreens(),
                                    settings: RouteSettings(
                                      arguments: {
                                        "otpType": PageType.forgot,
                                        "email": emailController.text,
                                      },
                                    ),
                                  ),
                                )
                              : null,
                        ),
                        const SizedBox(height: 30),
                        Padding(
                          padding: const EdgeInsets.fromLTRB(5, 0, 0, 10),
                          child: Align(
                            alignment: Alignment.centerLeft,
                            child: TextButton(
                              onPressed: () => Navigator.pushReplacement(
                                context,
                                MaterialPageRoute(
                                  builder: (context) => const LoginScreen(),
                                ),
                              ),
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
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
