'use client'

import { useState, useEffect } from 'react'
import { Search, Activity, Blocks, Zap, Clock, Hash, User } from 'lucide-react'
import Link from 'next/link'

interface Block {
  number: number
  hash: string
  timestamp: string
  gas_used: number
  gas_limit: number
  transaction_count: number
  miner: string
}

interface BlocksResponse {
  blocks: Block[]
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}

export default function Home() {
  const [searchQuery, setSearchQuery] = useState('')
  const [apiStatus, setApiStatus] = useState<'loading' | 'healthy' | 'error'>('loading')
  const [latestBlocks, setLatestBlocks] = useState<Block[]>([])
  const [blocksLoading, setBlocksLoading] = useState(true)

  useEffect(() => {
    // Check API health on component mount
    fetch('/api/v1/health')
      .then(res => res.json())
      .then(data => {
        setApiStatus(data.status === 'healthy' ? 'healthy' : 'error')
      })
      .catch(() => setApiStatus('error'))

    // Fetch latest blocks
    fetchLatestBlocks()
  }, [])

  const fetchLatestBlocks = async () => {
    try {
      const response = await fetch('/api/v1/blocks?page=1&limit=5')
      if (response.ok) {
        const data: BlocksResponse = await response.json()
        setLatestBlocks(data.blocks || [])
      }
    } catch (error) {
      console.error('Failed to fetch latest blocks:', error)
    } finally {
      setBlocksLoading(false)
    }
  }

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

      {/* Latest Blocks Section */}
      <div className="mb-12">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-semibold text-gray-900">Latest Blocks</h2>
          <Link href="/blocks" className="text-primary-600 hover:text-primary-800 text-sm font-medium">
            View all blocks →
          </Link>
        </div>
        
        {blocksLoading ? (
          <div className="card">
            <div className="flex justify-center items-center h-32">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
            </div>
          </div>
        ) : latestBlocks.length > 0 ? (
          <div className="card overflow-hidden">
            <div className="divide-y divide-gray-200">
              {latestBlocks.map((block) => (
                <div key={block.number} className="p-4 hover:bg-gray-50 transition-colors">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-4">
                      <div className="flex-shrink-0">
                        <div className="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center">
                          <Blocks className="w-5 h-5 text-primary-600" />
                        </div>
                      </div>
                      <div>
                        <div className="flex items-center space-x-2">
                          <Link
                            href={`/blocks/${block.number}`}
                            className="text-lg font-semibold text-primary-600 hover:text-primary-800"
                          >
                            Block {block.number.toLocaleString()}
                          </Link>
                          <span className="text-sm text-gray-500">
                            {new Date(block.timestamp).toLocaleTimeString()}
                          </span>
                        </div>
                        <div className="flex items-center space-x-4 mt-1 text-sm text-gray-600">
                          <div className="flex items-center">
                            <Hash className="w-3 h-3 mr-1" />
                            <span className="font-mono">
                              {block.hash.slice(0, 10)}...{block.hash.slice(-8)}
                            </span>
                          </div>
                          <div className="flex items-center">
                            <Activity className="w-3 h-3 mr-1" />
                            <span>{block.transaction_count} txns</span>
                          </div>
                          <div className="flex items-center">
                            <User className="w-3 h-3 mr-1" />
                            <span className="font-mono">
                              {block.miner.slice(0, 6)}...{block.miner.slice(-4)}
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="text-sm text-gray-600">Gas Used</div>
                      <div className="text-sm font-medium">
                        {((block.gas_used / block.gas_limit) * 100).toFixed(1)}%
                      </div>
                      <div className="w-16 bg-gray-200 rounded-full h-1 mt-1">
                        <div
                          className="bg-primary-600 h-1 rounded-full"
                          style={{
                            width: `${(block.gas_used / block.gas_limit) * 100}%`
                          }}
                        ></div>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        ) : (
          <div className="card">
            <div className="text-center py-8">
              <Blocks className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">No blocks available</h3>
              <p className="text-gray-600">
                Blocks will appear here once the blockchain ingestion service processes data.
              </p>
            </div>
          </div>
        )}
      </div>

      {/* Navigation Links */}
      <div className="mb-12">
        <h2 className="text-2xl font-semibold text-gray-900 mb-6 text-center">Explore the Platform</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Link href="/transactions" className="block p-4 bg-white rounded-lg border border-gray-200 hover:border-blue-300 hover:shadow-md transition-all">
            <div className="flex items-center mb-2">
              <Activity className="w-5 h-5 text-blue-600 mr-2" />
              <h3 className="font-semibold text-gray-800">Transactions</h3>
            </div>
            <p className="text-sm text-gray-600">Explore recent transactions</p>
          </Link>
          
          <Link href="/gas-analytics" className="block p-4 bg-white rounded-lg border border-gray-200 hover:border-blue-300 hover:shadow-md transition-all">
            <div className="flex items-center mb-2">
              <Zap className="w-5 h-5 text-blue-600 mr-2" />
              <h3 className="font-semibold text-gray-800">Gas Analytics</h3>
            </div>
            <p className="text-sm text-gray-600">Real-time gas prices and trends</p>
          </Link>
          
          <Link href="/transaction-flow" className="block p-4 bg-white rounded-lg border border-gray-200 hover:border-blue-300 hover:shadow-md transition-all">
            <div className="flex items-center mb-2">
              <Search className="w-5 h-5 text-blue-600 mr-2" />
              <h3 className="font-semibold text-gray-800">Transaction Flow</h3>
            </div>
            <p className="text-sm text-gray-600">Visualize transaction relationships</p>
          </Link>
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

    </div>
  )
}
