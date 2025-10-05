import 'package:flutter/material.dart';

class FixedTimeProvider extends ChangeNotifier {
  final List<TimeOfDay> _times = [
    const TimeOfDay(hour: 8, minute: 0),
    const TimeOfDay(hour: 13, minute: 0),
    const TimeOfDay(hour: 20, minute: 0),
  ];

  List<TimeOfDay> get times => _times;

  void addTimeOfDay() {
    _times.add(const TimeOfDay(hour: 8, minute: 0));
    notifyListeners();
  }

  void updateTime(int index, TimeOfDay newTime) {
    _times[index] = newTime;
    notifyListeners();
  }

  void removeTime(int index){
    _times.removeAt(index);
    notifyListeners();
  }

  String formatThaiTime(TimeOfDay time) {
  final hour = time.hour.toString().padLeft(2, '0');
  final minute = time.minute.toString().padLeft(2, '0');
  return '$hour:$minute น.';
}
}
