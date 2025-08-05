import 'package:get/get.dart';

import '../modules/home/bindings/home_binding.dart';
import '../modules/home/drug/bindings/drug_binding.dart';
import '../modules/home/drug/views/drug_view.dart';
import '../modules/home/notification/bindings/notification_binding.dart';
import '../modules/home/notification/views/notification_view.dart';
import '../modules/home/profile/bindings/profile_binding.dart';
import '../modules/home/profile/views/profile_view.dart';
import '../modules/home/today/bindings/today_binding.dart';
import '../modules/home/today/views/today_view.dart';
import '../modules/home/views/home_view.dart';
import '../modules/login/bindings/login_binding.dart';
import '../modules/login/forget_password/bindings/forget_password_binding.dart';
import '../modules/login/forget_password/new_password/bindings/new_password_binding.dart';
import '../modules/login/forget_password/new_password/views/new_password_view.dart';
import '../modules/login/forget_password/views/forget_password_view.dart';
import '../modules/login/views/login_view.dart';
import '../modules/otp/bindings/otp_binding.dart';
import '../modules/otp/views/otp_view.dart';
import '../modules/profile_setup/bindings/profile_setup_binding.dart';
import '../modules/profile_setup/views/profile_setup_view.dart';
import '../modules/register/bindings/register_binding.dart';
import '../modules/register/views/register_view.dart';
import '../modules/splash/bindings/splash_binding.dart';
import '../modules/splash/views/splash_view.dart';

part 'app_routes.dart';

class AppPages {
  AppPages._();

  static const INITIAL = Routes.SPLASH;

  static final routes = [
    GetPage(
      name: _Paths.HOME,
      page: () => const HomeScreen(),
      binding: HomeBinding(),
      children: [
        GetPage(
          name: _Paths.PROFILE,
          page: () => const ProfileScreen(),
          binding: ProfileBinding(),
        ),
        GetPage(
          name: _Paths.TODAY,
          page: () => const TodayScreen(),
          binding: TodayBinding(),
        ),
        GetPage(
          name: _Paths.NOTIFICATION,
          page: () => const NotificationScreen(),
          binding: NotificationBinding(),
        ),
        GetPage(
          name: _Paths.DRUG,
          page: () => const DrugScreen(),
          binding: DrugBinding(),
        ),
      ],
    ),
    GetPage(
      name: _Paths.LOGIN,
      page: () => const LoginScreen(),
      binding: LoginBinding(),
      children: [
        GetPage(
          name: _Paths.FORGET_PASSWORD,
          page: () => const ForgetPasswordView(),
          binding: ForgetPasswordBinding(),
          children: [
            GetPage(
              name: _Paths.NEW_PASSWORD,
              page: () => const NewPasswordView(),
              binding: NewPasswordBinding(),
            ),
          ],
        ),
      ],
    ),
    GetPage(
      name: _Paths.REGISTER,
      page: () => const RegisterScreen(),
      binding: RegisterBinding(),
    ),
    GetPage(
      name: _Paths.SPLASH,
      page: () => const SplashScreen(),
      binding: SplashBinding(),
    ),
    GetPage(
      name: _Paths.OTP,
      page: () => const OtpView(),
      binding: OtpBinding(),
    ),
    GetPage(
      name: _Paths.PROFILE_SETUP,
      page: () => const ProfileSetupView(),
      binding: ProfileSetupBinding(),
    ),
  ];
}
