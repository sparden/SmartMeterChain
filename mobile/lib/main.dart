import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:google_fonts/google_fonts.dart';
import 'blocs/auth_bloc.dart';
import 'services/api_service.dart';
import 'screens/login_screen.dart';
import 'screens/home_screen.dart';

void main() {
  runApp(const SmartMeterChainApp());
}

class SmartMeterChainApp extends StatelessWidget {
  const SmartMeterChainApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiRepositoryProvider(
      providers: [
        RepositoryProvider(create: (_) => ApiService()),
      ],
      child: BlocProvider(
        create: (ctx) => AuthBloc(ctx.read<ApiService>()),
        child: MaterialApp(
          title: 'SmartMeterChain',
          debugShowCheckedModeBanner: false,
          theme: ThemeData(
            colorScheme: ColorScheme.fromSeed(
              seedColor: const Color(0xFF2563EB),
              brightness: Brightness.light,
            ),
            textTheme: GoogleFonts.interTextTheme(),
            useMaterial3: true,
          ),
          home: BlocBuilder<AuthBloc, AuthState>(
            builder: (context, state) {
              if (state is Authenticated) return const HomeScreen();
              return const LoginScreen();
            },
          ),
        ),
      ),
    );
  }
}
