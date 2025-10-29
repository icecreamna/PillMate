import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_datetime_picker_plus/flutter_datetime_picker_plus.dart';
import 'package:frontend/providers/daily_weekly_provider.dart';
import 'package:provider/provider.dart';

class DailyWeeklyTimeWidget extends StatefulWidget {
  const DailyWeeklyTimeWidget({super.key});

  @override
  State<DailyWeeklyTimeWidget> createState() => _DailyWeeklyTimeWidgetState();
}

class _DailyWeeklyTimeWidgetState extends State<DailyWeeklyTimeWidget> {
  final _notiEveryController = TextEditingController();

  @override
  void dispose() {
    _notiEveryController.dispose();
    super.dispose();
  }

  Future<void> _pickTime(BuildContext context, int index) async {
    final provider = context.read<DailyWeeklyProvider>();
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
      showSecondsColumn: false,
      currentTime: initial,
      showTitleActions: true,
      locale: LocaleType.th,
      onConfirm: (time) {
        final picked = TimeOfDay(hour: time.hour, minute: time.minute);
        provider.updateTime(picked, index);
      },
    );
  }

  OutlineInputBorder _borderInput(Color color) {
    return OutlineInputBorder(
      borderRadius: BorderRadius.circular(13),
      borderSide: BorderSide(width: 1, color: color),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<DailyWeeklyProvider>(
      builder: (context, provider, _) {
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Text("แจ้งเตือนทุก", style: TextStyle(fontSize: 20)),
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
                    controller: _notiEveryController,
                    onChanged: (val) => context
                        .read<DailyWeeklyProvider>()
                        .setNotiEveryText(val),
                  ),
                ),
                const SizedBox(width: 10),
                const Text("วัน", style: TextStyle(fontSize: 20)),
                const SizedBox(width: 20),
                Visibility(
                  visible: provider.notiEveryError.isNotEmpty,
                  child: Text(
                    provider.notiEveryError,
                    style: const TextStyle(
                      fontSize: 12,
                      color: Color(0xFFFF0000),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 40),
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
          ],
        );
      },
    );
  }
}
