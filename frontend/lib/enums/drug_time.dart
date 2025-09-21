enum DrugTime {
  beforeMeal("ก่อนอาหาร"),
  afterMeal("หลังอาหาร"),
  withMeal("พร้อมอาหาร"),
  beforeSleep("ก่อนนอน");

  final String label ;

  const DrugTime(this.label);
}