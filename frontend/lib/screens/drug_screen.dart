import 'package:flutter/material.dart';
import 'package:frontend/utils/colors.dart' as color;

class DrugScreen extends StatelessWidget {
  const DrugScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 20, horizontal: 12),
      width: double.infinity,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
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
            style: TextStyle(fontSize: 32, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 10),
          const Text(
            "เพิ่มรายการยา กรุณากดปุ่ม +",
            style: TextStyle(fontSize: 20),
          ),
          const Spacer(),
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
    );
  }
}
