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
  String _selectedType = '';

  String get selectedType => _selectedType;

  String get startDate {
    if (_startDate == null) return '';
    final thaiYear = _startDate!.year + 543;
    final monthName = DateFormat.MMM('th').format(_startDate!);
    return '${_startDate!.day} $monthName $thaiYear';
  }

  String get endDate {
    if (_endDate == null) return '';
    final thaiYear = _endDate!.year + 543;
    final monthName = DateFormat.MMM('th').format(_endDate!);
    return '${_endDate!.day} $monthName $thaiYear';
  }

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
      firstDate: DateTime.now(),
      lastDate: _endDate ?? DateTime(2100),
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
      firstDate: _startDate ?? DateTime.now(),
      lastDate: DateTime(2100),
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
