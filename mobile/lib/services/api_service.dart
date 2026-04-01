import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiService {
  static const String baseUrl = 'http://10.0.2.2:3000/api/v1'; // Android emulator -> localhost
  final Dio _dio;
  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  ApiService() : _dio = Dio(BaseOptions(
    baseUrl: baseUrl,
    connectTimeout: const Duration(seconds: 10),
    receiveTimeout: const Duration(seconds: 10),
  )) {
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        final token = await _storage.read(key: 'token');
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        handler.next(options);
      },
    ));
  }

  Future<Map<String, dynamic>> login(String username, String password) async {
    final res = await _dio.post('/auth/login', data: {
      'username': username,
      'password': password,
    });
    final data = res.data['data'];
    await _storage.write(key: 'token', value: data['token']);
    return data;
  }

  Future<void> logout() async {
    await _storage.delete(key: 'token');
  }

  Future<Map<String, dynamic>> getProfile() async {
    final res = await _dio.get('/profile');
    return res.data['data'];
  }

  Future<List<dynamic>> getBills({int page = 1}) async {
    final res = await _dio.get('/bills', queryParameters: {'page': page});
    return res.data['data']['items'] ?? [];
  }

  Future<Map<String, dynamic>> verifyBill(String billId) async {
    final res = await _dio.get('/verify/bill/$billId');
    return res.data['data'];
  }

  Future<List<dynamic>> getReadings(String meterId, {int page = 1}) async {
    final res = await _dio.get('/meters/$meterId/readings', queryParameters: {'page': page});
    return res.data['data']['items'] ?? [];
  }

  Future<Map<String, dynamic>> getStats() async {
    final res = await _dio.get('/dashboard/stats');
    return res.data['data'];
  }

  Future<List<dynamic>> getDisputes({int page = 1}) async {
    final res = await _dio.get('/disputes', queryParameters: {'page': page});
    return res.data['data']['items'] ?? [];
  }

  Future<void> fileDispute(String billId, String reason) async {
    await _dio.post('/disputes', data: {
      'bill_id': billId,
      'reason': reason,
    });
  }
}
