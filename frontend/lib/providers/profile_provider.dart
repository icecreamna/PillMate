import 'package:flutter/material.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:intl/intl.dart';

class InfoUser {
  final String firstName;
  final String lastName;
  final String idCard;
  final String tel;

  InfoUser({
    required this.firstName,
    required this.lastName,
    required this.idCard,
    required this.tel,
  });
}

class InfoAppoinment {
  final DateTime dateTime;
  final String note;
  final DateTime hourMinute;

  InfoAppoinment({
    required this.dateTime,
    required this.note,
    required this.hourMinute,
  });
}

class ProfileProvider extends ChangeNotifier {
  final AuthService authService;
  ProfileProvider(this.authService);

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  InfoUser _user = InfoUser(
    firstName: "Kittabeth",
    lastName: "Chompoonich",
    idCard: "1739841323333",
    tel: "0864031301",
  );

  final InfoAppoinment _appoinment = InfoAppoinment(
    dateTime: DateTime.now(),
    note: "งดอาหาร 8 ชั่วโมง ก่อนเจาะเลือด",
    hourMinute: DateTime.now(),
  );

  InfoUser get user => _user;
  InfoAppoinment get appointment => _appoinment;

  String get appointmentDay {
    final thDate = DateFormat(
      "d MMMM yyyy",
      "th_TH",
    ).format(_appoinment.dateTime);
    final buddhistYear = _appoinment.dateTime.year + 543;
    return thDate.replaceAll('${_appoinment.dateTime.year}', '$buddhistYear');
  }

  String get appointmentHourMinute =>
      DateFormat("HH:mm").format(_appoinment.hourMinute);

  void updateUser(InfoUser newUser) {
    _user = newUser;
    notifyListeners();
  }

  Future<bool> logout() async {
    _isLoading = true;
    notifyListeners();

    try {
      final success = await authService.logout();
      return success;
    } catch (e) {
      debugPrint("❌ LogoutProvider error: $e");
      return false;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
