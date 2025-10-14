import 'package:flutter/material.dart';
import 'package:frontend/models/dose.dart';
import 'package:frontend/models/notification_info.dart';
import 'package:frontend/services/noti_service.dart';
import 'package:intl/intl.dart';

class AddNotificationProvider extends ChangeNotifier {
  final _notiService = NotiService();

  final String pageFrom;
  final Dose? dose;
  final String? keyName;
  final int? groupId;
  final List<String>? value;
  DateTime? _startDate;
  DateTime? _endDate;
  String _selectedType = '';
  bool _isLoading = false;

  List<NotiFormatModel> _formats = [];

  String get selectedType => _selectedType;
  List<NotiFormatModel> get formats => _formats;
  bool get isLoading => _isLoading;

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
    this.groupId,
    this.keyName,
    this.value,
  });

  void setSelectType(String typeId) {
    _selectedType = typeId;
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

  Future<void> loadNotiFormats() async {
    _isLoading = true;
    notifyListeners();
    try {
      _formats = await _notiService.getFormats();
      debugPrint(
        "✅ โหลด noti formats ${_formats.length} ${formats.map((e) => '{id: ${e.id}, name: ${e.formatName}').toList()}รายการสำเร็จ",
      );
    } catch (e) {
      debugPrint("Can't load notifomat $e");
      _formats = [];
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  String getTypeName(String id) {
    switch (id) {
      case "1":
        return "Fixed";
      case "2":
        return "Interval";
      case "3":
        return "DailyWeekly";
      case "4":
        return "Cycle";
      default:
        return "";
    }
  }

  Future<bool> addNotification(NotificationInfo info) async {
    _isLoading = true;
    notifyListeners();
    try {
      final int formatId = int.parse(selectedType);

      if (pageFrom == "group") {

        if (groupId == null) throw Exception("Group ID ไม่ถูกต้อง");

        await _notiService.addNotification(
          groupId: groupId,
          notiFormatId: formatId,
          info: info,
        );
        debugPrint("✅ เพิ่มแจ้งเตือนกลุ่มสำเร็จ (group_id=$groupId)");
      } else {
        final medId = int.tryParse(dose?.id ?? "");
        if (medId == null) throw Exception("Medicine ID ไม่ถูกต้อง");

        await _notiService.addNotification(
          myMedicineId: medId,
          notiFormatId: formatId,
          info: info,
        );
        debugPrint("✅ เพิ่มแจ้งเตือนยาเดี่ยวสำเร็จ (my_medicine_id=$medId)");
      }
      return true;
    } catch (e) {
      debugPrint("❌ เพิ่มแจ้งเตือนไม่สำเร็จ: $e");
      return false;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
