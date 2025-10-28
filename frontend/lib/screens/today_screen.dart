import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
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
  void initState() {
    super.initState();
    Future.microtask(
      () => Provider.of<TodayProvider>(context, listen: false).loadTodayData(),
    );
  }

  @override
  void dispose() {
    _saveSymptom.dispose();
    super.dispose();
  }

  //   List<Widget> buildDoseWidgets(DoseGroup d) {
  //   if (d.doses.length > 1) {
  //     return d.doses
  //         .map((dose) => Padding(
  //               padding: const EdgeInsets.only(top: 5, bottom: 5, left: 15),
  //               child: Text("• ${dose.name} ${dose.unit}"),
  //             ))
  //         .toList();
  //   } else {
  //     return [
  //       Padding(
  //         padding: const EdgeInsets.only(top: 5, bottom: 5, left: 15),
  //         child: Text("${d.doses.first.name} ${d.doses.first.unit}"),
  //       ),
  //     ];
  //   }
  // }

  @override
  Widget build(BuildContext context) {
    final p = context.watch<TodayProvider>();
    final items = p.doseSelect(p.selected);

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

    return Scaffold(
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        foregroundColor: Colors.white,
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              "ตารางกินยา",
              style: TextStyle(fontSize: 25, fontWeight: FontWeight.w700),
            ),
            const SizedBox(height: 2),
            Row(
              children: [
                // แตะเพื่อเปิดปฏิทินเลือกวัน
                InkWell(
                  onTap: () => p.pickDate(context),
                  child: Text(
                    p.dateLabel,
                    style: const TextStyle(fontSize: 20),
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
        child: items.isEmpty
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
                  final d = items[i];
                  final timeText = DateFormat('HH:mm').format(d.at.toLocal());
                  return Padding(
                    padding: const EdgeInsets.only(bottom: 5),
                    child: SizedBox(
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
                                                onPressed: () async {
                                                  await p.updateTakenStatus(
                                                    d,
                                                    true,
                                                    context,
                                                  );
                                                  Navigator.pop(context);
                                                },
                                                style: ElevatedButton.styleFrom(
                                                  shape: RoundedRectangleBorder(
                                                    borderRadius:
                                                        BorderRadius.circular(
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
                                                onPressed: () async {
                                                  await p.updateTakenStatus(
                                                    d,
                                                    false,
                                                    context,
                                                  );
                                                  Navigator.pop(context);
                                                },
                                                style: ElevatedButton.styleFrom(
                                                  shape: RoundedRectangleBorder(
                                                    borderRadius:
                                                        BorderRadius.circular(
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
                                        // SizedBox(
                                        //   width: 109,
                                        //   height: 37,
                                        //   child: ElevatedButton(
                                        //     onPressed: () {
                                        //       p.removeDose(d);
                                        //       Navigator.pop(context);
                                        //     },
                                        //     style: ElevatedButton.styleFrom(
                                        //       shape: RoundedRectangleBorder(
                                        //         borderRadius:
                                        //             BorderRadius.circular(10),
                                        //       ),
                                        //       backgroundColor: const Color(
                                        //         0xFFFFA100,
                                        //       ),
                                        //     ),
                                        //     child: const Text(
                                        //       "ลบออก",
                                        //       style: TextStyle(
                                        //         color: Colors.white,
                                        //         fontSize: 16,
                                        //       ),
                                        //     ),
                                        //   ),
                                        // ),
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
                                            style: TextStyle(
                                              fontSize: 24,
                                              fontWeight: FontWeight.normal,
                                              color: Colors.black,
                                              decoration: d.isTaken
                                                  ? TextDecoration.lineThrough
                                                  : null,
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
                                              _saveSymptom.text =
                                                  d.symptomNote ?? "";
                                              showDialog(
                                                context: context,
                                                builder: (context) {
                                                  return Dialog(
                                                    shape:
                                                        const RoundedRectangleBorder(
                                                          borderRadius:
                                                              BorderRadius.zero,
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
                                                                          BorderRadius.circular(
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
                                                                  onPressed: () async {
                                                                    if (_formKey
                                                                        .currentState!
                                                                        .validate()) {
                                                                      if (d.saveNote ==
                                                                          false) {
                                                                        //  ยังไม่มี — สร้างใหม่
                                                                        await p.createSymptom(
                                                                          dose:
                                                                              d,
                                                                          symptomNote: _saveSymptom
                                                                              .text
                                                                              .trim(),
                                                                          context:
                                                                              context,
                                                                        );
                                                                      } else {
                                                                        // มีแล้ว — อัปเดต
                                                                        await p.editSymptom(
                                                                          dose:
                                                                              d,
                                                                          symptomNote: _saveSymptom
                                                                              .text
                                                                              .trim(),
                                                                          context:
                                                                              context,
                                                                        );
                                                                      }
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
                                                                          BorderRadius.circular(
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
                                              color: d.saveNote
                                                  ? const Color(0xFFFFC800)
                                                  : const Color(0xFFA5A5A5),
                                            ),
                                          ),
                                        ],
                                      ),
                                      //  Column(
                                      //   crossAxisAlignment: CrossAxisAlignment.start,
                                      //   children: d.doses.map((dose) {
                                      //     return Text("${dose.name} ${dose.unit}");
                                      //   },).toList()
                                      //  )
                                      Column(
                                        crossAxisAlignment:
                                            CrossAxisAlignment.start,
                                        children: [
                                          if (d.doses.length > 1) ...[
                                            Text(
                                              d.nameGroup,
                                              style: const TextStyle(
                                                color: Colors.black,
                                                fontSize: 16,
                                              ),
                                            ),
                                            ...d.doses.map(
                                              (dose) => Padding(
                                                padding:
                                                    const EdgeInsetsGeometry.only(
                                                      top: 5,
                                                      bottom: 5,
                                                      left: 15,
                                                    ),
                                                child: Text(
                                                  "• ${dose.name} ${dose.amountPerTime} ${dose.unit}",
                                                  style: const TextStyle(
                                                    color: Colors.black,
                                                    fontSize: 16,
                                                    fontWeight:
                                                        FontWeight.normal,
                                                  ),
                                                ),
                                              ),
                                            ),
                                          ] else if (d.doses.length == 1) ...[
                                            Text(
                                              "${d.doses.first.name} ${d.doses.first.amountPerTime} ${d.doses.first.unit}",
                                              style: const TextStyle(
                                                color: Colors.black,
                                                fontSize: 16,
                                                fontWeight: FontWeight.normal,
                                              ),
                                            ),
                                          ],
                                        ],
                                      ),
                                      const SizedBox(height: 10),
                                    ],
                                  ),
                                ),
                                Column(
                                  crossAxisAlignment: CrossAxisAlignment.end,
                                  children: [
                                    Text(
                                      d.isTaken ? "กินแล้ว" : "ยังไม่กิน",
                                      style: TextStyle(
                                        color: d.isTaken
                                            ? color.AppColors.greenColor
                                            : color.AppColors.redColor,
                                        fontSize: 16,
                                        fontWeight: FontWeight.bold,
                                      ),
                                    ),
                                    const SizedBox(height: 15),
                                    //     Image.asset(
                                    //       "assets/images/pill.png",
                                    //       width: 33,
                                    //       height: 33,
                                    //     ),
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
                itemCount: items.length,
              ),
      ),
    );
  }
}