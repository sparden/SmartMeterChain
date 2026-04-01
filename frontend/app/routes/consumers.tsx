export default function ConsumersPage() {
  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-800">Consumers</h2>
      <p className="text-slate-500">Consumer management panel — view registered consumers and their linked meters.</p>
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-8 text-center text-slate-400">
        Consumer data loads from the backend API. Start the backend server to see consumer data.
      </div>
    </div>
  )
}
