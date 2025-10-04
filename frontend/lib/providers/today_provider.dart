import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

class Dose {
  final String name;
  // final String picture;
  final String unit;
  Dose({required this.name, required this.unit});
}

class DoseGroup {
  final String nameGroup;
  final String key;
  final List<Dose> doses;
  final DateTime at;
  final String instruction;
  bool saveNote;
  bool isTaken;

  DoseGroup({
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
  DateTime _selected = DateTime.now();

  DateTime get selected => _selected;

  String get dateLabel {
    final thDate = DateFormat("d MMMM yyyy", "th_TH").format(_selected);
    final buddhistYear = _selected.year + 543;
    return thDate.replaceAll('${_selected.year}', '$buddhistYear');
  }

  // String getTimeHourMinute(String key) => key.split("|").first;
  // String getInstruction(String key) => key.split("|").last;

  final List<DoseGroup> all = <DoseGroup>[
    DoseGroup(
      key: "08:00-หลังอาหาร",
      nameGroup: "กลุ่ม 1",
      at: DateTime(2025, 9, 14, 8, 0),
      instruction: "หลังอาหาร",
      doses: [Dose(name: "ยา1", unit: "1 เม็ด")],
    ),
    DoseGroup(
      key: "09:00-ก่อนอาหาร",
      nameGroup: "กลุ่ม 2",
      at: DateTime(2025, 9, 14, 9, 0),
      instruction: "ก่อนอาหาร",
      doses: [Dose(name: "ยา2", unit: "1 เม็ด")],
    ),
    DoseGroup(
      key: "13:00-หลังอาหาร",
      nameGroup: "กลุ่ม 3",
      at: DateTime(2025, 9, 14, 13, 0),
      instruction: "หลังอาหาร",
      doses: [
        Dose(name: "ยา3", unit: "1 เม็ด"),
        Dose(name: "แคปซูล1", unit: "1 เม็ด"),
        Dose(name: "แคปซูล2", unit: "1 เม็ด"),
        Dose(name: "ยาน้ำ1", unit: "1 ช้อน"),
        Dose(name: "ยาน้ำ2", unit: "1 ช้อน"),
      ],
    ),
    DoseGroup(
      nameGroup: "กลุ่ม 4",
      key: "13:00-ก่อนอาหาร",
      doses: [
        Dose(name: "ยา4", unit: "3 เม็ด"),
        Dose(name: "ยา5", unit: "5 ช้อน"),
        Dose(name: "ยาบ้า", unit: "100 ขวด"),
      ],
      at: DateTime(2025, 9, 14, 13, 0),
      instruction: "ก่อนอาหาร",
    ),
  ];
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
