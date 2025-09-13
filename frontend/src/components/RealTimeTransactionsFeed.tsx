'use client'

import React, { useState, useEffect } from 'react'
import Link from 'next/link'
import { useLatestTransaction } from '@/contexts/RealTimeContext'
import { formatDistanceToNow } from 'date-fns'
import { formatNumber, formatValue, formatAddress, formatGasPrice } from '@/utils/formatting'

interface Transaction {
  hash: string
  block_number: number
  from_address: string
  to_address?: string
  value: string
  gas_price: string
}

export default function RealTimeTransactionsFeed() {
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const latestTransaction = useLatestTransaction()

  // Add new transactions to the feed
  useEffect(() => {
    if (latestTransaction) {
      setTransactions(prevTxs => {
        // Check if transaction already exists
        const exists = prevTxs.some(tx => tx.hash === latestTransaction.hash)
        if (exists) return prevTxs

        // Add new transaction to the beginning and keep only last 10
        return [latestTransaction, ...prevTxs].slice(0, 10)
      })
    }
  }, [latestTransaction])

  if (transactions.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Latest Transactions</h3>
        <div className="text-gray-500 text-center py-8">
          Waiting for new transactions...
        </div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-medium text-gray-900">Latest Transactions</h3>
        <div className="flex items-center space-x-1">
          <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
          <span className="text-xs text-gray-500">Live Updates</span>
        </div>
      </div>
      
      <div className="space-y-3">
        {transactions.map((tx, index) => (
          <div
            key={tx.hash}
            className={`p-3 rounded-lg border transition-all duration-500 ${
              index === 0 && latestTransaction?.hash === tx.hash
                ? 'border-green-200 bg-green-50 animate-pulse'
                : 'border-gray-200 hover:border-gray-300'
            }`}
          >
            <div className="flex items-start justify-between">
              <div className="flex-1 min-w-0">
                <div className="flex items-center space-x-2 mb-2">
                  <div className="flex-shrink-0">
                    <div className="w-8 h-8 bg-purple-100 rounded-lg flex items-center justify-center">
                      <span className="text-purple-600 font-semibold text-xs">Tx</span>
                    </div>
                  </div>
                  <div className="flex-1 min-w-0">
                    <Link
                      href={`/transactions/${tx.hash}`}
                      className="text-sm font-medium text-gray-900 hover:text-blue-600 truncate block"
                    >
                      {formatAddress(tx.hash)}
                    </Link>
                    {index === 0 && latestTransaction?.hash === tx.hash && (
                      <span className="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800 mt-1">
                        NEW
                      </span>
                    )}
                  </div>
                </div>
                
                <div className="space-y-1">
                  <div className="flex items-center text-xs text-gray-500">
                    <span className="mr-2">From:</span>
                    <Link
                      href={`/addresses/${tx.from_address}`}
                      className="font-mono hover:text-blue-600"
                    >
                      {formatAddress(tx.from_address)}
                    </Link>
                  </div>
                  
                  {tx.to_address && (
                    <div className="flex items-center text-xs text-gray-500">
                      <span className="mr-2">To:</span>
                      <Link
                        href={`/addresses/${tx.to_address}`}
                        className="font-mono hover:text-blue-600"
                      >
                        {formatAddress(tx.to_address)}
                      </Link>
                    </div>
                  )}
                </div>
              </div>
              
              <div className="text-right ml-4">
                <div className="text-sm font-medium text-gray-900">
                  {formatValue(tx.value)}
                </div>
                <div className="text-xs text-gray-500">
                  Block #{formatNumber(tx.block_number)}
                </div>
                <div className="text-xs text-gray-400">
                  {formatGasPrice(tx.gas_price)}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
      
      <div className="mt-4 text-center">
        <Link
          href="/transactions"
          className="text-sm text-blue-600 hover:text-blue-800 font-medium"
        >
          View All Transactions →
        </Link>
      </div>
    </div>
  )
}
