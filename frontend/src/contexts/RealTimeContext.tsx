'use client'

import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react'
import { useWebSocket, WebSocketMessage, BlockUpdate, TransactionUpdate } from '@/hooks/useWebSocket'

interface RealTimeData {
  latestBlock: BlockUpdate | null
  latestTransaction: TransactionUpdate | null
  networkStats: any | null
  isConnected: boolean
  connectionError: string | null
}

interface RealTimeContextType extends RealTimeData {
  subscribe: (topic: string) => void
  unsubscribe: (topic: string) => void
  connect: () => void
  disconnect: () => void
}

const RealTimeContext = createContext<RealTimeContextType | undefined>(undefined)

interface RealTimeProviderProps {
  children: ReactNode
}

export function RealTimeProvider({ children }: RealTimeProviderProps) {
  const [latestBlock, setLatestBlock] = useState<BlockUpdate | null>(null)
  const [latestTransaction, setLatestTransaction] = useState<TransactionUpdate | null>(null)
  const [networkStats, setNetworkStats] = useState<any | null>(null)
  const [connectionError, setConnectionError] = useState<string | null>(null)

  // Get WebSocket URL from environment or default to localhost
  const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/api/v1/ws'

  const {
    isConnected,
    lastMessage,
    connect,
    disconnect,
    subscribe,
    unsubscribe,
  } = useWebSocket({
    url: wsUrl,
    autoConnect: true,
    reconnectAttempts: 5,
    reconnectInterval: 3000,
  })

  // Handle incoming WebSocket messages
  useEffect(() => {
    if (!lastMessage) return

    try {
      switch (lastMessage.type) {
        case 'block_update':
          console.log('Received block update:', lastMessage.data)
          setLatestBlock(lastMessage.data as BlockUpdate)
          setConnectionError(null)
          break

        case 'transaction_update':
          console.log('Received transaction update:', lastMessage.data)
          setLatestTransaction(lastMessage.data as TransactionUpdate)
          setConnectionError(null)
          break

        case 'network_stats':
          console.log('Received network stats:', lastMessage.data)
          setNetworkStats(lastMessage.data)
          setConnectionError(null)
          break

        default:
          console.log('Unknown message type:', lastMessage.type)
      }
    } catch (error) {
      console.error('Error processing WebSocket message:', error)
      setConnectionError('Failed to process real-time update')
    }
  }, [lastMessage])

  // Handle connection state changes
  useEffect(() => {
    if (isConnected) {
      setConnectionError(null)
      console.log('Real-time updates connected')
    } else {
      setConnectionError('Real-time updates disconnected')
      console.log('Real-time updates disconnected')
    }
  }, [isConnected])

  const contextValue: RealTimeContextType = {
    latestBlock,
    latestTransaction,
    networkStats,
    isConnected,
    connectionError,
    subscribe,
    unsubscribe,
    connect,
    disconnect,
  }

  return (
    <RealTimeContext.Provider value={contextValue}>
      {children}
    </RealTimeContext.Provider>
  )
}

export function useRealTime(): RealTimeContextType {
  const context = useContext(RealTimeContext)
  if (context === undefined) {
    throw new Error('useRealTime must be used within a RealTimeProvider')
  }
  return context
}

// Convenience hooks for specific data types
export function useLatestBlock(): BlockUpdate | null {
  const { latestBlock } = useRealTime()
  return latestBlock
}

export function useLatestTransaction(): TransactionUpdate | null {
  const { latestTransaction } = useRealTime()
  return latestTransaction
}

export function useNetworkStats(): any | null {
  const { networkStats } = useRealTime()
  return networkStats
}

export function useRealTimeConnection(): {
  isConnected: boolean
  connectionError: string | null
  connect: () => void
  disconnect: () => void
} {
  const { isConnected, connectionError, connect, disconnect } = useRealTime()
  return { isConnected, connectionError, connect, disconnect }
}
