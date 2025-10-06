import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

class Dose {
  final String name;
  final DateTime at;
  final String instruction;
  final String picture;
  final String unit;
   bool isTake;
  Dose({
    required this.name,
    required this.at,
    required this.instruction,
    required this.picture,
    required this.unit,
    this.isTake = false,
  });
}

class DoseGroup {
  final String key;
  final List<Dose> doses;
  bool isTaken;

  DoseGroup({required this.key, required this.doses, this.isTaken = false});
}

class TodayProvider extends ChangeNotifier {
  DateTime _selected = DateTime.now();

  DateTime get selected => _selected;

  String get dateLabel => DateFormat("MMM d, yyyy").format(_selected);

  String getTimeHourMinute(String key) => key.split("|").first;
  String getInstruction(String key) => key.split("|").last;

  final List<Dose> all = <Dose>[
    Dose(
      name: "ยา1",
      at: DateTime(2025, 9, 4, 8, 0),
      instruction: "หลังอาหาร",
      picture: "assets/images/pill.png",
      unit: "1 เม็ด",
      isTake: false
    ),
    Dose(
      name: "ยา2",
      at: DateTime(2025, 9, 4, 9, 0),
      instruction: "ก่อนอาหาร",
      picture: "assets/images/pill.png",
      unit: "1 เม็ด",
      isTake: false
    ),
    Dose(
      name: "ยา3",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      picture: "assets/images/pill.png",
      unit: "1 เม็ด",
      isTake: false
    ),
    Dose(
      name: "แคปซูล1",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      picture: "assets/images/capsule.png",
      unit: "1 เม็ด",
      isTake: true
    ),
    Dose(
      name: "แคปซูล2",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      picture: "assets/images/capsule.png",
      unit: "1 เม็ด",
      isTake: true
    ),
    Dose(
      name: "ยาน้ำ1",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      picture: "assets/images/syrup.png",
      unit: "1 ช้อน",
      isTake: true
    ),
    Dose(
      name: "ยาน้ำ2",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      picture: "assets/images/syrup.png",
      unit: "1 ช้อน",
      isTake: true
    ),
    Dose(
      name: "ยาทา1",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      picture: "assets/images/ointment.png",
      unit: "ทาบาง ๆ",
      isTake: true
    ),
    Dose(
      name: "ยาฉีด1",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      picture: "assets/images/vaccine.png",
      unit: "1 เข็ม",
      isTake: true
    ),
    Dose(
      name: "ยาฉีด2",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      picture: "assets/images/vaccine.png",
      unit: "1 เข็ม",
      isTake: true
    ),
    Dose(
      name: "ยาหยด1",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      picture: "assets/images/eye-drop 1.png",
      unit: "2 หยด",
      isTake: true
    ),
    Dose(
      name: "ยาหยด2",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      picture: "assets/images/eye-drop 1.png",
      unit: "2 หยด",
      isTake: true
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

  List<Dose> doseSelect(DateTime selectTime) {
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
  List<DoseGroup> get groupSelect {
    final Map<String, List<Dose>> group = {};
    for (var d in doseSelect(_selected)) {
      final key =
          "${d.at.hour}"
          ":"
          "${d.at.minute}"
          "|"
          "${d.instruction}";
      if (!group.containsKey(key)) {
        group[key] = [];
      }
      group[key]!.add(d);
    }
    return group.entries
        .map((e) => DoseGroup(key: e.key, doses: e.value))
        .toList();
  }

  Future<void> pickDate(BuildContext context) async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _selected,
      firstDate: DateTime(2000),
      lastDate: DateTime(2100),
      helpText: 'เลือกวันที่',
      cancelText: 'ยกเลิก',
      confirmText: 'ตกลง',
    );
    if (picked != null) {
      _selected = DateTime(picked.year, picked.month, picked.day);
      notifyListeners();
    }
  }

  void handleIsTaken(String field, Dose dose) {
    switch (field) {
      case "taken":
        // dose.isTaken = true;
        notifyListeners();
        break;
      case "not_taken":
        // dose.isTaken = false;
        notifyListeners();
        break;
      case "remove":
        all.remove(dose);
        notifyListeners();
        break;
    }
  }
}
