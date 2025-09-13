'use client'

import React, { useState, useEffect } from 'react'
import Link from 'next/link'
import { useLatestBlock } from '@/contexts/RealTimeContext'
import { formatDistanceToNow } from 'date-fns'
import { formatNumber, formatGas } from '@/utils/formatting'

interface Block {
  number: number
  hash: string
  transaction_count: number
  gas_used: number
  gas_limit: number
  timestamp: string
  miner: string
}

export default function RealTimeBlocksFeed() {
  const [blocks, setBlocks] = useState<Block[]>([])
  const latestBlock = useLatestBlock()

  // Add new blocks to the feed
  useEffect(() => {
    if (latestBlock) {
      setBlocks(prevBlocks => {
        // Check if block already exists
        const exists = prevBlocks.some(block => block.hash === latestBlock.hash)
        if (exists) return prevBlocks

        // Add new block to the beginning and keep only last 10
        return [latestBlock, ...prevBlocks].slice(0, 10)
      })
    }
  }, [latestBlock])

  if (blocks.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Latest Blocks</h3>
        <div className="text-gray-500 text-center py-8">
          Waiting for new blocks...
        </div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-medium text-gray-900">Latest Blocks</h3>
        <div className="flex items-center space-x-1">
          <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
          <span className="text-xs text-gray-500">Live Updates</span>
        </div>
      </div>
      
      <div className="space-y-3">
        {blocks.map((block, index) => (
          <div
            key={block.hash}
            className={`p-3 rounded-lg border transition-all duration-500 ${
              index === 0 && latestBlock?.hash === block.hash
                ? 'border-green-200 bg-green-50 animate-pulse'
                : 'border-gray-200 hover:border-gray-300'
            }`}
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <div className="flex-shrink-0">
                  <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
                    <span className="text-blue-600 font-semibold text-sm">
                      {block.number.toString().slice(-3)}
                    </span>
                  </div>
                </div>
                <div>
                  <div className="flex items-center space-x-2">
                    <Link
                      href={`/blocks/${block.number}`}
                      className="font-medium text-gray-900 hover:text-blue-600"
                    >
                      Block #{formatNumber(block.number)}
                    </Link>
                    {index === 0 && latestBlock?.hash === block.hash && (
                      <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
                        NEW
                      </span>
                    )}
                  </div>
                  <div className="text-sm text-gray-500">
                    {formatDistanceToNow(new Date(block.timestamp), { addSuffix: true })}
                  </div>
                </div>
              </div>
              <div className="text-right">
                <div className="text-sm font-medium text-gray-900">
                  {formatNumber(block.transaction_count)} txs
                </div>
                <div className="text-xs text-gray-500">
                  {formatGas(block.gas_used)} gas
                </div>
              </div>
            </div>
            <div className="mt-2 text-xs text-gray-400 font-mono">
              {block.hash.slice(0, 20)}...{block.hash.slice(-10)}
            </div>
          </div>
        ))}
      </div>
      
      <div className="mt-4 text-center">
        <Link
          href="/blocks"
          className="text-sm text-blue-600 hover:text-blue-800 font-medium"
        >
          View All Blocks →
        </Link>
      </div>
    </div>
  )
}
