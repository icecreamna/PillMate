import 'package:flutter/material.dart';
import 'package:frontend/providers/drug_provider.dart';

class AddSingleNotificationProvider extends ChangeNotifier{

late DoseTest _tempDose ; 

DoseTest get tempDose => _tempDose;

AddSingleNotificationProvider(DoseTest dose){
  _tempDose = dose;
}

void updatedTempDose(DoseTest newDose){ 
  _tempDose = newDose;
  notifyListeners();
}

}