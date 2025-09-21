import 'package:flutter/material.dart';
import 'package:frontend/enums/drug_form.dart';
import 'package:frontend/enums/drug_time.dart';

class AddEditProvider extends ChangeNotifier {
  final String pageFrom;
  DrugTime selectTime;
  DrugForm selectedForm;
  String? selectedUnit;

  AddEditProvider({
    required this.pageFrom,
    this.selectedForm = DrugForm.tablet,
    this.selectTime = DrugTime.beforeMeal,
  }) {
    selectedUnit = selectedForm.unit.first;
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
