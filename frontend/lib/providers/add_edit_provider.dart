import 'package:flutter/material.dart';
import 'package:frontend/enums/drug_form.dart';
import 'package:frontend/enums/drug_time.dart';

import '../models/dose.dart';

class AddEditProvider extends ChangeNotifier {
  final String pageFrom;
  DrugTime selectTime;
  DrugForm selectedForm;
  String? selectedUnit;

  Dose? editDose;

  AddEditProvider({
    required this.pageFrom,
    this.selectedForm = DrugForm.tablet,
    this.selectTime = DrugTime.beforeMeal,
    this.editDose
  }) {
    selectedUnit = selectedForm.unit.first;
    if(editDose != null){
      _loadDose(editDose!);
  }
  }

  void _loadDose(Dose dose){
    selectedUnit = dose.unit;
    selectTime = DrugTime.values.firstWhere((dt) => dt.label == dose.instruction,orElse: () => DrugTime.beforeMeal,);
    selectedForm = DrugForm.values.firstWhere((df) => df.image == dose.picture ,orElse: () => DrugForm.tablet,);
  }


  void setSelectForm(DrugForm drugSelect) {
    selectedForm = drugSelect;
    selectedUnit = drugSelect.unit.first;
    notifyListeners();
  }

  void setSelectTime(DrugTime newTime) {
    selectTime = newTime;
    notifyListeners();
  }

  void setUnit(String unit) {
    selectedUnit = unit;
    notifyListeners();
  }
}
