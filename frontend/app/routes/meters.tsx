import { useState, useEffect } from 'react'
import { api } from '../lib/api'

export default function MetersPage() {
  const [meters, setMeters] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getMeters(1, 50)
      .then((res) => setMeters(res.data?.items || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="text-slate-400 p-8">Loading meters...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-slate-800">Smart Meters</h2>
        <span className="text-sm text-slate-500">{meters.length} meters registered</span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {meters.map((meter: any) => (
          <div key={meter.meter_id} className="bg-white rounded-xl shadow-sm border border-slate-200 p-5">
            <div className="flex items-center justify-between mb-3">
              <span className="font-mono font-bold text-blue-600">{meter.meter_id}</span>
              <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
                meter.status === 'active' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
              }`}>
                {meter.status}
              </span>
            </div>
            <div className="space-y-2 text-sm text-slate-600">
              <p><span className="text-slate-400">Type:</span> {meter.meter_type}</p>
              <p><span className="text-slate-400">Location:</span> {meter.location}</p>
              <p><span className="text-slate-400">Consumer:</span> {meter.consumer_id}</p>
              <p><span className="text-slate-400">Last Reading:</span> {meter.last_reading?.toFixed(2) || '—'} kWh</p>
            </div>
            {meter.on_chain && (
              <div className="mt-3 flex items-center gap-1 text-xs text-blue-600">
                <span>🔗</span> On-chain verified
              </div>
            )}
          </div>
        ))}
      </div>

      {meters.length === 0 && (
        <div className="text-center text-slate-400 py-12">No meters registered yet. Start the simulator to add meters.</div>
      )}
    </div>
  )
}
