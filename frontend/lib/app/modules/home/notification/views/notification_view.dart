import 'package:flutter/material.dart';

class NotificationScreen extends StatelessWidget {
  const NotificationScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return  Container(
      child: const Center(
        child: Text("ไม่มีการแจ้งเตือน",style: TextStyle(
          fontSize: 20
        ),),
      ),
    );
  }
}