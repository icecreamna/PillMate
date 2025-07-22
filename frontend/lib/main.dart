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
    return  GetMaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'PillMate Demo',
      home: const SplashScreen(),
      theme: ThemeData(fontFamily: "Inter"), //หน้าโหลด
    );
    //ใช้ GetMaterail เพราะ ใช้ router get
  }
}
