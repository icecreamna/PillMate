class DrugTimeModel {
  final int id;
  final String name;

  DrugTimeModel({
    required this.id,
    required this.name,
  });

  factory DrugTimeModel.fromJson(Map<String, dynamic> json) {
    return DrugTimeModel(
      id: json['id'],
      name: json['instruction_name'],
    );
  }
}
