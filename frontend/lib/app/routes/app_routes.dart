part of 'app_pages.dart';
// DO NOT EDIT. This is code generated via package:get_cli/get_cli.dart

abstract class Routes {
  Routes._();
  static const HOME = _Paths.HOME;
  static const LOGIN = _Paths.LOGIN;
  static const REGISTER = _Paths.REGISTER;
  static const SPLASH = _Paths.SPLASH;
  static const PROFILE = _Paths.HOME + _Paths.PROFILE;
  static const TODAY = _Paths.HOME + _Paths.TODAY;
  static const NOTIFICATION = _Paths.HOME + _Paths.NOTIFICATION;
  static const DRUG = _Paths.HOME + _Paths.DRUG;
  static const FORGET_PASSWORD = _Paths.LOGIN + _Paths.FORGET_PASSWORD;
  static const OTP = _Paths.OTP;
  static const PROFILE_SETUP = _Paths.PROFILE_SETUP;
  static const NEW_PASSWORD =
      _Paths.LOGIN + _Paths.FORGET_PASSWORD + _Paths.NEW_PASSWORD;
}

abstract class _Paths {
  _Paths._();
  static const HOME = '/home';
  static const LOGIN = '/login';
  static const REGISTER = '/register';
  static const SPLASH = '/splash';
  static const PROFILE = '/profile';
  static const TODAY = '/today';
  static const NOTIFICATION = '/notification';
  static const DRUG = '/drug';
  static const FORGET_PASSWORD = '/forget-password';
  static const OTP = '/otp';
  static const PROFILE_SETUP = '/profile-setup';
  static const NEW_PASSWORD = '/new-password';
}
