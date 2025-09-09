'use client'

import { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { ArrowLeft, Wallet, Activity, Hash, Clock, Copy, CheckCircle, Tag, Shield, Zap, TrendingUp, Users } from 'lucide-react'
import { Address, Transaction, TransactionsResponse } from '@/types'
import { formatHash, formatNumber, formatTimestamp, formatValue, formatGasPrice } from '@/utils/formatting'
import { ErrorState, handleFetchError, handleNetworkError, createSuccessState } from '@/utils/errors'

export default function AddressDetailPage() {
  const params = useParams()
  const router = useRouter()
  const addressHash = params.address as string

  const [address, setAddress] = useState<Address | null>(null)
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)
  const [transactionsLoading, setTransactionsLoading] = useState(false)
  const [error, setError] = useState<ErrorState>(createSuccessState())
  const [copiedField, setCopiedField] = useState<string | null>(null)
  const [transactionPage, setTransactionPage] = useState(1)
  const [hasMoreTransactions, setHasMoreTransactions] = useState(true)

  const fetchAddress = async () => {
    try {
      setLoading(true)
      const response = await fetch(`/api/v1/addresses/${addressHash}`)
      
      if (!response.ok) {
        const errorState = await handleFetchError(response, 'fetch address')
        setError(errorState)
        return
      }
      
      const addressData: Address = await response.json()
      setAddress(addressData)
      setError(createSuccessState())
      
      // Fetch transactions for this address
      fetchAddressTransactions()
    } catch (err) {
      const errorState = handleNetworkError('fetch address')
      setError(errorState)
    } finally {
      setLoading(false)
    }
  }

  const fetchAddressTransactions = async (page = 1) => {
    try {
      setTransactionsLoading(true)
      const response = await fetch(`/api/v1/addresses/${addressHash}/transactions?page=${page}&limit=10`)
      
      if (response.ok) {
        const data: TransactionsResponse = await response.json()
        if (page === 1) {
          setTransactions(data.transactions || [])
        } else {
          setTransactions(prev => [...prev, ...(data.transactions || [])])
        }
        setHasMoreTransactions(data.pagination.page < data.pagination.total_pages)
      }
    } catch (err) {
      console.error('Failed to fetch address transactions:', err)
    } finally {
      setTransactionsLoading(false)
    }
  }

  const loadMoreTransactions = () => {
    if (!transactionsLoading && hasMoreTransactions) {
      const nextPage = transactionPage + 1
      setTransactionPage(nextPage)
      fetchAddressTransactions(nextPage)
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

  const getAddressType = (address: Address) => {
    if (address.is_contract) {
      return 'Smart Contract'
    }
    return 'Externally Owned Account (EOA)'
  }

  const getAddressTypeIcon = (address: Address) => {
    if (address.is_contract) {
      return <Shield className="w-5 h-5 text-purple-600" />
    }
    return <Wallet className="w-5 h-5 text-blue-600" />
  }

  const getTransactionDirection = (tx: Transaction, currentAddress: string) => {
    if (tx.from_address.toLowerCase() === currentAddress.toLowerCase()) {
      return 'out'
    } else if (tx.to_address?.toLowerCase() === currentAddress.toLowerCase()) {
      return 'in'
    }
    return 'unknown'
  }

  const getDirectionIcon = (direction: string) => {
    if (direction === 'out') {
      return <ArrowLeft className="w-4 h-4 text-red-500 rotate-45" />
    } else if (direction === 'in') {
      return <ArrowLeft className="w-4 h-4 text-green-500 -rotate-45" />
    }
    return <Activity className="w-4 h-4 text-gray-500" />
  }

  const getDirectionColor = (direction: string) => {
    if (direction === 'out') return 'text-red-600'
    if (direction === 'in') return 'text-green-600'
    return 'text-gray-600'
  }

  useEffect(() => {
    if (addressHash) {
      fetchAddress()
    }
  }, [addressHash])

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
          <h1 className="text-3xl font-bold text-gray-900">Address Details</h1>
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
          <h1 className="text-3xl font-bold text-gray-900">Address Details</h1>
        </div>
        <div className="card">
          <div className="text-center py-8">
            <p className="text-red-600 mb-4">Error: {error.message}</p>
            {error.details && (
              <p className="text-gray-600 mb-4 text-sm">{error.details}</p>
            )}
            <button
              onClick={fetchAddress}
              className="btn-primary"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>
    )
  }

  if (!address) {
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
          <h1 className="text-3xl font-bold text-gray-900">Address Not Found</h1>
        </div>
        <div className="card">
          <div className="text-center py-8">
            <p className="text-gray-600">Address {formatHash(addressHash)} was not found.</p>
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
          Back
        </button>
        <div className="flex items-center mb-2">
          {getAddressTypeIcon(address)}
          <h1 className="text-3xl font-bold text-gray-900 ml-3">Address Details</h1>
        </div>
        <p className="text-gray-600">
          {getAddressType(address)} • {address.label || 'Unlabeled Address'}
        </p>
      </div>

      {/* Address Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <div className="flex items-center">
            <Wallet className="w-8 h-8 text-blue-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">Balance</p>
              <p className="text-2xl font-bold text-gray-900">{formatValue(address.balance)} ETH</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <Activity className="w-8 h-8 text-green-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">Transactions</p>
              <p className="text-2xl font-bold text-gray-900">{formatNumber(address.transaction_count)}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <Hash className="w-8 h-8 text-purple-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">Nonce</p>
              <p className="text-2xl font-bold text-gray-900">{formatNumber(address.nonce)}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <TrendingUp className="w-8 h-8 text-orange-600 mr-3" />
            <div>
              <p className="text-sm text-gray-600">Type</p>
              <p className="text-lg font-bold text-gray-900">
                {address.is_contract ? 'Contract' : 'EOA'}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Address Details */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
        {/* Basic Information */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Address Information</h2>
          <div className="space-y-4">
            <div className="flex justify-between items-start">
              <span className="text-sm text-gray-600 font-medium">Address:</span>
              <div className="flex items-center">
                <span className="text-sm font-mono text-gray-900 mr-2">{formatHash(address.address)}</span>
                <button
                  onClick={() => copyToClipboard(address.address, 'address')}
                  className="text-gray-400 hover:text-gray-600"
                >
                  {copiedField === 'address' ? (
                    <CheckCircle className="w-4 h-4 text-green-600" />
                  ) : (
                    <Copy className="w-4 h-4" />
                  )}
                </button>
              </div>
            </div>

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Balance:</span>
              <span className="text-sm text-gray-900 font-mono">{formatValue(address.balance)} ETH</span>
            </div>

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Nonce:</span>
              <span className="text-sm text-gray-900">{formatNumber(address.nonce)}</span>
            </div>

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Address Type:</span>
              <span className="text-sm text-gray-900">{getAddressType(address)}</span>
            </div>

            {address.label && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-600 font-medium">Label:</span>
                <span className="text-sm text-gray-900">{address.label}</span>
              </div>
            )}

            {address.tags && address.tags.length > 0 && (
              <div className="flex justify-between items-start">
                <span className="text-sm text-gray-600 font-medium">Tags:</span>
                <div className="flex flex-wrap gap-1">
                  {address.tags.map((tag, index) => (
                    <span
                      key={index}
                      className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                    >
                      <Tag className="w-3 h-3 mr-1" />
                      {tag}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Contract Information */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">
            {address.is_contract ? 'Contract Information' : 'Activity Information'}
          </h2>
          <div className="space-y-4">
            {address.is_contract && address.contract_creator && (
              <div className="flex justify-between items-start">
                <span className="text-sm text-gray-600 font-medium">Creator:</span>
                <div className="flex items-center">
                  <Link
                    href={`/addresses/${address.contract_creator}`}
                    className="text-sm font-mono text-primary-600 hover:text-primary-800 mr-2"
                  >
                    {formatHash(address.contract_creator)}
                  </Link>
                  <button
                    onClick={() => copyToClipboard(address.contract_creator!, 'creator')}
                    className="text-gray-400 hover:text-gray-600"
                  >
                    {copiedField === 'creator' ? (
                      <CheckCircle className="w-4 h-4 text-green-600" />
                    ) : (
                      <Copy className="w-4 h-4" />
                    )}
                  </button>
                </div>
              </div>
            )}

            {address.is_contract && address.creation_transaction && (
              <div className="flex justify-between items-start">
                <span className="text-sm text-gray-600 font-medium">Creation Tx:</span>
                <div className="flex items-center">
                  <Link
                    href={`/transactions/${address.creation_transaction}`}
                    className="text-sm font-mono text-primary-600 hover:text-primary-800 mr-2"
                  >
                    {formatHash(address.creation_transaction)}
                  </Link>
                  <button
                    onClick={() => copyToClipboard(address.creation_transaction!, 'creation_tx')}
                    className="text-gray-400 hover:text-gray-600"
                  >
                    {copiedField === 'creation_tx' ? (
                      <CheckCircle className="w-4 h-4 text-green-600" />
                    ) : (
                      <Copy className="w-4 h-4" />
                    )}
                  </button>
                </div>
              </div>
            )}

            {address.first_seen_block && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-600 font-medium">First Seen:</span>
                <Link
                  href={`/blocks/${address.first_seen_block}`}
                  className="text-sm text-primary-600 hover:text-primary-800"
                >
                  Block {formatNumber(address.first_seen_block)}
                </Link>
              </div>
            )}

            {address.last_seen_block && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-600 font-medium">Last Seen:</span>
                <Link
                  href={`/blocks/${address.last_seen_block}`}
                  className="text-sm text-primary-600 hover:text-primary-800"
                >
                  Block {formatNumber(address.last_seen_block)}
                </Link>
              </div>
            )}

            <div className="flex justify-between">
              <span className="text-sm text-gray-600 font-medium">Total Transactions:</span>
              <span className="text-sm text-gray-900">{formatNumber(address.transaction_count)}</span>
            </div>

            {address.created_at && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-600 font-medium">First Indexed:</span>
                <span className="text-sm text-gray-900">{formatTimestamp(address.created_at)}</span>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Transaction History */}
      <div className="card">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-semibold text-gray-900">
            Transaction History ({formatNumber(address.transaction_count)})
          </h2>
          {address.transaction_count > 10 && (
            <Link
              href={`/transactions?address=${address.address}`}
              className="text-primary-600 hover:text-primary-800 text-sm font-medium"
            >
              View all transactions →
            </Link>
          )}
        </div>

        {transactionsLoading && transactions.length === 0 ? (
          <div className="flex justify-center items-center h-32">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          </div>
        ) : transactions.length > 0 ? (
          <div className="space-y-4">
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Direction
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Transaction Hash
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Block
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Counterparty
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
                  {transactions.map((tx) => {
                    const direction = getTransactionDirection(tx, address.address)
                    const counterparty = direction === 'out' ? tx.to_address : tx.from_address
                    
                    return (
                      <tr key={tx.hash} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center">
                            {getDirectionIcon(direction)}
                            <span className={`ml-2 text-sm font-medium ${getDirectionColor(direction)}`}>
                              {direction === 'out' ? 'OUT' : direction === 'in' ? 'IN' : 'SELF'}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <Link
                            href={`/transactions/${tx.hash}`}
                            className="text-primary-600 hover:text-primary-800 font-mono text-sm"
                          >
                            {formatHash(tx.hash)}
                          </Link>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <Link
                            href={`/blocks/${tx.block_number}`}
                            className="text-primary-600 hover:text-primary-800 text-sm"
                          >
                            {formatNumber(tx.block_number)}
                          </Link>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          {counterparty ? (
                            <Link
                              href={`/addresses/${counterparty}`}
                              className="text-gray-900 hover:text-primary-600 font-mono text-sm"
                            >
                              {formatHash(counterparty)}
                            </Link>
                          ) : (
                            <span className="text-gray-500 text-sm italic">Contract Creation</span>
                          )}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`text-sm font-mono ${getDirectionColor(direction)}`}>
                            {direction === 'out' ? '-' : direction === 'in' ? '+' : ''}
                            {formatValue(tx.value)} ETH
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {tx.gas_used ? formatNumber(tx.gas_used) : 'N/A'}
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>

            {hasMoreTransactions && (
              <div className="text-center pt-4">
                <button
                  onClick={loadMoreTransactions}
                  disabled={transactionsLoading}
                  className="btn-secondary"
                >
                  {transactionsLoading ? (
                    <div className="flex items-center">
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary-600 mr-2"></div>
                      Loading...
                    </div>
                  ) : (
                    'Load More Transactions'
                  )}
                </button>
              </div>
            )}
          </div>
        ) : (
          <div className="text-center py-8">
            <Activity className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No transactions found</h3>
            <p className="text-gray-600">This address has no transaction history.</p>
          </div>
        )}
      </div>
    </div>
  )
}
