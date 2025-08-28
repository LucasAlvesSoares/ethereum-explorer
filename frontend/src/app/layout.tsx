import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Ethereum Blockchain Explorer',
  description: 'Advanced Ethereum blockchain explorer with analytics',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <div className="min-h-screen bg-gray-50">
          <header className="bg-white shadow-sm border-b">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
              <div className="flex justify-between items-center h-16">
                <div className="flex items-center">
                  <h1 className="text-xl font-bold text-gray-900">
                    Ethereum Explorer
                  </h1>
                </div>
                <nav className="flex space-x-8">
                  <a href="/" className="text-gray-500 hover:text-gray-900">
                    Home
                  </a>
                  <a href="/blocks" className="text-gray-500 hover:text-gray-900">
                    Blocks
                  </a>
                  <a href="/transactions" className="text-gray-500 hover:text-gray-900">
                    Transactions
                  </a>
                  <a href="/analytics" className="text-gray-500 hover:text-gray-900">
                    Analytics
                  </a>
                </nav>
              </div>
            </div>
          </header>
          <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
            {children}
          </main>
        </div>
      </body>
    </html>
  )
}
