import 'package:flutter/material.dart';
import 'package:frontend/models/notification_info.dart';
import 'package:frontend/providers/add_notification_provider.dart';
import 'package:frontend/providers/cycle_time_provider.dart';
import 'package:frontend/providers/daily_weekly_provider.dart';
import 'package:frontend/providers/fixed_time_provider.dart';
import 'package:frontend/providers/interval_provider.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/widgets/interval_time_widget.dart';
import 'package:frontend/widgets/cycle_time_widget.dart';
import 'package:frontend/widgets/daily_weekly_time_widget.dart';
import 'package:frontend/widgets/fixed_time_widget.dart';
import 'package:provider/provider.dart';

class AddNotificationScreen extends StatefulWidget {
  const AddNotificationScreen({super.key});

  @override
  State<AddNotificationScreen> createState() => _AddNotificationScreenState();
}

class _AddNotificationScreenState extends State<AddNotificationScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<AddNotificationProvider>().loadNotiFormats();
    });
  }

  @override
  Widget build(BuildContext context) {
    final addN = context.watch<AddNotificationProvider>();
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => FixedTimeProvider()),
        ChangeNotifierProvider(create: (_) => IntervalProvider()),
        ChangeNotifierProvider(create: (_) => DailyWeeklyProvider()),
        ChangeNotifierProvider(create: (_) => CycleTimeProvider()),
      ],
      builder: (context, _) {
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
                          side: const BorderSide(
                            color: Colors.grey,
                            width: 0.5,
                          ),
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
                          side: const BorderSide(
                            color: Colors.grey,
                            width: 0.5,
                          ),
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
                                    if (addN.dose!.import) ...[
                                      Text(
                                        "วิธีการกิน: ${addN.dose!.note ?? ""}",
                                        style: const TextStyle(
                                          color: Colors.black,
                                          fontSize: 16,
                                        ),
                                      ),
                                      const SizedBox(height: 5),
                                      Text(
                                        "วันเริ่มทานยา: " +
                                                addN.dose!.startDate.toString() ??
                                            "",
                                        style: const TextStyle(
                                          color: Colors.black,
                                          fontSize: 16,
                                        ),
                                      ),
                                      const SizedBox(height: 5),
                                      Text(
                                        "วันหยุดทานยา: " +
                                                addN.dose!.endDate.toString() ??
                                            "",
                                        style: const TextStyle(
                                          color: Colors.black,
                                          fontSize: 16,
                                        ),
                                      ),
                                    ],
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
                                    addN.dose!.import
                                        ? "(โรงพยาบาล)"
                                        : "(ของฉัน)",
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
                  const Text(
                    "รูปแบบการแจ้งเตือน",
                    style: TextStyle(fontSize: 20),
                  ),
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
                      child: addN.isLoading
                          ? const Center(child: CircularProgressIndicator())
                          : DropdownButton<String>(
                              value: addN.selectedType.isEmpty
                                  ? null
                                  : addN.selectedType,
                              hint: const Text(
                                "เลือกรูปแบบการแจ้งเตือน",
                                style: TextStyle(
                                  fontSize: 18,
                                  color: Colors.black87,
                                ),
                              ),
                              items: addN.formats.map((f) {
                                return DropdownMenuItem<String>(
                                  value: f.id.toString(),
                                  child: Text(
                                    f.formatName,
                                    style: const TextStyle(fontSize: 20),
                                  ),
                                );
                              }).toList(),
                              onChanged: (value) => addN.setSelectType(value!),
                            ),
                    ),
                  ),
                  const SizedBox(height: 20),
                  if (addN.getTypeName(addN.selectedType) == "Fixed") ...[
                    const FixedTimeWidget(),
                  ],
                  if (addN.getTypeName(addN.selectedType) == "Interval") ...[
                    const IntervalTimeWidget(),
                  ],
                  if (addN.getTypeName(addN.selectedType) == "DailyWeekly") ...[
                    const DailyWeeklyTimeWidget(),
                  ],

                  if (addN.getTypeName(addN.selectedType) == "Cycle") ...[
                    const CycleTimeWidget(),
                  ],
                  Center(
                    child: ElevatedButton(
                      onPressed: () async {
                        bool valid = true;

                        if (addN.startDate.isEmpty || addN.endDate.isEmpty) {
                          valid = false;
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text(
                                "กรุณาเลือกวันเริ่มต้นและวันสิ้นสุด",
                              ),
                              backgroundColor: Colors.red,
                            ),
                          );
                        }
                        switch (addN.getTypeName(addN.selectedType)) {
                          case "Fixed":
                            final fixed = context.read<FixedTimeProvider>();
                            if (fixed.times.isEmpty) {
                              valid = false;
                              ScaffoldMessenger.of(context).showSnackBar(
                                const SnackBar(
                                  content: Text(
                                    "กรุณาเพิ่มเวลาอย่างน้อย 1 เวลา",
                                  ),
                                  backgroundColor: Colors.red,
                                ),
                              );
                            }
                            break;

                          case "Interval":
                            final interval = context.read<IntervalProvider>();
                            final ok1 = interval.validateHour(
                              interval.hourText ?? "",
                            );
                            final ok2 = interval.validateTake(
                              interval.takeText ?? "",
                            );
                            if (!ok1 || !ok2) valid = false;
                            break;

                          case "DailyWeekly":
                            final daily = context.read<DailyWeeklyProvider>();
                            if (!daily.validateNotiEvery(
                              daily.notiEveryText ?? "",
                            )) {
                              valid = false;
                            }
                            break;

                          case "Cycle":
                            final cycle = context.read<CycleTimeProvider>();
                            final ok1 = cycle.validateInTakeDays(
                              cycle.inTakeText ?? "",
                            );
                            final ok2 = cycle.validateBreakDays(
                              cycle.breakText ?? "",
                            );
                            if (!ok1 || !ok2) valid = false;
                            break;

                          default:
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(
                                content: Text("กรุณาเลือกรูปแบบการแจ้งเตือน"),
                                backgroundColor: Colors.red,
                              ),
                            );
                            valid = false;
                        }

                        if (!valid) return;

                        List<String>? selectedTimes;
                        String formatThaiTime(TimeOfDay t) {
                          final hour = t.hour.toString().padLeft(2, '0');
                          final minute = t.minute.toString().padLeft(2, '0');
                          return "$hour:$minute";
                        }

                        if (addN.getTypeName(addN.selectedType) == 'Fixed') {
                          selectedTimes = context
                              .read<FixedTimeProvider>()
                              .times
                              .map((t) => formatThaiTime(t))
                              .toList();
                        } else if (addN.getTypeName(addN.selectedType) ==
                            'DailyWeekly') {
                          selectedTimes = context
                              .read<DailyWeeklyProvider>()
                              .times
                              .map((t) => formatThaiTime(t))
                              .toList();
                        } else if (addN.getTypeName(addN.selectedType) ==
                            'Cycle') {
                          selectedTimes = context
                              .read<CycleTimeProvider>()
                              .times
                              .map((t) => formatThaiTime(t))
                              .toList();
                        } else if (addN.getTypeName(addN.selectedType) ==
                            'Interval') {
                          final start = context.read<IntervalProvider>().times;
                          selectedTimes = [formatThaiTime(start)];
                        }
                        final info = NotificationInfo(
                          type: addN.selectedType,
                          startDate: addN.startDate,
                          endDate: addN.endDate,
                          times: selectedTimes,
                          intervalHours:
                              (addN.getTypeName(addN.selectedType) ==
                                  'Interval')
                              ? int.tryParse(
                                  context.read<IntervalProvider>().hourText ??
                                      '',
                                )
                              : null,
                          intervalTake:
                              (addN.getTypeName(addN.selectedType) ==
                                  "Interval")
                              ? int.tryParse(
                                  context.read<IntervalProvider>().takeText ??
                                      '',
                                )
                              : null,
                          daysGap:
                              (addN.getTypeName(addN.selectedType) ==
                                  'DailyWeekly')
                              ? int.tryParse(
                                  context
                                          .read<DailyWeeklyProvider>()
                                          .notiEveryText ??
                                      '',
                                )
                              : null,
                          takeDays:
                              (addN.getTypeName(addN.selectedType) == 'Cycle')
                              ? int.tryParse(
                                  context
                                          .read<CycleTimeProvider>()
                                          .inTakeText ??
                                      '',
                                )
                              : null,
                          breakDays:
                              (addN.getTypeName(addN.selectedType) == 'Cycle')
                              ? int.tryParse(
                                  context.read<CycleTimeProvider>().breakText ??
                                      '',
                                )
                              : null,
                        );
                        final success = await addN.addNotification(info);
                        if (success && context.mounted) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text("✅ เพิ่มการแจ้งเตือนสำเร็จ"),
                            ),
                          );
                          Navigator.pop(context, true);
                        } else {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text("❌ เพิ่มการแจ้งเตือนล้มเหลว"),
                            ),
                          );
                        }
                      },
                      style: ElevatedButton.styleFrom(
                        elevation: 4,
                        backgroundColor: const Color(0xFF55FF00),
                        minimumSize: const Size(175, 38),
                        shape: BeveledRectangleBorder(
                          borderRadius: BorderRadius.circular(5),
                        ),
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
      },
    );
  }
}
