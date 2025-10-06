import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_datetime_picker_plus/flutter_datetime_picker_plus.dart';
import 'package:frontend/providers/interval_provider.dart';
import 'package:provider/provider.dart';

class IntervalTimeWidget extends StatefulWidget {
  const IntervalTimeWidget({super.key});

  @override
  State<IntervalTimeWidget> createState() => _IntervalTimeWidgetState();
}

class _IntervalTimeWidgetState extends State<IntervalTimeWidget> {
  final _formKey = GlobalKey<FormState>();

  final _hourController = TextEditingController();
  final _numberOfTakeController = TextEditingController();

  OutlineInputBorder _borderInput(Color color) {
    return OutlineInputBorder(
      borderRadius: BorderRadius.circular(13),
      borderSide: BorderSide(width: 1, color: color),
    );
  }

  Future<void> _pickTime(BuildContext context) async {
    final provider = context.read<IntervalProvider>();
    final current = provider.times;

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
        provider.updateTime(picked);
      },
    );
  }

  @override
  void dispose() {
    _hourController.dispose();
    _numberOfTakeController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final provider = context.watch<IntervalProvider>();
    return Form(
      key: _formKey,
      child: Column(
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
                  controller: _hourController,
                ),
              ),
              const SizedBox(width: 10),
              const Text("ชั่วโมง", style: TextStyle(fontSize: 20)),
              const SizedBox(width: 20),
              Visibility(
                visible: provider.errorHour.isNotEmpty,
                child: Text(
                  provider.errorHour,
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
              const Text("จำนวน", style: TextStyle(fontSize: 20)),
              const SizedBox(width: 58),
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
                  controller: _numberOfTakeController,
                ),
              ),
              const SizedBox(width: 10),
              const Text("ครั้ง/วัน", style: TextStyle(fontSize: 20)),
              const SizedBox(width: 20),
              Visibility(
                visible: provider.errorTake.isNotEmpty,
                child: Text(
                  provider.errorTake,
                  style: const TextStyle(
                    fontSize: 12,
                    color: Color(0xFFFF0000),
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 40),
          const Text("เวลาทานยา", style: TextStyle(fontSize: 20)),
          const SizedBox(height: 10),
          Row(
            children: [
              const Text("เวลาเริ่มต้น", style: TextStyle(fontSize: 20)),
              const SizedBox(width: 30),
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
                        provider.formatThaiTime(provider.times),
                        style: const TextStyle(fontSize: 16),
                      ),
                    ),
                    IconButton(
                      onPressed: () => _pickTime(context),
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
            ],
          ),
          const SizedBox(height: 60),
          Center(
            child: ElevatedButton(
              onPressed: () {
                final validHour = provider.validateHour(_hourController.text);
                final validTake = provider.validateTake(_numberOfTakeController.text);
                if (!validTake || !validHour) return;
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
      ),
    );
  }
}
