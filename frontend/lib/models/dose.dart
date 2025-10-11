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
    this.groupId
  });


  /// ✅ แปลงจาก JSON ที่มาจาก backend
  factory Dose.fromJson(Map<String, dynamic> json) {
    final formName = json["form_name"] ?? "-";
    return Dose(
      id: json["id"].toString(),
      name: json["med_name"] ?? "-",
      description: json["properties"] ?? "-",
      amountPerDose: json["amount_per_time"] ?? "-",
      frequency: json["times_per_day"] ?? "-",
      instruction: json["instruction_name"] ?? "-",
      unit: json["unit_name"] ?? "-",
      picture: _mapImage(formName),
      import: (json["source"] == "hospital"),

      formId: json["form_id"],
      unitId: json["unit_id"],
      instructionId: json["instruction_id"],
      groupId: json["group_id"]
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
}
