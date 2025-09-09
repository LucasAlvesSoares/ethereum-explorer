'use client'

import { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { ArrowLeft, Hash, User, Zap, Database, Activity, Copy, CheckCircle, XCircle, AlertCircle, FileText } from 'lucide-react'
import { Transaction } from '@/types'
import { formatHash, formatNumber, formatTimestamp, formatValue, formatGasPrice } from '@/utils/formatting'
import { ErrorState, handleFetchError, handleNetworkError, createSuccessState } from '@/utils/errors'

export default function TransactionDetailPage() {
  const params = useParams()
  const router = useRouter()
  const txHash = params.hash as string

  const [transaction, setTransaction] = useState<Transaction | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<ErrorState>(createSuccessState())
  const [copiedField, setCopiedField] = useState<string | null>(null)

  const fetchTransaction = async () => {
    try {
      setLoading(true)
      const response = await fetch(`/api/v1/transactions/${txHash}`)
      
      if (!response.ok) {
        const errorState = await handleFetchError(response, 'fetch transaction')
        setError(errorState)
        return
      }
      
      const txData: Transaction = await response.json()
      setTransaction(txData)
      setError(createSuccessState())
    } catch (err) {
      const errorState = handleNetworkError('fetch transaction')
      setError(errorState)
    } finally {
      setLoading(false)
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

  const getStatusIcon = (status: number | null) => {
    if (status === null) return <AlertCircle className="w-5 h-5 text-yellow-500" />
    if (status === 1) return <CheckCircle className="w-5 h-5 text-green-500" />
    return <XCircle className="w-5 h-5 text-red-500" />
  }

  const getStatusText = (status: number | null) => {
    if (status === null) return 'Pending'
    return status === 1 ? 'Success' : 'Failed'
  }

  const getStatusColor = (status: number | null) => {
    if (status === null) return 'text-yellow-700 bg-yellow-50 border-yellow-200'
    if (status === 1) return 'text-green-700 bg-green-50 border-green-200'
    return 'text-red-700 bg-red-50 border-red-200'
  }

  useEffect(() => {
    if (txHash) {
      fetchTransaction()
    }
  }, [txHash])

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
          <h1 className="text-3xl font-bold text-gray-900">Transaction Details</h1>
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
          <h1 className="text-3xl font-bold text-gray-900">Transaction Details</h1>
        </div>
        <div className="card">
          <div className="text-center py-8">
            <p className="text-red-600 mb-4">Error: {error.message}</p>
            {error.details && (
              <p className="text-gray-600 mb-4 text-sm">{error.details}</p>
            )}
            <button
              onClick={fetchTransaction}
              className="btn-primary"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>
    )
  }

  if (!transaction) {
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
          <h1 className="text-3xl font-bold text-gray-900">Transaction Not Found</h1>
        </div>
        <div className="card">
          <div className="text-center py-8">
            <p className="text-gray-600">Transaction {formatHash(txHash)} was not found.</p>
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
          Back to Transactions
        </button>
        <h1 className="text-3xl font-bold text-gray-900">Transaction Details</h1>
        <p className="text-gray-600 mt-2">Detailed information about this transaction</p>
      </div>

      {/* Transaction Status */}
      <div className={`card mb-8 border-l-4 ${getStatusColor(transaction.status)}`}>
        <div className="flex items-center">
          {getStatusIcon(transaction.status)}
          <div className="ml-3">
            <h3 className="text-lg font-medium">
              Transaction {getStatusText(transaction.status)}
            </h3>
            <p className="text-sm opacity-75">
              {transaction.status === 1 && 'This transaction was successfully executed.'}
              {transaction.status === 0 && 'This transaction failed during execution.'}
              {transaction.status === null && 'This transaction is pending confirmation.'}
            </p>
          </div>
        </div>
      </div>

      {/* Transaction Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <div className="flex items-center">
            <Hash className="w-8 h-8 text-primary-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">Block Number</p>
              <Link
                href={`/blocks/${transaction.block_number}`}
                className="text-2xl font-bold text-primary-600 hover:text-primary-800"
              >
                {formatNumber(transaction.block_number)}
              </Link>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <Database className="w-8 h-8 text-green-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">Transaction Index</p>
              <p className="text-2xl font-bold text-gray-900">{transaction.transaction_index}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <Zap className="w-8 h-8 text-yellow-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">Gas Used</p>
              <p className="text-2xl font-bold text-gray-900">
                {transaction.gas_used ? formatNumber(transaction.gas_used) : 'Pending'}
              </p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <Activity className="w-8 h-8 text-blue-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">Value</p>
              <p className="text-2xl font-bold text-gray-900">{formatValue(transaction.value)}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Transaction Details */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
        {/* Basic Information */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Transaction Information</h2>
          <div className="space-y-4">
            <div className="flex justify-between items-start">
              <span className="text-sm text-gray-600 font-medium">Transaction Hash:</span>
              <div className="flex items-center">
                <span className="text-sm font-mono text-gray-900 mr-2">{formatHash(transaction.hash)}</span>
                <button
                  onClick={() => copyToClipboard(transaction.hash, 'hash')}
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
              <span className="text-sm text-gray-600 font-medium">From:</span>
              <div className="flex items-center">
                <span className="text-sm font-mono text-gray-900 mr-2">{formatHash(transaction.from_address)}</span>
                <button
                  onClick={() => copyToClipboard(transaction.from_address, 'from')}
                  className="text-gray-400 hover:text-gray-600"
                >
                  {copiedField === 'from' ? (
                    <CheckCircle className="w-4 h-4 text-green-600" />
                  ) : (
                    <Copy className="w-4 h-4" />
                  )}
                </button>
              </div>
            </div>

            <div className="flex justify-between items-start">
              <span className="text-sm text-gray-600 font-medium">To:</span>
              <div className="flex items-center">
                {transaction.to_address ? (
                  <>
                    <span className="text-sm font-mono text-gray-900 mr-2">{formatHash(transaction.to_address)}</span>
                    <button
                      onClick={() => copyToClipboard(transaction.to_address!, 'to')}
                      className="text-gray-400 hover:text-gray-600"
                    >
                      {copiedField === 'to' ? (
                        <CheckCircle className="w-4 h-4 text-green-600" />
                      ) : (
                        <Copy className="w-4 h-4" />
                      )}
                    </button>
                  </>
                ) : (
                  <span className="text-sm text-gray-500 italic">Contract Creation</span>
                )}
              </div>
            </div>

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Value:</span>
              <span className="text-sm text-gray-900 font-mono">{formatValue(transaction.value)}</span>
            </div>

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Nonce:</span>
              <span className="text-sm text-gray-900">{transaction.nonce}</span>
            </div>

            {transaction.contract_address && (
              <div className="flex justify-between items-start">
                <span className="text-sm text-gray-600 font-medium">Contract Created:</span>
                <div className="flex items-center">
                  <span className="text-sm font-mono text-gray-900 mr-2">{formatHash(transaction.contract_address)}</span>
                  <button
                    onClick={() => copyToClipboard(transaction.contract_address!, 'contract')}
                    className="text-gray-400 hover:text-gray-600"
                  >
                    {copiedField === 'contract' ? (
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

        {/* Gas Information */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Gas Information</h2>
          <div className="space-y-4">
            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Gas Limit:</span>
              <span className="text-sm text-gray-900">{formatNumber(transaction.gas_limit)}</span>
            </div>

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Gas Used:</span>
              <span className="text-sm text-gray-900">
                {transaction.gas_used ? formatNumber(transaction.gas_used) : 'Pending'}
              </span>
            </div>

            {transaction.gas_used && (
              <div className="mt-4">
                <div className="flex justify-between text-sm text-gray-600 mb-2">
                  <span>Gas Usage</span>
                  <span>{((transaction.gas_used / transaction.gas_limit) * 100).toFixed(2)}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-3">
                  <div
                    className="bg-primary-600 h-3 rounded-full transition-all duration-300"
                    style={{
                      width: `${(transaction.gas_used / transaction.gas_limit) * 100}%`
                    }}
                  ></div>
                </div>
              </div>
            )}

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Gas Price:</span>
              <span className="text-sm text-gray-900">{formatGasPrice(transaction.gas_price)}</span>
            </div>

            {transaction.max_fee_per_gas && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-600 font-medium">Max Fee Per Gas:</span>
                <span className="text-sm text-gray-900">{formatGasPrice(transaction.max_fee_per_gas)}</span>
              </div>
            )}

            {transaction.max_priority_fee_per_gas && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-600 font-medium">Max Priority Fee:</span>
                <span className="text-sm text-gray-900">{formatGasPrice(transaction.max_priority_fee_per_gas)}</span>
              </div>
            )}

            {transaction.gas_used && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-600 font-medium">Transaction Fee:</span>
                <span className="text-sm text-gray-900 font-mono">
                  {formatValue((BigInt(transaction.gas_used) * BigInt(transaction.gas_price)).toString())}
                </span>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Input Data */}
      {transaction.input_data && transaction.input_data !== '0x' && (
        <div className="card mb-8">
          <h2 className="text-xl font-semibold text-gray-900 mb-6 flex items-center">
            <FileText className="w-5 h-5 mr-2" />
            Input Data
          </h2>
          <div className="bg-gray-50 rounded-lg p-4">
            <div className="flex justify-between items-start mb-2">
              <span className="text-sm text-gray-600 font-medium">Hex Data:</span>
              <button
                onClick={() => copyToClipboard(transaction.input_data, 'input')}
                className="text-gray-400 hover:text-gray-600"
              >
                {copiedField === 'input' ? (
                  <CheckCircle className="w-4 h-4 text-green-600" />
                ) : (
                  <Copy className="w-4 h-4" />
                )}
              </button>
            </div>
            <div className="bg-white rounded border p-3 max-h-40 overflow-y-auto">
              <code className="text-xs font-mono text-gray-800 break-all">
                {transaction.input_data}
              </code>
            </div>
            <p className="text-xs text-gray-500 mt-2">
              {Math.floor(transaction.input_data.length / 2)} bytes
            </p>
          </div>
        </div>
      )}

      {/* Additional Information */}
      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 mb-6">Additional Information</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-4">
            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Created At:</span>
              <span className="text-sm text-gray-900">{formatTimestamp(transaction.created_at)}</span>
            </div>

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Updated At:</span>
              <span className="text-sm text-gray-900">{formatTimestamp(transaction.updated_at)}</span>
            </div>

            {transaction.logs_bloom && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-600 font-medium">Logs Bloom:</span>
                <span className="text-sm text-gray-500">Present</span>
              </div>
            )}
          </div>

          <div className="space-y-4">
            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Transaction Type:</span>
              <span className="text-sm text-gray-900">
                {transaction.to_address ? 'Transfer' : 'Contract Creation'}
              </span>
            </div>

            {transaction.input_data && transaction.input_data !== '0x' && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-600 font-medium">Has Input Data:</span>
                <span className="text-sm text-gray-900">Yes ({Math.floor(transaction.input_data.length / 2)} bytes)</span>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
