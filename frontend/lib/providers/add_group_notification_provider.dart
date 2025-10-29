import 'package:flutter/widgets.dart';
import 'package:frontend/models/dose.dart';
import 'package:frontend/models/notification_info.dart';
import 'package:frontend/services/group_service.dart';
import 'package:frontend/services/noti_service.dart';

class AddGroupNotificationProvider extends ChangeNotifier {
  final GroupService _service = GroupService();
  final NotiService _notiService = NotiService();

  bool _isloading = false;
  String listError = "";
  String keyName;
  final int groupId;
  List<String> _value = [];
  List<Dose> memberList = [];

  List<String> get value => _value;
  bool get isLoading => _isloading;

  NotificationInfo? _savedNotification;
  NotificationInfo? get savedNotification => _savedNotification;

  AddGroupNotificationProvider({required this.keyName, required this.groupId}) {
    loadGroupDetail();
    loadNotification();
  }

  Future<void> loadGroupDetail() async {
    _isloading = true;
    notifyListeners();
    try {
      final detail = await _service.getGroupWithDetail(groupId: groupId);
      final members = detail!["members"] as List? ?? [];
      memberList = members.map((e) => Dose.fromJson(e)).toList();
      _value = memberList.map((m) => m.id).toList(); // ✅ เก็บ id ทั้งหมด
      debugPrint("✅ โหลดสมาชิกGroupDetail ${memberList.length} รายการสำเร็จ");
    } catch (e) {
      debugPrint("Cant load Detail $e");
    } finally {
      _isloading = false;
      notifyListeners();
    }
  }

  Future<bool> updateGroup({
    required int groupId,
    required String groupName,
    required List<String> medicineIds,
  }) async {
    _isloading = true;
    notifyListeners();

    try {
      final success = await _service.updateGroup(
        groupId: groupId,
        newGroupName: groupName,
        medicineIds: medicineIds,
      );
      if (success) {
        keyName = groupName;
        notifyListeners();
      }
      return success;
    } catch (e) {
      debugPrint("Error updated Group $e");
      return false;
    } finally {
      _isloading = false;
      notifyListeners();
    }
  }

  Future<bool> deleteGroup({required String groupId}) async {
    _isloading = true;
    notifyListeners();
    try {
      final success = await _service.deleteGroup(groupId: groupId);
      debugPrint("Delete success");
      return success;
    } catch (e) {
      debugPrint("Error delete Group $e");
      return false;
    } finally {
      _isloading = false;
      notifyListeners();
    }
  }

  Future<void> loadNotification() async {
    try {
      final info = await _notiService.getNotiInfo(
        type: "group",
        id: groupId.toString(),
      );
      _savedNotification = info;
      if (info != null) {
        print("✅ โหลด noti info ของกลุ่มสำเร็จ: ${info?.notiFormatName}");
      } else {
        print("❌ โหลด noti info ของกลุ่มไม่สำเร็จ: ${info?.notiFormatName}");
      }
    } catch (e) {
      print("❌ โหลด noti info กลุ่มล้มเหลว: $e");
      _savedNotification = null;
    }
    notifyListeners();
  }

  Future<bool> removeNoti() async {
    _isloading = true;
    notifyListeners();
    try {
      final success = await _notiService.removeNotification(
        id: _savedNotification!.id.toString(),
      );
      if (success) {
        _savedNotification = null;
        notifyListeners();
        return true;
      }
      return false;
    } catch (e) {
      debugPrint("Delete exception $e");
      return false;
    } finally {
      _isloading = false;
      notifyListeners();
    }
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
}
