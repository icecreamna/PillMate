import 'package:flutter/material.dart';

class RegrisScreen extends StatelessWidget {
  const RegrisScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Image.asset("assets/images/drugs.png")
            ],
        ),
      ),
    );
  }
}
