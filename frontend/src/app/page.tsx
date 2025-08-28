'use client'

import { useState, useEffect } from 'react'
import { Search, Activity, Blocks, Zap } from 'lucide-react'

export default function Home() {
  const [searchQuery, setSearchQuery] = useState('')
  const [apiStatus, setApiStatus] = useState<'loading' | 'healthy' | 'error'>('loading')

  useEffect(() => {
    // Check API health on component mount
    fetch('/api/v1/health')
      .then(res => res.json())
      .then(data => {
        setApiStatus(data.status === 'healthy' ? 'healthy' : 'error')
      })
      .catch(() => setApiStatus('error'))
  }, [])

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      // TODO: Implement search functionality
      console.log('Searching for:', searchQuery)
    }
  }

  return (
    <div className="px-4 sm:px-0">
      {/* Hero Section */}
      <div className="text-center mb-12">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          Ethereum Blockchain Explorer
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          Advanced analytics and real-time insights for the Ethereum blockchain
        </p>
        
        {/* Search Bar */}
        <form onSubmit={handleSearch} className="max-w-2xl mx-auto">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search by address, transaction hash, or block number..."
              className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
            <button
              type="submit"
              className="absolute right-2 top-1/2 transform -translate-y-1/2 btn-primary"
            >
              Search
            </button>
          </div>
        </form>
      </div>

      {/* API Status */}
      <div className="mb-8">
        <div className="card">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">System Status</h2>
            <div className="flex items-center space-x-2">
              <div className={`w-3 h-3 rounded-full ${
                apiStatus === 'healthy' ? 'bg-green-500' : 
                apiStatus === 'error' ? 'bg-red-500' : 'bg-yellow-500'
              }`} />
              <span className="text-sm text-gray-600">
                API: {apiStatus === 'loading' ? 'Checking...' : apiStatus}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Feature Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
        <div className="card">
          <div className="flex items-center mb-4">
            <Blocks className="w-8 h-8 text-primary-600 mr-3" />
            <h3 className="text-lg font-semibold text-gray-900">Block Explorer</h3>
          </div>
          <p className="text-gray-600">
            Browse blocks, transactions, and addresses with detailed information and real-time updates.
          </p>
        </div>

        <div className="card">
          <div className="flex items-center mb-4">
            <Activity className="w-8 h-8 text-primary-600 mr-3" />
            <h3 className="text-lg font-semibold text-gray-900">Advanced Analytics</h3>
          </div>
          <p className="text-gray-600">
            Analyze transaction flows, gas prices, and network statistics with interactive visualizations.
          </p>
        </div>

        <div className="card">
          <div className="flex items-center mb-4">
            <Zap className="w-8 h-8 text-primary-600 mr-3" />
            <h3 className="text-lg font-semibold text-gray-900">Real-time Data</h3>
          </div>
          <p className="text-gray-600">
            Get live updates on network activity, pending transactions, and block confirmations.
          </p>
        </div>
      </div>

      {/* Quick Stats */}
      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Network Overview</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="text-center">
            <div className="text-2xl font-bold text-primary-600">-</div>
            <div className="text-sm text-gray-600">Latest Block</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-primary-600">-</div>
            <div className="text-sm text-gray-600">Gas Price</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-primary-600">-</div>
            <div className="text-sm text-gray-600">Network Hash Rate</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-primary-600">-</div>
            <div className="text-sm text-gray-600">Active Addresses</div>
          </div>
        </div>
        <p className="text-sm text-gray-500 mt-4 text-center">
          Real-time data will be available once the backend is connected to an Ethereum node
        </p>
      </div>
    </div>
  )
}
