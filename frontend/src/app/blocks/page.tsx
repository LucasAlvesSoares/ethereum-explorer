'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { Clock, Hash, User, Zap } from 'lucide-react'

interface Block {
  number: number
  hash: string
  parent_hash: string
  timestamp: string
  gas_limit: number
  gas_used: number
  difficulty: string
  total_difficulty: string
  size: number
  transaction_count: number
  miner: string
  extra_data: string
  base_fee_per_gas?: string
  created_at: string
  updated_at: string
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

export default function BlocksPage() {
  const [blocks, setBlocks] = useState<Block[]>([])
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 20,
    total: 0,
    total_pages: 0
  })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchBlocks = async (page: number = 1) => {
    try {
      setLoading(true)
      const response = await fetch(`/api/v1/blocks?page=${page}&limit=20`)
      
      if (!response.ok) {
        throw new Error('Failed to fetch blocks')
      }
      
      const data: BlocksResponse = await response.json()
      setBlocks(data.blocks || []) // Ensure we always have an array
      setPagination(data.pagination || { page: 1, limit: 20, total: 0, total_pages: 0 })
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchBlocks()
  }, [])

  const formatTimestamp = (timestamp: string) => {
    return new Date(timestamp).toLocaleString()
  }

  const formatHash = (hash: string) => {
    return `${hash.slice(0, 10)}...${hash.slice(-8)}`
  }

  const formatNumber = (num: number) => {
    return num.toLocaleString()
  }

  const calculateGasUsedPercentage = (gasUsed: number, gasLimit: number) => {
    return ((gasUsed / gasLimit) * 100).toFixed(1)
  }

  if (loading && blocks.length === 0) {
    return (
      <div className="px-4 sm:px-0">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Blocks</h1>
          <p className="text-gray-600 mt-2">Latest blocks on the Ethereum blockchain</p>
        </div>
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="px-4 sm:px-0">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Blocks</h1>
          <p className="text-gray-600 mt-2">Latest blocks on the Ethereum blockchain</p>
        </div>
        <div className="card">
          <div className="text-center py-8">
            <p className="text-red-600 mb-4">Error: {error}</p>
            <button
              onClick={() => fetchBlocks(pagination.page)}
              className="btn-primary"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="px-4 sm:px-0">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Blocks</h1>
        <p className="text-gray-600 mt-2">Latest blocks on the Ethereum blockchain</p>
      </div>

      {/* Blocks Table */}
      <div className="card overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Block
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Age
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Transactions
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Miner
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Gas Used
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Size
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {blocks.map((block) => (
                <tr key={block.number} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex flex-col">
                      <Link
                        href={`/blocks/${block.number}`}
                        className="text-primary-600 hover:text-primary-800 font-medium"
                      >
                        {formatNumber(block.number)}
                      </Link>
                      <div className="flex items-center text-sm text-gray-500 mt-1">
                        <Hash className="w-3 h-3 mr-1" />
                        {formatHash(block.hash)}
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center text-sm text-gray-900">
                      <Clock className="w-4 h-4 mr-2 text-gray-400" />
                      {formatTimestamp(block.timestamp)}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <Link
                      href={`/transactions?block=${block.number}`}
                      className="text-primary-600 hover:text-primary-800"
                    >
                      {formatNumber(block.transaction_count)}
                    </Link>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <User className="w-4 h-4 mr-2 text-gray-400" />
                      <Link
                        href={`/addresses/${block.miner}`}
                        className="text-primary-600 hover:text-primary-800 text-sm"
                      >
                        {formatHash(block.miner)}
                      </Link>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex flex-col">
                      <div className="flex items-center">
                        <Zap className="w-4 h-4 mr-2 text-gray-400" />
                        <span className="text-sm text-gray-900">
                          {formatNumber(block.gas_used)}
                        </span>
                      </div>
                      <div className="text-xs text-gray-500">
                        {calculateGasUsedPercentage(block.gas_used, block.gas_limit)}% of {formatNumber(block.gas_limit)}
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-1 mt-1">
                        <div
                          className="bg-primary-600 h-1 rounded-full"
                          style={{
                            width: `${calculateGasUsedPercentage(block.gas_used, block.gas_limit)}%`
                          }}
                        ></div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {(block.size / 1024).toFixed(1)} KB
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {pagination.total_pages > 1 && (
          <div className="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
            <div className="flex-1 flex justify-between sm:hidden">
              <button
                onClick={() => fetchBlocks(pagination.page - 1)}
                disabled={pagination.page <= 1}
                className="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Previous
              </button>
              <button
                onClick={() => fetchBlocks(pagination.page + 1)}
                disabled={pagination.page >= pagination.total_pages}
                className="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Next
              </button>
            </div>
            <div className="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
              <div>
                <p className="text-sm text-gray-700">
                  Showing page <span className="font-medium">{pagination.page}</span> of{' '}
                  <span className="font-medium">{pagination.total_pages}</span> pages
                  {' '}({formatNumber(pagination.total)} total blocks)
                </p>
              </div>
              <div>
                <nav className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                  <button
                    onClick={() => fetchBlocks(pagination.page - 1)}
                    disabled={pagination.page <= 1}
                    className="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Previous
                  </button>
                  <button
                    onClick={() => fetchBlocks(pagination.page + 1)}
                    disabled={pagination.page >= pagination.total_pages}
                    className="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Next
                  </button>
                </nav>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
