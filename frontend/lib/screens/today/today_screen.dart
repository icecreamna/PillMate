import 'package:flutter/material.dart';
import 'package:frontend/utils/colors.dart' as color;

class TodayScreen extends StatelessWidget {
  const TodayScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      height: double.infinity,
      padding: const EdgeInsets.all(16),
      child: Column(
        children: <Widget>[
          const Center(child: Text("Today screen")),
          Expanded(
            child: Align(
              alignment: Alignment.bottomLeft,
                child: SizedBox(
                  width: 70,
                  height: 70,
                  child: RawMaterialButton(
                    onPressed: () {},
                    shape: const CircleBorder(),
                    fillColor: color.AppColors.backgroundColor1st,
                    splashColor: Colors.blueAccent.withOpacity(0.1), // สีคลื่นน้ำกระจาย
                    highlightColor: Colors.blueAccent.withOpacity(0.1), 
                    child: const Icon(Icons.add, color: Colors.white, size: 36),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
