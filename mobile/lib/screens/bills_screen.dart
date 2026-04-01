import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../services/api_service.dart';

class BillsScreen extends StatefulWidget {
  const BillsScreen({super.key});

  @override
  State<BillsScreen> createState() => _BillsScreenState();
}

class _BillsScreenState extends State<BillsScreen> {
  List<dynamic> bills = [];
  bool loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final data = await context.read<ApiService>().getBills();
      setState(() { bills = data; loading = false; });
    } catch (_) {
      setState(() => loading = false);
    }
  }

  Color _statusColor(String status) {
    switch (status) {
      case 'paid': return Colors.green;
      case 'pending': return Colors.orange;
      case 'disputed': return Colors.red;
      default: return Colors.grey;
    }
  }

  @override
  Widget build(BuildContext context) {
    if (loading) return const Center(child: CircularProgressIndicator());

    if (bills.isEmpty) {
      return const Center(child: Text('No bills yet', style: TextStyle(color: Colors.grey)));
    }

    return RefreshIndicator(
      onRefresh: _load,
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: bills.length,
        itemBuilder: (ctx, i) {
          final bill = bills[i];
          return Card(
            margin: const EdgeInsets.only(bottom: 12),
            child: ListTile(
              title: Text(bill['bill_id'] ?? '', style: const TextStyle(fontFamily: 'monospace', fontWeight: FontWeight.bold)),
              subtitle: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('${bill['units_used']?.toStringAsFixed(1) ?? '0'} kWh'),
                  Text('${bill['period_start']?.toString().substring(0, 10) ?? ''} to ${bill['period_end']?.toString().substring(0, 10) ?? ''}'),
                ],
              ),
              trailing: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text('INR ${bill['amount']?.toStringAsFixed(2) ?? '0'}',
                    style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                    decoration: BoxDecoration(
                      color: _statusColor(bill['status'] ?? '').withOpacity(0.1),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Text(bill['status'] ?? '',
                      style: TextStyle(color: _statusColor(bill['status'] ?? ''), fontSize: 12)),
                  ),
                ],
              ),
              isThreeLine: true,
            ),
          );
        },
      ),
    );
  }
}
