import 'package:flutter/material.dart';
import 'package:frontend/enums/drug_tab.dart';
import 'package:frontend/providers/add_edit_provider.dart';
import 'package:frontend/providers/add_group_provider.dart';
import 'package:frontend/providers/drug_provider.dart';
import 'package:frontend/screens/add_edit_screen.dart';
import 'package:frontend/screens/add_group_drug_screen.dart';
import 'package:frontend/screens/all_drug_screen.dart';
import 'package:frontend/screens/group_drug_screen.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/widgets/tab_button.dart';
import 'package:provider/provider.dart';

class DrugScreen extends StatefulWidget {
  const DrugScreen({super.key});

  @override
  State<DrugScreen> createState() => _DrugScreenState();
}

class _DrugScreenState extends State<DrugScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<DrugProvider>().loadMyMedicines();
      context.read<DrugProvider>().loadGroups();
    });
  }

  @override
  Widget build(BuildContext context) {
    final p = context.watch<DrugProvider>();

    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor2nd,
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        foregroundColor: Colors.white,
        title: const Text(
          "ยาของฉัน",
          style: TextStyle(fontSize: 25, fontWeight: FontWeight.w700),
        ),
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          await context.read<DrugProvider>().loadMyMedicines();
          await context.read<DrugProvider>().loadGroups();
        },
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 20, horizontal: 12),
          child: Stack(
            children: [
              SizedBox(
                width: double.infinity,
                child: p.doseAll.isEmpty
                    ? Column(
                        crossAxisAlignment: CrossAxisAlignment.center,
                        children: [
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Align(
                                alignment: Alignment.topRight,
                                child: SizedBox(
                                  width: 144,
                                  height: 45,
                                  child: ElevatedButton(
                                    style: ElevatedButton.styleFrom(
                                      backgroundColor: const Color(0xFFFFDF6A),
                                      padding: EdgeInsets.zero,
                                      shadowColor: Colors.black,
                                      elevation: 4,
                                      shape: const RoundedRectangleBorder(
                                        borderRadius: BorderRadius.zero,
                                      ),
                                    ),
                                    onPressed: () async {
                                      debugPrint("🟡 ปุ่มกดแล้ว เริ่ม sync...");
                                      await context
                                          .read<DrugProvider>()
                                          .syncHospitalMedicines(context);
                                      ScaffoldMessenger.of(
                                        context,
                                      ).showSnackBar(
                                        const SnackBar(
                                          content: Text(
                                            "✅ โหลดยาจากโรงพยาบาลเรียบร้อย",
                                          ),
                                        ),
                                      );
                                    },
                                    child: const Text(
                                      "เพิ่มยาจากโรงพยาบาล",
                                      style: TextStyle(
                                        color: Colors.black,
                                        fontSize: 16,
                                      ),
                                      maxLines: 1,
                                      overflow: TextOverflow.ellipsis,
                                    ),
                                  ),
                                ),
                              ),
                            ],
                          ),
                          const Spacer(),
                          Opacity(
                            opacity: 0.5,
                            child: Image.asset(
                              "assets/images/drugs.png",
                              width: 210,
                              height: 210,
                            ),
                          ),
                          const SizedBox(height: 50),
                          const Text(
                            "ไม่มีรายการยา",
                            style: TextStyle(
                              fontSize: 32,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 10),
                          const Text(
                            "เพิ่มรายการยา กรุณากดปุ่ม +",
                            style: TextStyle(fontSize: 20),
                          ),
                          const Spacer(),
                        ],
                      )
                    : Column(
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 9),
                            child: Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                TabButton(
                                  onTap: (t) =>
                                      context.read<DrugProvider>().setPage(t),
                                  selectPage: p.page,
                                ),
                                Expanded(
                                  child: p.page == DrugTab.all
                                      ? Align(
                                          alignment: Alignment.topRight,
                                          child: SizedBox(
                                            width: 144,
                                            height: 45,
                                            child: ElevatedButton(
                                              style: ElevatedButton.styleFrom(
                                                backgroundColor: const Color(
                                                  0xFFFFDF6A,
                                                ),
                                                padding: EdgeInsets.zero,
                                                elevation: 4,
                                                shadowColor: Colors.black,
                                                shape:
                                                    const RoundedRectangleBorder(
                                                      borderRadius:
                                                          BorderRadius.zero,
                                                    ),
                                              ),
                                              onPressed: () async {
                                                await context
                                                    .read<DrugProvider>()
                                                    .syncHospitalMedicines(
                                                      context,
                                                    );
    
                                              },
                                              child: const Text(
                                                "เพิ่มยาจากโรงพยาบาล",
                                                style: TextStyle(
                                                  color: Colors.black,
                                                  fontSize: 16,
                                                ),
                                                maxLines: 1,
                                                overflow: TextOverflow.ellipsis,
                                              ),
                                            ),
                                          ),
                                        )
                                      : Align(
                                          alignment: Alignment.topRight,
                                          child: SizedBox(
                                            width: 32,
                                            height: 45,
                                            child: RawMaterialButton(
                                              onPressed: () {
                                                Navigator.push(
                                                  context,
                                                  MaterialPageRoute(
                                                    builder: (_) => MultiProvider(
                                                      providers: [
                                                        ChangeNotifierProvider.value(
                                                          value: context
                                                              .read<
                                                                DrugProvider
                                                              >(),
                                                        ),
                                                        ChangeNotifierProvider(
                                                          create: (_) =>
                                                              AddGroupProvider(),
                                                        ),
                                                      ],
                                                      child:
                                                          const AddGroupDrug(),
                                                    ),
                                                  ),
                                                );
                                              },
                                              shape: const CircleBorder(),
                                              fillColor: const Color(
                                                0xFFFF92DB,
                                              ),
                                              child: const Icon(
                                                Icons.add,
                                                color: Colors.black,
                                                size: 28,
                                              ),
                                            ),
                                          ),
                                        ),
                                ),
                              ],
                            ),
                          ),
                          const SizedBox(height: 15),
                          Expanded(
                            child: p.page == DrugTab.all
                                ? const AllDrugScreen()
                                : const GroupDrugScreen(),
                          ),
                        ],
                      ),
              ),
              if (p.page == DrugTab.all) ...[
                Align(
                  alignment: Alignment.bottomLeft,
                  child: SizedBox(
                    width: 70,
                    height: 70,
                    child: RawMaterialButton(
                      onPressed: () async {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (_) => MultiProvider(
                              providers: [
                                ChangeNotifierProvider.value(
                                  value: context.read<DrugProvider>(),
                                ),
                                ChangeNotifierProvider(
                                  create: (_) =>
                                      AddEditProvider(pageFrom: "add"),
                                ),
                              ],
                              child: const AddEditView(),
                            ),
                          ),
                          // MaterialPageRoute(builder: (context) => const AddEditScreen(),settings: const RouteSettings(arguments: {
                          //   "pageType":"add"
                          // }))
                        );
                      },
                      shape: const CircleBorder(),
                      fillColor: Colors.transparent,
                      highlightColor: Colors.blueAccent.withOpacity(0.1),
                      splashColor: Colors.blueAccent.withOpacity(0.1),
                      child: const Icon(
                        Icons.add,
                        color: Colors.white,
                        size: 36,
                      ),
                    ),
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
