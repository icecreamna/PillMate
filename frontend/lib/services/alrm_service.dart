import 'dart:convert';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;
import 'package:frontend/services/auth_service.dart';

class AlarmService {
  static final _noti = FlutterLocalNotificationsPlugin();
  static const baseUrl = "http://10.0.2.2:8080";

  // ✅ เริ่มระบบ Notification
  static Future<void> init() async {
    const android = AndroidInitializationSettings('@mipmap/ic_launcher');
    const settings = InitializationSettings(android: android);
    await _noti.initialize(settings);
  }

  // ✅ แสดงแจ้งเตือน
  static Future<void> showNotification({
    required String title,
    required String body,
  }) async {
    print("เรียก shownotification แล้ว");
    const androidDetails = AndroidNotificationDetails(
      'pillmate_channel',
      'PillMate Notifications',
      channelDescription: 'การแจ้งเตือนการกินยา',
      importance: Importance.max,
      priority: Priority.high,
      playSound: true,
      enableVibration: true,
    );
    const details = NotificationDetails(android: androidDetails);
    print("รอเรียก noti.show ${details}");
    await _noti.show(
      DateTime.now().millisecondsSinceEpoch ~/ 1000,
      title,
      body,
      details,
    );
  }

  // ✅ ตรวจรายการที่ถึงเวลาแจ้งเตือน
  static Future<void> checkDueNow() async {
    try {
      String? token = AuthService.jwtToken;

      if (token == null) {
        print("🚫 alarm ไม่มี token — ลองโหลดจากไฟล์");
        token = await AuthService.loadTokenFromFile();
      }

      if (token == null) {
        print("❌ ไม่มี token ในระบบ (ยังไม่ login)");
        return;
      }

      print("🌐 เริ่มเรียก due-now ที่ ${DateTime.now()}");

      final res = await http.get(
        Uri.parse("$baseUrl/api/notify/due-now?window=1"),
        headers: {"Content-Type": "application/json", "Cookie": "jwt=$token"},
      );

      print("📦 STATUS: ${res.statusCode}");
      if (res.statusCode != 200) {
        print("❌ โหลด due-now ไม่สำเร็จ: ${res.body}");
        return;
      }

      final body = jsonDecode(res.body);
      final List data = body["data"] ?? [];

      if (data.isEmpty) {
        print("😴 ไม่มีรายการแจ้งเตือนในช่วงนี้");
        return;
      }

      for (var item in data) {
        // final String medName = item["med_name"] ?? "ยาไม่ทราบชื่อ";
        final String time = item["notify_time"] ?? "";

        print("💊รอเรียก shownotification");
        await showNotification(
          title: "ถึงเวลากินยาแล้ว 💊",
          body: "กินยาเวลา $time น.",
        );

        print("เรียก noti.show ไปแล้วว");
        // ✅ mark notify_status = true
        final id = item["id"];
        if (id != null) {
          await http.patch(
            Uri.parse("$baseUrl/api/noti-items/$id/notified"),
            headers: {
              "Content-Type": "application/json",
              "Cookie": "jwt=$token",
            },
            body: jsonEncode({"notified": true}),
          );
          print("🔔 อัปเดต notified = true สำหรับ ID $id");
        }
      }
    } catch (e) {
      print("❌ Error checkDueNow: $e");
    }
  }
}
