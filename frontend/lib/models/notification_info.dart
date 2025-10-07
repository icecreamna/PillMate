class NotificationInfo {
  final String type; // Fixed / Interval / DailyWeekly / Cycle
  final List<String>? times; // สำหรับ Fixed
  final int? intervalHours; // สำหรับ Interval
  final int? intervalTake;
  final int? daysGap; // สำหรับ DailyWeekly
  final int? takeDays; // สำหรับ Cycle
  final int? breakDays; // สำหรับ Cycle
  final String startDate;
  final String endDate;

  NotificationInfo({
    required this.type,
    this.times,
    this.intervalHours,
    this.intervalTake,
    this.daysGap,
    this.takeDays,
    this.breakDays,
    required this.startDate,
    required this.endDate,
  });
}
