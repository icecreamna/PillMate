import 'package:flutter/material.dart';
import 'package:frontend/enums/drug_tab.dart';

class DoseTest {
  final String name;
  final String numberOfTake;
  final String takePerDay;
  final String instruction;
  final String picture;
  final String description;
  final String unit;
  final bool import;

  DoseTest({
    required this.name,
    required this.description,
    required this.import,
    required this.numberOfTake,
    required this.takePerDay,
    required this.instruction,
    required this.unit,
    required this.picture,
  });
}

class DrugProvider extends ChangeNotifier {
  DrugTab _page = DrugTab.all;

  DrugTab get page => _page;

  final List<DoseTest> _all = [
    DoseTest(
      name: "Paracetamol",
      description: "แก้ปวดหัว",
      import: false,
      numberOfTake: "1",
      takePerDay: "3",
      instruction: "หลังอาหาร",
      unit: "เม็ด",
      picture: "assets/images/pill.png",
    ),
    DoseTest(
      name: "Paracetamol",
      description: "แก้ปวดหัว",
      import: true,
      numberOfTake: "1",
      takePerDay: "3",
      instruction: "ก่อนอาหาร",
      unit: "เม็ด",
      picture: "assets/images/pill.png",
    ),
    DoseTest(
      name: "ยาทา1",
      description: "แก้ระคายเคือง",
      import: false,
      numberOfTake: "1",
      takePerDay: "3",
      instruction: "หลังอาหาร",
      unit: "ช้อน",
      picture: "assets/images/ointment.png",
    ),
    DoseTest(
      name: "ยาน้ำ1",
      description: "ลดไข้",
      import: true,
      numberOfTake: "0.5",
      takePerDay: "4",
      instruction: "หลังอาหาร",
      unit: "ช้อน",
      picture: "assets/images/syrup.png",
    ),
    DoseTest(
      name: "ยาน้ำ12",
      description: "ลดไข้",
      import: true,
      numberOfTake: "2",
      takePerDay: "4",
      instruction: "หลังอาหาร",
      unit: "มิลลิต",
      picture: "assets/images/syrup.png",
    ),
  ];

  List<DoseTest> get doseAll => _all;

  void setPage(DrugTab selectPage) {
    _page = selectPage;
    notifyListeners();
  }

  void addDose(DoseTest newDose) {
    _all.add(newDose);
    print("มี ${_all.length}");
    notifyListeners();
  }
}
