import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../services/api_service.dart';

class VerifyScreen extends StatefulWidget {
  const VerifyScreen({super.key});

  @override
  State<VerifyScreen> createState() => _VerifyScreenState();
}

class _VerifyScreenState extends State<VerifyScreen> {
  final _ctrl = TextEditingController();
  Map<String, dynamic>? result;
  bool loading = false;
  String? error;

  Future<void> _verify() async {
    if (_ctrl.text.trim().isEmpty) return;
    setState(() { loading = true; error = null; result = null; });

    try {
      final data = await context.read<ApiService>().verifyBill(_ctrl.text.trim());
      setState(() { result = data; loading = false; });
    } catch (e) {
      setState(() { error = e.toString(); loading = false; });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Verify Bill')),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Icon(Icons.verified_user, size: 48, color: Color(0xFF2563EB)),
            const SizedBox(height: 16),
            const Text('Blockchain Bill Verification',
              textAlign: TextAlign.center,
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            const Text('Enter a Bill ID to verify its integrity against the blockchain record.',
              textAlign: TextAlign.center,
              style: TextStyle(color: Colors.grey, fontSize: 14)),
            const SizedBox(height: 24),

            TextField(
              controller: _ctrl,
              decoration: InputDecoration(
                labelText: 'Bill ID',
                hintText: 'e.g., BILL-abc12345',
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                prefixIcon: const Icon(Icons.receipt),
              ),
            ),
            const SizedBox(height: 16),

            FilledButton(
              onPressed: loading ? null : _verify,
              child: loading
                ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                : const Text('Verify on Blockchain'),
            ),

            if (error != null) ...[
              const SizedBox(height: 16),
              Text(error!, style: const TextStyle(color: Colors.red)),
            ],

            if (result != null) ...[
              const SizedBox(height: 24),
              Card(
                color: result!['verified'] == true ? Colors.green.shade50 : Colors.red.shade50,
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    children: [
                      Icon(
                        result!['verified'] == true ? Icons.check_circle : Icons.cancel,
                        size: 48,
                        color: result!['verified'] == true ? Colors.green : Colors.red,
                      ),
                      const SizedBox(height: 8),
                      Text(
                        result!['verified'] == true ? 'Verified' : 'Verification Failed',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                          color: result!['verified'] == true ? Colors.green : Colors.red,
                        ),
                      ),
                      if (result!['bill'] != null) ...[
                        const SizedBox(height: 12),
                        Text('Amount: INR ${result!['bill']['amount']?.toStringAsFixed(2) ?? '0'}'),
                        Text('Units: ${result!['bill']['units_used']?.toStringAsFixed(2) ?? '0'} kWh'),
                        Text('TX: ${result!['bill']['tx_id'] ?? 'N/A'}',
                          style: const TextStyle(fontFamily: 'monospace', fontSize: 10)),
                      ],
                    ],
                  ),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
