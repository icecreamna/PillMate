import 'package:flutter/material.dart';
import 'package:frontend/enums/drug_tab.dart';
import 'package:uuid/uuid.dart';


class DoseTest {
  final String id;
  final String name;
  final String amountPerDose;
  final String frequency;
  final String instruction;
  final String picture;
  final String description;
  final String unit;
  final bool import;

  DoseTest({
    required this.id,
    required this.name,
    required this.description,
    required this.import,
    required this.amountPerDose,
    required this.frequency,
    required this.instruction,
    required this.unit,
    required this.picture,
  });
  DoseTest copyWith({
    String? id,
    String? name,
    String? description,
    bool? import,
    String? amountPerDose,
    String? frequency,
    String? instruction,
    String? picture,
    String? unit,
  }) {
    return DoseTest(
      id: id ?? this.id,
      name: name ?? this.name,
      description: description ?? this.description,
      import: import ?? this.import,
      amountPerDose: amountPerDose ?? this.amountPerDose,
      frequency: frequency ?? this.frequency,
      instruction: instruction ?? this.instruction,
      picture: picture ?? this.picture,
      unit: unit ?? this.unit,
    );
  }
}

class DrugProvider extends ChangeNotifier {
  DrugTab _page = DrugTab.all;

  DrugTab get page => _page;

  final List<DoseTest> _all = [];

  List<DoseTest> get doseAll => _all;

  final _uuid = const Uuid();

  void setPage(DrugTab selectPage) {
    _page = selectPage;
    notifyListeners();
  }

  void addDose(DoseTest newDose) {
    final dose = newDose.copyWith(id:_uuid.v4());
    _all.add(dose);
    print("มี ${_all.length}");
    notifyListeners();
  }

  void removeDose(DoseTest dose){
    _all.removeWhere((d) => d.id == dose.id,);
    print("มี ${_all.length}");
    notifyListeners();
  }

  void updatedDose(DoseTest updateDose){  
    final index = _all.indexWhere((d) => d.id == updateDose.id);
    if(index != -1){
      _all[index] = updateDose;
      notifyListeners();
    }
  }
}
