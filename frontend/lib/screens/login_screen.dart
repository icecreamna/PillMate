import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:frontend/providers/login_provider.dart';
import 'package:frontend/screens/forget_password_screen.dart';
import 'package:frontend/screens/register_screen.dart';
import 'package:frontend/screens/splash_screen.dart';
import 'package:frontend/widgets/text_field_input.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/widgets/filled_button_custom.dart';
import 'package:provider/provider.dart';

class LoginScreen extends StatelessWidget {

  const LoginScreen({super.key});
   @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (_) => LoginProvider(),
      child: _LoginView(), // แยก View ออกมาเพื่อไม่ให้ provider ถูกสร้างซ้ำ
    );
  }
}

class _LoginView extends StatefulWidget{
  @override
  State<_LoginView> createState() => _LoginViewState();
}

class _LoginViewState extends State<_LoginView> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

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
                      key: _formKey,
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
                          Selector<LoginProvider, bool>(
                            selector: (_, p) => p.obsecurePassword,
                            builder: (_, obsecure, _) {
                              return TextFieldInput(
                                labelname: "Password",
                                controller: _passwordController,
                                inputFormatters: [
                                  FilteringTextInputFormatter.deny(
                                    RegExp(r'\s'),
                                  ),
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
                                  onPressed: () => context
                                      .read<LoginProvider>()
                                      .toggleObsecurePassword(),
                                  icon: Icon(
                                    color: Colors.black,
                                    obsecure
                                        ? Icons.visibility_off_outlined
                                        : Icons.visibility_outlined,
                                  ),
                                ),
                                isSuf: true,
                                hideText: obsecure,
                                textInputType: TextInputType.text,
                              );
                            },
                          ),
                          const SizedBox(height: 20),
                          Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 20.0,
                            ),
                            child: Align(
                              alignment: Alignment.centerRight,
                              child: TextButton(
                                onPressed: () => Navigator.pushReplacement(
                                  context,
                                  MaterialPageRoute(
                                    builder: (context) =>
                                        ForgetPasswordScreen(),
                                  ),
                                ),
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
                            onPressed: () => _formKey.currentState!.validate()
                                ? Navigator.pushReplacement(
                                    context,
                                    MaterialPageRoute(
                                      builder: (context) =>
                                          const SplashScreen(),
                                    ),
                                  )
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
                                onPressed: () => Navigator.pushReplacement(
                                  context,
                                  MaterialPageRoute(
                                    builder: (context) =>
                                        const RegisterScreen(),
                                  ),
                                ),
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
