import 'package:flutter/material.dart';

class RegisScreen extends StatelessWidget {
  const RegisScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [Image.asset("assets/images/drugs.png")],
        ),
      ),
    );
  }
}
