import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:frontend/providers/profile_setup_provider.dart';
import 'package:frontend/screens/register_screen.dart';
import 'package:frontend/services/profile_service.dart';
import 'package:frontend/widgets/filled_button_custom.dart';
import 'package:frontend/widgets/text_field_input.dart';
import 'package:provider/provider.dart';

import 'package:frontend/utils/colors.dart' as color;

class ProfileSetupScreen extends StatelessWidget {
  const ProfileSetupScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final args =
        ModalRoute.of(context)!.settings.arguments as Map<String, dynamic>;
    final int patientId = args["patient_id"];

    return ChangeNotifierProvider(
      create: (context) => ProfileSetupProvider(
        profileService: ProfileService(),
        patientId: patientId,
      ),
      child: _ProfileSetUpView(),
    );
  }
}

class _ProfileSetUpView extends StatefulWidget {
  @override
  State<_ProfileSetUpView> createState() => _ProfileSetupScreenState();
}

class _ProfileSetupScreenState extends State<_ProfileSetUpView> {
  final _formKey = GlobalKey<FormState>();
  @override
  Widget build(BuildContext context) {
    final p = context.watch<ProfileSetupProvider>();
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
                  child: Form(
                    key: _formKey,
                    child: Column(
                      children: [
                        const SizedBox(height: 30),
                        const Text(
                          "ข้อมูลผู้ใช้",
                          style: TextStyle(
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        TextFieldInput(
                          labelname: "ชื่อ",
                          textInputType: TextInputType.text,
                          controller: p.firstnameController,
                          inputFormatters: [
                            FilteringTextInputFormatter.allow(
                              RegExp(r'[a-zA-Zก-ฮะาิีึืุูเแโใไๅ์่้๊๋ั็็ฯๆฦฦำ]'),
                            ),
                          ],
                        ),
                        Consumer<ProfileSetupProvider>(
                          builder: (_, error, _) {
                            return Padding(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 27,
                              ),
                              child: Align(
                                alignment: Alignment.centerLeft,
                                child: Visibility(
                                  visible: error.hasError[0].isNotEmpty,
                                  child: Text(
                                    error.hasError[0],
                                    style: const TextStyle(
                                      fontSize: 12,
                                      color: Color(0xFFFF0000),
                                    ),
                                  ),
                                ),
                              ),
                            );
                          },
                        ),
                        const SizedBox(height: 5),
                        TextFieldInput(
                          labelname: "นามสกุล",
                          textInputType: TextInputType.text,
                          controller: p.lastnameController,
                          inputFormatters: [
                            FilteringTextInputFormatter.allow(
                              RegExp(r'[a-zA-Zก-ฮะาิีึืุูเแโใไๅ์่้๊๋ั็็ฯๆฦฦำ]'),
                            ),
                          ],
                        ),
                        Consumer<ProfileSetupProvider>(
                          builder: (_, error, _) {
                            return Padding(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 27,
                              ),
                              child: Align(
                                alignment: Alignment.centerLeft,
                                child: Visibility(
                                  visible: error.hasError[1].isNotEmpty,
                                  child: Text(
                                    error.hasError[1],
                                    style: const TextStyle(
                                      fontSize: 12,
                                      color: Color(0xFFFF0000),
                                    ),
                                  ),
                                ),
                              ),
                            );
                          },
                        ),
                        const SizedBox(height: 5),
                        TextFieldInput(
                          labelname: "เลขบัตรประชาชน",
                          textInputType: TextInputType.number,
                          controller: p.idcardController,
                          inputFormatters: [
                            FilteringTextInputFormatter.allow(RegExp(r'\d')),
                            LengthLimitingTextInputFormatter(13),
                          ],
                        ),
                        Consumer<ProfileSetupProvider>(
                          builder: (_, error, _) {
                            return Padding(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 27,
                              ),
                              child: Align(
                                alignment: Alignment.centerLeft,
                                child: Visibility(
                                  visible: error.hasError[2].isNotEmpty,
                                  child: Text(
                                    error.hasError[2],
                                    style: const TextStyle(
                                      fontSize: 12,
                                      color: Color(0xFFFF0000),
                                    ),
                                  ),
                                ),
                              ),
                            );
                          },
                        ),
                        const SizedBox(height: 5),
                        TextFieldInput(
                          labelname: "เบอร์โทรศัพท์",
                          textInputType: TextInputType.number,
                          controller: p.phoneController,
                          inputFormatters: [
                            FilteringTextInputFormatter.allow(RegExp(r'[0-9]')),
                            LengthLimitingTextInputFormatter(10),
                          ],
                        ),
                        Consumer<ProfileSetupProvider>(
                          builder: (_, error, _) {
                            return Padding(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 27,
                              ),
                              child: Align(
                                alignment: Alignment.centerLeft,
                                child: Visibility(
                                  visible: error.hasError[3].isNotEmpty,
                                  child: Text(
                                    error.hasError[3],
                                    style: const TextStyle(
                                      fontSize: 12,
                                      color: Color(0xFFFF0000),
                                    ),
                                  ),
                                ),
                              ),
                            );
                          },
                        ),
                        const Spacer(),
                        FilledButtonCustom(
                          text: p.isLoading ? "กำลังบันทึก..." : "ลงทะเบียน",
                          onPressed: p.isLoading
                              ? null
                              : () => p.checkSetUp(context),
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
                                  builder: (context) => const RegisterScreen(),
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
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
