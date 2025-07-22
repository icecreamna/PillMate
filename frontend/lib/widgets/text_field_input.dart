import 'package:flutter/material.dart';
import 'package:frontend/utils/colors.dart' as color;

class TextFieldInput extends StatelessWidget {
  final Widget? preIcon;
  final Widget? sufIcon;
  final String labelname;
  final bool isSuf;
  final bool hideText;
  final TextInputType textInputType;
  final FocusNode? focusNode;

  const TextFieldInput({
    super.key,
    required this.labelname,
    required this.preIcon,
    this.sufIcon,
    required this.textInputType,
    this.isSuf = false,
    this.hideText = false,
    this.focusNode,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 30),
      child: TextField(
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
