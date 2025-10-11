import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:frontend/models/drug_form.dart';

import 'package:frontend/providers/add_edit_provider.dart';
import 'package:frontend/providers/add_single_notification_provider.dart';
import 'package:frontend/providers/drug_provider.dart';
import 'package:provider/provider.dart';
import 'package:frontend/utils/colors.dart' as color;

import '../models/dose.dart';

// class AddEditScreen extends StatelessWidget {
//   const AddEditScreen({super.key});

//   @override
//   Widget build(BuildContext context) {
//     final args =
//         ModalRoute.of(context)!.settings.arguments as Map<String, dynamic>;
//     final String pageFrom = args['pageType'];
//     return MultiProvider(
//       providers: [
//         ChangeNotifierProvider(create:(context) => DrugProvider(),),
//         ChangeNotifierProvider(create:(context) => AddEditProvider(pageFrom: pageFrom),),
//       ],
//       child: _AddEditView(),
//     );
//   }
// }

// UnderlineInputBorder _borderUnderLine(Color color) {
//   return UnderlineInputBorder(borderSide: BorderSide(width: 1, color: color));
// }

class AddEditView extends StatefulWidget {
  // final Dose? dose;
  const AddEditView({super.key});

  @override
  State<AddEditView> createState() => AddEditViewState();
}

class AddEditViewState extends State<AddEditView> {
  OutlineInputBorder _borderInput(Color color) {
    return OutlineInputBorder(
      borderRadius: BorderRadius.circular(12),
      borderSide: BorderSide(color: color),
    );
  }

  Text _headText(String text) {
    return Text(
      text,
      style: const TextStyle(
        color: Colors.black,
        fontWeight: FontWeight.normal,
        fontSize: 24,
        letterSpacing: 0,
      ),
    );
  }

  Container _container(double width, double height, List<Widget> children) {
    return Container(
      width: width,
      height: height,
      decoration: BoxDecoration(
        color: Colors.white,
        border: Border.all(color: Colors.black, width: 1),
        borderRadius: BorderRadius.circular(15),
      ),
      child: Padding(
        padding: const EdgeInsets.all(10.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: children,
        ),
      ),
    );
  }

  final _formKey = GlobalKey<FormState>();

  final _nameDrugController = TextEditingController();
  final _descriptionController = TextEditingController();
  final _amountPerTimeController = TextEditingController();
  final _frequencyController = TextEditingController();

  bool _initializedText = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    final addP = context.read<AddEditProvider>();

