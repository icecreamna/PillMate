import 'package:flutter/material.dart';
// import 'package:flutter/services.dart';
import 'package:frontend/utils/colors.dart' as color ;

class FilledButtonCustom extends StatelessWidget {
  final String text;
  final void Function()? onPressed;

  const FilledButtonCustom({
    super.key,
    required this.text,
    required this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    return FilledButton(
      style: FilledButton.styleFrom(
        minimumSize: const Size(288, 49),
        backgroundColor: color.AppColors.buttonColor,
      ),
      onPressed: onPressed,
      child: Text(
        text,
        style: const TextStyle(color: Colors.white, fontSize: 24),
      ),
    );
  }
}
