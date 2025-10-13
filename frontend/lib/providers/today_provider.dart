import 'package:flutter/material.dart';
import 'package:frontend/services/today_service.dart';
import 'package:intl/intl.dart';

class DoseSingle {
  final String name;
  // final String picture;
  final String amountPerTime;
  final String unit;
  DoseSingle({
    required this.name,
    required this.amountPerTime,
    required this.unit,
  });
}

class DoseGroup {
  final int? id;
  final int? groupId;
  final String nameGroup;
  final String key;
  final List<DoseSingle> doses;
  final DateTime at;
  final String instruction;
  bool saveNote;
  bool isTaken;

  DoseGroup({
    this.id,
    this.groupId,
    required this.nameGroup,
    required this.key,
    required this.doses,
    required this.at,
    required this.instruction,
    this.saveNote = false,
    this.isTaken = false,
  });
}

class TodayProvider extends ChangeNotifier {
  final TodayService _todayService = TodayService();

  bool _isLoading = false;

  DateTime _selected = DateTime.now();

  DateTime get selected => _selected;
  bool get isLoading => _isLoading;

  String get dateLabel {
    final thDate = DateFormat("d MMMM yyyy", "th_TH").format(_selected);
    final buddhistYear = _selected.year + 543;
    return thDate.replaceAll('${_selected.year}', '$buddhistYear');
  }

  // String getTimeHourMinute(String key) => key.split("|").first;
  // String getInstruction(String key) => key.split("|").last;

  List<DoseGroup> all = <DoseGroup>[];
  // List<Dose> get doseSelect {
  //   return all
  //       .where(
  //         (d) =>
  //             d.at.year == _selected.year &&
  //             d.at.month == _selected.month &&
  //             d.at.day == _selected.day,
  //       )
  //       .toList();
  // }//function
  Future<void> loadTodayData() async {
    _isLoading = true;
    notifyListeners();
    try {
      final data = await _todayService.fetchTodayNoti(_selected);
      all = data;
      notifyListeners();
    } catch (e) {
      print("❌ โหลดข้อมูลวันนี้ไม่สำเร็จ: $e");
    } finally {
      await Future.delayed(const Duration(milliseconds: 1000));
      _isLoading = false;
      notifyListeners();
    }
  }

  List<DoseGroup> doseSelect(DateTime? selectTime) {
    if (selectTime == null) return [];
    return all
        .where(
          (d) =>
              d.at.year == selectTime.year &&
              d.at.month == selectTime.month &&
              d.at.day == selectTime.day,
        )
        .toList();
  }

  //function kub

  Future<void> pickDate(BuildContext context) async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _selected,
      firstDate: DateTime(2000),
      lastDate: DateTime(2100),
      helpText: 'เลือกวันที่',
      cancelText: 'ยกเลิก',
      confirmText: 'ตกลง',
      locale: const Locale('th', 'TH'),
    );
    if (picked != null) {
      _selected = DateTime(picked.year, picked.month, picked.day);
      notifyListeners();
    }
  }

  void setIsTaken(bool taken, DoseGroup doseGroup) {
    doseGroup.isTaken = taken;
    notifyListeners();
  }

  void removeDose(DoseGroup doseGroup) {
    all.remove(doseGroup);
    notifyListeners();
  }

  void setNote(bool saveNote, DoseGroup doseGroup) {
    doseGroup.saveNote = saveNote;
    notifyListeners();
  }
}
