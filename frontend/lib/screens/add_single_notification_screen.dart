import 'package:flutter/material.dart';
import 'package:frontend/models/notification_info.dart';
import 'package:frontend/providers/add_edit_provider.dart';
import 'package:frontend/providers/add_notification_provider.dart';
import 'package:frontend/providers/add_single_notification_provider.dart';

import 'package:frontend/screens/add_edit_screen.dart';
import 'package:frontend/screens/add_notification_screen.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/providers/drug_provider.dart';
import 'package:provider/provider.dart';
import '../models/dose.dart';

class AddSingleNotification extends StatelessWidget {
  final Dose dose;
  const AddSingleNotification({super.key, required this.dose});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (_) => AddSingleNotificationProvider(dose),
      child: _AddSingleNotificationView(),
    );
  }
}

class _AddSingleNotificationView extends StatefulWidget {
  @override
  State<_AddSingleNotificationView> createState() =>
      _AddSingleNotificationViewState();
}

class _AddSingleNotificationViewState
    extends State<_AddSingleNotificationView> {
  @override
  Widget build(BuildContext context) {
    final addS = context.watch<AddSingleNotificationProvider>();
    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor2nd,
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        foregroundColor: Colors.white,
        title: const Text(
          "รายาการยา",
          style: TextStyle(fontWeight: FontWeight.bold, fontSize: 25),
        ),
      ),
      body: Padding(
        padding: const EdgeInsets.symmetric(vertical: 20, horizontal: 12),
        child: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(
                width: 384,
                child: Card(
                  color: addS.tempDose.import
                      ? const Color(0xFFFFF5D0)
                      : Colors.white,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                    side: const BorderSide(color: Colors.grey, width: 0.5),
                  ),
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(10, 13, 16, 0),
                    child: Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              addS.tempDose.name,
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 20,
                              ),
                            ),
                            const SizedBox(height: 5),
                            Text(
                              addS.tempDose.description,
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 16,
                              ),
                            ),
                            Text(
                              "ครั้งละ " +
                                  addS.tempDose.amountPerDose +
                                  addS.tempDose.unit +
                                  " " +
                                  "วันละ " +
                                  addS.tempDose.frequency +
                                  " " +
                                  "ครั้ง",
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 16,
                              ),
                            ),
                            Text(
                              addS.tempDose.instruction,
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 16,
                              ),
                            ),
                            const SizedBox(height: 10),
                            if (!addS.tempDose.import) ...[
                              SizedBox(
                                width: 95,
                                height: 35,
                                child: ElevatedButton(
                                  onPressed: () {
                                    Navigator.push(
                                      context,
                                      MaterialPageRoute(
                                        builder: (_) => MultiProvider(
                                          providers: [
                                            ChangeNotifierProvider.value(
                                              value: context
                                                  .read<
                                                    AddSingleNotificationProvider
                                                  >(),
                                            ),
                                            ChangeNotifierProvider(
                                              create: (_) => AddEditProvider(
                                                pageFrom: "edit",
                                                editDose: addS.tempDose,
                                              ),
                                            ),
                                          ],
                                          child: const AddEditView(),
                                        ),
                                      ),
                                    );
                                  },
                                  style: ElevatedButton.styleFrom(
                                    backgroundColor: const Color(0xFF55FF00),
                                    elevation: 4,
                                    padding: const EdgeInsets.symmetric(
                                      vertical: 4,
                                    ),
                                    shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(5),
                                    ),
                                  ),
                                  child: const Text(
                                    "แก้ไขข้อมูลยา",
                                    style: TextStyle(
                                      fontSize: 16,
                                      color: Colors.black,
                                    ),
                                    maxLines: 1,
                                    overflow: TextOverflow.visible,
                                  ),
                                ),
                              ),
                            ],
                            const SizedBox(height: 15),
                          ],
                        ),
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.end,
                          children: [
                            Image.asset(
                              addS.tempDose.picture,
                              width: 40,
                              height: 40,
                            ),
                            const SizedBox(height: 40),
                            Text(
                              addS.tempDose.import ? "(โรงพยาบาล)" : "(ของฉัน)",
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 16,
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 30),
              const Text("การแจ้งเตือน", style: TextStyle(fontSize: 20)),
              const SizedBox(height: 15),
              if (addS.savedNotification == null) ...[
                const Text(
                  "ยังไม่มีข้อมูลการแจ้งเตือน",
                  style: TextStyle(color: Color(0xFF959595), fontSize: 20),
                ),
              ] else ...[
                if (addS.savedNotification!.type == "Fixed") ...[
                  Text(
                    "- เวลา: ${addS.savedNotification!.times?.join(', ')} (ทุกวัน)",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                  Text(
                    "- เริ่ม: ${addS.savedNotification!.startDate}",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                  Text(
                    "- สิ้นสุด: ${addS.savedNotification!.endDate}",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                ] else if (addS.savedNotification!.type == "Interval") ...[
                  Text(
                    "- ทุก ${addS.savedNotification!.intervalHours} ชั่วโมง ${addS.savedNotification!.intervalTake} ครั้ง/วัน",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                  Text(
                    "- เริ่ม: ${addS.savedNotification!.times?.join(', ')} ของวันที่ ${addS.savedNotification!.startDate}",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                  Text(
                    "- สิ้นสุด: ${addS.savedNotification!.endDate}",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                ] else if (addS.savedNotification!.type == "DailyWeekly") ...[
                  Text(
                    "- เวลา: ${addS.savedNotification!.times?.join(', ')} (ทุก ${addS.savedNotification!.daysGap} วัน)",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                  Text(
                    "- เริ่ม: ${addS.savedNotification!.startDate}",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                  Text(
                    "- สิ้นสุด: ${addS.savedNotification!.endDate}",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                ] else ...[
                  Text(
                    "- เวลา: ${addS.savedNotification!.times?.join(', ')}",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                  Text(
                    "- กิน ${addS.savedNotification!.takeDays} วัน ต่อเนื่อง พัก ${addS.savedNotification!.breakDays} วัน",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                  Text(
                    "- เริ่ม: ${addS.savedNotification!.startDate}",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                  Text(
                    "- สิ้นสุด: ${addS.savedNotification!.endDate}",
                    style: const TextStyle(
                      color: Color(0xFF959595),
                      fontSize: 20,
                    ),
                  ),
                ],
              ],
              const SizedBox(height: 90),
              SizedBox(
                width: 120,
                height: 35,
                child: ElevatedButton(
                  onPressed: () async {
                    if (addS.savedNotification == null) {
                      final info = await Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (_) => MultiProvider(
                            providers: [
                              ChangeNotifierProvider(
                                create: (_) => AddNotificationProvider(
                                  pageFrom: "single",
                                  dose: addS.tempDose,
                                ),
                              ),
                            ],
                            child: const AddNotificationScreen(),
                          ),
                        ),
                      );
                      if (info != null && info is NotificationInfo) {
                        debugPrint(
                          "ได้ข้อมูลกลับจาก AddNotificationScreen มา Singleแล้ว",
                        );
                        debugPrint("ชนิด: ${info.type}");
                        debugPrint("เวลา: ${info.times}");
                        debugPrint("วันเริ่มต้น: ${info.startDate}");
                        debugPrint("วันสิ้นสุด: ${info.endDate}");

                        addS.saveNotification(info);
                      }
                    } else {
                      addS.clearNotification();
                    }
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: addS.savedNotification == null
                        ? const Color(0xFF55FF00)
                        : const Color(0xFFFF8080),
                    elevation: 4,
                    padding: const EdgeInsets.symmetric(vertical: 4),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(5),
                    ),
                  ),
                  child: Text(
                    addS.savedNotification == null
                        ? "เพิ่มการแจ้งเตือน"
                        : "ลบการแจ้งเตือน",
                    style: const TextStyle(color: Colors.black, fontSize: 16),
                    maxLines: 1,
                  ),
                ),
              ),
              const SizedBox(height: 90),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  if (!addS.tempDose.import) ...[
                    SizedBox(
                      width: 181,
                      height: 70,
                      child: ElevatedButton(
                        onPressed: () {
                          context.read<DrugProvider>().removeDose(
                            addS.tempDose,
                          );
                          Navigator.pop(context);
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFFFF0000),
                          elevation: 4,
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(5),
                          ),
                        ),
                        child: const Text(
                          "ลบรายการยา",
                          style: TextStyle(color: Colors.white, fontSize: 24),
                        ),
                      ),
                    ),
                  ],
                  const SizedBox(width: 15),
                  if (addS.tempDose.import) ...[
                    Center(
                      child: SizedBox(
                        width: 181,
                        height: 70,
                        child: ElevatedButton(
                          onPressed: () {
                            context.read<DrugProvider>().updatedDose(
                              addS.tempDose,
                            );
                            Navigator.pop(context);
                          },
                          style: ElevatedButton.styleFrom(
                            backgroundColor: const Color(0xFF94B4C1),
                            elevation: 4,
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(5),
                            ),
                          ),
                          child: const Text(
                            "บันทึก",
                            style: TextStyle(color: Colors.white, fontSize: 24),
                          ),
                        ),
                      ),
                    ),
                  ] else ...[
                    SizedBox(
                      width: 181,
                      height: 70,
                      child: ElevatedButton(
                        onPressed: () {
                          context.read<DrugProvider>().updatedDose(
                            addS.tempDose,
                          );
                          Navigator.pop(context);
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF94B4C1),
                          elevation: 4,
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(5),
                          ),
                        ),
                        child: const Text(
                          "บันทึก",
                          style: TextStyle(color: Colors.white, fontSize: 24),
                        ),
                      ),
                    ),
                  ],
                ],
              ),
              const SizedBox(height: 70),
            ],
          ),
        ),
      ),
    );
  }
}
