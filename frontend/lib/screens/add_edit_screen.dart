import 'package:flutter/material.dart';
import 'package:frontend/providers/add_edit_provider.dart';
import 'package:provider/provider.dart';
import 'package:frontend/utils/colors.dart' as color;

class AddEditScreen extends StatelessWidget {
  const AddEditScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (context) => AddEditProvider(),
      child: _AddEditView(),
    );
  }
}

// UnderlineInputBorder _borderUnderLine(Color color) {
//   return UnderlineInputBorder(borderSide: BorderSide(width: 1, color: color));
// }

Text _headText(String text) {
  return Text(
    text,
    style: const TextStyle(
      color: Colors.black,
      fontWeight: FontWeight.normal,
      fontSize: 24,
      letterSpacing: 0,
    ),
  );
}

Container _container(double width, double height, List<Widget> children) {
  return Container(
    width: width,
    height: height,
    decoration: BoxDecoration(
      color: Colors.white,
      border: Border.all(color: Colors.black, width: 1),
      borderRadius: BorderRadius.circular(15),
    ),
    child: Padding(
      padding: const EdgeInsets.all(10.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: children,
      ),
    ),
  );
}

class _AddEditView extends StatefulWidget {
  @override
  State<_AddEditView> createState() => _AddEditViewState();
}

class _AddEditViewState extends State<_AddEditView> {
  @override
  Widget build(BuildContext context) {
    final args =
        ModalRoute.of(context)!.settings.arguments as Map<String, dynamic>;
    final String pageFrom = args['pageType'];
    
    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor2nd,
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        foregroundColor: Colors.white,
        title: Text(
          pageFrom == "add" ? "เพิ่มรายการยา" : "แก้ไขรายการยา",
          style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 25),
        ),
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsetsGeometry.all(12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              _container(384, 125, [
                _headText("ชื่อยา"),
                TextFormField(
                  decoration: const InputDecoration(
                    hint: Text("ใส่ชื่อยาหรือยี่ห้อยา"),
                    enabledBorder: UnderlineInputBorder(
                      borderSide: BorderSide(width: 1, color: Colors.black),
                    ),
                    focusedBorder: UnderlineInputBorder(
                      borderSide: BorderSide(width: 1, color: Colors.black),
                    ),
                    contentPadding: EdgeInsets.only(bottom: 0),
                  ),
                ),
              ]),
              const SizedBox(height: 15),
              _container(384, 125, [
                _headText("สรรพคุณ"),
                TextFormField(
                  decoration: const InputDecoration(
                    hint: Text("ใส่สรรพคุณยา, รักษาอาการ "),
                    enabledBorder: UnderlineInputBorder(
                      borderSide: BorderSide(width: 1, color: Colors.black),
                    ),
                    focusedBorder: UnderlineInputBorder(
                      borderSide: BorderSide(width: 1, color: Colors.black),
                    ),
                    contentPadding: EdgeInsets.all(0),
                  ),
                ),
              ]),
              const SizedBox(height: 15),
              _container(384, 191, [_headText("ลักษณะของยา"),List]),
              const SizedBox(height: 15),
              _container(384, 203, [_headText("การใช้ยา")]),
              const SizedBox(height: 15),
              _container(384, 223, [_headText("ช่วงเวลาใช้ยา")]),
              const SizedBox(height: 15),
              ElevatedButton(
                onPressed: () {},
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF94B4C1),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(7),
                  ),
                ),
                child: const Text(
                  "บันทึก",
                  style: TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.normal,
                    fontSize: 24,
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
