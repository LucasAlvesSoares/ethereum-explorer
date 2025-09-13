import { useEffect, useRef, useState, useCallback } from 'react'

export interface WebSocketMessage {
  type: string
  data: any
}

export interface BlockUpdate {
  number: number
  hash: string
  transaction_count: number
  gas_used: number
  gas_limit: number
  timestamp: string
  miner: string
}

export interface TransactionUpdate {
  hash: string
  block_number: number
  from_address: string
  to_address?: string
  value: string
  gas_price: string
}

export interface UseWebSocketOptions {
  url: string
  autoConnect?: boolean
  reconnectAttempts?: number
  reconnectInterval?: number
}

export interface UseWebSocketReturn {
  isConnected: boolean
  lastMessage: WebSocketMessage | null
  sendMessage: (message: any) => void
  connect: () => void
  disconnect: () => void
  subscribe: (topic: string) => void
  unsubscribe: (topic: string) => void
}

export function useWebSocket({
  url,
  autoConnect = true,
  reconnectAttempts = 5,
  reconnectInterval = 3000,
}: UseWebSocketOptions): UseWebSocketReturn {
  const [isConnected, setIsConnected] = useState(false)
  const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const reconnectCountRef = useRef(0)
  const subscribedTopicsRef = useRef<Set<string>>(new Set())

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      return
    }

    try {
      wsRef.current = new WebSocket(url)

      wsRef.current.onopen = () => {
        console.log('WebSocket connected')
        setIsConnected(true)
        reconnectCountRef.current = 0

        // Re-subscribe to previously subscribed topics
        subscribedTopicsRef.current.forEach(topic => {
          wsRef.current?.send(JSON.stringify({
            type: 'subscribe',
            topic,
          }))
        })
      }

      wsRef.current.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)
          setLastMessage(message)
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      wsRef.current.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason)
        setIsConnected(false)

        // Attempt to reconnect if not a manual disconnect and we haven't exceeded attempts
        if (event.code !== 1000 && reconnectCountRef.current < reconnectAttempts) {
          reconnectCountRef.current++
          console.log(`Attempting to reconnect (${reconnectCountRef.current}/${reconnectAttempts})...`)
          
          reconnectTimeoutRef.current = setTimeout(() => {
            connect()
          }, reconnectInterval)
        } else if (reconnectCountRef.current >= reconnectAttempts) {
          console.error('Max reconnection attempts reached')
        }
      }

      wsRef.current.onerror = (error) => {
        console.error('WebSocket error:', error)
      }
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error)
    }
  }, [url, reconnectAttempts, reconnectInterval])

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }

    if (wsRef.current) {
      wsRef.current.close(1000, 'Manual disconnect')
      wsRef.current = null
    }

    setIsConnected(false)
    reconnectCountRef.current = 0
  }, [])

  const sendMessage = useCallback((message: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message))
    } else {
      console.warn('WebSocket is not connected')
    }
  }, [])

  const subscribe = useCallback((topic: string) => {
    subscribedTopicsRef.current.add(topic)
    sendMessage({
      type: 'subscribe',
      topic,
    })
  }, [sendMessage])

  const unsubscribe = useCallback((topic: string) => {
    subscribedTopicsRef.current.delete(topic)
    sendMessage({
      type: 'unsubscribe',
      topic,
    })
  }, [sendMessage])

  useEffect(() => {
    if (autoConnect) {
      connect()
    }

    return () => {
      disconnect()
    }
  }, [autoConnect, connect, disconnect])

  return {
    isConnected,
    lastMessage,
    sendMessage,
    connect,
    disconnect,
    subscribe,
    unsubscribe,
  }
}
