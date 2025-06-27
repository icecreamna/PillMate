import 'package:flutter/material.dart';
import 'package:frontend/utils/colors.dart' as color;
import 'package:frontend/widgets/text_field_input.dart';

class LoginScreen extends StatelessWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: color.AppColors.backgroundColor1st,
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              SizedBox(
                width: double.infinity,
                child: Image.asset(
                  "assets/images/drugs.png",
                  height: 191,
                  width: 191,
                ),
              ),
              const SizedBox(height: 15),
              Container(
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(20),
                  color: Colors.white,
                ),
                width: double.infinity,
                height: 514,
                child: const Column(
                  children: [
                    SizedBox(height: 12),
                    Text(
                      "Welcome",
                      style: TextStyle(
                        color: Color(0xFF0B0870),
                        fontWeight: FontWeight.bold,
                        fontSize: 24,
                      ),
                    ),
                    SizedBox(height: 10),
                    Text(
                      "Please log in to continue",
                      style: TextStyle(
                        color: Color(0xFF0B0870),
                        fontWeight: FontWeight.bold,
                        fontSize: 13,
                      ),
                    ),
                    TextFieldInput(
                      labelname: "User",
                      preIcon: Icon(Icons.lock_open_outlined),
                      sufIcon: Icon(Icons.visibility_off_outlined),
                      isEye: true,
                    ),
                   
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
