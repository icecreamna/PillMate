import 'dart:convert';
import 'package:frontend/models/drug_form.dart';
import 'package:frontend/models/drug_time.dart';
import 'package:http/http.dart' as http;

class InitialDataService {
  static const baseUrl = "http://10.0.2.2:8080";

  Future<List<DrugFormModel>> fetchDrugForms() async {
    final res = await http.get(Uri.parse("$baseUrl/forms?with_relations=true"));
    if (res.statusCode == 200) {
      final List data = jsonDecode(res.body);
      return data.map((e) => DrugFormModel.fromJson(e)).toList();
    } else {
      throw Exception("Failed to load forms: ${res.body}");
    }
  }

  Future<List<DrugTimeModel>> fetchDrugTimes() async {
    final res = await http.get(Uri.parse("$baseUrl/instructions"));
    if (res.statusCode == 200) {
      final List data = jsonDecode(res.body);
      return data.map((e) => DrugTimeModel.fromJson(e)).toList();
    } else {
      throw Exception("Failed to load instructions: ${res.body}");
    }
  }
}
