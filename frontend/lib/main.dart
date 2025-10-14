import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:android_alarm_manager_plus/android_alarm_manager_plus.dart';
import 'package:frontend/services/alrm_service.dart';
import 'package:frontend/screens/splash_screen.dart';
import 'package:frontend/services/permission_service.dart';
import 'package:permission_handler/permission_handler.dart';

@pragma('vm:entry-point')
Future<void> backgroundTask() async {
  WidgetsFlutterBinding.ensureInitialized();
  print("🚀 backgroundTask RUNNING at ${DateTime.now()}");
  await AlarmService.init();
  await AlarmService.checkDueNow(); // จะใช้ http ต้อง ensureInitialized ก่อน
}

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // ✅ ขอสิทธิ์ก่อนตั้ง alarm
  await PermissionService.requestNotificationPermission();

  final exactPerm = await Permission.scheduleExactAlarm.request();
  if (!exactPerm.isGranted) {
    debugPrint("⚠️ ผู้ใช้ยังไม่อนุญาต SCHEDULE_EXACT_ALARM");
  }

  await AndroidAlarmManager.initialize();
  await AlarmService.init();

  // ✅ ตั้ง Alarm ให้ทำงานทุก 1 นาที
  await AndroidAlarmManager.periodic(
    const Duration(minutes: 1),
    123,
    backgroundTask,
    wakeup: true,
    rescheduleOnReboot: true,
    exact: true,
  );

  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: "PillMate",
      home: const SplashScreen(),
      theme: ThemeData(fontFamily: "NotoSansThai"),
      debugShowCheckedModeBanner: false,
      localizationsDelegates: const [
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: const [Locale('th', 'TH'), Locale('en', 'US')],
      locale: const Locale('th', 'TH'),
    );
  }
}
