import 'package:flutter/material.dart';
import 'package:frontend/providers/today_provider.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';

class TodayScreen extends StatelessWidget {
  const TodayScreen({super.key});
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
        child: p.doseSelect.isEmpty
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
                  final d = p.doseSelect[i];
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
                                        ],
                                      ),
                                      const SizedBox(height: 7),
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
                                    Image.asset(
                                      d.picture,
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
                itemCount: p.doseSelect.length,
              ),
      ),
    );
  }
}
