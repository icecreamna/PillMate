import 'package:flutter/material.dart';
import 'package:frontend/providers/add_notification_provider.dart';
import 'package:frontend/providers/fixed_time_provider.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/widgets/Interval_time_widget.dart';
import 'package:frontend/widgets/cycle_time_widget.dart';
import 'package:frontend/widgets/daily_weekly_time_widget.dart';
import 'package:frontend/widgets/fixed_time_widget.dart';
import 'package:provider/provider.dart';

class AddNotificationScreen extends StatelessWidget {
  const AddNotificationScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final addN = context.watch<AddNotificationProvider>();
    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor2nd,
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        foregroundColor: Colors.white,
        title: const Text(
          "เพิ่มรายการแจ้งเตือน",
          style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
        ),
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              if (addN.pageFrom == "group") ...[
                Container(
                  width: 384,
                  height: 80,
                  margin: const EdgeInsets.only(bottom: 15),
                  child: Card(
                    color: Colors.white,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                      side: const BorderSide(color: Colors.grey, width: 0.5),
                    ),
                    child: Padding(
                      padding: const EdgeInsets.symmetric(
                        vertical: 5,
                        horizontal: 10,
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            addN.keyName!,
                            style: const TextStyle(
                              color: Colors.black,
                              fontSize: 20,
                            ),
                          ),
                          const SizedBox(height: 7),
                          Text(
                            "${addN.value!.length} รายการ",
                            style: const TextStyle(
                              color: Color(0xFF454545),
                              fontSize: 16,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ] else ...[
                Container(
                  width: 384,
                  margin: const EdgeInsets.only(bottom: 15),
                  child: Card(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                      side: const BorderSide(color: Colors.grey, width: 0.5),
                    ),
                    color: addN.dose!.import
                        ? const Color(0xFFFFF5D0)
                        : Colors.white,
                    child: Padding(
                      padding: const EdgeInsets.fromLTRB(10, 13, 16, 0),
                      child: Row(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  addN.dose!.name,
                                  style: const TextStyle(
                                    color: Colors.black,
                                    fontSize: 20,
                                  ),
                                ),
                                const SizedBox(height: 5),
                                Text(
                                  addN.dose!.description,
                                  style: const TextStyle(
                                    color: Colors.black,
                                    fontSize: 16,
                                  ),
                                ),
                                Text(
                                  "ครั้งละ " +
                                      addN.dose!.amountPerDose +
                                      addN.dose!.unit +
                                      " " +
                                      "วันละ " +
                                      addN.dose!.frequency +
                                      " " +
                                      "ครั้ง",
                                  style: const TextStyle(
                                    color: Colors.black,
                                    fontSize: 16,
                                  ),
                                ),
                                Text(
                                  addN.dose!.instruction,
                                  style: const TextStyle(
                                    color: Colors.black,
                                    fontSize: 16,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          Column(
                            crossAxisAlignment: CrossAxisAlignment.end,
                            children: [
                              Image.asset(
                                addN.dose!.picture,
                                width: 40,
                                height: 40,
                              ),
                              const SizedBox(height: 40),
                              Text(
                                addN.dose!.import ? "(โรงพยาบาล)" : "(ของฉัน)",
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
              ],
              const Row(
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  Text(
                    "วันเริ่มต้น",
                    style: TextStyle(
                      fontSize: 20,
                      wordSpacing: 0,
                      letterSpacing: 0,
                    ),
                  ),
                  SizedBox(width: 140),
                  Text(
                    "วันสิ้นสุด",
                    style: TextStyle(
                      fontSize: 20,
                      wordSpacing: 0,
                      letterSpacing: 0,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 6),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    width: 171,
                    height: 50,
                    decoration: BoxDecoration(
                      color: Colors.white,
                      border: Border.all(color: Colors.grey, width: 1),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                      children: [
                        Text(
                          addN.startDate,
                          style: const TextStyle(fontSize: 16),
                        ),
                        IconButton(
                          onPressed: () => addN.pickStartDate(context),
                          padding: EdgeInsets.zero,
                          icon: const Icon(
                            Icons.calendar_month_rounded,
                            color: Color(0xFF29AEDE),
                            size: 32,
                          ),
                        ),
                      ],
                    ),
                  ),
                  Container(
                    width: 171,
                    height: 50,
                    decoration: BoxDecoration(
                      color: Colors.white,
                      border: Border.all(color: Colors.grey, width: 1),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                      children: [
                        Text(
                          addN.endDate,
                          style: const TextStyle(fontSize: 16),
                        ),
                        IconButton(
                          onPressed: () => addN.pickEndDate(context),
                          padding: EdgeInsets.zero,
                          icon: const Icon(
                            Icons.calendar_month_rounded,
                            color: Color(0xFF29AEDE),
                            size: 32,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 20),
              const Text("รูปแบบการแจ้งเตือน", style: TextStyle(fontSize: 20)),
              const SizedBox(height: 15),
              Container(
                width: 384,
                height: 50,
                padding: const EdgeInsets.symmetric(horizontal: 15),
                decoration: BoxDecoration(
                  color: Colors.white,
                  border: Border.all(color: Colors.black, width: 0.9),
                  borderRadius: BorderRadius.circular(5),
                ),
                child: DropdownButtonHideUnderline(
                  child: DropdownButton<String>(
                    value: addN.selectedType.isEmpty ? null : addN.selectedType,
                    hint: const Text(
                      "เลือกรูปแบบการแจ้งเตือน",
                      style: TextStyle(fontSize: 18, color: Colors.black87),
                    ),
                    items: const [
                      DropdownMenuItem(
                        value: "Fixed",
                        child: Text(
                          "เวลาเฉพาะ (Fixed Times)",
                          style: TextStyle(fontSize: 20),
                        ),
                      ),
                      DropdownMenuItem(
                        value: "Interval",
                        child: Text(
                          "ทุกกี่ชั่วโมง (Interval)",
                          style: TextStyle(fontSize: 20),
                        ),
                      ),
                      DropdownMenuItem(
                        value: "DailyWeekly",
                        child: Text(
                          "วันเว้นวัน / ทุกกี่วัน",
                          style: TextStyle(fontSize: 20),
                        ),
                      ),
                      DropdownMenuItem(
                        value: "Cycle",
                        child: Text(
                          "ทานต่อเนื่อง/พักยา (Cycle)",
                          style: TextStyle(fontSize: 20),
                        ),
                      ),
                    ],
                    onChanged: (value) => addN.setSelectType(value!),
                  ),
                ),
              ),
              const SizedBox(height: 20),
              if (addN.selectedType == "Fixed") ...[
                ChangeNotifierProvider(
                  create: (_) => FixedTimeProvider(),
                  child: const FixedTimeWidget(),
                ),
              ],
              if (addN.selectedType == "Interval") const IntervalTimeWidget(),
              if (addN.selectedType == "DailyWeekly")
                const DailyWeeklyTimeWidget(),
              if (addN.selectedType == "Cycle") const CycleTimeWidget(),
              const SizedBox(height: 60,) ,
              Center(
                child: ElevatedButton(
                  onPressed: () {
                    Navigator.pop(context);
                  },
                  style: ElevatedButton.styleFrom(
                    elevation: 4,
                    backgroundColor: const Color(0xFF55FF00),
                    minimumSize: const Size(175, 38),
                    shape: BeveledRectangleBorder(
                      borderRadius: BorderRadius.circular(5)
                    )
                  ),
                  child: const Text(
                    "เพิ่มการแจ้งเตือน",
                    style: TextStyle(color: Colors.white, fontSize: 24),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
