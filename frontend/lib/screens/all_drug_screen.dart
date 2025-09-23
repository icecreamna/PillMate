import 'package:flutter/material.dart';
import 'package:frontend/providers/drug_provider.dart';
import 'package:frontend/screens/add_single_notification.dart';
import 'package:provider/provider.dart';

class AllDrugScreen extends StatelessWidget {
  const AllDrugScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final p = context.watch<DrugProvider>();
    return ListView.builder(
      itemBuilder: (_, index) {
        final d = p.doseAll[index];
        return Padding(
          padding: const EdgeInsets.only(bottom: 10),
          child: SizedBox(
            width: 384,
            child: Card(
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
                side: const BorderSide(color: Colors.grey, width: 0.5),
              ),
              color: d.import ? const Color(0xFFFFF5D0) : Colors.white,
              child: InkWell(
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (_) => ChangeNotifierProvider.value(
                        value: context.read<DrugProvider>(),
                        child: AddSingleNotification(dose: d),
                      ),
                    ),
                  );
                },
                child: Padding(
                  padding: const EdgeInsets.fromLTRB(10, 13, 16, 0),
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              d.name,
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 20,
                              ),
                            ),
                            const SizedBox(height: 5),
                            Text(
                              d.description,
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 16,
                              ),
                            ),
                            Text(
                              "ครั้งละ " +
                                  d.amountPerDose +
                                  d.unit +
                                  " " +
                                  "วันละ " +
                                  d.frequency +
                                  " " +
                                  "ครั้ง",
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 16,
                              ),
                            ),
                            Text(
                              d.instruction,
                              style: const TextStyle(
                                color: Colors.black,
                                fontSize: 16,
                              ),
                            ),
                          ],
                        ),
                      ),
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.end,
                        children: [
                          Image.asset(d.picture, width: 40, height: 40),
                          const SizedBox(height: 40),
                          Text(
                            d.import ? "(โรงพยาบาล)" : "(ของฉัน)",
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
          ),
        );
      },
      itemCount: p.doseAll.length,
    );
  }
}
