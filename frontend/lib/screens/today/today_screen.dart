import 'package:flutter/material.dart';
import 'package:frontend/utils/colors.dart' as color;

class TodayScreen extends StatelessWidget {
  const TodayScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: <Widget>[
          const Spacer(),
          Opacity(
            opacity: 0.5,
            child: Image.asset(
              "assets/images/drugs.png",
              height: 210,
              width: 210,
            ),
          ),
          const SizedBox(height: 50),
          const Text(
            "วันนี้ไม่ต้องกินยา",
            style: TextStyle(fontWeight: FontWeight.bold, fontSize: 32),
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
                    splashColor: Colors.blueAccent.withOpacity(0.1,), // สีคลื่นน้ำกระจาย
                    highlightColor: Colors.blueAccent.withOpacity(0.1),
                    child: const Icon(Icons.add, color: Colors.white, size: 36),
                  ),
                ),
            ),
        ],
      ),
    );
  }
}