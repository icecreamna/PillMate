import 'package:flutter/material.dart';

class IntervalProvider extends ChangeNotifier {
  TimeOfDay _times = const TimeOfDay(hour: 8, minute: 0);

  TimeOfDay get times => _times;

  String _errorHour = "";
  String _errorTake = "";

  String get errorHour => _errorHour;
  String get errorTake => _errorTake;

  bool validateHour(String hour) {
    if (hour.trim().isEmpty) {
      _errorHour = "กรุณากรอกค่า";
      notifyListeners();
      return false;
    } else if (int.parse(hour.trim()) < 1) {
      _errorHour = "เวลาควรมากกว่า 0";
      notifyListeners();
      return false;
    } else if (int.parse(hour.trim()) > 24) {
      _errorHour = "เวลาควร <= 24";
      notifyListeners();
      return false;
    }
    _errorHour = "";
    notifyListeners();
    return true;
  }

  bool validateTake(String takePerDay) {
    if (takePerDay.trim().isEmpty) {
      _errorTake = "กรุณากรอกค่า";
      notifyListeners();
      return false;
    } else if (int.parse(takePerDay.trim()) < 1) {
      _errorTake = "ใส่ค่ามากกว่า 0";
      notifyListeners();
      return false;
    }
    _errorTake = "";
    notifyListeners();
    return true;
  }

  void updateTime(TimeOfDay newTime) {
    _times = newTime;
    notifyListeners();
  }

  String formatThaiTime(TimeOfDay time) {
    final hour = time.hour.toString().padLeft(2, '0');
    final minute = time.minute.toString().padLeft(2, '0');
    return '$hour:$minute น.';
  }
}