    // ✅ เซ็ตค่า text หลัง loadInitialData เสร็จ และเฉพาะตอน edit เท่านั้น
    if (!_initializedText &&
        !addP.isLoading &&
        addP.pageFrom == "edit" &&
        addP.editDose != null) {
      final d = addP.editDose!;
      _nameDrugController.text = d.name;
      _descriptionController.text = d.description;
      _amountPerTimeController.text = d.amountPerDose;
      _frequencyController.text = d.frequency;
      _initializedText = true;
    }
  }

  @override
  void dispose() {
    _nameDrugController.dispose();
    _descriptionController.dispose();
    _amountPerTimeController.dispose();
    _frequencyController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final addP = context.watch<AddEditProvider>();

    if (addP.isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor2nd,
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        foregroundColor: Colors.white,
        title: Text(
          addP.pageFrom == "add" ? "เพิ่มรายการยา" : "แก้ไขรายการยา",
          style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 25),
        ),
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                _container(384, 125, [
                  _headText("ชื่อยา"),
                  TextFormField(
                    inputFormatters: [
                      FilteringTextInputFormatter.allow(RegExp(r'[a-zA-Zก-ฮ]')),
                    ],
                    decoration: const InputDecoration(
                      hint: Text("ใส่ชื่อยาหรือยี่ห้อยา"),
                      enabledBorder: UnderlineInputBorder(
                        borderSide: BorderSide(width: 1, color: Colors.black),
                      ),
                      focusedBorder: UnderlineInputBorder(
                        borderSide: BorderSide(width: 1, color: Colors.black),
                      ),
                      contentPadding: EdgeInsets.only(bottom: 0),
                    ),
                    controller: _nameDrugController,
                    validator: (name) {
                      if (name == null || name.trim().isEmpty) {
                        return "กรุณากรอกค่า";
                      }
                      return null;
                    },
                  ),
                ]),
                const SizedBox(height: 15),
                _container(384, 125, [
                  _headText("สรรพคุณ"),
                  TextFormField(
                    inputFormatters: [
                      FilteringTextInputFormatter.allow(RegExp(r'[a-zA-Zก-ฮ]')),
                    ],
                    decoration: const InputDecoration(
                      hint: Text("ใส่สรรพคุณยา, รักษาอาการ "),
                      enabledBorder: UnderlineInputBorder(
                        borderSide: BorderSide(width: 1, color: Colors.black),
                      ),
                      focusedBorder: UnderlineInputBorder(
                        borderSide: BorderSide(width: 1, color: Colors.black),
                      ),
                      contentPadding: EdgeInsets.all(0),
                    ),
                    validator: (des) {
                      if (des == null || des.trim().isEmpty) {
                        return "กรุณากรอกค่า";
                      }
                      return null;
                    },
                    controller: _descriptionController,
                  ),
                ]),
                const SizedBox(height: 15),
                _container(384, 191, [
                  _headText("ลักษณะของยา"),
                  const SizedBox(height: 10),
                  SizedBox(
                    width: 360,
                    height: 109,
                    child: ListView(
                      scrollDirection: Axis.horizontal,
                      children: addP.forms.map((df) {
                        return GestureDetector(
                          onTap: () {
                            _amountPerTimeController.clear();
                            _frequencyController.clear();
                            addP.setSelectForm(df);
                          },
                          child: Container(
                            width: 77,
                            height: 109,
                            margin: const EdgeInsets.only(
                              right: 5,
                              left: 3,
                              top: 3,
                              bottom: 3,
                            ),
                            padding: const EdgeInsets.all(12),
                            decoration: BoxDecoration(
                              color: addP.selectedForm == df
                                  ? const Color(0xFF84E8FF)
                                  : const Color(0xFFD9F8FF),
                              borderRadius: BorderRadius.circular(12),
                              boxShadow: [
                                BoxShadow(
                                  color: Colors.black.withOpacity(0.3),
                                  blurRadius: 8,
                                  offset: const Offset(0, 4),
                                ),
                              ],
                            ),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.center,
                              children: [
                                Image.asset(df.image, height: 40, width: 40),
                                const Spacer(),
                                Text(df.name),
                              ],
                            ),
                          ),
                        );
                      }).toList(),
                    ),
                  ),
                ]),
                const SizedBox(height: 15),
                _container(384, 203, [
                  _headText("การใช้ยา"),
                  const SizedBox(height: 17),
                  Row(
                    children: [
                      SizedBox(
                        width: 171,
                        height: 50,
                        child: TextFormField(
                          inputFormatters: [
                            FilteringTextInputFormatter.allow(RegExp(r'[0-9]')),
                          ],
                          decoration: InputDecoration(
                            labelText: "ครั้งละ",
                            border: _borderInput(Colors.grey),
                            enabledBorder: _borderInput(Colors.grey),
                            focusedBorder: _borderInput(Colors.grey),
                            disabledBorder: _borderInput(Colors.grey),
                            errorBorder: _borderInput(Colors.grey),
                            focusedErrorBorder: _borderInput(Colors.grey),
                            errorStyle: const TextStyle(height: 0, fontSize: 0),
                          ),
                          validator: (amount) {
                            if (amount == null || amount.trim().isEmpty) {
                              return "กรุณากรอกค่า";
                            }
                            return null;
                          },
                          controller: _amountPerTimeController,
                        ),
                      ),
                      const SizedBox(width: 15),
                      SizedBox(
                        width: 171,
                        height: 50,
                        child: DropdownButtonFormField<DrugUnitModel>(
                          value: addP.selectedUnit,
                          decoration: InputDecoration(
                            label: const Text("หน่วย"),
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: const BorderSide(color: Colors.grey),
                            ),
                          ),
                          items: (addP.selectedForm?.units ?? []).map((unit) {
                            final unitName = unit.name ?? "ไม่ระบุหน่วย";
                            return DropdownMenuItem(
                              value: unit,
                              child: Text(unitName),
                            );
                          }).toList(),
                          onChanged: (unit) {
                            if (unit != null) {
                              if (unit == addP.selectedUnit) return;
                              _amountPerTimeController.clear();
                              _frequencyController.clear();
                              addP.setUnit(unit);
                            }
                            return;
                          },
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 17),
                  Row(
                    children: [
                      SizedBox(
                        width: 171,
                        height: 50,
                        child: TextFormField(
                          inputFormatters: [
                            FilteringTextInputFormatter.allow(RegExp(r'[0-9]')),
                          ],
                          decoration: InputDecoration(
                            labelText: "วันละ",
                            enabledBorder: _borderInput(Colors.grey),
                            focusedBorder: _borderInput(Colors.grey),
                            disabledBorder: _borderInput(Colors.grey),
                            errorBorder: _borderInput(Colors.grey),
                            focusedErrorBorder: _borderInput(Colors.grey),
                            errorStyle: const TextStyle(height: 0, fontSize: 0),
                          ),
                          validator: (fr) {
                            if (fr == null || fr.trim().isEmpty) {
                              return "กรุณากรอกค่า";
                            }
                            return null;
                          },
                          controller: _frequencyController,
                        ),
                      ),
                      const SizedBox(width: 20),
                      const Text(
                        "ครั้ง",
                        style: TextStyle(
                          color: Colors.black,
                          fontSize: 16,
                          fontWeight: FontWeight.normal,
                        ),
                      ),
                    ],
                  ),
                ]),
                const SizedBox(height: 15),
                _container(384, 223, [
                  _headText("ช่วงเวลาใช้ยา"),
                  const SizedBox(height: 18),
                  SizedBox(
                    height: 149,
                    child: GridView.count(
                      crossAxisCount: 2,
                      childAspectRatio: 3,
                      physics: const NeverScrollableScrollPhysics(),
                      children: addP.times.map((dt) {
                        return GestureDetector(
                          onTap: () {
                            addP.setSelectTime(dt);
                          },
                          child: Container(
                            width: 171,
                            height: 50,
                            margin: const EdgeInsets.only(
                              left: 3,
                              bottom: 8,
                              right: 3,
                              top: 3,
                            ),
                            decoration: BoxDecoration(
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(color: Colors.black, width: 1),
                              boxShadow: [
                                BoxShadow(
                                  color: Colors.black.withOpacity(0.3),
                                  blurRadius: 8,
                                  offset: const Offset(0, 4),
                                ),
                              ],
                              color: addP.selectTime == dt
                                  ? const Color(0xFF84E8FF)
                                  : Colors.white,
                            ),
                            child: Center(
                              child: Text(
                                dt.name,
                                style: const TextStyle(
                                  fontWeight: FontWeight.bold,
                                  fontSize: 20,
                                ),
                              ),
                            ),
                          ),
                        );
                      }).toList(),
                    ),
                  ),
                ]),
                const SizedBox(height: 15),
                ElevatedButton(
                  onPressed: () async {
                    if (!_formKey.currentState!.validate()) {
                      return;
                    }
                    if (addP.pageFrom == "edit" && addP.editDose != null) {
                      final success = await addP.editMedicine(
                        id: int.parse(addP.editDose!.id),
                        name: _nameDrugController.text.trim(),
                        properties: _descriptionController.text.trim(),
                        amountPerTime: _amountPerTimeController.text.trim(),
                        timePerDay: _frequencyController.text.trim(),
                        selectedFormId:
                            addP.selectedForm?.id ?? addP.editDose!.formId ?? 0,
                        selectedUnitId:
                            addP.selectedUnit?.id ?? addP.editDose!.unitId ?? 0,
                        selectTimeId:
                            addP.selectTime?.id ??
                            addP.editDose!.instructionId ??
                            0,
                      );
                      if (success && context.mounted) {
                        ScaffoldMessenger.of(context).showSnackBar(
                          const SnackBar(content: Text("✅ แก้ไขยาเรียบร้อย")),
                        );
                        Navigator.pop(context,true);
                      } else {
                        ScaffoldMessenger.of(context).showSnackBar(
                          const SnackBar(content: Text("❌ แก้ไขยาไม่สำเร็จ")),
                        );
                      }
                    } else {
                      final success = await addP.addMedicine(
                        name: _nameDrugController.text.trim(),
                        properties: _descriptionController.text.trim(),
                        amountPerTime: _amountPerTimeController.text.trim(),
                        timePerDay: _frequencyController.text.trim(),
                      );

                      if (success && context.mounted) {
                        _amountPerTimeController.clear();
                        _frequencyController.clear();
                        _nameDrugController.clear();
                        _descriptionController.clear();
                        ScaffoldMessenger.of(context).showSnackBar(
                          const SnackBar(content: Text("✅ เพิ่มยาเรียบร้อย")),
                        );
                        Navigator.pop(context);
                      } else {
                        ScaffoldMessenger.of(context).showSnackBar(
                          const SnackBar(content: Text("❌ เพิ่มยาไม่สำเร็จ")),
                        );
                      }
                    }
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF94B4C1),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(7),
                    ),
                  ),
                  child: const Text(
                    "บันทึก",
                    style: TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.normal,
                      fontSize: 24,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
