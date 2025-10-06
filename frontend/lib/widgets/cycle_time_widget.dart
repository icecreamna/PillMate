import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_datetime_picker_plus/flutter_datetime_picker_plus.dart';
import 'package:frontend/providers/cycle_time_provider.dart';
import 'package:provider/provider.dart';

class CycleTimeWidget extends StatefulWidget {
  const CycleTimeWidget({super.key});

  @override
  State<CycleTimeWidget> createState() => _CycleTimeWidgetState();
}

class _CycleTimeWidgetState extends State<CycleTimeWidget> {
  final _breakDaysController = TextEditingController();
  final _inTakeDaysController = TextEditingController();

  @override
  void dispose() {
    // TODO: implement dispose
    _breakDaysController.dispose();
    _inTakeDaysController.dispose();
    super.dispose();
  }

  OutlineInputBorder _borderInput(Color color) {
    return OutlineInputBorder(
      borderRadius: BorderRadius.circular(13),
      borderSide: BorderSide(width: 1, color: color),
    );
  }

  Future<void> _pickTime(BuildContext context, int index) async {
    final provider = context.read<CycleTimeProvider>();
    final current = provider.times[index];

    final now = DateTime.now();
    final initial = DateTime(
      now.year,
      now.month,
      now.day,
      current.hour,
      current.minute,
    );

    DatePicker.showTimePicker(
      context,
      currentTime: initial,
      showTitleActions: true,
      showSecondsColumn: false,
      locale: LocaleType.th,
      onConfirm: (time) {
        final picked = TimeOfDay(hour: time.hour, minute: time.minute);
        provider.updateTime(index, picked);
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<CycleTimeProvider>(
      builder: (context, provider, _) {
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Text("ทานต่อเนื่อง", style: TextStyle(fontSize: 20)),
                const SizedBox(width: 10),
                Container(
                  width: 81,
                  height: 50,
                  color: Colors.white,
                  child: TextFormField(
                    inputFormatters: [
                      FilteringTextInputFormatter.allow(RegExp('[0-9]')),
                      LengthLimitingTextInputFormatter(2),
                    ],
                    keyboardType: TextInputType.number,
                    decoration: InputDecoration(
                      fillColor: Colors.white,
                      enabledBorder: _borderInput(const Color(0xFF8D8D8D)),
                      focusedBorder: _borderInput(const Color(0xFF8D8D8D)),
                      errorBorder: _borderInput(const Color(0xFF8D8D8D)),
                      focusedErrorBorder: _borderInput(const Color(0xFF8D8D8D)),
                    ),
                    controller: _inTakeDaysController,
                  ),
                ),
                const SizedBox(width: 10),
                const Text("วัน", style: TextStyle(fontSize: 20)),
                const SizedBox(width: 20),
                Visibility(
                  visible: provider.inTakeDaysError.isNotEmpty,
                  child: Text(
                    provider.inTakeDaysError,
                    style: const TextStyle(
                      fontSize: 12,
                      color: Color(0xFFFF0000),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 10),
            Row(
              children: [
                const Text("พักทาน", style: TextStyle(fontSize: 20)),
                const SizedBox(width: 50),
                Container(
                  width: 81,
                  height: 50,
                  color: Colors.white,
                  child: TextFormField(
                    inputFormatters: [
                      FilteringTextInputFormatter.allow(RegExp('[0-9]')),
                      LengthLimitingTextInputFormatter(2),
                    ],
                    keyboardType: TextInputType.number,
                    decoration: InputDecoration(
                      fillColor: Colors.white,
                      enabledBorder: _borderInput(const Color(0xFF8D8D8D)),
                      focusedBorder: _borderInput(const Color(0xFF8D8D8D)),
                      errorBorder: _borderInput(const Color(0xFF8D8D8D)),
                      focusedErrorBorder: _borderInput(const Color(0xFF8D8D8D)),
                    ),
                    controller: _breakDaysController,
                  ),
                ),
                const SizedBox(width: 10),
                const Text("วัน", style: TextStyle(fontSize: 20)),
                const SizedBox(width: 20),
                Visibility(
                  visible: provider.breakDayError.isNotEmpty,
                  child: Text(
                    provider.breakDayError,
                    style: const TextStyle(
                      fontSize: 12,
                      color: Color(0xFFFF0000),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 20),
            Row(
              children: [
                const Expanded(
                  child: Text("เวลาใช้ยา", style: TextStyle(fontSize: 20)),
                ),
                ElevatedButton(
                  onPressed: () => provider.addTimeOfDay(),
                  style: ElevatedButton.styleFrom(
                    minimumSize: const Size(70, 33),
                    elevation: 4,
                    backgroundColor: const Color(0xFF29AEDE),
                    shape: BeveledRectangleBorder(
                      borderRadius: BorderRadius.circular(5),
                    ),
                  ),
                  child: const Text.rich(
                    TextSpan(
                      children: [
                        WidgetSpan(
                          alignment: PlaceholderAlignment.middle,
                          child: Text(
                            "+  ",
                            style: TextStyle(fontSize: 20, color: Colors.white),
                          ),
                        ),
                        TextSpan(
                          text: "เพิ่มเวลา",
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 20,
                            letterSpacing: 0,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 5),
            ...List.generate(provider.times.length, (index) {
              final t = provider.times[index];
              return Row(
                children: [
                  Column(
                    children: [
                      Row(
                        children: [
                          Container(
                            width: 171,
                            height: 50,
                            margin: const EdgeInsets.only(bottom: 10),
                            padding: const EdgeInsets.symmetric(horizontal: 20),
                            decoration: BoxDecoration(
                              color: Colors.white,
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(width: 1, color: Colors.grey),
                            ),
                            child: Row(
                              children: [
                                Expanded(
                                  child: Text(
                                    provider.formatThaiTime(t),
                                    style: const TextStyle(fontSize: 16),
                                  ),
                                ),
                                IconButton(
                                  onPressed: () => _pickTime(context, index),
                                  padding: EdgeInsets.zero,
                                  icon: const Icon(
                                    Icons.timer_outlined,
                                    color: Color(0xFF29AEDE),
                                    size: 36,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          if (index != 0) ...[
                            Center(
                              child: IconButton(
                                onPressed: () => provider.removeTime(index),
                                icon: const Icon(
                                  Icons.remove_circle_outline,
                                  size: 40,
                                  color: Color(0xFFFF0000),
                                ),
                              ),
                            ),
                          ],
                        ],
                      ),
                    ],
                  ),
                ],
              );
            }),
            const SizedBox(height: 60),
            Center(
              child: ElevatedButton(
                onPressed: () {
                  final breakDayValid = provider.validateBreakDays(_breakDaysController.text);
                  final inTakeDayValid = provider.validateInTakeDays(_inTakeDaysController.text);
                  if(!breakDayValid || !inTakeDayValid) return ;
                  Navigator.pop(context);
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
        );
      },
    );
  }
}
