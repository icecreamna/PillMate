import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

class Dose {
  final String name;
  final DateTime at;
  final String instruction;
  final String picture;
  bool isTaken;
  Dose({
    required this.name,
    required this.at,
    required this.instruction,
    required this.isTaken,
    required this.picture,
  });
}

class TodayProvider extends ChangeNotifier {
  DateTime _selected = DateTime.now();

  DateTime get selected => _selected;

  String get dateLabel => DateFormat("MMM d, yyyy").format(_selected);

  final List<Dose> all = <Dose>[
    Dose(
      name: "ยา1",
      at: DateTime(2025, 9, 4, 8, 0),
      instruction: "หลังอาหาร",
      isTaken: false,
      picture: "assets/images/pill.png",
    ),
    Dose(
      name: "ยา2",
      at: DateTime(2025, 9, 4, 9, 0),
      instruction: "ก่อนอาหาร",
      isTaken: false,
      picture: "assets/images/pill.png",
    ),
    Dose(
      name: "ยา3",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      isTaken: false,
      picture: "assets/images/pill.png",
    ),
    Dose(
      name: "แคปซูล1",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      isTaken: false,
      picture: "assets/images/capsule.png",
    ),
    Dose(
      name: "แคปซูล2",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      isTaken: false,
      picture: "assets/images/capsule.png",
    ),
    Dose(
      name: "ยาน้ำ1",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      isTaken: false,
      picture: "assets/images/syrup.png",
    ),
    Dose(
      name: "ยาน้ำ2",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      isTaken: false,
      picture: "assets/images/syrup.png",
    ),
    Dose(
      name: "ยาทา1",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      isTaken: false,
      picture: "assets/images/ointment.png",
    ),
    Dose(
      name: "ยาฉีด1",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      isTaken: false,
      picture: "assets/images/vaccine.png",
    ),
    Dose(
      name: "ยาฉีด2",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      isTaken: false,
      picture: "assets/images/vaccine.png",
    ),
    Dose(
      name: "ยาหยด",
      at: DateTime(2025, 9, 4, 13, 0),
      instruction: "หลังอาหาร",
      isTaken: false,
      picture: "assets/images/eye-drop 1.png",
    ),
  ];
  List<Dose> get doseSelect {
    return all
        .where(
          (d) =>
              d.at.year == _selected.year &&
              d.at.month == _selected.month &&
              d.at.day == _selected.day,
        )
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
        dose.isTaken = true;
        notifyListeners();
        break;
      case "not_taken":
        dose.isTaken = false;
        notifyListeners();
        break;
      case "remove":
        all.remove(dose);
        notifyListeners();
        break;
    }
  }
}
