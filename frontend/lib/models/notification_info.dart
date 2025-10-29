class NotificationInfo {
  final int? id;
  final int? myMedicineId;
  final int? groupId;
  final int? notiFormatId;
  final String? notiFormatName; // เช่น "ทานต่อเนื่อง/พักยา (Cycle)"
  final String? type; // Fixed, Interval, DailyWeekly, Cycle
  final String? startDate;
  final String? endDate;
  final List<String>? times;
  final int? intervalHours;
  final int? intervalTake; // times_per_day
  final int? daysGap; // interval_day
  final int? takeDays;
  final int? breakDays;

  NotificationInfo({
    this.id,
    this.myMedicineId,
    this.groupId,
    this.notiFormatId,
    this.notiFormatName,
    this.type,
    this.startDate,
    this.endDate,
    this.times,
    this.intervalHours,
    this.intervalTake,
    this.daysGap,
    this.takeDays,
    this.breakDays,
  });

  factory NotificationInfo.fromJson(Map<String, dynamic> json) {
    // 🧩 แปลง cycle_pattern เช่น [1,3]
    int? takeDays;
    int? breakDays;
    if (json["cycle_pattern"] != null && json["cycle_pattern"] is List) {
      final list = json["cycle_pattern"] as List;
      if (list.length >= 2) {
        takeDays = list[0];
        breakDays = list[1];
      }
    }

    // 🧠 หาประเภท (type) จากชื่อ format หรือ id
    final formatName = json["NotiFormatName"] ?? json["noti_format_name"];
    String? type;
    if (formatName != null) {
      final name = formatName.toString();
      if (name.contains("เวลา") || name.contains("Fixed"))
        type = "Fixed";
      else if (name.contains("ชั่วโมง") || name.contains("Interval"))
        type = "Interval";
      else if (name.contains("วัน") || name.contains("Every"))
        type = "DailyWeekly";
      else if (name.contains("พักยา") || name.contains("Cycle"))
        type = "Cycle";
    } else if (json["noti_format_id"] != null) {
      switch (json["noti_format_id"]) {
        case 1:
          type = "Fixed";
          break;
        case 2:
          type = "Interval";
          break;
        case 3:
          type = "DailyWeekly";
          break;
        case 4:
          type = "Cycle";
          break;
      }
    }

    return NotificationInfo(
      id: json["id"],
      myMedicineId: json["my_medicine_id"],
      groupId: json["group_id"],
      notiFormatId: json["noti_format_id"],
      notiFormatName: formatName,
      type: type,
      startDate: json["start_date"],
      endDate: json["end_date"],
      times: json["times"] != null ? List<String>.from(json["times"]) : [],
      intervalHours: json["interval_hours"],
      intervalTake: json["times_per_day"],
      daysGap: json["interval_day"],
      takeDays: takeDays,
      breakDays: breakDays,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      "id": id,
      "my_medicine_id": myMedicineId,
      "group_id": groupId,
      "noti_format_id": notiFormatId,
      "noti_format_name": notiFormatName,
      "start_date": startDate,
      "end_date": endDate,
      "times": times,
      "interval_hours": intervalHours,
      "times_per_day": intervalTake,
      "interval_day": daysGap,
      "cycle_pattern": [takeDays, breakDays],
    };
  }
}
class NotiFormatModel {
  final int id;
  final String formatName;

  NotiFormatModel({required this.id, required this.formatName});

  factory NotiFormatModel.fromJson(Map<String, dynamic> json) {
    return NotiFormatModel(
      id: json["id"],
      formatName: json["format_name"] ?? "-",
    );
  }
}
