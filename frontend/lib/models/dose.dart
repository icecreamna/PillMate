import 'package:intl/intl.dart';

class Dose {
  final String id;
  final String name;
  final String description;
  final String amountPerDose;
  final String frequency;
  final String instruction;
  final String unit;
  final String picture;
  final bool import; // hospital = true, manual = false
  final String? startDate;
  final String? endDate;
  final String? note;

  final int? formId;
  final int? unitId;
  final int? instructionId;
  final int? groupId;

  Dose({
    required this.id,
    required this.name,
    required this.description,
    required this.amountPerDose,
    required this.frequency,
    required this.instruction,
    required this.unit,
    required this.picture,
    required this.import,
    this.formId,
    this.unitId,
    this.instructionId,
    this.groupId,
    this.startDate,
    this.endDate,
    this.note,
  });

  /// ✅ แปลงจาก JSON ที่มาจาก backend
  factory Dose.fromJson(Map<String, dynamic> json) {
    final formName = json["form_name"] ?? "-";
    final startDate = json["start_date"] ?? "-";
    final endDate = json["end_date"] ?? "-";
    return Dose(
      id: json["id"].toString(),
      name: json["med_name"] ?? "-",
      description: json["properties"] ?? "-",
      amountPerDose: json["amount_per_time"] ?? "-",
      frequency: json["times_per_day"] ?? "-",
      instruction: json["instruction_name"] ?? "-",
      unit: json["unit_name"] ?? "-",
      startDate: _toThaiDate(startDate),
      endDate: _toThaiDate(endDate),
      note: json["note"] ?? "",
      picture: _mapImage(formName),
      import: (json["source"] == "hospital"),
      formId: json["form_id"],
      unitId: json["unit_id"],
      instructionId: json["instruction_id"],
      groupId: json["group_id"],
    );
  }

  /// ✅ Map ชื่อ form เป็นภาพไอคอน
  static String _mapImage(String formName) {
    switch (formName) {
      case "ยาเม็ด":
        return "assets/images/pill.png";
      case "แคปซูล":
        return "assets/images/capsule.png";
      case "ยาน้ำ":
        return "assets/images/syrup.png";
      case "ยาใช้ทา":
        return "assets/images/ointment.png";
      case "ยาฉีด":
        return "assets/images/vaccine.png";
      case "ยาใช้หยด":
        return "assets/images/eye-drop 1.png";
      default:
        return "assets/images/pill.png";
    }
  }

  static String? _toThaiDate(dynamic dateStr) {
    if (dateStr == null || dateStr == "") return null;
    try {
      DateTime utc = DateTime.parse(dateStr).toUtc();
      DateTime thaiTime = utc.add(const Duration(hours: 7));

      // แปลงปีเป็น พ.ศ.
      int buddhistYear = thaiTime.year + 543;

      // จัดรูปแบบเป็นวัน/เดือน/ปี (พ.ศ.)
      String day = thaiTime.day.toString().padLeft(2, '0');
      String month = thaiTime.month.toString().padLeft(2, '0');
      return "$day/$month/$buddhistYear";
    } catch (e) {
      return dateStr.toString(); // ถ้าแปลงไม่ได้ให้คืนค่าเดิม
    }
  }
}
