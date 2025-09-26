import 'package:flutter/material.dart';
import 'package:frontend/models/dose.dart';

class AddGroupProvider extends ChangeNotifier {

final Map<String,List<Dose>> _groupSelected = {};
bool _doseSelect = false;

bool get doseSelected => _doseSelect;


void toggleCheckbox(){
  _doseSelect = !_doseSelect;
}

}
