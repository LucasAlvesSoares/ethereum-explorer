'use client'

import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { ArrowRight, Hash, User, Zap, CheckCircle, XCircle } from 'lucide-react'

interface Transaction {
  hash: string
  block_number: number
  transaction_index: number
  from_address: string
  to_address: string | null
  value: string
  gas_limit: number
  gas_used: number | null
  gas_price: string
  max_fee_per_gas?: string
  max_priority_fee_per_gas?: string
  nonce: number
  input_data: string
  status: number | null
  contract_address: string | null
  logs_bloom: string | null
  created_at: string
  updated_at: string
}

interface TransactionsResponse {
  transactions: Transaction[]
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}

export default function TransactionsPage() {
  const searchParams = useSearchParams()
  const blockNumber = searchParams.get('block')
  
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 20,
    total: 0,
    total_pages: 0
  })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchTransactions = async (page: number = 1) => {
    try {
      setLoading(true)
      let url = `/api/v1/transactions?page=${page}&limit=20`
      if (blockNumber) {
        url += `&block=${blockNumber}`
      }
      
      const response = await fetch(url)
      
      if (!response.ok) {
        throw new Error('Failed to fetch transactions')
      }
      
      const data: TransactionsResponse = await response.json()
      setTransactions(data.transactions || []) // Ensure we always have an array
      setPagination(data.pagination || { page: 1, limit: 20, total: 0, total_pages: 0 })
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchTransactions()
  }, [blockNumber])

  const formatHash = (hash: string) => {
    return `${hash.slice(0, 10)}...${hash.slice(-8)}`
  }

  const formatNumber = (num: number) => {
    return num.toLocaleString()
  }

  const formatValue = (value: string) => {
    // Convert wei to ETH (1 ETH = 10^18 wei)
    const eth = parseFloat(value) / Math.pow(10, 18)
    if (eth === 0) return '0 ETH'
    if (eth < 0.001) return '<0.001 ETH'
    return `${eth.toFixed(6)} ETH`
  }

  const formatGasPrice = (gasPrice: string) => {
    // Convert wei to gwei (1 gwei = 10^9 wei)
    const gwei = parseFloat(gasPrice) / Math.pow(10, 9)
    return `${gwei.toFixed(2)} gwei`
  }

  const getStatusIcon = (status: number | null) => {
    if (status === null) return null
    if (status === 1) {
      return <CheckCircle className="w-4 h-4 text-green-500" />
    } else {
      return <XCircle className="w-4 h-4 text-red-500" />
    }
  }

  const getStatusText = (status: number | null) => {
    if (status === null) return 'Pending'
    return status === 1 ? 'Success' : 'Failed'
  }

  if (loading && transactions.length === 0) {
    return (
      <div className="px-4 sm:px-0">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">
            {blockNumber ? `Block ${blockNumber} Transactions` : 'Transactions'}
          </h1>
          <p className="text-gray-600 mt-2">
            {blockNumber 
              ? `All transactions in block ${blockNumber}`
              : 'Latest transactions on the Ethereum blockchain'
            }
          </p>
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
          <h1 className="text-3xl font-bold text-gray-900">
            {blockNumber ? `Block ${blockNumber} Transactions` : 'Transactions'}
          </h1>
          <p className="text-gray-600 mt-2">
            {blockNumber 
              ? `All transactions in block ${blockNumber}`
              : 'Latest transactions on the Ethereum blockchain'
            }
          </p>
        </div>
        <div className="card">
          <div className="text-center py-8">
            <p className="text-red-600 mb-4">Error: {error}</p>
            <button
              onClick={() => fetchTransactions(pagination.page)}
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
        <h1 className="text-3xl font-bold text-gray-900">
          {blockNumber ? `Block ${blockNumber} Transactions` : 'Transactions'}
        </h1>
        <p className="text-gray-600 mt-2">
          {blockNumber 
            ? `All transactions in block ${blockNumber}`
            : 'Latest transactions on the Ethereum blockchain'
          }
        </p>
        {blockNumber && (
          <div className="mt-4">
            <Link
              href="/transactions"
              className="text-primary-600 hover:text-primary-800 text-sm"
            >
              ← View all transactions
            </Link>
          </div>
        )}
      </div>

      {/* Transactions Table */}
      <div className="card overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Transaction Hash
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Block
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  From → To
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Value
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Gas
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {transactions.map((tx) => (
                <tr key={tx.hash} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <Hash className="w-4 h-4 mr-2 text-gray-400" />
                      <Link
                        href={`/transactions/${tx.hash}`}
                        className="text-primary-600 hover:text-primary-800 font-mono text-sm"
                      >
                        {formatHash(tx.hash)}
                      </Link>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <Link
                      href={`/blocks/${tx.block_number}`}
                      className="text-primary-600 hover:text-primary-800"
                    >
                      {formatNumber(tx.block_number)}
                    </Link>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center space-x-2">
                      <div className="flex items-center">
                        <User className="w-3 h-3 mr-1 text-gray-400" />
                        <Link
                          href={`/addresses/${tx.from_address}`}
                          className="text-primary-600 hover:text-primary-800 font-mono text-xs"
                        >
                          {formatHash(tx.from_address)}
                        </Link>
                      </div>
                      <ArrowRight className="w-3 h-3 text-gray-400" />
                      <div className="flex items-center">
                        <User className="w-3 h-3 mr-1 text-gray-400" />
                        {tx.to_address ? (
                          <Link
                            href={`/addresses/${tx.to_address}`}
                            className="text-primary-600 hover:text-primary-800 font-mono text-xs"
                          >
                            {formatHash(tx.to_address)}
                          </Link>
                        ) : (
                          <span className="text-gray-500 text-xs">Contract Creation</span>
                        )}
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className="text-sm text-gray-900 font-mono">
                      {formatValue(tx.value)}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex flex-col">
                      <div className="flex items-center">
                        <Zap className="w-4 h-4 mr-2 text-gray-400" />
                        <span className="text-sm text-gray-900">
                          {tx.gas_used ? formatNumber(tx.gas_used) : 'Pending'}
                        </span>
                      </div>
                      <div className="text-xs text-gray-500">
                        {formatGasPrice(tx.gas_price)}
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      {getStatusIcon(tx.status)}
                      <span className={`ml-2 text-sm ${
                        tx.status === 1 ? 'text-green-700' : 
                        tx.status === 0 ? 'text-red-700' : 'text-yellow-700'
                      }`}>
                        {getStatusText(tx.status)}
                      </span>
                    </div>
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
                onClick={() => fetchTransactions(pagination.page - 1)}
                disabled={pagination.page <= 1}
                className="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Previous
              </button>
              <button
                onClick={() => fetchTransactions(pagination.page + 1)}
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
                  {' '}({formatNumber(pagination.total)} total transactions)
                </p>
              </div>
              <div>
                <nav className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                  <button
                    onClick={() => fetchTransactions(pagination.page - 1)}
                    disabled={pagination.page <= 1}
                    className="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Previous
                  </button>
                  <button
                    onClick={() => fetchTransactions(pagination.page + 1)}
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
