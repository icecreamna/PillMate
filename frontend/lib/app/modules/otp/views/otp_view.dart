import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:frontend/app/widgets/filled_button_custom.dart';

import 'package:get/get.dart';
import 'package:pin_code_fields/pin_code_fields.dart';

import '../controllers/otp_controller.dart';
import '../../../utils//colors.dart' as color;

class OtpView extends GetView<OtpController> {
  const OtpView({super.key});

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
                              child: const Row(
                                crossAxisAlignment: CrossAxisAlignment.center,
                                children: [
                                  Icon(
                                    Icons.email_outlined,
                                    color: Colors.black,
                                    size: 30,
                                  ),
                                  SizedBox(width: 10),
                                  Expanded(
                                    child: Text(
                                      "Kittabeth554@gmail.com",
                                      style: TextStyle(
                                        color: Colors.black,
                                        fontSize: 16,
                                      ),
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
                              controller: controller.otpController,
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
                                Obx(() {
                                  final isCount =
                                      controller.countdown.value > 0;
                                  return TextButton(
                                    style: TextButton.styleFrom(
                                      padding: const EdgeInsets.symmetric(horizontal: 0),
                                      minimumSize: const Size(0, 0)
                                    ),
                                    onPressed: isCount
                                        ? null
                                        : controller.sendOtp,
                                    child: Text(
                                      isCount
                                          ? "ขออีกครั้งใน ${controller.countdown.value} วินาที"
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
                                }),
                              ],
                            ),
                            Obx(
                              () => Align(
                                alignment: Alignment.centerLeft,
                                child: Visibility(
                                  visible: controller.errorOtp.isNotEmpty,
                                  child: Text(
                                    controller.errorOtp.value,
                                    style: const TextStyle(
                                      fontSize: 14,
                                      color: Color(0xFFFF0000),
                                    ),
                                  ),
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 5,),
                      const Spacer(),
                      FilledButtonCustom(
                        text: "ยืนยันอีเมล",
                        onPressed: controller.validateOtp,
                      ),
                      const SizedBox(height: 30),
                      Padding(
                        padding: const EdgeInsets.fromLTRB(5, 0, 0, 10),
                        child: Align(
                          alignment: Alignment.centerLeft,
                          child: TextButton(
                            onPressed: () => controller.goBackScreen(),
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
