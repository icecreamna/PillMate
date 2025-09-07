import 'package:flutter/material.dart';
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

   InfoUser _user = InfoUser(
    firstName: "Kittabeth",
    lastName: "Chompoonich",
    idCard: "1739841323333",
    tel: "0864031301",
  );

  final InfoAppoinment _appoinment = InfoAppoinment(
    dateTime: DateTime.now(),
    note: "งดอาหาร 8 ชั่วโมง ก่อนเจาะเลือด",
    hourMinute: DateTime.now(),
  );

  InfoUser get user => _user;
  InfoAppoinment get appointment => _appoinment;

  String get appointmentDay =>
      DateFormat("MMM d, yyyy").format(_appoinment.dateTime);

  String get appointmentHourMinute =>
      DateFormat("HH:mm").format(_appoinment.hourMinute);

  void updateUser(InfoUser newUser){
    _user = newUser ;
    notifyListeners();
  }
}
