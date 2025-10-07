import 'package:flutter/widgets.dart';
import 'package:frontend/models/notification_info.dart';

class AddGroupNotificationProvider extends ChangeNotifier {
  String listError = "";
  late String _keyName;
  late List<String> _value;

  String get keyName => _keyName;
  List<String> get value => _value;

  NotificationInfo? _savedNotification;
  NotificationInfo? get savedNotification => _savedNotification;

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

  void saveNotification(NotificationInfo info) {
    _savedNotification = info;
    notifyListeners();
  }

  void clearNotification() {
  _savedNotification = null;
  notifyListeners();
}
}
