import 'package:flutter/material.dart';
import 'package:frontend/enums/drug_tab.dart';
import 'package:uuid/uuid.dart';
import '../models/dose.dart';

class DrugProvider extends ChangeNotifier {
  DrugTab _page = DrugTab.all;

  DrugTab get page => _page;

  final List<Dose> _all = [
    Dose(
      id: "1",
      name: "paracetamol",
      description: "fattu",
      import: true,
      amountPerDose: "1",
      frequency: "3",
      instruction: "หลังอาหาร",
      unit: "เม็ด",
      picture: "assets/images/pill.png",
    ),
  ];

  List<Dose> get doseAll => _all;

  final _uuid = const Uuid();

  void setPage(DrugTab selectPage) {
    _page = selectPage;
    notifyListeners();
  }

  void addDose(Dose newDose) {
    final dose = newDose.copyWith(id: _uuid.v4());
    _all.add(dose);
    print("มี ${_all.length}");
    notifyListeners();
  }

  void removeDose(Dose dose) {
    _all.removeWhere((d) => d.id == dose.id);
    print("มี ${_all.length}");
    notifyListeners();
  }

  void updatedDose(Dose updateDose) {
    final index = _all.indexWhere((d) => d.id == updateDose.id);
    if (index != -1) {
      _all[index] = updateDose;
      notifyListeners();
    }
  }

  
}
