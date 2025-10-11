import 'package:flutter/material.dart';
import 'package:frontend/services/group_service.dart';

class AddGroupProvider extends ChangeNotifier {
  final GroupService _service = GroupService();

  String listError = "";

  bool _isLoading = false;
  List<Map<String, dynamic>> _groups = [];
  List<String> _selectedList = [];

  bool get isLoading => _isLoading;
  List<Map<String, dynamic>> get groups => _groups;
  List<String> get selectedList => _selectedList;

  Future<bool> addGroups(String groupName, List<String> medicineIds) async {
    _isLoading = true;
    notifyListeners();
    try {
      final success = await _service.createGroup(
        groupName: groupName,
        medicineIds: medicineIds,
      );
      if (success) {
        return true;
      } else {
        debugPrint("Provider can't created Group");
        return false;
      }
    } catch (e) {
      debugPrint("Provider Catch created Group $e");
      return false;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  void setSelectedList(List<String> list) {
    _selectedList = list;
    notifyListeners();
  }

  void removeSelectedList(String id){
    _selectedList.removeWhere((s) => s == id,);
    notifyListeners();
  }

  void setListError() {
    listError = "กรุณาใส่ยามากกว่า 1 ตัว";
    notifyListeners();
  }

  void clearListError() {
    listError = "";
    notifyListeners();
  }
}
