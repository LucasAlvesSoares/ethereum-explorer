'use client'

import React from 'react'
import { Inter } from 'next/font/google'
import Link from 'next/link'
import ModeSelector from '@/components/ModeSelector'
import RealTimeStatus from '@/components/RealTimeStatus'
import { RealTimeProvider } from '@/contexts/RealTimeContext'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <RealTimeProvider>
          <div className="min-h-screen bg-gray-50">
            <header className="bg-white shadow-sm border-b">
              <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex justify-between items-center h-16">
                  <div className="flex items-center">
                    <Link href="/" className="text-xl font-bold text-gray-900 hover:text-gray-700">
                      Ethereum Explorer
                    </Link>
                  </div>
                  <div className="flex items-center space-x-8">
                    <RealTimeStatus />
                    <nav className="flex space-x-8">
                      <Link href="/" className="text-gray-500 hover:text-gray-900">
                        Home
                      </Link>
                      <Link href="/blocks" className="text-gray-500 hover:text-gray-900">
                        Blocks
                      </Link>
                      <Link href="/transactions" className="text-gray-500 hover:text-gray-900">
                        Transactions
                      </Link>
                      <Link href="/transaction-flow" className="text-gray-500 hover:text-gray-900">
                        Flow Analysis
                      </Link>
                      <Link href="/mev-analytics" className="text-gray-500 hover:text-gray-900">
                        MEV Analytics
                      </Link>
                    </nav>
                    <ModeSelector />
                  </div>
                </div>
              </div>
            </header>
            <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
              {children}
            </main>
          </div>
        </RealTimeProvider>
      </body>
    </html>
  )
}
