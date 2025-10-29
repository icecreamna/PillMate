import 'package:flutter/material.dart';
import 'package:frontend/models/notification_info.dart';
import 'package:frontend/services/noti_service.dart';
import '../models/dose.dart';

class AddSingleNotificationProvider extends ChangeNotifier {
  final NotiService _notiService = NotiService();

  late Dose _tempDose;

  NotificationInfo? _savedNotification;
  NotificationInfo? get savedNotification => _savedNotification;
  bool _isLoading = false;


  Dose get tempDose => _tempDose;
  bool get isLoading => _isLoading;

  AddSingleNotificationProvider(Dose dose) {
    _tempDose = dose;
    loadNotification();
  }

  Future<void> loadNotification() async {
    _isLoading = true;
    try {
      final info = await _notiService.getNotiInfo(
        type: "medicine",
        id: _tempDose.id,
      );
      _savedNotification = info;
      notifyListeners();
      if (info != null) {
        print("✅ โหลด noti info ของยาเดี่ยวสำเร็จ: ${info?.notiFormatName}");
      } else {
        print("❌ ไม่มี noti info ของยาเดี่ยว: ${info?.notiFormatName}");
      }
    } catch (e) {
      print("❌ โหลด noti info ล้มเหลว: $e");
      _savedNotification = null;
    }finally {
      _isLoading = false;
      notifyListeners();
    }
    notifyListeners();
  }

  Future<bool> removeNoti() async {
    _isLoading = true;
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
      _isLoading = false;
      notifyListeners();
    }
  }

  void updatedTempDose(Dose newDose) {
    _tempDose = newDose;
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
