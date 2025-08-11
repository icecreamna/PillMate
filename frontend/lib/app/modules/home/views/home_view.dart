import 'package:flutter/material.dart';
import 'package:frontend/app/modules/home/today/views/today_view.dart';
import 'package:frontend/app/modules/home/drug/views/drug_view.dart';
import 'package:frontend/app/modules/home/notification/views/notification_view.dart';
import 'package:frontend/app/modules/home/profile/views/profile_view.dart';
import 'package:frontend/app/utils/colors.dart' as color;

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _selectedIndex = 0;

  static final List<Widget> _screenOptions = const [
    TodayScreen(),
    DrugScreen(),
    NotificationScreen(),
    ProfileScreen(),
  ];

  void onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor2nd,
      appBar: AppBar(
        backgroundColor: color.AppColors.backgroundColor1st,
        title: Text(
          _selectedIndex == 0
              ? "ตารางกินยา"
              : _selectedIndex == 1
              ? "ยาของฉัน"
              : _selectedIndex == 2
              ? "การแจ้งเตือน"
              : "ข้อมูลผู้ใช้",
          style: const TextStyle(
            color: Colors.white,
            fontSize: 32,
            fontWeight: FontWeight.bold,
          ),
          textAlign: TextAlign.center,
        ),
      ),
      body: _screenOptions.elementAt(_selectedIndex),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _selectedIndex,
        onTap: onItemTapped,
        type: BottomNavigationBarType.fixed,
        selectedItemColor: color.AppColors.buttonColor,
        unselectedItemColor: const Color(0xFF454545),
        showUnselectedLabels: true,
        items: const <BottomNavigationBarItem>[
          BottomNavigationBarItem(
            label: "วันนี้",
            icon: Icon(Icons.calendar_month_outlined),
          ),
          BottomNavigationBarItem(
            label: "ยา",
            icon: Icon(Icons.medication_outlined),
          ),
          BottomNavigationBarItem(
            label: "แจ้งเตือน",
            icon: Icon(Icons.notifications_none_sharp),
          ),
          BottomNavigationBarItem(
            label: "ผู้ใช้",
            icon: Icon(Icons.person_outline),
          ),
        ],
      ),
    );
  }
}