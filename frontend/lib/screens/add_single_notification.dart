import 'package:flutter/material.dart';
import 'package:frontend/providers/add_edit_provider.dart';
import 'package:frontend/screens/add_edit_screen.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/providers/drug_provider.dart';
import 'package:provider/provider.dart';

class AddSingleNotification extends StatelessWidget {
  final DoseTest dose;

  const AddSingleNotification({super.key, required this.dose});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor2nd,
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        foregroundColor: Colors.white,
        title: const Text(
          "รายาการยา",
          style: TextStyle(fontWeight: FontWeight.bold, fontSize: 25),
        ),
      ),
      body: Padding(
        padding: const EdgeInsets.symmetric(vertical: 20, horizontal: 12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            SizedBox(
              width: 384,
              child: Card(
                color: dose.import ? const Color(0xFFFFF5D0) : Colors.white,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                  side: const BorderSide(color: Colors.grey, width: 0.5),
                ),
                child: Padding(
                  padding: const EdgeInsets.fromLTRB(10, 13, 16, 0),
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              dose.name,
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 20,
                              ),
                            ),
                            const SizedBox(height: 5),
                            Text(
                              dose.description,
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 16,
                              ),
                            ),
                            Text(
                              "ครั้งละ " +
                                  dose.numberOfTake +
                                  dose.unit +
                                  " " +
                                  "วันละ " +
                                  dose.takePerDay +
                                  " " +
                                  "ครั้ง",
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 16,
                              ),
                            ),
                            Text(
                              dose.instruction,
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 16,
                              ),
                            ),
                            const SizedBox(height: 10),
                            if (!dose.import) ...[
                              SizedBox(
                                width: 95,
                                height: 35,
                                child: ElevatedButton(
                                  onPressed: () {
                                    Navigator.push(
                                      context,
                                      MaterialPageRoute(
                                        builder: (_) => MultiProvider(
                                          providers: [
                                            ChangeNotifierProvider.value(
                                              value: context
                                                  .read<DrugProvider>(),
                                            ),
                                            ChangeNotifierProvider(
                                              create: (_) => AddEditProvider(
                                                pageFrom: "edit",
                                                editDose: dose
                                              ),
                                            ),
                                          ],
                                          child: const AddEditView(),
                                        ),
                                      ),
                                    );
                                  },
                                  style: ElevatedButton.styleFrom(
                                    backgroundColor: const Color(0xFF55FF00),
                                    padding: const EdgeInsets.symmetric(
                                      vertical: 4,
                                    ),
                                    shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(5),
                                    ),
                                  ),
                                  child: const Text(
                                    "แก้ไขข้อมูลยา",
                                    style: TextStyle(
                                      fontSize: 16,
                                      color: Colors.black,
                                    ),
                                    maxLines: 1,
                                    overflow: TextOverflow.visible,
                                  ),
                                ),
                              ),
                            ],
                            const SizedBox(height: 15),
                          ],
                        ),
                      ),
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.end,
                        children: [
                          Image.asset(dose.picture, width: 40, height: 40),
                          const SizedBox(height: 40),
                          Text(
                            dose.import ? "(โรงพยาบาล)" : "(ของฉัน)",
                            style: const TextStyle(
                              color: Colors.black,
                              fontSize: 16,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
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
                onPressed: () {},
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF55FF00),
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
          ],
        ),
      ),
    );
  }
}
