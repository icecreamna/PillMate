import 'package:flutter/material.dart';
import 'package:frontend/providers/add_group_notification_provider.dart';
import 'package:frontend/providers/add_notification_provider.dart';
// import 'package:frontend/providers/add_group_provider.dart';
import 'package:frontend/providers/drug_provider.dart';
import 'package:frontend/screens/add_notification_screen.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:provider/provider.dart';

class AddGroupNotificationScreen extends StatelessWidget {
  const AddGroupNotificationScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Builder(
      builder: (context) {
        return _AddGroupNotificationView();
      },
    );
  }
}

class _AddGroupNotificationView extends StatefulWidget {
  @override
  State<_AddGroupNotificationView> createState() =>
      _AddGroupNotificationScreenState();
}

class _AddGroupNotificationScreenState
    extends State<_AddGroupNotificationView> {
  final _formKey = GlobalKey<FormState>();

  UnderlineInputBorder _inputBorder(Color c) {
    return UnderlineInputBorder(borderSide: BorderSide(color: c, width: 1));
  }

  @override
  Widget build(BuildContext context) {
    final dp = context.watch<DrugProvider>();
    // final agp = context.watch<AddGroupProvider>();
    final addG = context.watch<AddGroupNotificationProvider>();
    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor2nd,
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        foregroundColor: Colors.white,
        title: const Text(
          "กลุ่มยา",
          style: TextStyle(fontWeight: FontWeight.bold, fontSize: 25),
        ),
      ),
      body: Padding(
        padding: const EdgeInsets.symmetric(vertical: 20, horizontal: 12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              width: 384,
              height: 125,
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: Colors.grey, width: 1),
              ),
              child: Padding(
                padding: const EdgeInsets.all(12.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text("ชื่อกลุ่มยา", style: TextStyle(fontSize: 20)),
                    Form(
                      key: _formKey,
                      child: TextFormField(
                        enabled: false,
                        decoration: InputDecoration(
                          hint: Text(
                            addG.keyName,
                            style: const TextStyle(
                              color: Colors.black,
                              fontSize: 20,
                            ),
                          ),
                          border: _inputBorder(Colors.grey),
                          enabledBorder: _inputBorder(Colors.grey),
                          errorBorder: _inputBorder(Colors.grey),
                          focusedBorder: _inputBorder(Colors.grey),
                          disabledBorder: _inputBorder(Colors.grey),
                          focusedErrorBorder: _inputBorder(Colors.grey),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 25),
            InkWell(
              onTap: () async {
                final chosen = await showDialog<List<String>>(
                  context: context,
                  barrierDismissible: false,
                  builder: (context) {
                    final myDrugs = dp.doseAll.where((d) => !d.import).toList();
                    final hospitalDrugs = dp.doseAll
                        .where((d) => d.import)
                        .toList();
                    final prevChosen = addG.value;
                    List<bool> selectedMy = List.generate(
                      myDrugs.length,
                      (i) => prevChosen.contains(myDrugs[i].id),
                    );
                    List<bool> selectedHospital = List.generate(
                      hospitalDrugs.length,
                      (i) => prevChosen.contains(hospitalDrugs[i].id),
                    );
                    return Dialog(
                      child: Container(
                        width: 328,
                        height: 753,
                        padding: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 8,
                        ),
                        decoration: BoxDecoration(
                          color: Colors.white,
                          borderRadius: BorderRadius.circular(13),
                        ),
                        child: StatefulBuilder(
                          builder: (context, setState) {
                            return Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Row(
                                  children: [
                                    IconButton(
                                      padding: EdgeInsets.zero,
                                      onPressed: () => Navigator.pop(context),
                                      icon: const Icon(
                                        Icons.close,
                                        size: 40,
                                        color: Colors.black,
                                      ),
                                    ),
                                    const SizedBox(width: 27),
                                    const Text(
                                      "เลือกรายการยา",
                                      style: TextStyle(
                                        color: Colors.black,
                                        fontSize: 24,
                                      ),
                                    ),
                                  ],
                                ),
                                const Padding(
                                  padding: EdgeInsets.symmetric(
                                    horizontal: 12,
                                    vertical: 15,
                                  ),
                                  child: Text(
                                    "รายการยา(ของฉัน)",
                                    style: TextStyle(
                                      color: Colors.black,
                                      fontSize: 20,
                                    ),
                                  ),
                                ),
                                Expanded(
                                  child: ListView.builder(
                                    itemCount: myDrugs.length,
                                    itemBuilder: (context, index) {
                                      final drug = myDrugs[index];
                                      return CheckboxListTile(
                                        title: Text(
                                          drug.name,
                                          style: const TextStyle(
                                            color: Colors.black,
                                            fontSize: 16,
                                          ),
                                        ),
                                        value: selectedMy[index],
                                        onChanged: (val) {
                                          setState(() {
                                            selectedMy[index] = val ?? false;
                                          });
                                        },
                                        controlAffinity:
                                            ListTileControlAffinity.leading,
                                      );
                                    },
                                  ),
                                ),
                                const Padding(
                                  padding: EdgeInsets.symmetric(
                                    horizontal: 12,
                                    vertical: 15,
                                  ),
                                  child: Text(
                                    "รายการยา(โรงพยาบาล)",
                                    style: TextStyle(
                                      color: Colors.black,
                                      fontSize: 20,
                                    ),
                                  ),
                                ),
                                Expanded(
                                  child: ListView.builder(
                                    itemCount: hospitalDrugs.length,
                                    itemBuilder: (context, index) {
                                      final drug = hospitalDrugs[index];
                                      return CheckboxListTile(
                                        title: Text(
                                          drug.name,
                                          style: const TextStyle(
                                            color: Colors.black,
                                            fontSize: 16,
                                          ),
                                        ),
                                        value: selectedHospital[index],
                                        onChanged: (val) {
                                          setState(() {
                                            selectedHospital[index] =
                                                val ?? false;
                                          });
                                        },
                                        controlAffinity:
                                            ListTileControlAffinity.leading,
                                      );
                                    },
                                  ),
                                ),
                                Center(
                                  child: SizedBox(
                                    width: 237,
                                    height: 50,
                                    child: ElevatedButton(
                                      onPressed: () {
                                        final chosen = <String>[];
                                        for (
                                          int i = 0;
                                          i < selectedMy.length;
                                          i++
                                        ) {
                                          if (selectedMy[i]) {
                                            chosen.add(myDrugs[i].id);
                                          }
                                        }
                                        for (
                                          int i = 0;
                                          i < selectedHospital.length;
                                          i++
                                        ) {
                                          if (selectedHospital[i]) {
                                            chosen.add(hospitalDrugs[i].id);
                                          }
                                        }
                                        Navigator.pop(context, chosen);
                                        debugPrint("มี id $chosen");
                                      },
                                      style: ElevatedButton.styleFrom(
                                        elevation: 4,
                                        backgroundColor: const Color(
                                          0xFF03B200,
                                        ),
                                        shape: BeveledRectangleBorder(
                                          borderRadius: BorderRadius.circular(
                                            5,
                                          ),
                                        ),
                                      ),
                                      child: const Text(
                                        "ยืนยัน",
                                        style: TextStyle(
                                          color: Colors.white,
                                          fontSize: 24,
                                        ),
                                      ),
                                    ),
                                  ),
                                ),
                                const SizedBox(height: 20),
                              ],
                            );
                          },
                        ),
                      ),
                    );
                  },
                );
                if (chosen != null) {
                  addG.setSelectedList(chosen);
                }
              },
              child: Container(
                width: 384,
                height: 42,
                padding: const EdgeInsets.symmetric(horizontal: 12),
                decoration: BoxDecoration(
                  color: const Color(0xFFD9D9D9),
                  borderRadius: BorderRadius.circular(7),
                ),
                child: const Row(
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    Expanded(
                      child: Text(
                        "เลือกรายการยา",
                        style: TextStyle(fontSize: 16),
                      ),
                    ),
                    Icon(Icons.arrow_forward_outlined, size: 32),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 10),
            Visibility(
              visible: addG.listError.isNotEmpty,
              child: Text(
                addG.listError,
                style: const TextStyle(color: Color(0xFFFF0000), fontSize: 12),
              ),
            ),
            const SizedBox(height: 20),
            SizedBox(
              width: double.infinity,
              height: 150,
              child: ListView.builder(
                itemCount: addG.value.length,
                itemBuilder: (context, index) {
                  final selectId = addG.value[index];
                  final dose = dp.doseAll.firstWhere((d) => d.id == selectId);
                  return Column(
                    children: [
                      Container(
                        width: 384,
                        height: 40,
                        padding: const EdgeInsets.symmetric(
                          vertical: 3,
                          horizontal: 9,
                        ),
                        decoration: const UnderlineTabIndicator(
                          borderSide: BorderSide(width: 1, color: Colors.grey),
                        ),
                        child: Row(
                          children: [
                            Expanded(
                              child: Text(
                                dose.name,
                                style: const TextStyle(
                                  color: Colors.black,
                                  fontSize: 20,
                                ),
                                maxLines: 1,
                                overflow: TextOverflow.ellipsis,
                              ),
                            ),
                            SizedBox(
                              width: 32,
                              height: 32,
                              child: RawMaterialButton(
                                fillColor: const Color(0xFFFF0000),
                                onPressed: () {
                                  addG.removeSelected(selectId);
                                },
                                shape: const CircleBorder(),
                                child: const Icon(
                                  Icons.remove,
                                  size: 28,
                                  color: Colors.black,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 15),
                    ],
                  );
                },
              ),
            ),
            const SizedBox(height: 30),
            const Text("การแจ้งเตือน", style: TextStyle(fontSize: 20)),
            const SizedBox(height: 15),
            const Text(
              "ยังไม่มีข้อมูลการแจ้งเตือน",
              style: TextStyle(color: Color(0xFF959595), fontSize: 20),
            ),
            const SizedBox(height: 90),
            SizedBox(
              width: 120,
              height: 35,
              child: ElevatedButton(
                onPressed: () {
                  bool hasError = false;
                  if (addG.value.length < 2) {
                    hasError = true;
                    addG.setListError();
                  } else {
                    addG.clearListError();
                  }
                  if (hasError) return;
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (_) => MultiProvider(
                        providers: [
                          ChangeNotifierProvider(
                            create: (_) => AddNotificationProvider(
                              pageFrom: "group",
                              keyName: addG.keyName,
                              value: addG.value,
                            ),
                          ),
                        ],
                        child: const AddNotificationScreen(),
                      ),
                    ),
                  );
                },
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF55FF00),
                  elevation: 4,
                  padding: const EdgeInsets.symmetric(vertical: 4),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(5),
                  ),
                ),
                child: const Text(
                  "เพิ่มการแจ้งเตือน",
                  style: TextStyle(color: Colors.black, fontSize: 16),
                  maxLines: 1,
                ),
              ),
            ),
            const Spacer(),
            Row(
              children: [
                Container(
                  width: 181,
                  height: 70,
                  margin: const EdgeInsets.only(bottom: 60),
                  child: ElevatedButton(
                    onPressed: () {
                      dp.removeGroup(addG.keyName);
                      Navigator.pop(context);
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFFFF0000),
                      elevation: 4,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(5),
                      ),
                    ),
                    child: const Text(
                      "ลบรายการยา",
                      style: TextStyle(color: Colors.white, fontSize: 24),
                    ),
                  ),
                ),
                const SizedBox(width: 15),
                Container(
                  width: 181,
                  height: 70,
                  margin: const EdgeInsets.only(bottom: 60),
                  child: ElevatedButton(
                    onPressed: () {
                      bool hasError = false;
                      if (addG.value.length < 2) {
                        hasError = true;
                        addG.setListError();
                      } else {
                        addG.clearListError();
                      }
                      if (hasError) return;

                      dp.updatedDoseGroup(addG.keyName, addG.value);
                      Navigator.pop(context);
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFF94B4C1),
                      elevation: 4,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(5),
                      ),
                    ),
                    child: const Text(
                      "บันทึก",
                      style: TextStyle(color: Colors.white, fontSize: 24),
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}