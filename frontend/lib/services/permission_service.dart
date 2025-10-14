import 'package:flutter/material.dart';
import 'package:permission_handler/permission_handler.dart';

class PermissionService {
  static Future<void> requestNotificationPermission([BuildContext? context]) async {
    final status = await Permission.notification.request();

    if (status.isDenied && context != null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text("กรุณาอนุญาตการแจ้งเตือนเพื่อให้ PillMate ทำงานได้เต็มที่"),
          backgroundColor: Colors.orange,
        ),
      );
    }

    if (status.isPermanentlyDenied) {
      await openAppSettings();
    }
  }
}
