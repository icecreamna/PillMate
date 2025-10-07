import 'package:flutter/material.dart';
import 'package:frontend/models/notification_info.dart';
import '../models/dose.dart';

class AddSingleNotificationProvider extends ChangeNotifier {
  late Dose _tempDose;

  NotificationInfo? _savedNotification;
  NotificationInfo? get savedNotification => _savedNotification;

  Dose get tempDose => _tempDose;

  AddSingleNotificationProvider(Dose dose) {
    _tempDose = dose;
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
