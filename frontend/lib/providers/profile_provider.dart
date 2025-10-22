import 'package:flutter/material.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:frontend/services/profile_service.dart';
import 'package:intl/intl.dart';

class InfoUser {
  final String firstName;
  final String lastName;
  final String idCard;
  final String tel;
  final String? patientCode;

  InfoUser({
    required this.firstName,
    required this.lastName,
    required this.idCard,
    required this.tel,
    this.patientCode,
  });
}

class InfoAppoinment {
  final DateTime dateTime;
  final String note;
  final DateTime hourMinute;

  InfoAppoinment({
    required this.dateTime,
    required this.note,
    required this.hourMinute,
  });
}

class ProfileProvider extends ChangeNotifier {
  final AuthService authService;
  final ProfileService profileService;
  ProfileProvider({required this.authService, required this.profileService}) {
    loadProfile();
    loadAppointment();
  }

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  InfoUser? _user;
  InfoUser? get user => _user;

  InfoAppoinment? _appoinment;

  InfoAppoinment? get appointment => _appoinment;

  String get appointmentDay {
    if (_appoinment == null) return "-";
    final thDate = DateFormat(
      "d MMMM yyyy",
      "th_TH",
    ).format(_appoinment!.dateTime);
    final buddhistYear = _appoinment!.dateTime.year + 543;
    return thDate.replaceAll('${_appoinment!.dateTime.year}', '$buddhistYear');
  }

  String get appointmentHourMinute {
    if (_appoinment == null) return "-";
    return DateFormat("HH:mm").format(_appoinment!.hourMinute);
  }

  Future<void> loadProfile() async {
    _isLoading = true;
    notifyListeners();

    try {
      final data = await profileService.fetchProfile();
      if (data != null) {
        _user = InfoUser(
          firstName: data["first_name"] ?? "",
          lastName: data["last_name"] ?? "",
          idCard: data["id_card_number"] ?? "",
          tel: data["phone_number"] ?? "",
          patientCode: data["patient_code"] ?? "ไม่มีรหัสผู้ป่วย",
        );
        debugPrint("✅ Profile loaded: ${_user!.firstName} ${_user!.lastName}");
      } else {
        debugPrint("⚠️ Profile is null (no data returned)");
      }
    } catch (e) {
      throw Exception("provider cannot load profile info ${e.toString()}");
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> loadAppointment() async {
    _isLoading = true;
    notifyListeners();
    try {
      final data = await profileService.fetchAppointment();
      if (data != null) {
        final date = DateTime.parse(data["appointment_date"]);
        final time = DateFormat("HH:mm").parse(data["appointment_time"]);

        _appoinment = InfoAppoinment(
          dateTime: date,
          note: data['note'],
          hourMinute: time,
        );
        debugPrint("✅ Appointment loaded");
      } else {
        debugPrint("✅ Appointment failed: (no data)");
      }
    } catch (e) {
      debugPrint("Provider fetchappoinment cache $e");
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> updatedInfoProfile(
    InfoUser newInfo,
    BuildContext context,
  ) async {
    _isLoading = true;
    notifyListeners();

    try {
      final success = await profileService.updateInfoProfile(
        idCard: newInfo.idCard,
        firstName: newInfo.firstName,
        lastName: newInfo.lastName,
        tel: newInfo.tel,
      );
      if (success) {
        _user = newInfo;
        notifyListeners();
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text("✅ อัปเดตข้อมูลสำเร็จ"),
            backgroundColor: Colors.green,
          ),
        );
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text("❌ อัปเดตข้อมูลไม่สำเร็จ"),
            backgroundColor: Colors.red,
          ),
        );
      }
    } catch (e) {
      debugPrint("provider catch can't updated $e");
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text("⚠️ เกิดข้อผิดพลาดระหว่างอัปเดต"),
          backgroundColor: Colors.orange,
        ),
      );
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<bool> logout() async {
    _isLoading = true;
    notifyListeners();

    try {
      final success = await authService.logout();
      return success;
    } catch (e) {
      debugPrint("❌ LogoutProvider error: $e");
      return false;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
