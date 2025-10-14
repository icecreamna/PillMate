import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:frontend/enums/page_type.dart';
import 'package:frontend/providers/register_provider.dart';
import 'package:frontend/screens/login_screen.dart';
import 'package:frontend/screens/otp_screens.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/widgets/filled_button_custom.dart';
import 'package:frontend/widgets/text_field_input.dart';
import 'package:provider/provider.dart';

class RegisterScreen extends StatelessWidget {
  const RegisterScreen({super.key});

  @override
  Widget build(BuildContext context) {
    // TODO: implement build
    return ChangeNotifierProvider(
      create: (_) => RegisterProvider(AuthService()),
      child: _RegisterView(),
    );
  }
}

class _RegisterView extends StatefulWidget {
  @override
  State<_RegisterView> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends State<_RegisterView> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _emailController;
  late final TextEditingController _passwordController;
  late final TextEditingController _confirmPasswordController;

  @override
  void initState() {
    super.initState();
    _emailController = TextEditingController();
    _passwordController = TextEditingController();
    _confirmPasswordController = TextEditingController();
  }

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  Future<void> _onSubmit() async {
    if (!_formKey.currentState!.validate()) return;
    final provider = context.read<RegisterProvider>();

    try {
      final res =await provider.register(
        email: _emailController.text.trim(),
        password: _passwordController.text,
      );

      final patienId = res["patient_id"];
      if (!mounted) return;
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(
          builder: (context) => const OtpScreens(),
          settings: RouteSettings(
            arguments: {
              "otpType": PageType.register,
              "email": _emailController.text,
              "patient_id":patienId
            },
          ),
        ),
      );
    } catch (e) {
      if (!mounted) return;
      showDialog(
        context: context,
        builder: (_) => AlertDialog(
          title: const Text("สมัครสมาชิกไม่สำเร็จ"),
          content: Text(e.toString()),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text("ปิด"),
            ),
          ],
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final isLoading = context.watch<RegisterProvider>().isLoading;

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
                    key: _formKey,
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
                            final emailRegex = RegExp(
                              r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$',
                            );
                            if (!emailRegex.hasMatch(email.trim())) {
                              return "รูปแบบอีเมลไม่ถูกต้อง";
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
                        Consumer<RegisterProvider>(
                          builder: (context, register, _) {
                            return Column(
                              children: [
                                TextFieldInput(
                                  labelname: "Password",
                                  controller: _passwordController,
                                  textInputType: TextInputType.text,
                                  hideText: register.obsecurePassword,
                                  inputFormatters: [
                                    FilteringTextInputFormatter.deny(
                                      RegExp(r'\s'),
                                    ),
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
                                        register.toggleObsecurePassword(),
                                    icon: Icon(
                                      register.obsecurePassword
                                          ? Icons.visibility_off_outlined
                                          : Icons.visibility_outlined,
                                      color: Colors.black,
                                    ),
                                  ),
                                ),
                                const SizedBox(height: 15),

                                TextFieldInput(
                                  labelname: "Confirm Password",
                                  controller: _confirmPasswordController,
                                  textInputType: TextInputType.text,
                                  hideText: register.obsecureConfirmPassword,
                                  inputFormatters: [
                                    FilteringTextInputFormatter.deny(
                                      RegExp(r'\s'),
                                    ),
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
                                    onPressed: () => register
                                        .toggleObsecureConfirmPassword(),
                                    icon: Icon(
                                      register.obsecureConfirmPassword
                                          ? Icons.visibility_off_outlined
                                          : Icons.visibility_outlined,
                                      color: Colors.black,
                                    ),
                                  ),
                                ),
                              ],
                            );
                          },
                        ),
                        const Spacer(),
                        FilledButtonCustom(
                          text: isLoading ? "กำลังส่ง..." : "ถัดไป",
                          onPressed: () => isLoading ? null : _onSubmit(),
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
                                Navigator.pushReplacement(
                                  context,
                                  MaterialPageRoute(
                                    builder: (context) => const LoginScreen(),
                                  ),
                                );
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
