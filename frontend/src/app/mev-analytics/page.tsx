'use client'

import { useState, useEffect } from 'react'
import { formatEther, formatUnits } from '@/utils/formatting'

interface MEVActivity {
  id: string
  transaction_hash: string
  block_number: number
  mev_type: string
  value_extracted: string
  gas_used: number
  addresses_involved: string[]
  confidence_score: number
  timestamp: string
}

interface MEVStats {
  total_mev_extracted_24h: string
  mev_transactions_count: number
  average_mev_per_transaction: string
  most_active_mev_addresses: Array<{
    address: string
    mev_count: number
    total_extracted: string
  }>
  mev_by_type: Array<{
    type: string
    count: number
    total_value: string
  }>
}

export default function MEVAnalytics() {
  const [mevActivities, setMevActivities] = useState<MEVActivity[]>([])
  const [mevStats, setMevStats] = useState<MEVStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [selectedType, setSelectedType] = useState<string>('all')

  useEffect(() => {
    fetchMEVData()
  }, [selectedType])

  const fetchMEVData = async () => {
    try {
      setLoading(true)
      
      // Fetch MEV activities
      const activitiesParams = selectedType !== 'all' ? `?type=${selectedType}` : ''
      const activitiesResponse = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/mev-analytics/activities${activitiesParams}`)
      
      if (activitiesResponse.ok) {
        const activitiesData = await activitiesResponse.json()
        setMevActivities(activitiesData.data || [])
      }

      // Fetch MEV statistics
      const statsResponse = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/mev-analytics/stats`)
      if (statsResponse.ok) {
        const statsData = await statsResponse.json()
        setMevStats(statsData.data)
      }
    } catch (error) {
      console.error('Error fetching MEV data:', error)
    } finally {
      setLoading(false)
    }
  }

  const getMEVTypeColor = (type: string) => {
    const colors: { [key: string]: string } = {
      'sandwich_attack': 'bg-red-100 text-red-800',
      'front_running': 'bg-orange-100 text-orange-800',
      'back_running': 'bg-yellow-100 text-yellow-800',
      'arbitrage': 'bg-blue-100 text-blue-800',
      'liquidation': 'bg-purple-100 text-purple-800',
      'suspicious': 'bg-gray-100 text-gray-800'
    }
    return colors[type] || 'bg-gray-100 text-gray-800'
  }

  const getConfidenceColor = (score: number) => {
    if (score >= 0.8) return 'text-red-600 font-semibold'
    if (score >= 0.6) return 'text-orange-600 font-medium'
    if (score >= 0.4) return 'text-yellow-600'
    return 'text-gray-600'
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="bg-white p-6 rounded-lg shadow">
                <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                <div className="h-6 bg-gray-200 rounded w-1/2"></div>
              </div>
            ))}
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="h-6 bg-gray-200 rounded w-1/3 mb-4"></div>
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-4 bg-gray-200 rounded mb-2"></div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">MEV Analytics</h1>
        <p className="text-gray-600">
          Monitor Maximal Extractable Value (MEV) activities and analyze suspicious transactions
        </p>
      </div>

      {/* MEV Statistics */}
      {mevStats && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-sm font-medium text-gray-500 mb-2">24h MEV Extracted</h3>
            <p className="text-2xl font-bold text-gray-900">
              {formatEther(mevStats.total_mev_extracted_24h)} ETH
            </p>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-sm font-medium text-gray-500 mb-2">MEV Transactions</h3>
            <p className="text-2xl font-bold text-gray-900">
              {mevStats.mev_transactions_count.toLocaleString()}
            </p>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-sm font-medium text-gray-500 mb-2">Avg MEV per Tx</h3>
            <p className="text-2xl font-bold text-gray-900">
              {formatEther(mevStats.average_mev_per_transaction)} ETH
            </p>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-sm font-medium text-gray-500 mb-2">Active MEV Bots</h3>
            <p className="text-2xl font-bold text-gray-900">
              {mevStats.most_active_mev_addresses.length}
            </p>
          </div>
        </div>
      )}

      {/* MEV by Type */}
      {mevStats && mevStats.mev_by_type.length > 0 && (
        <div className="bg-white p-6 rounded-lg shadow mb-8">
          <h2 className="text-xl font-bold text-gray-900 mb-4">MEV Activity by Type</h2>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
            {mevStats.mev_by_type.map((type) => (
              <div key={type.type} className="text-center">
                <div className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${getMEVTypeColor(type.type)} mb-2`}>
                  {type.type.replace('_', ' ').toUpperCase()}
                </div>
                <p className="text-sm text-gray-600">{type.count} transactions</p>
                <p className="text-sm font-medium text-gray-900">
                  {formatEther(type.total_value)} ETH
                </p>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Filter and Activities */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-6 border-b border-gray-200">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
            <h2 className="text-xl font-bold text-gray-900 mb-4 sm:mb-0">Recent MEV Activities</h2>
            <div className="flex items-center space-x-4">
              <select
                value={selectedType}
                onChange={(e) => setSelectedType(e.target.value)}
                className="border border-gray-300 rounded-md px-3 py-2 bg-white text-sm"
              >
                <option value="all">All Types</option>
                <option value="sandwich_attack">Sandwich Attacks</option>
                <option value="front_running">Front Running</option>
                <option value="back_running">Back Running</option>
                <option value="arbitrage">Arbitrage</option>
                <option value="liquidation">Liquidations</option>
                <option value="suspicious">Suspicious</option>
              </select>
            </div>
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Transaction
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Type
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Value Extracted
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Confidence
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Block
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Time
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {mevActivities.length > 0 ? (
                mevActivities.map((activity) => (
                  <tr key={activity.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <a
                        href={`/transactions/${activity.transaction_hash}`}
                        className="text-blue-600 hover:text-blue-800 font-mono text-sm"
                      >
                        {activity.transaction_hash.slice(0, 12)}...{activity.transaction_hash.slice(-8)}
                      </a>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getMEVTypeColor(activity.mev_type)}`}>
                        {activity.mev_type.replace('_', ' ')}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {formatEther(activity.value_extracted)} ETH
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <span className={getConfidenceColor(activity.confidence_score)}>
                        {(activity.confidence_score * 100).toFixed(1)}%
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      <a
                        href={`/blocks/${activity.block_number}`}
                        className="text-blue-600 hover:text-blue-800"
                      >
                        {activity.block_number.toLocaleString()}
                      </a>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(activity.timestamp).toLocaleString()}
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={6} className="px-6 py-4 text-center text-gray-500">
                    No MEV activities found
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Most Active MEV Addresses */}
      {mevStats && mevStats.most_active_mev_addresses.length > 0 && (
        <div className="bg-white p-6 rounded-lg shadow mt-8">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Most Active MEV Addresses</h2>
          <div className="space-y-4">
            {mevStats.most_active_mev_addresses.slice(0, 5).map((address, index) => (
              <div key={address.address} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                <div className="flex items-center space-x-4">
                  <div className="flex-shrink-0">
                    <div className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center text-sm font-medium text-gray-600">
                      #{index + 1}
                    </div>
                  </div>
                  <div>
                    <a
                      href={`/addresses/${address.address}`}
                      className="text-blue-600 hover:text-blue-800 font-mono text-sm"
                    >
                      {address.address.slice(0, 12)}...{address.address.slice(-8)}
                    </a>
                    <p className="text-sm text-gray-500">{address.mev_count} MEV transactions</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-sm font-medium text-gray-900">
                    {formatEther(address.total_extracted)} ETH
                  </p>
                  <p className="text-sm text-gray-500">Total extracted</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
