import 'package:flutter/material.dart';
import 'package:frontend/utils/colors.dart' as color;

class TextFieldInput extends StatelessWidget {
  final Icon preIcon;
  final Icon sufIcon;
  final String labelname;
  final bool isEye;
  final bool hideText;

  const TextFieldInput({
    super.key,
    required this.labelname,
    required this.preIcon,
    required this.sufIcon,
    this.isEye = false,
    this.hideText = false,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 30),
      child: TextField(
        decoration: InputDecoration(
          prefixIcon: preIcon,
          suffixIcon: isEye ? sufIcon : null,
          labelText: labelname,
          labelStyle: const TextStyle(color: Colors.grey),
          enabledBorder: UnderlineInputBorder(
            borderSide: BorderSide(color: color.AppColors.buttonColor),
          ),
          focusedBorder: UnderlineInputBorder(
            borderSide: BorderSide(color: color.AppColors.buttonColor)
          )
        ),
        obscureText: hideText,
      ),
    );
  }
}
