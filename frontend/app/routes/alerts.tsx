import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { api } from '~/lib/api'

export const Route = createFileRoute('/alerts')({
  component: AlertsPage,
})

function AlertsPage() {
  const [alerts, setAlerts] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getAlerts(1)
      .then((res) => setAlerts(res.data?.items || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  const severityColor = (s: string) => {
    switch (s) {
      case 'critical': return 'bg-red-100 text-red-700 border-red-200'
      case 'high': return 'bg-orange-100 text-orange-700 border-orange-200'
      case 'medium': return 'bg-yellow-100 text-yellow-700 border-yellow-200'
      default: return 'bg-blue-100 text-blue-700 border-blue-200'
    }
  }

  if (loading) return <div className="text-slate-400 p-8">Loading alerts...</div>

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-800">Tamper Alerts</h2>

      <div className="space-y-3">
        {alerts.map((alert: any) => (
          <div key={alert.ID} className={`rounded-xl border p-5 ${severityColor(alert.severity)}`}>
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-3">
                <span className="font-mono font-bold">{alert.meter_id}</span>
                <span className="px-2 py-0.5 rounded text-xs font-medium uppercase">{alert.alert_type}</span>
              </div>
              <span className="text-xs">{alert.detected_at?.slice(0, 16)}</span>
            </div>
            <p className="text-sm">{alert.description}</p>
            <div className="mt-2 text-xs opacity-75">
              Value: {alert.reading_value?.toFixed(2)} | Expected: {alert.expected_min?.toFixed(2)} — {alert.expected_max?.toFixed(2)}
            </div>
          </div>
        ))}
        {alerts.length === 0 && (
          <div className="text-center text-slate-400 py-12">No tamper alerts detected. System is healthy.</div>
        )}
      </div>
    </div>
  )
}
