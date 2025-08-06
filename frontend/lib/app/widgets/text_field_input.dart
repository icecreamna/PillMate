import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:frontend/app/utils/colors.dart' as color;

class TextFieldInput extends StatelessWidget {
  final Widget? preIcon;
  final Widget? sufIcon;
  final String labelname;
  final bool isSuf;
  final bool hideText;
  final TextInputType textInputType;
  final List<TextInputFormatter>? inputFormatters;

  const TextFieldInput({
    super.key,
    required this.labelname,
    this.preIcon,
    this.sufIcon,
    required this.textInputType,
    this.isSuf = false,
    this.hideText = false,
    this.inputFormatters,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 30),
      child: TextField(
        inputFormatters: inputFormatters,
        keyboardType: textInputType,
        obscureText: hideText,
        decoration: InputDecoration(
          prefixIcon: preIcon,
          suffixIcon: isSuf ? sufIcon : null,
          labelText: labelname,
          labelStyle: const TextStyle(color: Colors.grey),
          enabledBorder: UnderlineInputBorder(
            borderSide: BorderSide(width: 1.5,color: color.AppColors.buttonColor),
          ),
          focusedBorder: UnderlineInputBorder(
            borderSide: BorderSide(width: 1.5,color: color.AppColors.buttonColor),
          ),
        ),
      ),
    );
  }
}