import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom'
import { useState } from 'react'
import { clearToken, getUser } from '../lib/api'

export default function RootLayout() {
  const location = useLocation()
  const path = location.pathname

  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar currentPath={path} />
      <main className="flex-1 overflow-y-auto bg-slate-50">
        <Header />
        <div className="p-6">
          <Outlet />
        </div>
      </main>
    </div>
  )
}

function Sidebar({ currentPath }: { currentPath: string }) {
  const user = getUser()
  const role = user?.role || 'admin'

  const allLinks = [
    { to: '/', label: 'Dashboard', icon: '📊', roles: ['admin', 'consumer', 'regulator'] },
    { to: '/meters', label: 'Meters', icon: '⚡', roles: ['admin', 'regulator'] },
    { to: '/billing', label: 'Billing', icon: '💰', roles: ['admin', 'consumer'] },
    { to: '/consumers', label: 'Consumers', icon: '👥', roles: ['admin'] },
    { to: '/disputes', label: 'Disputes', icon: '⚖️', roles: ['admin', 'consumer'] },
    { to: '/tariffs', label: 'Tariffs', icon: '📋', roles: ['admin', 'regulator'] },
    { to: '/alerts', label: 'Alerts', icon: '🔔', roles: ['admin', 'regulator'] },
    { to: '/explorer', label: 'Blockchain', icon: '🔗', roles: ['admin', 'consumer', 'regulator'] },
  ]

  const links = allLinks.filter((l) => l.roles.includes(role))

  return (
    <aside className="w-64 bg-slate-900 text-white flex flex-col">
      <div className="p-4 border-b border-slate-700">
        <h1 className="text-lg font-bold">SmartMeterChain</h1>
        <p className="text-xs text-slate-400 mt-1">Blockchain Smart Meter Ecosystem</p>
      </div>
      <nav className="flex-1 p-3 space-y-1">
        {links.map((link) => (
          <Link
            key={link.to}
            to={link.to}
            className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-colors ${
              currentPath === link.to
                ? 'bg-blue-600 text-white'
                : 'text-slate-300 hover:bg-slate-800 hover:text-white'
            }`}
          >
            <span>{link.icon}</span>
            {link.label}
          </Link>
        ))}
      </nav>
      {role === 'regulator' && (
        <div className="px-4 py-2 text-xs text-yellow-400 border-t border-slate-700">
          Read-only access (Regulator)
        </div>
      )}
      <div className="p-4 border-t border-slate-700 text-xs text-slate-500">
        Powered by Hyperledger Fabric
      </div>
    </aside>
  )
}

function Header() {
  const navigate = useNavigate()
  const [user] = useState(() => getUser())

  const handleLogout = () => {
    clearToken()
    navigate('/login')
  }

  const roleColor = (role: string) => {
    switch (role) {
      case 'admin': return 'bg-blue-100 text-blue-700'
      case 'consumer': return 'bg-green-100 text-green-700'
      case 'regulator': return 'bg-purple-100 text-purple-700'
      default: return 'bg-slate-100 text-slate-700'
    }
  }

  return (
    <header className="h-14 border-b border-slate-200 bg-white flex items-center justify-between px-6">
      <div />
      <div className="flex items-center gap-4">
        <span className="text-sm text-slate-600">{user?.name || 'User'}</span>
        <span className={`px-2 py-0.5 text-xs rounded-full font-medium ${roleColor(user?.role)}`}>
          {user?.role || 'unknown'}
        </span>
        <button
          onClick={handleLogout}
          className="px-3 py-1.5 text-xs font-medium text-red-600 bg-red-50 rounded-lg hover:bg-red-100 transition-colors"
        >
          Logout
        </button>
      </div>
    </header>
  )
}
