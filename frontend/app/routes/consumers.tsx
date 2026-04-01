import { useState, useEffect } from 'react'
import { api } from '../lib/api'

export default function ConsumersPage() {
  const [consumers, setConsumers] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getConsumers(1)
      .then((res) => setConsumers(res.data?.items || res.data || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="text-slate-400 p-8">Loading consumers...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-slate-800">Consumers</h2>
        <span className="text-sm text-slate-500">{consumers.length} registered</span>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-slate-50 text-slate-600">
            <tr>
              <th className="text-left p-4 font-medium">ID</th>
              <th className="text-left p-4 font-medium">Username</th>
              <th className="text-left p-4 font-medium">Name</th>
              <th className="text-left p-4 font-medium">Email</th>
              <th className="text-left p-4 font-medium">Phone</th>
              <th className="text-center p-4 font-medium">Role</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {consumers.map((c: any) => (
              <tr key={c.ID || c.id || c.username} className="hover:bg-slate-50">
                <td className="p-4 font-mono text-slate-500">{c.ID || c.id}</td>
                <td className="p-4 font-medium text-blue-600">{c.username}</td>
                <td className="p-4">{c.name || '—'}</td>
                <td className="p-4 text-slate-500">{c.email}</td>
                <td className="p-4 text-slate-500">{c.phone || '—'}</td>
                <td className="p-4 text-center">
                  <span className="px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-700">
                    {c.role}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {consumers.length === 0 && (
          <div className="text-center text-slate-400 py-12">No consumers found.</div>
        )}
      </div>
    </div>
  )
}
