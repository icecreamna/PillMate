import 'package:flutter/material.dart';
import 'package:flutter_datetime_picker_plus/flutter_datetime_picker_plus.dart';
import 'package:frontend/providers/fixed_time_provider.dart';
import 'package:provider/provider.dart';

class FixedTimeWidget extends StatelessWidget {
  const FixedTimeWidget({super.key});

  Future<void> _pickTime(BuildContext context, int index) async {
    final provider = context.read<FixedTimeProvider>();
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
      showTitleActions: true,
      currentTime: initial,
      locale: LocaleType.th,
      onConfirm: (DateTime date) {
        final picked = TimeOfDay(hour: date.hour, minute: date.minute);
        provider.updateTime(index, picked);
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<FixedTimeProvider>(
      builder: (context, provider, _) {
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
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
          ],
        );
      },
    );
  }
}
