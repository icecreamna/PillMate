import 'package:flutter/material.dart';
import 'package:frontend/enums/drug_tab.dart';
import 'package:frontend/services/group_service.dart';
import 'package:frontend/services/medicine_service.dart';
import '../models/dose.dart';

class DrugProvider extends ChangeNotifier {
  final _medicineService = MedicineService();
  final _groupService = GroupService();

  DrugTab _page = DrugTab.all;

  List<Dose> _myMedicines = [];

  final Map<String, List<String>> _groups = {};
  final Map<String, int> _groupsId = {};

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  DrugTab get page => _page;
  Map<String, List<String>> get groups => _groups;
  Map<String, int> get groupsId => _groupsId;
  List<Dose> get doseAll => _myMedicines;

  Future<void> loadMyMedicines() async {
    _isLoading = true;
    notifyListeners();

    try {
      final medicines = await _medicineService.getMyMedicines();
      _myMedicines = medicines!;
      debugPrint("✅ โหลดยา ${medicines.length} รายการสำเร็จ");
    } catch (e) {
      debugPrint("❌ โหลดยาไม่สำเร็จ: $e");
      _myMedicines = [];
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<bool> removeMedicine({required int id}) async {
    _isLoading = true;
    notifyListeners();
    try {
      final success = await _medicineService.deleteMedicineInfo(id: id);
      return success;
    } catch (e) {
      debugPrint("❌ ลบไม่สำเร็จ: $e");
      return false;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> loadGroups() async {
    _isLoading = true;
    notifyListeners();
    try {
      final data = await _groupService.getGroups();
      _groups.clear();
      _groupsId.clear();
      debugPrint("✅ โหลดกลุ่มยา ${data.length} รายการสำเร็จ");
      for (final g in data) {
        final groupData = g["group"] ?? {};
        final name = groupData["group_name"] ?? "-";
        final id = groupData["id"];
        final count = g["member_count"] ?? 0;
        _groups[name] = List.filled(count, "");
        _groupsId[name] = id;
      }
    } catch (e) {
      debugPrint("Provider can't load Groups");
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  void setPage(DrugTab selectPage) {
    _page = selectPage;
    notifyListeners();
  }

  void updatedDose(Dose updateDose) {
    final index = _myMedicines.indexWhere((d) => d.id == updateDose.id);
    if (index != -1) {
      _myMedicines[index] = updateDose;
      notifyListeners();
    }
  }

  void addGroup(String groupName, List<String> listDrugIds) {
    _groups[groupName] = List.from(listDrugIds);
    notifyListeners();
  }

  void updatedDoseGroup(String groupName, List<String> listDrugIds) {
    _groups[groupName] = List.from(listDrugIds);
    notifyListeners();
  }

  void removeGroup(String groupName) {
    _groups.removeWhere((key, _) => key == groupName);
    notifyListeners();
  }
}
