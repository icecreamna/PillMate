import 'package:flutter/material.dart';
import 'package:frontend/models/dose.dart';
import 'package:intl/intl.dart';

class AddNotificationProvider extends ChangeNotifier {
  final String pageFrom;
  final Dose? dose;
  final String? keyName;
  final List<String>? value;
  DateTime? _startDate;
  DateTime? _endDate;
  String? _selectedType;

  String? get selectedType => _selectedType;
  String get startDate => _startDate != null
      ? DateFormat('d MMM yyyy', 'th').format(_startDate!)
      : '';
  String get endDate =>
      _endDate != null ? DateFormat('d MMM yyyy', 'th').format(_endDate!) : '';

  AddNotificationProvider({
    required this.pageFrom,
    this.dose,
    this.keyName,
    this.value,
  });

  void setSelectType(String type) {
    _selectedType = type;
    notifyListeners();
  }

  Future<void> pickStartDate(BuildContext context) async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _startDate ?? DateTime.now(),
      firstDate: DateTime(2000),
      lastDate: DateTime(2100),
      helpText: 'เลือกวันเริ่มต้น',
      cancelText: 'ยกเลิก',
      confirmText: 'ตกลง',
      locale: const Locale('th', 'TH'),
    );
    if (picked != null) {
      _startDate = DateTime(picked.year, picked.month, picked.day);
      notifyListeners();
    }
  }

  Future<void> pickEndDate(BuildContext context) async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _endDate ?? (_startDate ?? DateTime.now()),
      firstDate: DateTime(2000),
      lastDate: DateTime(2100),
      helpText: 'เลือกวันสิ้นสุด',
      cancelText: 'ยกเลิก',
      confirmText: 'ตกลง',
      locale: const Locale('th', 'TH'),
    );
    if (picked != null) {
      _endDate = DateTime(picked.year, picked.month, picked.day);
      notifyListeners();
    }
  }
}
