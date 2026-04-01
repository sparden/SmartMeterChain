import { useState, useEffect } from 'react'
import { api } from '../lib/api'

export default function BillingPage() {
  const [bills, setBills] = useState<any[]>([])
  const [meters, setMeters] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState({ meter_id: '', period_start: '', period_end: '' })
  const [submitting, setSubmitting] = useState(false)
  const [msg, setMsg] = useState({ text: '', type: '' })

  const loadBills = () => {
    api.getBills(undefined, 1)
      .then((res) => setBills(res.data?.items || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadBills()
    api.getMeters(1, 50).then((res) => setMeters(res.data?.items || [])).catch(() => {})
  }, [])

  const handleGenerate = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)
    setMsg({ text: '', type: '' })
    try {
      const res = await api.generateBill({
        meter_id: formData.meter_id,
        period_start: formData.period_start,
        period_end: formData.period_end,
      })
      setMsg({ text: `Bill ${res.data?.bill_id || ''} generated! Amount: INR ${res.data?.amount?.toFixed(2) || '0'}`, type: 'success' })
      setShowForm(false)
      setFormData({ meter_id: '', period_start: '', period_end: '' })
      loadBills()
    } catch (err: any) {
      setMsg({ text: err.message || 'Failed to generate bill', type: 'error' })
    } finally {
      setSubmitting(false)
    }
  }

  const handlePay = async (billId: string) => {
    try {
      await api.payBill(billId)
      setMsg({ text: `Bill ${billId} marked as paid`, type: 'success' })
      loadBills()
    } catch (err: any) {
      setMsg({ text: err.message, type: 'error' })
    }
  }

  const statusColor = (status: string) => {
    switch (status) {
      case 'paid': return 'bg-green-100 text-green-700'
      case 'pending': return 'bg-yellow-100 text-yellow-700'
      case 'disputed': return 'bg-red-100 text-red-700'
      case 'overdue': return 'bg-orange-100 text-orange-700'
      default: return 'bg-slate-100 text-slate-700'
    }
  }

  if (loading) return <div className="text-slate-400 p-8">Loading bills...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-slate-800">Billing</h2>
        <button
          onClick={() => setShowForm(!showForm)}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700"
        >
          {showForm ? 'Cancel' : 'Generate Bill'}
        </button>
      </div>

      {msg.text && (
        <div className={`p-3 rounded-lg text-sm ${msg.type === 'success' ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'}`}>
          {msg.text}
        </div>
      )}

      {showForm && (
        <form onSubmit={handleGenerate} className="bg-white rounded-xl shadow-sm border border-slate-200 p-5 space-y-4">
          <h3 className="font-semibold text-slate-700">Generate New Bill</h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Meter</label>
              <select
                value={formData.meter_id}
                onChange={(e) => setFormData({ ...formData, meter_id: e.target.value })}
                className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
                required
              >
                <option value="">Select meter...</option>
                {meters.map((m: any) => (
                  <option key={m.meter_id} value={m.meter_id}>{m.meter_id} ({m.meter_type})</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Period Start</label>
              <input
                type="date"
                value={formData.period_start}
                onChange={(e) => setFormData({ ...formData, period_start: e.target.value })}
                className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-600 mb-1">Period End</label>
              <input
                type="date"
                value={formData.period_end}
                onChange={(e) => setFormData({ ...formData, period_end: e.target.value })}
                className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
                required
              />
            </div>
          </div>
          <button
            type="submit"
            disabled={submitting}
            className="px-5 py-2 bg-green-600 text-white rounded-lg text-sm font-medium hover:bg-green-700 disabled:opacity-50"
          >
            {submitting ? 'Generating...' : 'Generate Bill'}
          </button>
        </form>
      )}

      <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-slate-50 text-slate-600">
            <tr>
              <th className="text-left p-4 font-medium">Bill ID</th>
              <th className="text-left p-4 font-medium">Meter</th>
              <th className="text-left p-4 font-medium">Period</th>
              <th className="text-right p-4 font-medium">Units (kWh)</th>
              <th className="text-right p-4 font-medium">Amount (INR)</th>
              <th className="text-center p-4 font-medium">Status</th>
              <th className="text-center p-4 font-medium">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {bills.map((bill: any) => (
              <tr key={bill.bill_id} className="hover:bg-slate-50">
                <td className="p-4 font-mono text-blue-600 text-xs">{bill.bill_id}</td>
                <td className="p-4">{bill.meter_id}</td>
                <td className="p-4 text-slate-500">
                  {bill.period_start?.slice(0, 10)} to {bill.period_end?.slice(0, 10)}
                </td>
                <td className="p-4 text-right">{bill.units_used?.toFixed(2)}</td>
                <td className="p-4 text-right font-medium">{bill.amount?.toFixed(2)}</td>
                <td className="p-4 text-center">
                  <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${statusColor(bill.status)}`}>
                    {bill.status}
                  </span>
                </td>
                <td className="p-4 text-center space-x-2">
                  {bill.status === 'pending' && (
                    <button
                      onClick={() => handlePay(bill.bill_id)}
                      className="px-2 py-1 bg-green-100 text-green-700 rounded text-xs font-medium hover:bg-green-200"
                    >
                      Pay
                    </button>
                  )}
                  {bill.bill_hash && <span className="text-green-600" title="Verified on blockchain">🔗</span>}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {bills.length === 0 && (
          <div className="text-center text-slate-400 py-12">No bills generated yet. Click "Generate Bill" to create one.</div>
        )}
      </div>
    </div>
  )
}
