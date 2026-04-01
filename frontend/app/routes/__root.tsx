import { createRootRoute, Outlet, Link, useRouter } from '@tanstack/react-router'
import { useState } from 'react'
import '../styles.css'

export const Route = createRootRoute({
  component: RootLayout,
})

function RootLayout() {
  const router = useRouter()
  const path = router.state.location.pathname
  const isLoginPage = path === '/login'

  if (isLoginPage) {
    return <Outlet />
  }

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
  const links = [
    { to: '/', label: 'Dashboard', icon: '📊' },
    { to: '/meters', label: 'Meters', icon: '⚡' },
    { to: '/billing', label: 'Billing', icon: '💰' },
    { to: '/consumers', label: 'Consumers', icon: '👥' },
    { to: '/disputes', label: 'Disputes', icon: '⚖️' },
    { to: '/tariffs', label: 'Tariffs', icon: '📋' },
    { to: '/alerts', label: 'Alerts', icon: '🔔' },
    { to: '/explorer', label: 'Blockchain', icon: '🔗' },
  ]

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
      <div className="p-4 border-t border-slate-700 text-xs text-slate-500">
        Powered by Hyperledger Fabric
      </div>
    </aside>
  )
}

function Header() {
  const [user] = useState(() => {
    if (typeof window !== 'undefined') {
      const stored = localStorage.getItem('smc_user')
      return stored ? JSON.parse(stored) : null
    }
    return null
  })

  return (
    <header className="h-14 border-b border-slate-200 bg-white flex items-center justify-between px-6">
      <div />
      <div className="flex items-center gap-4">
        <span className="text-sm text-slate-600">{user?.name || 'Admin'}</span>
        <span className="px-2 py-0.5 text-xs rounded-full bg-blue-100 text-blue-700 font-medium">
          {user?.role || 'admin'}
        </span>
      </div>
    </header>
  )
}
