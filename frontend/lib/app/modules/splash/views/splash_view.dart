import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:get/route_manager.dart';
import 'package:frontend/app/utils/colors.dart' as color;

import '../../../routes/app_pages.dart';

class SplashScreen extends StatelessWidget {
  const SplashScreen({super.key});
  
  @override
  Widget build(BuildContext context) {
    Future.delayed(const Duration(milliseconds: 1500), () {
      Get.offAllNamed(Routes.LOGIN);
    });

    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor1st,
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            SizedBox(
              height: 280,
              child: Stack(
                clipBehavior: Clip.none,
                alignment: Alignment.center,
                children: [
                  SvgPicture.asset(
                    "assets/images/clock.svg",
                    colorFilter: const ColorFilter.mode(
                      Colors.white,
                      BlendMode.srcIn,
                    ),
                    height: 190,
                    width: 200,
                  ),
                  Positioned(
                    bottom: -20,
                    left: -70,
                    child: Image.asset(
                      "assets/images/drugs.png",
                      height: 153,
                      width: 153,
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 120),
            const Text(
              "PillMate",
              style: TextStyle(color: Colors.white, fontSize: 48),
            ),
          ],
        ),
      ),
    );
  }
}