import 'package:flutter/material.dart';

class CycleTimeProvider extends ChangeNotifier {
  final List<TimeOfDay> _times = [
    const TimeOfDay(hour: 8, minute: 0),
    const TimeOfDay(hour: 13, minute: 0),
    const TimeOfDay(hour: 20, minute: 0),
  ];

  String _breakDaysError = "";
  String _inTakeDaysError = "";

  List<TimeOfDay> get times => _times;
  String get breakDayError => _breakDaysError;
  String get inTakeDaysError => _inTakeDaysError;

  void addTimeOfDay() {
    _times.add(const TimeOfDay(hour: 8, minute: 0));
    notifyListeners();
  }

  void updateTime(int index, TimeOfDay newTime) {
    _times[index] = newTime;
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

  bool validateBreakDays(String bd) {
    if (bd.trim().isEmpty) {
      _breakDaysError = "กรุณากรอกค่า";
      notifyListeners();
      return false;
    }else if(int.parse(bd.trim()) < 1){
      _breakDaysError = "ค่าต้องมากกว่า 0";
      notifyListeners();
      return false;
    }
    _breakDaysError = "";
    notifyListeners();
    return true;
  }

  bool validateInTakeDays(String bd) {
    if (bd.trim().isEmpty) {
      
      _inTakeDaysError = "กรุณากรอกค่า";
      notifyListeners();
      return false;
    }else if(int.parse(bd.trim()) < 1){
      _inTakeDaysError = "ค่าต้องมากกว่า 0";
      notifyListeners();
      return false;
    }
    _inTakeDaysError = "" ;
    notifyListeners();
    return true;
  }
}
