import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_svg/svg.dart';
import 'package:frontend/providers/profile_provider.dart';
import 'package:frontend/screens/login_screen.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:provider/provider.dart';

class ProfileScreen extends StatelessWidget {
  ProfileScreen({super.key});
  final _formKey = GlobalKey<FormState>();

  OutlineInputBorder _border(Color color) {
    return OutlineInputBorder(
      borderRadius: BorderRadius.circular(13),
      borderSide: BorderSide(width: 1, color: color),
    );
  }

  @override
  Widget build(BuildContext context) {
    final p = context.watch<ProfileProvider>();
    if (p.isLoading) {
      Scaffold(
        backgroundColor: color.AppColors.backgroundColor1st,
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SizedBox(
                height: 280,
                child: Stack(
                  clipBehavior: Clip.none,
                  alignment: Alignment.center,
                  children: [
                    SvgPicture.asset(
                      "assets/images/clock.svg",
                      colorFilter: const ColorFilter.mode(
                        Colors.white,
                        BlendMode.srcIn,
                      ),
                      height: 190,
                      width: 200,
                    ),
                    Positioned(
                      bottom: -20,
                      left: -70,
                      child: Image.asset(
                        "assets/images/drugs.png",
                        height: 153,
                        width: 153,
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 120),
              const Text(
                "PillMate",
                style: TextStyle(color: Colors.white, fontSize: 48),
              ),
            ],
          ),
        ),
      );
    }

    if (p.user == null) {
      return const Scaffold(body: Center(child: Text("ไม่พบข้อมูลผู้ใช้ ❌")));
    }
    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor2nd,
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        foregroundColor: Colors.white,
        title: const Text(
          "ข้อมูลผู้ใช้",
          style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold),
        ),
      ),
      body: Container(
        width: double.infinity,
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              width: 389,
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(14),
                border: Border.all(color: const Color(0xFFFF92DB), width: 1.5),
              ),
              child: Column(
                children: [
                  Container(
                    width: double.infinity,
                    height: 64,
                    padding: const EdgeInsets.all(12),
                    decoration: const BoxDecoration(
                      color: Color(0xFFFF92DB),
                      borderRadius: BorderRadius.only(
                        topLeft: Radius.circular(12),
                        topRight: Radius.circular(12),
                      ),
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          "ข้อมูลทั่วไป",
                          style: TextStyle(color: Colors.black, fontSize: 24),
                        ),
                        IconButton(
                          onPressed: () {
                            showDialog(
                              context: context,
                              barrierDismissible: false,
                              builder: (context) {
                                final idCardController = TextEditingController(
                                  text: p.user!.idCard,
                                );
                                final firstNameController =
                                    TextEditingController(
                                      text: p.user!.firstName,
                                    );
                                final lastNameController =
                                    TextEditingController(
                                      text: p.user!.lastName,
                                    );
                                final telController = TextEditingController(
                                  text: p.user!.tel,
                                );
                                return SingleChildScrollView(
                                  child: Dialog(
                                    child: Form(
                                      key: _formKey,
                                      child: Container(
                                        padding: const EdgeInsets.all(12),
                                        decoration: BoxDecoration(
                                          color: Colors.white,
                                          borderRadius: BorderRadius.circular(
                                            13,
                                          ),
                                        ),
                                        width: 328,
                                        height: 753,
                                        child: Column(
                                          crossAxisAlignment:
                                              CrossAxisAlignment.start,
                                          children: [
                                            Center(
                                              child: Row(
                                                mainAxisAlignment:
                                                    MainAxisAlignment.center,
                                                children: [
                                                  const Expanded(
                                                    child: Text(
                                                      "ระบุข้อมูลผู้ใช้",
                                                      style: TextStyle(
                                                        color: Colors.black,
                                                        fontSize: 24,
                                                      ),
                                                    ),
                                                  ),
                                                  IconButton(
                                                    onPressed: () {
                                                      Navigator.pop(context);
                                                    },
                                                    icon: const Icon(
                                                      size: 40,
                                                      CupertinoIcons.xmark,
                                                    ),
                                                  ),
                                                ],
                                              ),
                                            ),
                                            const SizedBox(height: 25),
                                            const Text(
                                              "เลขบัตรประชาชน",
                                              style: TextStyle(
                                                color: Colors.black,
                                                fontSize: 20,
                                              ),
                                            ),
                                            const SizedBox(height: 10),
                                            TextFormField(
                                              inputFormatters: [
                                                FilteringTextInputFormatter.allow(
                                                  RegExp(r'[0-9]'),
                                                ),
                                                LengthLimitingTextInputFormatter(
                                                  13,
                                                ),
                                              ],
                                              keyboardType:
                                                  TextInputType.number,
                                              decoration: InputDecoration(
                                                enabledBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                                focusedBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                                errorBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                                focusedErrorBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                              ),
                                              validator: (cn) {
                                                if (cn == null ||
                                                    cn.trim().isEmpty) {
                                                  return "กรุณากรอกค่า";
                                                } else if (cn.length != 13) {
                                                  return "กรุณากรอกให้ครบ 13 ตัว";
                                                }
                                                return null;
                                              },
                                              controller: idCardController,
                                            ),
                                            const SizedBox(height: 20),
                                            const Text(
                                              "ชื่อ",
                                              style: TextStyle(
                                                color: Colors.black,
                                                fontSize: 20,
                                              ),
                                            ),
                                            const SizedBox(height: 10),
                                            TextFormField(
                                              inputFormatters: [
                                                FilteringTextInputFormatter.allow(
                                                  RegExp(
                                                    r'[a-zA-Zก-ฮะาิีึืุูเแโใไๅ์่้๊๋ั็็ฯๆฦฦำ]',
                                                  ),
                                                ),
                                              ],
                                              keyboardType: TextInputType.text,
                                              decoration: InputDecoration(
                                                enabledBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                                focusedBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                                errorBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                                focusedErrorBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                              ),
                                              validator: (fname) =>
                                                  fname == null ||
                                                      fname.trim().isEmpty
                                                  ? "กรุณากรอกค่า"
                                                  : null,
                                              controller: firstNameController,
                                            ),
                                            const SizedBox(height: 20),
                                            const Text(
                                              "นามสกุล",
                                              style: TextStyle(
                                                color: Colors.black,
                                                fontSize: 20,
                                              ),
                                            ),
                                            const SizedBox(height: 10),
                                            TextFormField(
                                              inputFormatters: [
                                                FilteringTextInputFormatter.allow(
                                                  RegExp(
                                                    r'[a-zA-Zก-ฮะาิีึืุูเแโใไๅ์่้๊๋ั็็ฯๆฦฦำ]',
                                                  ),
                                                ),
                                              ],
                                              keyboardType: TextInputType.text,
                                              decoration: InputDecoration(
                                                enabledBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                                focusedBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                                errorBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                                focusedErrorBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                              ),
                                              validator: (lname) =>
                                                  lname == null ||
                                                      lname.trim().isEmpty
                                                  ? "กรุณากรอกค่า"
                                                  : null,
                                              controller: lastNameController,
                                            ),
                                            const SizedBox(height: 20),
                                            const Text(
                                              "เบอร์โทรศัพท์",
                                              style: TextStyle(
                                                color: Colors.black,
                                                fontSize: 20,
                                              ),
                                            ),
                                            const SizedBox(height: 10),
                                            TextFormField(
                                              inputFormatters: [
                                                FilteringTextInputFormatter.allow(
                                                  RegExp('[0-9]'),
                                                ),
                                                LengthLimitingTextInputFormatter(
                                                  10,
                                                ),
                                              ],
                                              keyboardType: TextInputType.phone,
                                              decoration: InputDecoration(
                                                enabledBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                                focusedBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                                errorBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                                focusedErrorBorder: _border(
                                                  const Color(0xFF8D8D8D),
                                                ),
                                              ),
                                              validator: (tel) {
                                                if (tel == null ||
                                                    tel.trim().isEmpty) {
                                                  return "กรุณากรอกค่า";
                                                } else if (tel.length != 10) {
                                                  return "กรุณากรอกให้ครบ 10 ตัว";
                                                }
                                                return null;
                                              },
                                              controller: telController,
                                            ),
                                            const Spacer(),
                                            Center(
                                              child: ElevatedButton(
                                                style: ElevatedButton.styleFrom(
                                                  shape: RoundedRectangleBorder(
                                                    borderRadius:
                                                        BorderRadius.circular(
                                                          5,
                                                        ),
                                                  ),
                                                  minimumSize: const Size(
                                                    237,
                                                    50,
                                                  ),
                                                  shadowColor: Colors.black,
                                                  elevation: 4,
                                                  backgroundColor: const Color(
                                                    0xFF03B200,
                                                  ),
                                                ),
                                                onPressed: () {
                                                  if (_formKey.currentState!
                                                      .validate()) {
                                                    p.updatedInfoProfile(
                                                      InfoUser(
                                                        firstName:
                                                            firstNameController
                                                                .text
                                                                .trim(),
                                                        lastName:
                                                            lastNameController
                                                                .text
                                                                .trim(),
                                                        idCard: idCardController
                                                            .text
                                                            .trim(),
                                                        tel: telController.text
                                                            .trim(),
                                                      ),
                                                      context,
                                                    );
                                                    Navigator.pop(context);
                                                  }
                                                },
                                                child: Text(
                                                  p.isLoading
                                                      ? "กำลังบันทึก..."
                                                      : "บันทึกข้อมูล",
                                                  style: const TextStyle(
                                                    color: Colors.white,
                                                    fontSize: 24,
                                                  ),
                                                ),
                                              ),
                                            ),
                                            const SizedBox(height: 15),
                                          ],
                                        ),
                                      ),
                                    ),
                                  ),
                                );
                              },
                            );
                          },
                          icon: const Icon(Icons.edit_square, size: 26),
                          color: Colors.black,
                        ),
                      ],
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.all(12),
                    width: double.infinity,
                    height: 178,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          "เลขบัตรประชาชน :"
                          " ${p.user!.idCard}",
                          style: const TextStyle(
                            color: Colors.black,
                            fontSize: 20,
                          ),
                        ),
                        const SizedBox(height: 10),
                        Text(
                          "ชื่อ :"
                          " ${p.user!.firstName}",
                          style: const TextStyle(
                            color: Colors.black,
                            fontSize: 20,
                          ),
                        ),
                        const SizedBox(height: 10),
                        Text(
                          "นามสกุล :"
                          " ${p.user!.lastName}",
                          style: const TextStyle(
                            color: Colors.black,
                            fontSize: 20,
                          ),
                        ),
                        const SizedBox(height: 10),
                        Text(
                          "เบอร์โทรศัพท์ :"
                          " ${p.user!.tel}",
                          style: const TextStyle(
                            color: Colors.black,
                            fontSize: 20,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 20),
            Container(
              width: 389,
              decoration: BoxDecoration(
                color: Colors.white,
                border: Border.all(width: 1.5, color: const Color(0xFF87F07B)),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Column(
                children: [
                  Container(
                    padding: const EdgeInsets.all(12),
                    width: double.infinity,
                    height: 64,
                    decoration: const BoxDecoration(
                      color: Color(0xFF87F07B),
                      borderRadius: BorderRadius.only(
                        topLeft: Radius.circular(12),
                        topRight: Radius.circular(12),
                      ),
                    ),
                    child: const Text(
                      "ข้อมูลนัดพบแพทย์",
                      style: TextStyle(fontSize: 24, color: Colors.black),
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.all(12),
                    width: double.infinity,
                    height: 178,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          p.appointmentDay == null
                              ? "ไม่มีใบนัดถัดไป"
                              : "วันที่ :"
                                    " ${p.appointmentDay}",
                          style: const TextStyle(
                            color: Colors.black,
                            fontSize: 20,
                          ),
                        ),
                        const SizedBox(height: 10),
                        Text(
                          p.appointmentHourMinute == null
                              ? "-"
                              : "เวลา :"
                                    " ${p.appointmentHourMinute}",
                          style: const TextStyle(
                            color: Colors.black,
                            fontSize: 20,
                          ),
                        ),
                        const SizedBox(height: 10),
                        Text(
                          "Note : ${p.appointment?.note ?? '-'}",
                          style: const TextStyle(
                            color: Colors.black,
                            fontSize: 20,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            const Spacer(),
            Center(
              child: ElevatedButton(
                onPressed: p.isLoading
                    ? null
                    : () async {
                        final success = await p.logout();
                        if (success && context.mounted) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text("ออกจากระบบสำเร็จ ✅"),
                              backgroundColor: Colors.green,
                              duration: Duration(seconds: 2),
                            ),
                          );
                          Navigator.pushAndRemoveUntil(
                            context,
                            MaterialPageRoute(
                              builder: (_) => const LoginScreen(),
                            ),
                            (route) => false,
                          );
                        } else {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text("ออกจากระบบไม่สำเร็จ ❌"),
                              backgroundColor: Colors.red,
                            ),
                          );
                        }
                      },
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.grey[300],
                  minimumSize: const Size(389, 50),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(5),
                  ),
                ),
                child: const Text(
                  "ออกจากระบบ",
                  style: TextStyle(color: Colors.black, fontSize: 24),
                ),
              ),
            ),
            const SizedBox(height: 13),
          ],
        ),
      ),
    );
  }
}
