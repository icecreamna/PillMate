import 'package:flutter/material.dart';
import 'package:frontend/providers/add_group_provider.dart';
import 'package:frontend/providers/drug_provider.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:provider/provider.dart';

class AddGroupDrug extends StatelessWidget {
  const AddGroupDrug({super.key});

  UnderlineInputBorder _inputBorder(Color c) {
    return UnderlineInputBorder(borderSide: BorderSide(color: c, width: 1));
  }

  @override
  Widget build(BuildContext context) {
    final dp = context.watch<DrugProvider>();
    final agp = context.watch<AddGroupProvider>();

    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor2nd,
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        foregroundColor: Colors.white,
        title: const Text(
          "สร้างกลุ่มยา",
          style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold),
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
                padding: const EdgeInsets.all(12),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text("ชื่อกลุ่มยา", style: TextStyle(fontSize: 20)),
                    const SizedBox(height: 15),
                    TextFormField(
                      decoration: InputDecoration(
                        hint: const Text(
                          "กลุ่ม 1 ",
                          style: TextStyle(
                            color: Color(0xFF959595),
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
                    final prevChosen = agp.selectedList;
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
                  agp.setSelectedList(chosen);
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
            const SizedBox(height: 20),
            SizedBox(
              width: double.infinity,
              height: 300,
              child: ListView.builder(
                itemCount: agp.selectedList.length,
                itemBuilder: (context, index) {
                  final selectId = agp.selectedList[index];
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
                                  agp.removeSelected(selectId);
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
            const Spacer(),
            Container(
              width: 384,
              height: 70,
              margin: const EdgeInsets.only(bottom: 60),
              child: ElevatedButton(
                onPressed: () {},
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF94B4C1),
                  elevation: 4,
                  shape: BeveledRectangleBorder(
                    borderRadius: BorderRadius.circular(5),
                  ),
                ),
                child: const Text(
                  "สร้าง",
                  style: TextStyle(color: Colors.white, fontSize: 24),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
