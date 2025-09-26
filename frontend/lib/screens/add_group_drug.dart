import 'package:flutter/material.dart';
import 'package:frontend/providers/drug_provider.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:provider/provider.dart';

class AddGroupDrug extends StatelessWidget {
  const AddGroupDrug({super.key});

  UnderlineInputBorder _inputBorder(Color color) {
    return UnderlineInputBorder(borderSide: BorderSide(color: color, width: 1));
  }

  @override
  Widget build(BuildContext context) {
    final dp = context.watch<DrugProvider>();
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
              onTap: () {
                showDialog(
                  context: context,
                  barrierDismissible: false,
                  builder: (context) {
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
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              mainAxisAlignment: MainAxisAlignment.start,
                              children: [
                                IconButton(
                                  padding: const EdgeInsets.all(0),
                                  onPressed: () {
                                    Navigator.pop(context);
                                  },
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
                            ...dp.doseAll.where((d) => !d.import).toList().map((
                              dp,
                            ) {
                              return CheckboxListTile(
                                title: Text(dp.name),
                                value: true,
                                onChanged: (value) {},
                              );
                            }),
                          ],
                        ),
                      ),
                    );
                  },
                );
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
          ],
        ),
      ),
    );
  }
}
