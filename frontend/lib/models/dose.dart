class Dose {
  final String id;
  final String name;
  final String amountPerDose;
  final String frequency;
  final String instruction;
  final String picture;
  final String description;
  final String unit;
  final bool import;

  Dose({
    required this.id,
    required this.name,
    required this.description,
    required this.import,
    required this.amountPerDose,
    required this.frequency,
    required this.instruction,
    required this.unit,
    required this.picture,
  });
  Dose copyWith({
    String? id,
    String? name,
    String? description,
    bool? import,
    String? amountPerDose,
    String? frequency,
    String? instruction,
    String? picture,
    String? unit,
  }) {
    return Dose(
      id: id ?? this.id,
      name: name ?? this.name,
      description: description ?? this.description,
      import: import ?? this.import,
      amountPerDose: amountPerDose ?? this.amountPerDose,
      frequency: frequency ?? this.frequency,
      instruction: instruction ?? this.instruction,
      picture: picture ?? this.picture,
      unit: unit ?? this.unit,
    );
  }
}
