import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:frontend/enums/page_type.dart';
import 'package:frontend/providers/otp_provider.dart';
import 'package:frontend/screens/forget_password_screen.dart';
import 'package:frontend/screens/new_password_screen.dart';
import 'package:frontend/screens/profile_setup_screen.dart';
import 'package:frontend/screens/register_screen.dart';

import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/widgets/filled_button_custom.dart';
import 'package:pin_code_fields/pin_code_fields.dart';
import 'package:provider/provider.dart';

class OtpScreens extends StatelessWidget {
  const OtpScreens({super.key});

  @override
  Widget build(BuildContext context) {
    final args =
        ModalRoute.of(context)!.settings.arguments as Map<String, dynamic>;
    final PageType otpType = args['otpType'] as PageType;
    final String email = args['email'] as String;

    return ChangeNotifierProvider(
      create: (_) {
        final p = OtpProvider(otpType: otpType, emailText: email);
        p.init();
        return p;
      },
      child: _OtpView(),
    );
  }
}

class _OtpView extends StatefulWidget {
  @override
  State<_OtpView> createState() => _OtpViewState();
}

class _OtpViewState extends State<_OtpView> {
  String _otp = "";
  @override
  void initState() {
    super.initState();
  }

  @override
  void dispose() {
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
                    width: 191,
                    height: 191,
                  ),
                ),
                const SizedBox(height: 15),
                Container(
                  width: double.infinity,
                  height: 514,
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(20),
                    color: Colors.white,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 45.0),
                        child: Column(
                          children: [
                            const SizedBox(height: 30),
                            const Text(
                              "ยืนยันอีเมล",
                              style: TextStyle(
                                color: Colors.black,
                                fontWeight: FontWeight.bold,
                                fontSize: 24,
                              ),
                            ),
                            const SizedBox(height: 12),
                            const Text(
                              "กรุณาใส่รหัส OTP 6 หลักที่ส่งไปยัง",
                              style: TextStyle(fontSize: 16),
                            ),
                            const SizedBox(height: 15),
                            Container(
                              width: 290,
                              padding: const EdgeInsets.symmetric(
                                vertical: 5.0,
                                horizontal: 14,
                              ),
                              decoration: BoxDecoration(
                                borderRadius: BorderRadius.circular(30),
                                border: Border.all(
                                  color: Colors.black,
                                  width: 1,
                                ),
                                color: Colors.white,
                              ),
                              child: Row(
                                crossAxisAlignment: CrossAxisAlignment.center,
                                children: [
                                  const Icon(
                                    Icons.email_outlined,
                                    color: Colors.black,
                                    size: 30,
                                  ),
                                  const SizedBox(width: 10),
                                  Expanded(
                                    child: Selector<OtpProvider, String>(
                                      builder: (_, emailShow, _) {
                                        return Text(
                                          emailShow,
                                          style: const TextStyle(
                                            color: Colors.black,
                                            fontSize: 16,
                                          ),
                                        );
                                      },
                                      selector: (_, p) => p.emailText,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            const SizedBox(height: 40),
                            const Align(
                              alignment: Alignment.centerLeft,
                              child: Text(
                                "รหัสผ่านมีอายุการใช้งาน 3 นาที",
                                style: TextStyle(
                                  color: Color(0xFFFF0000),
                                  fontSize: 16,
                                ),
                              ),
                            ),
                            const SizedBox(height: 15),
                            PinCodeTextField(
                              appContext: context,
                              length: 6,
                              onCompleted: (v) => _otp = v,
                              textStyle: const TextStyle(
                                fontSize: 20,
                                color: Colors.black,
                                fontWeight: FontWeight.normal,
                              ),
                              inputFormatters: [
                                FilteringTextInputFormatter.digitsOnly,
                              ],
                              cursorHeight: 22,
                              cursorColor: Colors.black,
                              pinTheme: PinTheme(
                                shape: PinCodeFieldShape.box,
                                inactiveColor: Colors.black,
                                activeColor: Colors.black,
                                selectedColor: Colors.blue,
                                fieldHeight: 50,
                                borderRadius: BorderRadius.circular(8),
                              ),
                            ),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.start,
                              children: [
                                const Text(
                                  "ไม่ได้รับรหัส?",
                                  style: TextStyle(
                                    color: Colors.black,
                                    fontWeight: FontWeight.normal,
                                    fontSize: 16,
                                  ),
                                ),
                                const SizedBox(width: 4),
                                Consumer<OtpProvider>(
                                  builder: (_, otpSend, _) {
                                    final isCount = otpSend.countdown > 0;
                                    return TextButton(
                                      style: TextButton.styleFrom(
                                        padding: const EdgeInsets.symmetric(
                                          horizontal: 0,
                                        ),
                                        minimumSize: const Size(0, 0),
                                      ),
                                      onPressed: isCount
                                          ? null
                                          : otpSend.sendOtp,
                                      child: Text(
                                        isCount
                                            ? "ขออีกครั้งใน ${otpSend.countdown} วินาที"
                                            : "ขอรหัส OTP อีกครั้ง",
                                        style: TextStyle(
                                          fontSize: 16,
                                          fontWeight: FontWeight.normal,
                                          color: isCount
                                              ? const Color(0xFF00C907)
                                              : const Color(0xFF0873FF),
                                        ),
                                      ),
                                    );
                                  },
                                ),
                              ],
                            ),
                            Consumer<OtpProvider>(
                              builder: (_, errorOtp, _) {
                                return Align(
                                  alignment: Alignment.centerLeft,
                                  child: Visibility(
                                    visible: errorOtp.errorOtp.isNotEmpty,
                                    child: Text(
                                      errorOtp.errorOtp,
                                      style: const TextStyle(
                                        fontSize: 14,
                                        color: Color(0xFFFF0000),
                                      ),
                                    ),
                                  ),
                                );
                              },
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 5),
                      const Spacer(),
                      FilledButtonCustom(
                        text: "ยืนยันอีเมล",
                        onPressed: () {
                          final p = context.read<OtpProvider>();
                          final ok = p.validateOtp(_otp);
                          if (!ok) return;
                          switch (p.otpType) {
                            case PageType.register:
                              Navigator.pushReplacement(
                                context,
                                MaterialPageRoute(
                                  builder: (context) =>
                                      const ProfileSetupScreen(),
                                ),
                              );
                              break;
                            case PageType.forgot:
                              Navigator.pushReplacement(
                                context,
                                MaterialPageRoute(
                                  builder: (context) =>
                                      const NewPasswordScreen(),
                                  settings: RouteSettings(
                                    arguments: {"email": p.emailText},
                                  ),
                                ),
                              );
                              break;
                          }
                        },
                      ),
                      const SizedBox(height: 30),
                      Padding(
                        padding: const EdgeInsets.fromLTRB(5, 0, 0, 10),
                        child: Align(
                          alignment: Alignment.centerLeft,
                          child: TextButton(
                            onPressed: () {
                              final p = context.read<OtpProvider>();
                              switch (p.otpType) {
                                case PageType.register:
                                  Navigator.pushReplacement(
                                    context,
                                    MaterialPageRoute(
                                      builder: (context) =>
                                          const RegisterScreen(),
                                    ),
                                  );
                                  break;
                                case PageType.forgot:
                                  Navigator.pushReplacement(
                                    context,
                                    MaterialPageRoute(
                                      builder: (context) =>
                                          const ForgetPasswordScreen(),
                                    ),
                                  );
                                  break;
                              }
                            },
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
                                      color: Colors.black,
                                      fontSize: 20,
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
