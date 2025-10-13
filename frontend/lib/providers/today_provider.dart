import 'package:flutter/material.dart';
import 'package:frontend/services/today_service.dart';
import 'package:intl/intl.dart';

class DoseSingle {
  final String name;
  // final String picture;
  final String amountPerTime;
  final String unit;
  DoseSingle({
    required this.name,
    required this.amountPerTime,
    required this.unit,
  });
}

class DoseGroup {
  final int? notiSingleId;
  final int? medicineId;
  final int? groupId;
  final List<int>? notiGroupIds;
  int? symptomId;
  final String nameGroup;
  final String key;
  final List<DoseSingle> doses;
  final DateTime at;
  final String instruction;
  bool saveNote;
  bool isTaken;
  String? symptomNote;

  DoseGroup({
    this.notiSingleId,
    this.medicineId,
    this.groupId,
    this.notiGroupIds,
    required this.nameGroup,
    required this.key,
    required this.doses,
    required this.at,
    required this.instruction,
    this.saveNote = false,
    this.isTaken = false,
    this.symptomNote = "",
    this.symptomId,
  });
}

class TodayProvider extends ChangeNotifier {
  final TodayService _todayService = TodayService();

  bool _isLoading = false;

  DateTime _selected = DateTime.now();

  DateTime get selected => _selected;
  bool get isLoading => _isLoading;

  String get dateLabel {
    final thDate = DateFormat("d MMMM yyyy", "th_TH").format(_selected);
    final buddhistYear = _selected.year + 543;
    return thDate.replaceAll('${_selected.year}', '$buddhistYear');
  }

  // String getTimeHourMinute(String key) => key.split("|").first;
  // String getInstruction(String key) => key.split("|").last;

