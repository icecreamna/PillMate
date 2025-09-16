import 'package:flutter/material.dart';
import 'package:frontend/providers/drug_provider.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/widgets/tab_button.dart';
import 'package:provider/provider.dart';

class DrugScreen extends StatelessWidget {
  const DrugScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final p = context.watch<DrugProvider>();

    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor2nd,
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        title: const Text(
          "ยาของฉัน",
          style: TextStyle(
            fontSize: 25,
            fontWeight: FontWeight.w700,
            color: Colors.white,
          ),
        ),
      ),
      body: Padding(
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
                                    elevation: 3,
                                    shape: const RoundedRectangleBorder(
                                      borderRadius: BorderRadius.zero,
                                    ),
                                  ),
                                  onPressed: () {},
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
                          padding: const EdgeInsetsGeometry.symmetric(
                            horizontal: 7,
                          ),
                          child: Row(
                            children: [
                              TabButton(
                                onTap: (t) => context.read<DrugProvider>().setPage(t),
                                selectPage: p.page,
                              ),
                              Expanded(
                                child: Align(
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
                                        elevation: 3,
                                        shadowColor: Colors.black,
                                        shape: const RoundedRectangleBorder(
                                          borderRadius: BorderRadius.zero,
                                        ),
                                      ),
                                      onPressed: () {},
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
                              ),
                            ],
                          ),
                        ),
                        const SizedBox(height: 15),
                        Expanded(
                          child: ListView.builder(
                            itemBuilder: (_, index) {
                              final d = p.doseAll[index];
                              return Padding(
                                padding: const EdgeInsetsGeometry.only(
                                  bottom: 10,
                                ),
                                child: SizedBox(
                                  width: 384,
                                  height: 135,
                                  child: Card(
                                    shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(12),
                                      side: const BorderSide(
                                        color: Colors.grey,
                                        width: 0.5,
                                      ),
                                    ),
                                    color: d.import
                                        ? const Color(0xFFFFF5D0)
                                        : Colors.white,
                                    child: InkWell(
                                      onTap: () {},
                                      child: Padding(
                                        padding: const EdgeInsets.fromLTRB(
                                          10,
                                          13,
                                          16,
                                          0,
                                        ),
                                        child: Row(
                                          crossAxisAlignment:
                                              CrossAxisAlignment.start,
                                          children: [
                                            Expanded(
                                              child: Column(
                                                crossAxisAlignment:
                                                    CrossAxisAlignment.start,
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
                                                    d.drugIndication,
                                                    style: const TextStyle(
                                                      color: Colors.black,
                                                      fontSize: 16,
                                                    ),
                                                  ),
                                                  Text(
                                                    "ครั้งละ " +
                                                        d.numberOfTake +
                                                        "เม็ด" +
                                                        " " +
                                                        "วันละ " +
                                                        d.takePerDay +
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
                                              crossAxisAlignment:
                                                  CrossAxisAlignment.end,
                                              children: [
                                                Image.asset(
                                                  d.picture,
                                                  width: 40,
                                                  height: 40,
                                                ),
                                                const SizedBox(height: 40),
                                                Text(
                                                  d.import
                                                      ? "(โรงพยาบาล)"
                                                      : "(ของฉัน)",
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
                          ),
                        ),
                      ],
                    ),
            ),
            Align(
              alignment: Alignment.bottomLeft,
              child: SizedBox(
                width: 70,
                height: 70,
                child: RawMaterialButton(
                  onPressed: () {},
                  shape: const CircleBorder(),
                  fillColor: color.AppColors.backgroundColor1st,
                  highlightColor: Colors.blueAccent.withOpacity(0.1),
                  splashColor: Colors.blueAccent.withOpacity(0.1),
                  child: const Icon(Icons.add, color: Colors.white, size: 36),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
