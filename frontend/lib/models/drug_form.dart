class DrugFormModel {
  static String _mapImage(String name) {
    switch (name) {
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

  final int id;
  final String name;
  final String image; // ✅ เพิ่ม field image
  final List<DrugUnitModel> units;

  DrugFormModel({
    required this.id,
    required this.name,
    required this.units,
    required this.image,
  });

  factory DrugFormModel.fromJson(Map<String, dynamic> json) {
    final name = json["form_name"] ?? "-";
    return DrugFormModel(
      id: json["id"],
      name: name,
      image: _mapImage(name), // ✅ map รูปตรงนี้เลย
      units:
          (json["units"] as List<dynamic>?)
              ?.map((e) => DrugUnitModel.fromJson(e))
              .toList() ??
          [],
    );
  }
}

class DrugUnitModel {
  final int id;
  final String name;

  DrugUnitModel({required this.id, required this.name});

  factory DrugUnitModel.fromJson(Map<String, dynamic> json) {
    return DrugUnitModel(id: json["id"], name: json["unit_name"]);
  }
}
