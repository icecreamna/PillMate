import 'package:flutter/material.dart';

class AddGroupProvider extends ChangeNotifier {
  List<String> _selectedList = [];

  List<String> get selectedList => _selectedList;

  void setSelectedList(List<String> list) {
    _selectedList = list;
    debugPrint("เลือกมา${list.length}");
    notifyListeners();
  }

  void removeSelected(String id) {
    _selectedList.removeWhere((d) => d == id,);
    notifyListeners();
  }
}
