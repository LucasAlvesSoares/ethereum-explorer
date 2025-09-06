'use client'

import { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { ArrowLeft, Clock, Hash, User, Zap, Database, Activity, Copy, CheckCircle } from 'lucide-react'
import { Block, Transaction } from '@/types'
import { formatHash, formatNumber, formatTimestamp, formatBytes, formatPercentage, formatValue, formatGasPrice } from '@/utils/formatting'
import { ErrorState, handleFetchError, handleNetworkError, createSuccessState } from '@/utils/errors'

export default function BlockDetailPage() {
  const params = useParams()
  const router = useRouter()
  const blockId = params.id as string

  const [block, setBlock] = useState<Block | null>(null)
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)
  const [transactionsLoading, setTransactionsLoading] = useState(false)
  const [error, setError] = useState<ErrorState>(createSuccessState())
  const [copiedField, setCopiedField] = useState<string | null>(null)

  const fetchBlock = async () => {
    try {
      setLoading(true)
      const response = await fetch(`/api/v1/blocks/${blockId}`)
      
      if (!response.ok) {
        const errorState = await handleFetchError(response, 'fetch block')
        setError(errorState)
        return
      }
      
      const blockData: Block = await response.json()
      setBlock(blockData)
      setError(createSuccessState())
      
      // Fetch transactions for this block
      fetchBlockTransactions()
    } catch (err) {
      const errorState = handleNetworkError('fetch block')
      setError(errorState)
    } finally {
      setLoading(false)
    }
  }

  const fetchBlockTransactions = async () => {
    try {
      setTransactionsLoading(true)
      const response = await fetch(`/api/v1/transactions?block=${blockId}&limit=10`)
      
      if (response.ok) {
        const data = await response.json()
        setTransactions(data.transactions || [])
      }
    } catch (err) {
      console.error('Failed to fetch block transactions:', err)
    } finally {
      setTransactionsLoading(false)
    }
  }

  const copyToClipboard = async (text: string, field: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopiedField(field)
      setTimeout(() => setCopiedField(null), 2000)
    } catch (err) {
      console.error('Failed to copy to clipboard:', err)
    }
  }

  useEffect(() => {
    if (blockId) {
      fetchBlock()
    }
  }, [blockId])

  if (loading) {
    return (
      <div className="px-4 sm:px-0">
        <div className="mb-8">
          <button
            onClick={() => router.back()}
            className="flex items-center text-primary-600 hover:text-primary-800 mb-4"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back
          </button>
          <h1 className="text-3xl font-bold text-gray-900">Block Details</h1>
        </div>
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        </div>
      </div>
    )
  }

  if (error.isError) {
    return (
      <div className="px-4 sm:px-0">
        <div className="mb-8">
          <button
            onClick={() => router.back()}
            className="flex items-center text-primary-600 hover:text-primary-800 mb-4"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back
          </button>
          <h1 className="text-3xl font-bold text-gray-900">Block Details</h1>
        </div>
        <div className="card">
          <div className="text-center py-8">
            <p className="text-red-600 mb-4">Error: {error.message}</p>
            {error.details && (
              <p className="text-gray-600 mb-4 text-sm">{error.details}</p>
            )}
            <button
              onClick={fetchBlock}
              className="btn-primary"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>
    )
  }

  if (!block) {
    return (
      <div className="px-4 sm:px-0">
        <div className="mb-8">
          <button
            onClick={() => router.back()}
            className="flex items-center text-primary-600 hover:text-primary-800 mb-4"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back
          </button>
          <h1 className="text-3xl font-bold text-gray-900">Block Not Found</h1>
        </div>
        <div className="card">
          <div className="text-center py-8">
            <p className="text-gray-600">Block {blockId} was not found.</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="px-4 sm:px-0">
      {/* Header */}
      <div className="mb-8">
        <button
          onClick={() => router.back()}
          className="flex items-center text-primary-600 hover:text-primary-800 mb-4"
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back to Blocks
        </button>
        <h1 className="text-3xl font-bold text-gray-900">Block {formatNumber(block.number)}</h1>
        <p className="text-gray-600 mt-2">Detailed information about this block</p>
      </div>

      {/* Block Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <div className="flex items-center">
            <Database className="w-8 h-8 text-primary-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">Block Height</p>
              <p className="text-2xl font-bold text-gray-900">{formatNumber(block.number)}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <Activity className="w-8 h-8 text-green-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">Transactions</p>
              <p className="text-2xl font-bold text-gray-900">{formatNumber(block.transaction_count)}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <Zap className="w-8 h-8 text-yellow-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">Gas Used</p>
              <p className="text-2xl font-bold text-gray-900">
                {formatPercentage(block.gas_used, block.gas_limit)}
              </p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <Database className="w-8 h-8 text-blue-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">Size</p>
              <p className="text-2xl font-bold text-gray-900">{formatBytes(block.size)}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Block Details */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
        {/* Basic Information */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Block Information</h2>
          <div className="space-y-4">
            <div className="flex justify-between items-start">
              <span className="text-sm text-gray-600 font-medium">Block Hash:</span>
              <div className="flex items-center">
                <span className="text-sm font-mono text-gray-900 mr-2">{formatHash(block.hash)}</span>
                <button
                  onClick={() => copyToClipboard(block.hash, 'hash')}
                  className="text-gray-400 hover:text-gray-600"
                >
                  {copiedField === 'hash' ? (
                    <CheckCircle className="w-4 h-4 text-green-600" />
                  ) : (
                    <Copy className="w-4 h-4" />
                  )}
                </button>
              </div>
            </div>

            <div className="flex justify-between items-start">
              <span className="text-sm text-gray-600 font-medium">Parent Hash:</span>
              <div className="flex items-center">
                <Link
                  href={`/blocks/${parseInt(block.number.toString()) - 1}`}
                  className="text-sm font-mono text-primary-600 hover:text-primary-800 mr-2"
                >
                  {formatHash(block.parent_hash)}
                </Link>
                <button
                  onClick={() => copyToClipboard(block.parent_hash, 'parent_hash')}
                  className="text-gray-400 hover:text-gray-600"
                >
                  {copiedField === 'parent_hash' ? (
                    <CheckCircle className="w-4 h-4 text-green-600" />
                  ) : (
                    <Copy className="w-4 h-4" />
                  )}
                </button>
              </div>
            </div>

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Timestamp:</span>
              <span className="text-sm text-gray-900">{formatTimestamp(block.timestamp)}</span>
            </div>

            <div className="flex justify-between items-start">
              <span className="text-sm text-gray-600 font-medium">Miner:</span>
              <div className="flex items-center">
                <span className="text-sm font-mono text-gray-900 mr-2">{formatHash(block.miner)}</span>
                <button
                  onClick={() => copyToClipboard(block.miner, 'miner')}
                  className="text-gray-400 hover:text-gray-600"
                >
                  {copiedField === 'miner' ? (
                    <CheckCircle className="w-4 h-4 text-green-600" />
                  ) : (
                    <Copy className="w-4 h-4" />
                  )}
                </button>
              </div>
            </div>

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Difficulty:</span>
              <span className="text-sm text-gray-900">{block.difficulty}</span>
            </div>

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Total Difficulty:</span>
              <span className="text-sm text-gray-900">{block.total_difficulty}</span>
            </div>
          </div>
        </div>

        {/* Gas Information */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Gas Information</h2>
          <div className="space-y-4">
            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Gas Used:</span>
              <span className="text-sm text-gray-900">{formatNumber(block.gas_used)}</span>
            </div>

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Gas Limit:</span>
              <span className="text-sm text-gray-900">{formatNumber(block.gas_limit)}</span>
            </div>

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Gas Utilization:</span>
              <span className="text-sm text-gray-900">
                {formatPercentage(block.gas_used, block.gas_limit)}
              </span>
            </div>

            {/* Gas Usage Bar */}
            <div className="mt-4">
              <div className="flex justify-between text-sm text-gray-600 mb-2">
                <span>Gas Usage</span>
                <span>{formatPercentage(block.gas_used, block.gas_limit)}</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-3">
                <div
                  className="bg-primary-600 h-3 rounded-full transition-all duration-300"
                  style={{
                    width: formatPercentage(block.gas_used, block.gas_limit)
                  }}
                ></div>
              </div>
            </div>

            {block.base_fee_per_gas && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-600 font-medium">Base Fee:</span>
                <span className="text-sm text-gray-900">{formatGasPrice(block.base_fee_per_gas)}</span>
              </div>
            )}

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Block Size:</span>
              <span className="text-sm text-gray-900">{formatBytes(block.size)}</span>
            </div>

            {block.extra_data && block.extra_data !== '0x' && (
              <div className="flex justify-between items-start">
                <span className="text-sm text-gray-600 font-medium">Extra Data:</span>
                <div className="flex items-center">
                  <span className="text-sm font-mono text-gray-900 mr-2 max-w-32 truncate">
                    {block.extra_data}
                  </span>
                  <button
                    onClick={() => copyToClipboard(block.extra_data, 'extra_data')}
                    className="text-gray-400 hover:text-gray-600"
                  >
                    {copiedField === 'extra_data' ? (
                      <CheckCircle className="w-4 h-4 text-green-600" />
                    ) : (
                      <Copy className="w-4 h-4" />
                    )}
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Block Transactions */}
      <div className="card">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-semibold text-gray-900">
            Transactions ({formatNumber(block.transaction_count)})
          </h2>
          {block.transaction_count > 10 && (
            <Link
              href={`/transactions?block=${block.number}`}
              className="text-primary-600 hover:text-primary-800 text-sm font-medium"
            >
              View all transactions →
            </Link>
          )}
        </div>

        {transactionsLoading ? (
          <div className="flex justify-center items-center h-32">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          </div>
        ) : transactions.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Transaction Hash
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    From
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    To
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Value
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Gas Used
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {transactions.slice(0, 10).map((tx) => (
                  <tr key={tx.hash} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <Link
                        href={`/transactions/${tx.hash}`}
                        className="text-primary-600 hover:text-primary-800 font-mono text-sm"
                      >
                        {formatHash(tx.hash)}
                      </Link>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className="text-gray-900 font-mono text-sm">
                        {formatHash(tx.from_address)}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className="text-gray-900 font-mono text-sm">
                        {tx.to_address ? formatHash(tx.to_address) : 'Contract Creation'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {formatValue(tx.value)} ETH
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {tx.gas_used ? formatNumber(tx.gas_used) : 'N/A'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="text-center py-8">
            <Activity className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No transactions found</h3>
            <p className="text-gray-600">This block contains no transactions.</p>
          </div>
        )}
      </div>
    </div>
  )
}
