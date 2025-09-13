'use client'

import React from 'react'
import { useRealTimeConnection } from '@/contexts/RealTimeContext'

export default function RealTimeStatus() {
  const { isConnected, connectionError, connect, disconnect } = useRealTimeConnection()

  return (
    <div className="flex items-center space-x-2">
      <div className="flex items-center space-x-1">
        <div
          className={`w-2 h-2 rounded-full ${
            isConnected ? 'bg-green-500' : 'bg-red-500'
          }`}
        />
        <span className="text-xs text-gray-600">
          {isConnected ? 'Live' : 'Offline'}
        </span>
      </div>
      
      {connectionError && (
        <div className="flex items-center space-x-1">
          <button
            onClick={connect}
            className="text-xs text-blue-600 hover:text-blue-800 underline"
            title={connectionError}
          >
            Reconnect
          </button>
        </div>
      )}
    </div>
  )
}
