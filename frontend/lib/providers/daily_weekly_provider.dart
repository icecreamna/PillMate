import 'package:flutter/material.dart';

class DailyWeeklyProvider extends ChangeNotifier {
  final List<TimeOfDay> _times = [
    const TimeOfDay(hour: 8, minute: 0),
    const TimeOfDay(hour: 13, minute: 0),
    const TimeOfDay(hour: 20, minute: 0),
  ];

  String? notiEveryText;

  String _notiEveryError = "";

  List<TimeOfDay> get times => _times;

  String get notiEveryError => _notiEveryError;

  void setNotiEveryText(String val) {
    notiEveryText = val;
    notifyListeners();
  }

  void updateTime(TimeOfDay newTime, int index) {
    _times[index] = newTime;
    notifyListeners();
  }

  void addTimeOfDay() {
    _times.add(const TimeOfDay(hour: 8, minute: 0));
    notifyListeners();
  }

  void removeTime(int index) {
    _times.removeAt(index);
    notifyListeners();
  }

  String formatThaiTime(TimeOfDay time) {
    final hour = time.hour.toString().padLeft(2, '0');
    final minute = time.minute.toString().padLeft(2, '0');
    return '$hour:$minute น.';
  }

  bool validateNotiEvery(String ne) {
    if (ne.trim().isEmpty) {
      _notiEveryError = "กรุณากรอกค่า";
      notifyListeners();
      return false;
    } else if (int.parse(ne.trim()) < 1) {
      _notiEveryError = "เวลาควรมากกว่า 0";
      notifyListeners();
      return false;
    }
    _notiEveryError = "";
    notifyListeners();
    return true;
  }
}
