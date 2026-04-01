import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../blocs/auth_bloc.dart';
import '../services/api_service.dart';
import 'bills_screen.dart';
import 'consumption_screen.dart';
import 'disputes_screen.dart';
import 'verify_screen.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _currentIndex = 0;

  final _screens = const [
    _DashboardTab(),
    BillsScreen(),
    ConsumptionScreen(),
    DisputesScreen(),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('SmartMeterChain'),
        actions: [
          IconButton(
            icon: const Icon(Icons.verified),
            tooltip: 'Verify Bill',
            onPressed: () => Navigator.push(context,
              MaterialPageRoute(builder: (_) => const VerifyScreen())),
          ),
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () => context.read<AuthBloc>().add(LogoutRequested()),
          ),
        ],
      ),
      body: _screens[_currentIndex],
      bottomNavigationBar: NavigationBar(
        selectedIndex: _currentIndex,
        onDestinationSelected: (i) => setState(() => _currentIndex = i),
        destinations: const [
          NavigationDestination(icon: Icon(Icons.dashboard), label: 'Home'),
          NavigationDestination(icon: Icon(Icons.receipt_long), label: 'Bills'),
          NavigationDestination(icon: Icon(Icons.show_chart), label: 'Usage'),
          NavigationDestination(icon: Icon(Icons.gavel), label: 'Disputes'),
        ],
      ),
    );
  }
}

class _DashboardTab extends StatefulWidget {
  const _DashboardTab();

  @override
  State<_DashboardTab> createState() => _DashboardTabState();
}

class _DashboardTabState extends State<_DashboardTab> {
  Map<String, dynamic>? stats;
  bool loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final data = await context.read<ApiService>().getStats();
      setState(() { stats = data; loading = false; });
    } catch (_) {
      setState(() => loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (loading) return const Center(child: CircularProgressIndicator());

    return RefreshIndicator(
      onRefresh: _load,
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          _StatCard('Active Meters', '${stats?['meters']?['active'] ?? 0}', Icons.electric_meter, Colors.blue),
          _StatCard('Today\'s Readings', '${stats?['readings']?['today'] ?? 0}', Icons.data_usage, Colors.green),
          _StatCard('Open Disputes', '${stats?['disputes']?['open'] ?? 0}', Icons.warning, Colors.orange),
          _StatCard('Tamper Alerts', '${stats?['alerts']?['unacknowledged'] ?? 0}', Icons.security, Colors.red),
          const SizedBox(height: 16),
          Card(
            child: ListTile(
              leading: const Icon(Icons.link, color: Color(0xFF2563EB)),
              title: const Text('Blockchain Secured'),
              subtitle: const Text('All data is immutably recorded on Hyperledger Fabric'),
            ),
          ),
        ],
      ),
    );
  }
}

class _StatCard extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;
  final Color color;

  const _StatCard(this.title, this.value, this.icon, this.color);

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: color.withOpacity(0.1),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, color: color, size: 28),
            ),
            const SizedBox(width: 16),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(value, style: Theme.of(context).textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.bold)),
                Text(title, style: Theme.of(context).textTheme.bodySmall?.copyWith(color: Colors.grey)),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
