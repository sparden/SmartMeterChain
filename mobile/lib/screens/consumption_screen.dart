import 'package:flutter/material.dart';

class ConsumptionScreen extends StatelessWidget {
  const ConsumptionScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return const Center(
      child: Padding(
        padding: EdgeInsets.all(24),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.show_chart, size: 64, color: Colors.blue),
            SizedBox(height: 16),
            Text('Consumption Analytics',
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
            SizedBox(height: 8),
            Text('View your daily and monthly electricity consumption patterns.',
              textAlign: TextAlign.center,
              style: TextStyle(color: Colors.grey)),
            SizedBox(height: 24),
            Text('Start the backend server and simulator\nto see live consumption data.',
              textAlign: TextAlign.center,
              style: TextStyle(color: Colors.grey, fontSize: 12)),
          ],
        ),
      ),
    );
  }
}
