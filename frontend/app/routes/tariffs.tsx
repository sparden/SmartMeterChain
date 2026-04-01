import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { api } from '~/lib/api'

export const Route = createFileRoute('/tariffs')({
  component: TariffsPage,
})

function TariffsPage() {
  const [tariffs, setTariffs] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getTariffs()
      .then((res) => setTariffs(res.data || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="text-slate-400 p-8">Loading tariffs...</div>

  const grouped = tariffs.reduce((acc: any, t: any) => {
    if (!acc[t.category]) acc[t.category] = []
    acc[t.category].push(t)
    return acc
  }, {})

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-800">Tariff Slabs</h2>

      {Object.entries(grouped).map(([category, slabs]: [string, any]) => (
        <div key={category} className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
          <div className="bg-slate-50 px-5 py-3 border-b border-slate-200">
            <h3 className="font-semibold text-slate-700 capitalize">{category}</h3>
          </div>
          <table className="w-full text-sm">
            <thead>
              <tr className="text-slate-500">
                <th className="text-left p-4 font-medium">Slab Range (kWh)</th>
                <th className="text-right p-4 font-medium">Rate/Unit (INR)</th>
                <th className="text-right p-4 font-medium">Fixed Charge (INR)</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {slabs.map((slab: any) => (
                <tr key={slab.tariff_id}>
                  <td className="p-4">{slab.slab_start} — {slab.slab_end >= 99999 ? '∞' : slab.slab_end}</td>
                  <td className="p-4 text-right font-medium">{slab.rate_per_unit?.toFixed(2)}</td>
                  <td className="p-4 text-right">{slab.fixed_charge?.toFixed(2)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ))}
    </div>
  )
}
