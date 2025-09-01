import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

class Dose {
  final String name;
  final DateTime at;
  final String instruction;
  bool isTaken;
  Dose({
    required this.name,
    required this.at,
    required this.instruction,
    required this.isTaken,
  });
}

class TodayProvider extends ChangeNotifier {
  DateTime _selected = DateTime.now();

  DateTime get selected => _selected;

  String get dateLabel => DateFormat("MMM d, yyyy").format(_selected);

  final List<Dose> all = <Dose>[
    Dose(
      name: "ยา1",
      at: DateTime(2025, 9, 1, 8, 0),
      instruction: "หลังอาหาร",
      isTaken: false,
    ),
    Dose(
      name: "ยา2",
      at: DateTime(2025, 9, 1, 9, 0),
      instruction: "ก่อนอาหาร",
      isTaken: false,
    ),
    Dose(
      name: "ยา3",
      at: DateTime(2025, 9, 2, 13, 0),
      instruction: "หลังอาหาร",
      isTaken: false,
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
