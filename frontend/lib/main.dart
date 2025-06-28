import 'package:flutter/material.dart';
import 'package:frontend/screens/splash/splash_screen.dart';
import 'package:get/get.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return const GetMaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'PillMate Demo',
      home: SplashScreen(), //หน้าโหลด
    );
    //ใช้ GetMaterail เพราะ ใช้ router get
  }
}
