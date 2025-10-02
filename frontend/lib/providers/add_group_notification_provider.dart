import 'package:flutter/widgets.dart';

class AddGroupNotificationProvider extends ChangeNotifier {
  String listError = "";
  late String _keyName;
  late List<String> _value;

  String get keyName => _keyName;
  List<String> get value => _value;

  AddGroupNotificationProvider(String key, List<String> value) {
    _keyName = key;
    _value = List.from(value);
  }
  void setSelectedList(List<String> list) {
    _value = list;
    debugPrint("เลือกมา${list.length}");
    notifyListeners();
  }

  void removeSelected(String id) {
    _value.removeWhere((d) => d == id);
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
