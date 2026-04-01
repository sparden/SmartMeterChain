import { useState, useEffect } from 'react'
import { api } from '../lib/api'

export default function AlertsPage() {
  const [alerts, setAlerts] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [msg, setMsg] = useState({ text: '', type: '' })

  const loadAlerts = () => {
    api.getAlerts(1)
      .then((res) => setAlerts(res.data?.items || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  useEffect(() => { loadAlerts() }, [])

  const handleAcknowledge = async (id: number) => {
    try {
      await api.acknowledgeAlert(id)
      setMsg({ text: `Alert #${id} acknowledged`, type: 'success' })
      loadAlerts()
    } catch (err: any) {
      setMsg({ text: err.message, type: 'error' })
    }
  }

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
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-slate-800">Tamper Alerts</h2>
        <div className="flex items-center gap-3">
          <span className="text-sm text-slate-500">{alerts.length} alerts</span>
          <button
            onClick={loadAlerts}
            className="px-3 py-1.5 bg-slate-200 text-slate-600 rounded-lg text-xs font-medium hover:bg-slate-300"
          >
            Refresh
          </button>
        </div>
      </div>

      {msg.text && (
        <div className={`p-3 rounded-lg text-sm ${msg.type === 'success' ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'}`}>
          {msg.text}
        </div>
      )}

      <div className="space-y-3">
        {alerts.map((alert: any) => (
          <div key={alert.ID} className={`rounded-xl border p-5 ${severityColor(alert.severity)}`}>
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-3">
                <span className="font-mono font-bold">{alert.meter_id}</span>
                <span className="px-2 py-0.5 rounded text-xs font-medium uppercase">{alert.alert_type}</span>
                <span className={`px-2 py-0.5 rounded text-xs font-medium uppercase ${
                  alert.severity === 'critical' ? 'bg-red-600 text-white' : ''
                }`}>
                  {alert.severity}
                </span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-xs">{alert.detected_at?.slice(0, 16)}</span>
                {!alert.acknowledged && (
                  <button
                    onClick={() => handleAcknowledge(alert.ID)}
                    className="px-2.5 py-1 bg-white border border-current rounded text-xs font-medium hover:opacity-80"
                  >
                    Acknowledge
                  </button>
                )}
                {alert.acknowledged && (
                  <span className="px-2 py-0.5 bg-green-200 text-green-800 rounded text-xs">Acknowledged</span>
                )}
              </div>
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
