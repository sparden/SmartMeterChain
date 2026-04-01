import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { api } from '~/lib/api'

export const Route = createFileRoute('/billing')({
  component: BillingPage,
})

function BillingPage() {
  const [bills, setBills] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getBills(undefined, 1)
      .then((res) => setBills(res.data?.items || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

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
      <h2 className="text-2xl font-bold text-slate-800">Billing</h2>

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
              <th className="text-center p-4 font-medium">Verified</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {bills.map((bill: any) => (
              <tr key={bill.bill_id} className="hover:bg-slate-50">
                <td className="p-4 font-mono text-blue-600">{bill.bill_id}</td>
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
                <td className="p-4 text-center">
                  {bill.bill_hash ? <span className="text-green-600">🔗</span> : '—'}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {bills.length === 0 && (
          <div className="text-center text-slate-400 py-12">No bills generated yet.</div>
        )}
      </div>
    </div>
  )
}
