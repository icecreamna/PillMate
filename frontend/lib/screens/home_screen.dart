import 'package:flutter/material.dart';
import 'package:frontend/providers/drug_provider.dart';
import 'package:frontend/providers/home_provider.dart';
import 'package:frontend/providers/today_provider.dart';
import 'package:frontend/screens/drug_screen.dart';
import 'package:frontend/screens/notification_screen.dart';
import 'package:frontend/screens/profile_screen.dart';
import 'package:frontend/screens/today_screen.dart';

import 'package:frontend/utils/colors.dart' as color;
import 'package:provider/provider.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (_) => HomeProvider(),
      child: _HomeView(), // แยก View ออกมาเพื่อไม่ให้ provider ถูกสร้างซ้ำ
    );
  }
}

class _HomeView extends StatefulWidget {
  @override
  State<_HomeView> createState() => _HomeViewState();
}

class _HomeViewState extends State<_HomeView> {
  final List<Widget> _screens = [
    ChangeNotifierProvider(
      create: (_) => TodayProvider(),
      child: const TodayScreen(),
    ),
    ChangeNotifierProvider(
      create: (_) => DrugProvider(),
      child: const DrugScreen(),
    ),
    const NotificationScreen(),
    const ProfileScreen(),
  ];

  @override
  Widget build(BuildContext context) {
    final app = context.watch<HomeProvider>();
    return Scaffold(
      body: _screens[app.selectIndex],
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: app.selectIndex,
        onTap: app.setSelectIndex,
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
