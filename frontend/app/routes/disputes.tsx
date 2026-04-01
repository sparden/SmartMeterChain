import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { api } from '~/lib/api'

export const Route = createFileRoute('/disputes')({
  component: DisputesPage,
})

function DisputesPage() {
  const [disputes, setDisputes] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getDisputes(1)
      .then((res) => setDisputes(res.data?.items || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="text-slate-400 p-8">Loading disputes...</div>

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-800">Disputes</h2>

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
          </div>
        ))}
        {disputes.length === 0 && (
          <div className="text-center text-slate-400 py-12">No disputes filed yet.</div>
        )}
      </div>
    </div>
  )
}
