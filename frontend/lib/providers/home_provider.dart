import 'package:flutter/material.dart';

class HomeProvider extends ChangeNotifier {
  int _selectedIndex = 0;

  int get selectIndex => _selectedIndex;

  void setSelectIndex(int index) {
      _selectedIndex = index;
      notifyListeners();
  }
}
