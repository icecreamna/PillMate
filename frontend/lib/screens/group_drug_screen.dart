import 'package:flutter/material.dart';
import 'package:frontend/providers/add_group_notification_provider.dart';
import 'package:frontend/providers/add_group_provider.dart';
import 'package:frontend/providers/drug_provider.dart';
import 'package:frontend/screens/add_group_notification_screen.dart';
import 'package:provider/provider.dart';
// import 'package:frontend/providers/drug_provider.dart';
// import 'package:provider/provider.dart';

class GroupDrugScreen extends StatelessWidget {
  const GroupDrugScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final dp = context.watch<DrugProvider>();
    final keys = dp.groups.keys.toList();
    return ListView.builder(
      itemCount: keys.length,
      itemBuilder: (context, index) {
        final key = keys[index];
        final value = dp.groups[key]!;
        return Padding(
          padding: const EdgeInsets.only(bottom: 10),
          child: SizedBox(
            width: 384,
            height: 80,
            child: Card(
              color: Colors.white,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
                side: const BorderSide(color: Colors.grey, width: 0.5),
              ),
              child: InkWell(
                onTap: () {
                  final groupId = dp.groupsId[key]; // 🆔 เก็บ id ของกลุ่ม

                  if (groupId == null) {
                    debugPrint("❌ ไม่พบ groupId สำหรับกลุ่ม $key");
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text("ไม่พบข้อมูลกลุ่มนี้")),
                    );
                    return; // ❌ หยุดไม่ไปหน้าใหม่
                  }
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (_) => MultiProvider(
                        providers: [
                          ChangeNotifierProvider.value(
                            value: context.read<DrugProvider>(),
                          ),
                          ChangeNotifierProvider(
                            create: (_) => AddGroupProvider(),
                          ),
                          ChangeNotifierProvider(
                            create: (_) => AddGroupNotificationProvider(
                              groupId: groupId,
                              keyName: key,
                            ),
                          ),
                        ],
                        child: const AddGroupNotificationScreen(),
                      ),
                    ),
                  );
                },
                child: Padding(
                  padding: const EdgeInsets.symmetric(
                    vertical: 5,
                    horizontal: 10,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        key,
                        style: const TextStyle(
                          color: Colors.black,
                          fontSize: 20,
                        ),
                      ),
                      const SizedBox(height: 7),
                      Text(
                        "${value.length} รายการ",
                        style: const TextStyle(
                          color: Color(0xFF454545),
                          fontSize: 16,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
        );
      },
    );
  }
}
