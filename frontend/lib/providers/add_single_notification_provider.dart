import 'package:flutter/material.dart';
import '../models/dose.dart';

class AddSingleNotificationProvider extends ChangeNotifier{

late Dose _tempDose ; 

Dose get tempDose => _tempDose;

AddSingleNotificationProvider(Dose dose){
  _tempDose = dose;
}

void updatedTempDose(Dose newDose){ 
  _tempDose = newDose;
  notifyListeners();
}

}