import 'package:flutter_bloc/flutter_bloc.dart';
import '../services/api_service.dart';

// Events
abstract class AuthEvent {}

class LoginRequested extends AuthEvent {
  final String username;
  final String password;
  LoginRequested(this.username, this.password);
}

class LogoutRequested extends AuthEvent {}

class CheckAuth extends AuthEvent {}

// States
abstract class AuthState {}

class AuthInitial extends AuthState {}

class AuthLoading extends AuthState {}

class Authenticated extends AuthState {
  final Map<String, dynamic> user;
  Authenticated(this.user);
}

class AuthError extends AuthState {
  final String message;
  AuthError(this.message);
}

class Unauthenticated extends AuthState {}

// Bloc
class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final ApiService _api;

  AuthBloc(this._api) : super(AuthInitial()) {
    on<LoginRequested>(_onLogin);
    on<LogoutRequested>(_onLogout);
    on<CheckAuth>(_onCheckAuth);
  }

  Future<void> _onLogin(LoginRequested event, Emitter<AuthState> emit) async {
    emit(AuthLoading());
    try {
      final user = await _api.login(event.username, event.password);
      emit(Authenticated(user));
    } catch (e) {
      emit(AuthError('Invalid credentials'));
    }
  }

  Future<void> _onLogout(LogoutRequested event, Emitter<AuthState> emit) async {
    await _api.logout();
    emit(Unauthenticated());
  }

  Future<void> _onCheckAuth(CheckAuth event, Emitter<AuthState> emit) async {
    try {
      final profile = await _api.getProfile();
      emit(Authenticated(profile));
    } catch (_) {
      emit(Unauthenticated());
    }
  }
}
