import 'package:flutter/material.dart';

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
    required this.isEye,
    this.hideText = false,
  });

  @override
  Widget build(BuildContext context) {
    return TextField(
      decoration: InputDecoration(
        prefixIcon: preIcon,
        suffixIcon: isEye ? sufIcon : null,
        labelText: labelname,
      ),
    );
  }
}
