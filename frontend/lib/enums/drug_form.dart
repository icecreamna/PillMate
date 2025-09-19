enum DrugForm {
  tablet("ยาเม็ด","assets/images/pill.png",["เม็ด"]),
  capsule("แคปซูล","assets/images/capsule.png",["แคปซูล"]),
  syrup("ยาน้ำ","assets/images/syrup.png",["ช้อนชา", "ช้อนโต๊ะ", "มิลลิลิตร", "cc"]),
  ointment("ยาใช้ทา","assets/images/ointment.png",["กรัม", "มิลลิลิตร", "หลอด", "ช้อนชา"]),
  injection("ยาฉีด","assets/images/vaccine.png",["ยูนิต", "cc", "มิลลิลิตร", "มิลลิกรัม", "ไมโครกรัม"]),
  drop("ยาใช้หยด","assets/images/eye-drop 1.png",["หยด", "มิลลิลิตร", "cc"]);

  final String label;
  final String image;
  final List<String> unit;

  const DrugForm(this.label,this.image,this.unit);
}