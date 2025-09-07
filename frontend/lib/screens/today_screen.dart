import 'package:flutter/material.dart';
import 'package:frontend/providers/today_provider.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';

class TodayScreen extends StatefulWidget {
  const TodayScreen({super.key});

  @override
  State<TodayScreen> createState() => _TodayScreenState();
}

class _TodayScreenState extends State<TodayScreen> {
  final _formKey = GlobalKey<FormState>();
  final _saveSymptom = TextEditingController();

  @override
  void dispose() {
    _saveSymptom.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final p = context.watch<TodayProvider>();

    return Scaffold(
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              "ตารางกินยา",
              style: TextStyle(
                fontSize: 25,
                color: Colors.white,
                fontWeight: FontWeight.w700,
              ),
            ),
            const SizedBox(height: 2),
            Row(
              children: [
                // แตะเพื่อเปิดปฏิทินเลือกวัน
                InkWell(
                  onTap: () => p.pickDate(context),
                  child: Text(
                    p.dateLabel,
                    style: const TextStyle(fontSize: 20, color: Colors.white),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
      backgroundColor: color.AppColors.backgroundColor2nd,
      body: Container(
        width: double.infinity,
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 20),
        child: p.doseSelect(p.selected).isEmpty
            ? Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: <Widget>[
                  const Spacer(),
                  Opacity(
                    opacity: 0.5,
                    child: Image.asset(
                      "assets/images/drugs.png",
                      height: 210,
                      width: 210,
                    ),
                  ),
                  const SizedBox(height: 50),
                  const Text(
                    "วันนี้ไม่ต้องกินยา",
                    style: TextStyle(fontWeight: FontWeight.bold, fontSize: 32),
                  ),
                  const SizedBox(height: 10),
                  const Spacer(),
                ],
              )
            : ListView.builder(
                itemBuilder: (_, i) {
                  final d = p.doseSelect(p.selected)[i];
                  final timeText = DateFormat('HH:mm').format(d.at.toLocal());
                  return Padding(
                    padding: const EdgeInsets.only(bottom: 5),
                    child: SizedBox(
                      height: 100,
                      child: Card(
                        color: Colors.white,
                        child: InkWell(
                          onTap: () {
                            showDialog(
                              context: context,
                              builder: (context) {
                                return Dialog(
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(0),
                                  ),
                                  child: Container(
                                    padding: const EdgeInsets.all(10),
                                    height: 178,
                                    width: 200,
                                    color: Colors.white,
                                    child: Column(
                                      crossAxisAlignment:
                                          CrossAxisAlignment.center,
                                      children: [
                                        const SizedBox(height: 15),
                                        Row(
                                          crossAxisAlignment:
                                              CrossAxisAlignment.start,
                                          mainAxisSize: MainAxisSize.min,
                                          children: [
                                            SizedBox(
                                              height: 62,
                                              width: 120,
                                              child: ElevatedButton(
                                                onPressed: () {
                                                  p.handleIsTaken("taken", d);
                                                  Navigator.pop(context);
                                                },
                                                style: ElevatedButton.styleFrom(
                                                  shape: RoundedRectangleBorder(
                                                    borderRadius:
                                                        BorderRadiusGeometry.circular(
                                                          10,
                                                        ),
                                                  ),
                                                  backgroundColor: color
                                                      .AppColors
                                                      .greenColor,
                                                ),
                                                child: const Text(
                                                  "กินแล้ว",
                                                  style: TextStyle(
                                                    color: Colors.white,
                                                    fontSize: 24,
                                                  ),
                                                ),
                                              ),
                                            ),
                                            const SizedBox(width: 10),
                                            SizedBox(
                                              height: 62,
                                              width: 130,
                                              child: ElevatedButton(
                                                onPressed: () {
                                                  p.handleIsTaken(
                                                    "not_taken",
                                                    d,
                                                  );
                                                  Navigator.pop(context);
                                                },
                                                style: ElevatedButton.styleFrom(
                                                  shape: RoundedRectangleBorder(
                                                    borderRadius:
                                                        BorderRadiusGeometry.circular(
                                                          10,
                                                        ),
                                                  ),
                                                  backgroundColor:
                                                      color.AppColors.redColor,
                                                ),
                                                child: const Text(
                                                  "ยังไม่กิน",
                                                  style: TextStyle(
                                                    color: Colors.white,
                                                    fontSize: 24,
                                                  ),
                                                ),
                                              ),
                                            ),
                                          ],
                                        ),
                                        const SizedBox(height: 30),
                                        SizedBox(
                                          width: 109,
                                          height: 37,
                                          child: ElevatedButton(
                                            onPressed: () {
                                              p.handleIsTaken("remove", d);
                                              Navigator.pop(context);
                                            },
                                            style: ElevatedButton.styleFrom(
                                              shape: RoundedRectangleBorder(
                                                borderRadius:
                                                    BorderRadiusGeometry.circular(
                                                      10,
                                                    ),
                                              ),
                                              backgroundColor: const Color(
                                                0xFFFFA100,
                                              ),
                                            ),
                                            child: const Text(
                                              "ลบออก",
                                              style: TextStyle(
                                                color: Colors.white,
                                                fontSize: 16,
                                              ),
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                  ),
                                );
                              },
                            );
                          },
                          child: Padding(
                            padding: const EdgeInsets.fromLTRB(10, 15, 16, 0),
                            child: Row(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Expanded(
                                  child: Column(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                      Row(
                                        children: [
                                          Text(
                                            timeText + " น.",
                                            style: const TextStyle(
                                              fontSize: 24,
                                              fontWeight: FontWeight.normal,
                                              color: Colors.black,
                                            ),
                                          ),
                                          const SizedBox(width: 7),
                                          Text(
                                            "(${d.instruction})",
                                            style: const TextStyle(
                                              color: Colors.black,
                                              fontSize: 16,
                                              fontWeight: FontWeight.normal,
                                            ),
                                          ),
                                          const SizedBox(width: 50),
                                          IconButton(
                                            onPressed: () {
                                              showDialog(
                                                context: context,
                                                builder: (context) {
                                                  return Dialog(
                                                    shape:
                                                        const RoundedRectangleBorder(
                                                          borderRadius:
                                                              BorderRadiusGeometry
                                                                  .zero,
                                                        ),
                                                    child: Form(
                                                      key: _formKey,
                                                      child: Container(
                                                        color: const Color(
                                                          0xFFFFE78E,
                                                        ),
                                                        width: 329,
                                                        height: 287,
                                                        padding:
                                                            const EdgeInsets.all(
                                                              10,
                                                            ),
                                                        child: Column(
                                                          crossAxisAlignment:
                                                              CrossAxisAlignment
                                                                  .start,
                                                          children: [
                                                            Expanded(
                                                              child: TextFormField(
                                                                controller:
                                                                    _saveSymptom,
                                                                validator: (e) {
                                                                  if (e ==
                                                                          null ||
                                                                      e
                                                                          .trim()
                                                                          .isEmpty) {
                                                                    return "กรุณากรอกอาการ";
                                                                  }
                                                                  return null;
                                                                },
                                                                decoration: const InputDecoration(
                                                                  hint: Text(
                                                                    "กรอกอาการ",
                                                                  ),
                                                                  enabledBorder: UnderlineInputBorder(
                                                                    borderSide: BorderSide(
                                                                      width: 0,
                                                                      color: Color(
                                                                        0xFFFFE78E,
                                                                      ),
                                                                    ),
                                                                  ),
                                                                  focusedBorder: UnderlineInputBorder(
                                                                    borderSide: BorderSide(
                                                                      width: 0,
                                                                      color: Color(
                                                                        0xFFFFE78E,
                                                                      ),
                                                                    ),
                                                                  ),
                                                                ),
                                                                maxLines: null,
                                                                maxLength: 200,
                                                              ),
                                                            ),
                                                            Row(
                                                              mainAxisAlignment:
                                                                  MainAxisAlignment
                                                                      .end,
                                                              children: [
                                                                ElevatedButton(
                                                                  onPressed: () {
                                                                    Navigator.pop(
                                                                      context,
                                                                    );
                                                                    _saveSymptom
                                                                        .clear();
                                                                  },
                                                                  style: ElevatedButton.styleFrom(
                                                                    shape: RoundedRectangleBorder(
                                                                      borderRadius:
                                                                          BorderRadiusGeometry.circular(
                                                                            5,
                                                                          ),
                                                                    ),
                                                                    backgroundColor:
                                                                        const Color(
                                                                          0xFF000000,
                                                                        ),
                                                                    minimumSize:
                                                                        const Size(
                                                                          109,
                                                                          37,
                                                                        ),
                                                                  ),
                                                                  child: const Text(
                                                                    "ยกเลิก",
                                                                    style: TextStyle(
                                                                      color: Colors
                                                                          .white,
                                                                      fontSize:
                                                                          16,
                                                                    ),
                                                                  ),
                                                                ),
                                                                const SizedBox(
                                                                  width: 10,
                                                                ),
                                                                ElevatedButton(
                                                                  onPressed: () {
                                                                    if (_formKey
                                                                        .currentState!
                                                                        .validate()) {
                                                                      Navigator.pop(
                                                                        context,
                                                                      );
                                                                      _saveSymptom
                                                                          .clear();
                                                                    }
                                                                  },
                                                                  style: ElevatedButton.styleFrom(
                                                                    shape: RoundedRectangleBorder(
                                                                      borderRadius:
                                                                          BorderRadiusGeometry.circular(
                                                                            5,
                                                                          ),
                                                                    ),
                                                                    backgroundColor:
                                                                        const Color(
                                                                          0xFFA3A3A3,
                                                                        ),
                                                                    minimumSize:
                                                                        const Size(
                                                                          109,
                                                                          37,
                                                                        ),
                                                                  ),
                                                                  child: const Text(
                                                                    "บันทึกอาการ",
                                                                    style: TextStyle(
                                                                      color: Colors
                                                                          .white,
                                                                      fontSize:
                                                                          16,
                                                                    ),
                                                                  ),
                                                                ),
                                                              ],
                                                            ),
                                                          ],
                                                        ),
                                                      ),
                                                    ),
                                                  );
                                                },
                                              );
                                            },
                                            icon: Icon(
                                              Icons.note_add_outlined,
                                              size: 32,
                                              color: d.isTake
                                                  ? const Color(0xFFFFC800)
                                                  : const Color(0xFFA5A5A5),
                                            ),
                                          ),
                                        ],
                                      ),
                                      const Text(
                                        "paracetamol",
                                        style: TextStyle(
                                          color: Colors.black,
                                          fontSize: 16,
                                          fontWeight: FontWeight.normal,
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                                Column(
                                  crossAxisAlignment: CrossAxisAlignment.end,
                                  children: [
                                    // Text(
                                    //   d.isTaken ? "กินแล้ว" : "ยังไม่กิน",
                                    //   style: TextStyle(
                                    //     color: d.isTaken
                                    //         ? color.AppColors.greenColor
                                    //         : color.AppColors.redColor,
                                    //     fontSize: 16,
                                    //     fontWeight: FontWeight.bold,
                                    //   ),
                                    // ),
                                    const SizedBox(height: 15),
                                    Image.asset(
                                      "assets/images/pill.png",
                                      width: 33,
                                      height: 33,
                                    ),
                                  ],
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ),
                  );
                },
                itemCount: p.doseSelect(p.selected).length,
              ),
      ),
    );
  }
}
