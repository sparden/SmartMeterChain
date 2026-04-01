import { useState, useEffect } from 'react'
import { api } from '../lib/api'
import {
  XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  LineChart, Line,
} from 'recharts'

function StatCard({ title, value, subtitle, color }: {
  title: string; value: string | number; subtitle?: string; color: string
}) {
  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-5">
      <p className="text-sm text-slate-500">{title}</p>
      <p className={`text-3xl font-bold mt-1 ${color}`}>{value}</p>
      {subtitle && <p className="text-xs text-slate-400 mt-1">{subtitle}</p>}
    </div>
  )
}

export default function Dashboard() {
  const [stats, setStats] = useState<any>(null)
  const [trend, setTrend] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([api.getStats(), api.getConsumptionTrend(30)])
      .then(([s, t]) => {
        setStats(s.data)
        setTrend(t.data || [])
      })
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return <div className="flex items-center justify-center h-64 text-slate-400">Loading dashboard...</div>
  }

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-800">Dashboard</h2>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard title="Active Meters" value={stats?.meters?.active || 0} subtitle="Total registered" color="text-blue-600" />
        <StatCard title="Today's Readings" value={stats?.readings?.today || 0} subtitle="Ingested today" color="text-green-600" />
        <StatCard title="Revenue" value={`INR ${(stats?.billing?.total_revenue || 0).toLocaleString()}`} subtitle="Total collected" color="text-emerald-600" />
        <StatCard title="Open Disputes" value={stats?.disputes?.open || 0} subtitle="Pending resolution" color="text-amber-600" />
      </div>

      {/* Tamper Alerts Summary */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <StatCard title="Total Alerts" value={stats?.alerts?.total_alerts || 0} color="text-red-600" />
        <StatCard title="Unacknowledged" value={stats?.alerts?.unacknowledged || 0} color="text-orange-600" />
        <StatCard title="Critical" value={stats?.alerts?.critical || 0} color="text-red-700" />
      </div>

      {/* Consumption Trend Chart */}
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-5">
        <h3 className="text-lg font-semibold text-slate-700 mb-4">Consumption Trend (30 Days)</h3>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={trend}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" tick={{ fontSize: 12 }} />
            <YAxis tick={{ fontSize: 12 }} />
            <Tooltip />
            <Line type="monotone" dataKey="avg_value" stroke="#2563eb" strokeWidth={2} name="Avg kWh" />
            <Line type="monotone" dataKey="readings" stroke="#16a34a" strokeWidth={2} name="Readings" />
          </LineChart>
        </ResponsiveContainer>
      </div>

      {/* Blockchain Badge */}
      <div className="bg-gradient-to-r from-blue-600 to-indigo-700 rounded-xl p-5 text-white">
        <div className="flex items-center gap-3">
          <span className="text-2xl">🔗</span>
          <div>
            <p className="font-semibold">Blockchain Secured</p>
            <p className="text-sm text-blue-100">All meter readings and bills are immutably recorded on Hyperledger Fabric</p>
          </div>
        </div>
      </div>
    </div>
  )
}