  List<DoseGroup> all = <DoseGroup>[];
  // List<Dose> get doseSelect {
  //   return all
  //       .where(
  //         (d) =>
  //             d.at.year == _selected.year &&
  //             d.at.month == _selected.month &&
  //             d.at.day == _selected.day,
  //       )
  //       .toList();
  // }//function
  Future<void> loadTodayData() async {
    _isLoading = true;
    notifyListeners();
    try {
      final data = await _todayService.fetchTodayNoti(_selected);
      all = data;

      await loadAllSymptomsAndMap();
      notifyListeners();
    } catch (e) {
      print("❌ โหลดข้อมูลวันนี้ไม่สำเร็จ: $e");
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> updateTakenStatus(
    DoseGroup doseGroup,
    bool taken,
    BuildContext context,
  ) async {
    _isLoading = true;
    notifyListeners();
    try {
      if (doseGroup.notiSingleId != null) {
        // ✅ แบบเดี่ยว
        final success = await _todayService.markTaken(
          notiItemId: doseGroup.notiSingleId!,
          taken: taken,
        );
        if (success) {
          doseGroup.isTaken = taken;
          notifyListeners();
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text(
                "อัปเดตสถานะการกิน เสร็จสิ้น",
                style: TextStyle(color: Colors.white),
              ),
              backgroundColor: Colors.green,
              behavior: SnackBarBehavior.floating,
              duration: Duration(seconds: 2),
            ),
          );
        } else {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text(
                "อัปเดตสถานะการกิน ไม่เสร็จสิ้น",
                style: TextStyle(color: Colors.white),
              ),
              backgroundColor: Colors.red,
              behavior: SnackBarBehavior.floating,
              duration: Duration(seconds: 2),
            ),
          );
        }
      } else if (doseGroup.notiGroupIds != null) {
        bool allSuccess = true;
        for (final id in doseGroup.notiGroupIds!) {
          print("ส่ง notiGroupId ${id}");
          final success = await _todayService.markTaken(
            notiItemId: id,
            taken: taken,
          );
          if (!success) allSuccess = false;
        }
        if (allSuccess) {
          doseGroup.isTaken = taken;
          notifyListeners();
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text(
                "อัปเดตสถานะการกิน เสร็จสิ้น",
                style: TextStyle(color: Colors.white),
              ),
              backgroundColor: Colors.green,
              behavior: SnackBarBehavior.floating,
              duration: Duration(seconds: 2),
            ),
          );
        } else {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text(
                "อัปเดตสถานะการกิน ไม่เสร็จสิ้น",
                style: TextStyle(color: Colors.white),
              ),
              backgroundColor: Colors.red,
              behavior: SnackBarBehavior.floating,
              duration: Duration(seconds: 2),
            ),
          );
        }
      }
    } catch (e) {
      print("Error $e");
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> createSymptom({
    DoseGroup? dose,
    required String symptomNote,
    required BuildContext context,
  }) async {
    _isLoading = true;
    notifyListeners();
    try {
      final notiId =
          dose?.notiSingleId ??
          (dose?.notiGroupIds != null && dose!.notiGroupIds!.isNotEmpty
              ? dose.notiGroupIds!.first
              : null);

      if (notiId == null) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text("❌ ไม่พบรหัสแจ้งเตือน (notiItemId)")),
        );
        return;
      }

      final created = await _todayService.createSymptom(
        symptomNote: symptomNote,
        notiItemId: notiId,
        groupId: dose?.groupId,
        myMedicineId: dose?.medicineId,
      );
      if (created != null) {
        dose?.saveNote = true;
        dose?.symptomNote = symptomNote;
        dose!.symptomId = created["id"];
        print("idคือออออ${dose.symptomId}");
        notifyListeners();
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text(
              "บันทึกอาการสำเร็จ",
              style: TextStyle(color: Colors.white),
            ),
            backgroundColor: Colors.green,
            behavior: SnackBarBehavior.floating,
            duration: Duration(seconds: 2),
          ),
        );
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text(
              "อัปเดตสถานะการกิน ไม่เสร็จสิ้น",
              style: TextStyle(color: Colors.white),
            ),
            backgroundColor: Colors.red,
            behavior: SnackBarBehavior.floating,
            duration: Duration(seconds: 2),
          ),
        );
      }
    } catch (e) {
      print("Error exception $e");
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> editSymptom({
    DoseGroup? dose,
    required String symptomNote,
    required BuildContext context,
  }) async {
    _isLoading = true;
    notifyListeners();
    try {
      if (dose!.symptomId == null) {
        throw Exception("❌ ไม่มี Symptom ID สำหรับอัปเดต");
      }

      final ok = await _todayService.updateSymptom(
        symptomId: dose.symptomId!,
        symptomNote: symptomNote,
      );

      if (ok) {
        dose.symptomNote = symptomNote;
        notifyListeners();

        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text(
              "อัปเดตอาการสำเร็จ",
              style: TextStyle(color: Colors.white),
            ),
            backgroundColor: Colors.green,
          ),
        );
      } else {
        throw Exception("อัปเดตไม่สำเร็จ");
      }
    } catch (e) {
      print("❌ แก้ไขอาการล้มเหลว: $e");
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
            "เกิดข้อผิดพลาด: $e",
            style: const TextStyle(color: Colors.white),
          ),
          backgroundColor: Colors.red,
        ),
      );
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> loadAllSymptomsAndMap() async {
    try {
      final list = await _todayService.fetchAllSymptoms();

      // เลือก "อาการล่าสุด" ต่อ noti_item_id
      // (ถ้า API ส่ง created_at มาเป็น RFC3339 ก็ parse ได้ตรง ๆ)
      final Map<int, Map<String, dynamic>> latestByNoti = {};
      for (final s in list) {
        final int notiId = s["noti_item_id"];
        final String createdAt = s["created_at"] ?? "";
        final prev = latestByNoti[notiId];
        if (prev == null) {
          latestByNoti[notiId] = s;
        } else {
          final prevTime =
              DateTime.tryParse(prev["created_at"] ?? "") ??
              DateTime.fromMillisecondsSinceEpoch(0);
          final curTime =
              DateTime.tryParse(createdAt) ??
              DateTime.fromMillisecondsSinceEpoch(0);
          if (curTime.isAfter(prevTime)) latestByNoti[notiId] = s;
        }
      }

      // แมพเข้า DoseGroup (single ใช้ notiSingleId, group ใช้ notiGroupIds.first)
      for (final d in all) {
        final notiId =
            d.notiSingleId ??
            (d.notiGroupIds != null && d.notiGroupIds!.isNotEmpty
                ? d.notiGroupIds!.first
                : null);
        if (notiId == null) continue;

        final match = latestByNoti[notiId];
        if (match != null) {
          d.saveNote = true;
          d.symptomNote = match["symptom_note"] ?? "";
        }
      }
    } catch (e) {
      print("❌ โหลด/แมพ symptoms ไม่สำเร็จ: $e");
    }
  }

  List<DoseGroup> doseSelect(DateTime? selectTime) {
    if (selectTime == null) return [];
    return all
        .where(
          (d) =>
              d.at.year == selectTime.year &&
              d.at.month == selectTime.month &&
              d.at.day == selectTime.day,
        )
        .toList();
  }

  //function kub

  Future<void> pickDate(BuildContext context) async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _selected,
      firstDate: DateTime(2000),
      lastDate: DateTime(2100),
      helpText: 'เลือกวันที่',
      cancelText: 'ยกเลิก',
      confirmText: 'ตกลง',
      locale: const Locale('th', 'TH'),
    );
    if (picked != null) {
      _selected = DateTime(picked.year, picked.month, picked.day);
      notifyListeners();
    }
  }

  void removeDose(DoseGroup doseGroup) {
    all.remove(doseGroup);
    notifyListeners();
  }

  void setNote(bool saveNote, DoseGroup doseGroup) {
    doseGroup.saveNote = saveNote;
    notifyListeners();
  }
}
