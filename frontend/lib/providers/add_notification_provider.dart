import 'package:flutter/material.dart';
import 'package:frontend/models/dose.dart';

class AddNotificationProvider extends ChangeNotifier {
  final String pageFrom;
  final Dose? dose;
  final String? keyName;
  final List<String>? value;

  AddNotificationProvider({
    required this.pageFrom,
    this.dose,
    this.keyName,
    this.value,
  });
}
