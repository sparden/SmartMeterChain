import { useState, useEffect } from 'react'
import { api, getUser } from '../lib/api'

export default function DisputesPage() {
  const [disputes, setDisputes] = useState<any[]>([])
  const [bills, setBills] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [billId, setBillId] = useState('')
  const [reason, setReason] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [msg, setMsg] = useState({ text: '', type: '' })
  const [resolveId, setResolveId] = useState('')
  const [resolution, setResolution] = useState('')
  const user = getUser()
  const isAdmin = user?.role === 'admin'

  const loadDisputes = () => {
    api.getDisputes(1)
      .then((res) => setDisputes(res.data?.items || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadDisputes()
    api.getBills(undefined, 1).then((res) => setBills(res.data?.items || [])).catch(() => {})
  }, [])

  const handleFile = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)
    setMsg({ text: '', type: '' })
    try {
      await api.fileDispute(billId, reason)
      setMsg({ text: 'Dispute filed successfully', type: 'success' })
      setShowForm(false)
      setBillId('')
      setReason('')
      loadDisputes()
    } catch (err: any) {
      setMsg({ text: err.message, type: 'error' })
    } finally {
      setSubmitting(false)
    }
  }

  const handleResolve = async (disputeId: string) => {
    if (!resolution.trim()) return
    try {
      await api.resolveDispute(disputeId, resolution)
      setMsg({ text: `Dispute ${disputeId} resolved`, type: 'success' })
      setResolveId('')
      setResolution('')
      loadDisputes()
    } catch (err: any) {
      setMsg({ text: err.message, type: 'error' })
    }
  }

  if (loading) return <div className="text-slate-400 p-8">Loading disputes...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-slate-800">Disputes</h2>
        <button
          onClick={() => setShowForm(!showForm)}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700"
        >
          {showForm ? 'Cancel' : 'File Dispute'}
        </button>
      </div>

      {msg.text && (
        <div className={`p-3 rounded-lg text-sm ${msg.type === 'success' ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'}`}>
          {msg.text}
        </div>
      )}

      {showForm && (
        <form onSubmit={handleFile} className="bg-white rounded-xl shadow-sm border border-slate-200 p-5 space-y-4">
          <h3 className="font-semibold text-slate-700">File New Dispute</h3>
          <div>
            <label className="block text-sm font-medium text-slate-600 mb-1">Bill ID</label>
            {bills.length > 0 ? (
              <select
                value={billId}
                onChange={(e) => setBillId(e.target.value)}
                className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
                required
              >
                <option value="">Select bill...</option>
                {bills.map((b: any) => (
                  <option key={b.bill_id} value={b.bill_id}>{b.bill_id} — {b.meter_id} (INR {b.amount?.toFixed(2)})</option>
                ))}
              </select>
            ) : (
              <input
                type="text"
                value={billId}
                onChange={(e) => setBillId(e.target.value)}
                placeholder="Enter Bill ID"
                className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none"
                required
              />
            )}
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-600 mb-1">Reason</label>
            <textarea
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="Describe the issue with this bill..."
              rows={3}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none resize-none"
              required
            />
          </div>
          <button
            type="submit"
            disabled={submitting}
            className="px-5 py-2 bg-orange-600 text-white rounded-lg text-sm font-medium hover:bg-orange-700 disabled:opacity-50"
          >
            {submitting ? 'Filing...' : 'File Dispute'}
          </button>
        </form>
      )}

      <div className="space-y-4">
        {disputes.map((d: any) => (
          <div key={d.dispute_id} className="bg-white rounded-xl shadow-sm border border-slate-200 p-5">
            <div className="flex items-center justify-between mb-3">
              <span className="font-mono font-bold text-slate-700">{d.dispute_id}</span>
              <span className={`px-2.5 py-0.5 rounded-full text-xs font-medium ${
                d.status === 'open' ? 'bg-yellow-100 text-yellow-700' :
                d.status === 'resolved' ? 'bg-green-100 text-green-700' :
                'bg-red-100 text-red-700'
              }`}>
                {d.status}
              </span>
            </div>
            <p className="text-sm text-slate-600 mb-2"><strong>Bill:</strong> {d.bill_id}</p>
            <p className="text-sm text-slate-600 mb-2"><strong>Reason:</strong> {d.reason}</p>
            {d.resolution && (
              <p className="text-sm text-green-700"><strong>Resolution:</strong> {d.resolution}</p>
            )}
            <p className="text-xs text-slate-400 mt-3">Filed: {d.filed_at?.slice(0, 10)}</p>

            {d.status === 'open' && isAdmin && (
              <div className="mt-3 pt-3 border-t border-slate-100">
                {resolveId === d.dispute_id ? (
                  <div className="flex gap-2">
                    <input
                      type="text"
                      value={resolution}
                      onChange={(e) => setResolution(e.target.value)}
                      placeholder="Enter resolution..."
                      className="flex-1 px-3 py-1.5 border border-slate-300 rounded-lg text-sm outline-none"
                    />
                    <button
                      onClick={() => handleResolve(d.dispute_id)}
                      className="px-3 py-1.5 bg-green-600 text-white rounded-lg text-xs font-medium hover:bg-green-700"
                    >
                      Resolve
                    </button>
                    <button
                      onClick={() => { setResolveId(''); setResolution('') }}
                      className="px-3 py-1.5 bg-slate-200 text-slate-600 rounded-lg text-xs font-medium hover:bg-slate-300"
                    >
                      Cancel
                    </button>
                  </div>
                ) : (
                  <button
                    onClick={() => setResolveId(d.dispute_id)}
                    className="px-3 py-1.5 bg-blue-100 text-blue-700 rounded-lg text-xs font-medium hover:bg-blue-200"
                  >
                    Resolve Dispute
                  </button>
                )}
              </div>
            )}
          </div>
        ))}
        {disputes.length === 0 && (
          <div className="text-center text-slate-400 py-12">No disputes filed yet. Click "File Dispute" to create one.</div>
        )}
      </div>
    </div>
  )
}
