import 'package:flutter/material.dart';
// import 'package:frontend/providers/drug_provider.dart';
// import 'package:provider/provider.dart';

class GroupDrugScreen extends StatelessWidget {
  const GroupDrugScreen({super.key});

  @override
  Widget build(BuildContext context) {
    // final p = context.read<DrugProvider>();

    return Padding(
      padding: const EdgeInsets.all(5),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: 384,
            height: 125,
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: Colors.grey,width: 1)
            ),
          ),
        ],
      ),
    );
  }
}
