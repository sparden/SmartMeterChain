import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { api } from '~/lib/api'

export const Route = createFileRoute('/explorer')({
  component: ExplorerPage,
})

function ExplorerPage() {
  const [billId, setBillId] = useState('')
  const [result, setResult] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleVerify = async () => {
    if (!billId.trim()) return
    setLoading(true)
    setError('')
    setResult(null)

    try {
      const res = await api.verifyBill(billId.trim())
      setResult(res.data)
    } catch (err: any) {
      setError(err.message || 'Verification failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-800">Blockchain Explorer</h2>

      {/* Bill Verification */}
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <h3 className="text-lg font-semibold text-slate-700 mb-4">Verify Bill on Blockchain</h3>
        <div className="flex gap-3">
          <input
            type="text"
            value={billId}
            onChange={(e) => setBillId(e.target.value)}
            placeholder="Enter Bill ID (e.g., BILL-abc12345)"
            className="flex-1 px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
          />
          <button
            onClick={handleVerify}
            disabled={loading}
            className="px-6 py-2.5 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 disabled:opacity-50"
          >
            {loading ? 'Verifying...' : 'Verify'}
          </button>
        </div>

        {error && <p className="mt-3 text-red-600 text-sm">{error}</p>}

        {result && (
          <div className="mt-4 p-4 rounded-lg bg-slate-50 border">
            <div className="flex items-center gap-2 mb-3">
              {result.verified ? (
                <span className="text-green-600 font-semibold">Verified on Blockchain</span>
              ) : (
                <span className="text-red-600 font-semibold">Verification Failed</span>
              )}
            </div>
            {result.bill && (
              <div className="grid grid-cols-2 gap-2 text-sm">
                <p><span className="text-slate-400">Bill ID:</span> {result.bill.bill_id}</p>
                <p><span className="text-slate-400">Meter:</span> {result.bill.meter_id}</p>
                <p><span className="text-slate-400">Units:</span> {result.bill.units_used?.toFixed(2)} kWh</p>
                <p><span className="text-slate-400">Amount:</span> INR {result.bill.amount?.toFixed(2)}</p>
                <p><span className="text-slate-400">TX ID:</span> <span className="font-mono text-xs">{result.bill.tx_id}</span></p>
                <p><span className="text-slate-400">Hash:</span> <span className="font-mono text-xs">{result.bill.bill_hash?.slice(0, 16)}...</span></p>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Info Section */}
      <div className="bg-gradient-to-r from-slate-800 to-slate-900 rounded-xl p-6 text-white">
        <h3 className="text-lg font-semibold mb-3">How Blockchain Verification Works</h3>
        <div className="space-y-2 text-sm text-slate-300">
          <p>1. Every meter reading is hashed (SHA-256) and submitted to Hyperledger Fabric</p>
          <p>2. Bills are generated from on-chain readings using slab-based tariffs</p>
          <p>3. Bill hash = SHA-256(bill_id + units_used + amount)</p>
          <p>4. Verification compares off-chain hash with on-chain record</p>
          <p>5. Any tampering is immediately detectable and logged</p>
        </div>
      </div>
    </div>
  )
}
