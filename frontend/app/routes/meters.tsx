import { useState, useEffect } from 'react'
import { api } from '../lib/api'

export default function MetersPage() {
  const [meters, setMeters] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState({ meter_id: '', consumer_id: '', location: '', meter_type: 'domestic' })
  const [submitting, setSubmitting] = useState(false)
  const [msg, setMsg] = useState({ text: '', type: '' })

  const loadMeters = () => {
    api.getMeters(1, 50)
      .then((res) => setMeters(res.data?.items || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  useEffect(() => { loadMeters() }, [])

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)
    setMsg({ text: '', type: '' })
    try {
      await api.registerMeter(formData)
      setMsg({ text: `Meter ${formData.meter_id} registered successfully`, type: 'success' })
      setShowForm(false)
      setFormData({ meter_id: '', consumer_id: '', location: '', meter_type: 'domestic' })
      loadMeters()
    } catch (err: any) {
      setMsg({ text: err.message, type: 'error' })
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) return <div className="text-slate-400 p-8">Loading meters...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-slate-800">Smart Meters</h2>
        <div className="flex items-center gap-3">
          <span className="text-sm text-slate-500">{meters.length} meters registered</span>
          <button
            onClick={() => setShowForm(!showForm)}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700"
          >
            {showForm ? 'Cancel' : 'Register Meter'}
          </button>
        </div>
      </div>

      {msg.text && (
        <div className={`p-3 rounded-lg text-sm ${msg.type === 'success' ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'}`}>
          {msg.text}
        </div>
      )}

      {showForm && (
        <form onSubmit={handleRegister} className="bg-white rounded-xl shadow-sm border border-slate-200 p-5 space-y-4">
          <h3 className="font-semibold text-slate-700">Register New Meter</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Meter ID</label>
              <input
                type="text"
                value={formData.meter_id}
                onChange={(e) => setFormData({ ...formData, meter_id: e.target.value })}
                placeholder="e.g., SM-DEL-003"
                className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Consumer ID</label>
              <input
                type="text"
                value={formData.consumer_id}
                onChange={(e) => setFormData({ ...formData, consumer_id: e.target.value })}
                placeholder="e.g., consumer1"
                className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Location</label>
              <input
                type="text"
                value={formData.location}
                onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                placeholder="e.g., New Delhi, Sector 10"
                className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Meter Type</label>
              <select
                value={formData.meter_type}
                onChange={(e) => setFormData({ ...formData, meter_type: e.target.value })}
                className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
              >
                <option value="domestic">Domestic</option>
                <option value="commercial">Commercial</option>
                <option value="industrial">Industrial</option>
              </select>
            </div>
          </div>
          <button
            type="submit"
            disabled={submitting}
            className="px-5 py-2 bg-green-600 text-white rounded-lg text-sm font-medium hover:bg-green-700 disabled:opacity-50"
          >
            {submitting ? 'Registering...' : 'Register Meter'}
          </button>
        </form>
      )}

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
        <div className="text-center text-slate-400 py-12">No meters registered yet. Click "Register Meter" to add one.</div>
      )}
    </div>
  )
}
