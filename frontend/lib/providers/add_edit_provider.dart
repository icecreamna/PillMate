import 'package:flutter/material.dart';
import 'package:frontend/models/drug_form.dart';
import 'package:frontend/models/drug_time.dart';
import 'package:frontend/services/initialdata_service.dart';
import 'package:frontend/services/medicine_service.dart';

import '../models/dose.dart';

class AddEditProvider extends ChangeNotifier {
  final String pageFrom;
  DrugTimeModel? selectTime;
  DrugFormModel? selectedForm;
  DrugUnitModel? selectedUnit;
  Dose? editDose;

  bool isLoading = false;
  List<DrugFormModel> forms = [];
  List<DrugTimeModel> times = [];

  final _service = InitialDataService();
  final _medicineService = MedicineService();

  AddEditProvider({
    required this.pageFrom,
    this.selectedForm,
    this.selectTime,
    this.editDose,
  }) {
    loadInitialData();
    if (forms.isNotEmpty) {
      selectedForm ??= forms.first;
      if (selectedForm!.units.isNotEmpty) {
        selectedUnit ??= selectedForm!.units.first;
      }
    }
  }

  // void _loadDose(Dose dose) {
  //   selectedUnit = dose.unit;
  // }
  void setSelectForm(DrugFormModel drugSelect) {
    selectedForm = drugSelect;
    selectedUnit = (drugSelect.units.isNotEmpty
        ? drugSelect.units.first
        : null);
    notifyListeners();
  }

  void setSelectTime(DrugTimeModel newTime) {
    selectTime = newTime;
    notifyListeners();
  }

  void setUnit(DrugUnitModel unit) {
    selectedUnit = unit;
    notifyListeners();
  }

  Future<void> loadInitialData() async {
    isLoading = true;
    notifyListeners();

    try {
      forms = await _service.fetchDrugForms();
      times = await _service.fetchDrugTimes();

      if (pageFrom == "edit" && editDose != null) {
        selectedForm = forms.firstWhere(
          (f) => f.id == editDose!.formId,
          orElse: () => forms.first,
        );
        selectedUnit = selectedForm!.units.firstWhere(
          (u) => u.id == editDose!.unitId,
          orElse: () => selectedForm!.units.first,
        );
        selectTime = times.firstWhere(
          (t) => t.id == editDose!.instructionId,
          orElse: () => times.first,
        );
      } else {
        if (forms.isNotEmpty) {
          selectedForm ??= forms.first;
          selectedUnit = (forms.first.units.isNotEmpty
              ? forms.first.units.first
              : null);
        }
        if (times.isNotEmpty) {
          selectTime ??= times.first;
        }
      }
    } catch (e) {
      debugPrint("❌ Failed to load initial data: $e");
    }

    isLoading = false;
    notifyListeners();
  }

  Future<bool> addMedicine({
    required String name,
    required String properties,
    required String amountPerTime,
    required String timePerDay,
  }) async {
    if (selectedForm == null || selectTime == null || selectedUnit == null) {
      debugPrint("❌ Missing form/unit/instruction");
      return false;
    }

    final success = await _medicineService.addMedicineInfo(
      medName: name,
      genericName: "-",
      properties: properties,
      formId: selectedForm!.id,
      unitId: selectedUnit!.id,
      instructionId: selectTime!.id,
      amountPerTime: amountPerTime,
      timePerDay: timePerDay,
    );
    return success;
  }

  Future<bool> editMedicine({
    required int id,
    required String name,
    required String properties,
    required String amountPerTime,
    required String timePerDay,
    required int selectedFormId,
    required int selectTimeId,
    required int selectedUnitId,
  }) async {
    if (selectedFormId == null ||
        selectTimeId == null ||
        selectedUnitId == null) {
      debugPrint("❌ Missing form/unit/instruction");
      return false;
    }

    final success = await _medicineService.updatedMedicineInfo(
      id: id,
      medName: name,
      genericName: "-",
      properties: properties,
      formId: selectedFormId,
      unitId: selectedUnitId,
      instructionId: selectTimeId,
      amountPerTime: amountPerTime,
      timePerDay: timePerDay,
    );
    return success;
  }
}
