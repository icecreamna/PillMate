import 'package:flutter/material.dart';
import 'package:frontend/enums/drug_tab.dart';

class TabButton extends StatelessWidget {
  final Function(DrugTab) onTap;
  final DrugTab selectPage;

  const TabButton({super.key, required this.onTap, required this.selectPage});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 154,
      height: 37,
      decoration: BoxDecoration(
        border: BoxBorder.all(color: Colors.grey, width: 1),
        borderRadius: BorderRadius.circular(20),
        color: Colors.transparent,
      ),
      child: Row(
        children: [
          Expanded(
            child: GestureDetector(
              onTap: () => onTap(DrugTab.all),
              child: Container(
                decoration: BoxDecoration(
                  color: selectPage == DrugTab.all
                      ? const Color(0xFF84E8FF)
                      : Colors.transparent,
                  border: Border.all(
                    color: selectPage == DrugTab.all
                        ? Colors.black
                        : Colors.transparent,
                    width: 1,
                  ),
                  borderRadius: BorderRadius.circular(20),
                ),
                alignment: Alignment.center,
                child: const Text(
                  "ทั้งหมด",
                  style: TextStyle(
                    color: Colors.black,
                    fontSize: 16,
                    fontWeight: FontWeight.normal,
                  ),
                ),
              ),
            ),
          ),
          Expanded(
            child: GestureDetector(
              onTap: () => onTap(DrugTab.group),
              child: Container(
                decoration: BoxDecoration(
                  color: selectPage == DrugTab.group
                      ? const Color(0xFF84E8FF)
                      : Colors.transparent,
                  border: Border.all(
                    color: selectPage == DrugTab.group
                        ? Colors.black
                        : Colors.transparent,
                    width: 1,
                  ),
                  borderRadius: BorderRadius.circular(20),
                ),
                alignment: Alignment.center,
                child: const Text(
                  "รวมกลุ่ม",
                  style: TextStyle(
                    color: Colors.black,
                    fontSize: 16,
                    fontWeight: FontWeight.normal,
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
