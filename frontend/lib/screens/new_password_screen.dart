import 'package:flutter/material.dart';
import 'package:frontend/providers/new_password_provider.dart';
import 'package:frontend/screens/login_screen.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/widgets/filled_button_custom.dart';
import 'package:frontend/widgets/text_field_input.dart';
import 'package:provider/provider.dart';

class NewPasswordScreen extends StatelessWidget {
  const NewPasswordScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final args =
        ModalRoute.of(context)!.settings.arguments as Map<String, dynamic>;
    final String email = args['email'];
    final int patientId = args["patient_id"];

    return ChangeNotifierProvider(
      create: (_) => NewPasswordProvider(
        authService: AuthService(),
        email: email,
        patientId: patientId,
      ),
      child: _NewPasswordView(),
    );
  }
}

class _NewPasswordView extends StatefulWidget {
  @override
  State<_NewPasswordView> createState() => _NewPasswordScreenState();
}

class _NewPasswordScreenState extends State<_NewPasswordView> {
  final _formKey = GlobalKey<FormState>();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();

  @override
  void dispose() {
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final npp = context.read<NewPasswordProvider>();
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
                  child: Form(
                    key: _formKey,
                    child: Column(
                      children: [
                        const SizedBox(height: 30),
                        const Text(
                          "สร้างรหัสผ่านใหม่",
                          style: TextStyle(
                            color: Colors.black,
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 30),
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 42),
                          child: Row(
                            children: [
                              const Icon(
                                Icons.email_outlined,
                                color: Colors.black,
                              ),
                              const SizedBox(width: 15),
                              Text(
                                npp.email,
                                style: const TextStyle(
                                  color: Colors.black,
                                  fontSize: 16,
                                ),
                              ),
                            ],
                          ),
                        ),
                        const SizedBox(height: 15),
                        Consumer<NewPasswordProvider>(
                          builder: (_, pw, _) {
                            return TextFieldInput(
                              labelname: "Password",
                              textInputType: TextInputType.text,
                              hideText: pw.obsecurePassword,
                              controller: _passwordController,
                              validator: (password) {
                                if (password == null ||
                                    password.trim().isEmpty) {
                                  return "กรุณากรอกค่า";
                                } else if (password.length < 6) {
                                  return "รหัสควรมีความยาวมากกว่าเท่ากับ 6 ตัว";
                                }
                                return null;
                              },
                              preIcon: const Icon(
                                Icons.lock_outline,
                                color: Colors.black,
                              ),
                              isSuf: true,
                              sufIcon: IconButton(
                                onPressed: () => pw.toggleObsecurePassword(),
                                icon: Icon(
                                  pw.obsecurePassword
                                      ? Icons.visibility_off_outlined
                                      : Icons.visibility_outlined,
                                  color: Colors.black,
                                ),
                              ),
                            );
                          },
                        ),
                        const SizedBox(height: 15),
                        Consumer<NewPasswordProvider>(
                          builder: (_, cfm, _) {
                            return TextFieldInput(
                              labelname: "Confirm Password",
                              textInputType: TextInputType.text,
                              hideText: cfm.obsecureConfirmPassword,
                              controller: _confirmPasswordController,
                              validator: (cfmp) {
                                if (cfmp == null || cfmp.trim().isEmpty) {
                                  return "กรุณากรอกค่า";
                                } else if (cfmp != _passwordController.text) {
                                  return "รหัสไม่ตรงกัน";
                                }
                                return null;
                              },
                              preIcon: const Icon(
                                Icons.lock_outline,
                                color: Colors.black,
                              ),
                              isSuf: true,
                              sufIcon: IconButton(
                                onPressed: () =>
                                    cfm.toggleObsecureConfirmPassword(),
                                icon: Icon(
                                  cfm.obsecureConfirmPassword
                                      ? Icons.visibility_off_outlined
                                      : Icons.visibility_outlined,
                                  color: Colors.black,
                                ),
                              ),
                            );
                          },
                        ),
                        const Spacer(),
                        FilledButtonCustom(
                          text: "ตกลง",
                          onPressed: () async {
                            if (!_formKey.currentState!.validate()) return;

                            final provider = context
                                .read<NewPasswordProvider>();
                            final res = await provider.resetPassword(
                              _passwordController.text.trim(),
                            );

                            if (!mounted) return;
                            if (res != null) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                SnackBar(
                                  content: Text(
                                    "อัปเดตรหัสผ่านของ ${provider.email} เสร็จสิ้น",
                                    style: const TextStyle(color: Colors.white),
                                  ),
                                  backgroundColor: Colors.green,
                                  behavior: SnackBarBehavior.floating,
                                  duration: const Duration(seconds: 2),
                                ),
                              );
                              await Future.delayed(const Duration(seconds: 1));
                              Navigator.pushReplacement(
                                context,
                                MaterialPageRoute(
                                  builder: (_) => const LoginScreen(),
                                ),
                              );
                            }
                          },
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
                                          color: Colors.black,
                                          fontSize: 20,
                                        ),
                                      ),
                                    ),
                                    TextSpan(
                                      text: "กลับ",
                                      style: TextStyle(
                                        color: Colors.black,
                                        letterSpacing: 0,
                                        fontSize: 20,
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
