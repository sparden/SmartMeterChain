import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../services/api_service.dart';

class DisputesScreen extends StatefulWidget {
  const DisputesScreen({super.key});

  @override
  State<DisputesScreen> createState() => _DisputesScreenState();
}

class _DisputesScreenState extends State<DisputesScreen> {
  List<dynamic> disputes = [];
  bool loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final data = await context.read<ApiService>().getDisputes();
      setState(() { disputes = data; loading = false; });
    } catch (_) {
      setState(() => loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (loading) return const Center(child: CircularProgressIndicator());

    if (disputes.isEmpty) {
      return const Center(child: Text('No disputes filed', style: TextStyle(color: Colors.grey)));
    }

    return RefreshIndicator(
      onRefresh: _load,
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: disputes.length,
        itemBuilder: (ctx, i) {
          final d = disputes[i];
          final isOpen = d['status'] == 'open';
          return Card(
            margin: const EdgeInsets.only(bottom: 12),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(d['dispute_id'] ?? '',
                        style: const TextStyle(fontFamily: 'monospace', fontWeight: FontWeight.bold)),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          color: isOpen ? Colors.orange.withOpacity(0.1) : Colors.green.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Text(d['status'] ?? '',
                          style: TextStyle(color: isOpen ? Colors.orange : Colors.green, fontSize: 12)),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Text('Bill: ${d['bill_id'] ?? ''}', style: const TextStyle(color: Colors.grey)),
                  const SizedBox(height: 4),
                  Text(d['reason'] ?? ''),
                  if (d['resolution'] != null) ...[
                    const SizedBox(height: 8),
                    Text('Resolution: ${d['resolution']}',
                      style: const TextStyle(color: Colors.green)),
                  ],
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}
