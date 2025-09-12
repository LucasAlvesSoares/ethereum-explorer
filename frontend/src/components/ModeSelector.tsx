'use client';

import { useState, useEffect } from 'react';

interface ModeSelectorProps {
  className?: string;
}

export default function ModeSelector({ className = '' }: ModeSelectorProps) {
  const [currentMode, setCurrentMode] = useState<'demo' | 'live'>('demo');
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    // Check current mode from environment or API
    fetchCurrentMode();
  }, []);

  const fetchCurrentMode = async () => {
    try {
      const response = await fetch('/api/v1/health');
      const data = await response.json();
      // In a real implementation, you might get mode info from the health check
      // For now, we'll check if ethereum is connected to determine mode
      setCurrentMode(data.ethereum === 'connected' ? 'live' : 'demo');
    } catch (error) {
      console.warn('Failed to fetch current mode:', error);
      setCurrentMode('demo');
    }
  };

  const handleModeChange = async (newMode: 'demo' | 'live') => {
    if (newMode === currentMode || isLoading) return;

    setIsLoading(true);
    try {
      // In a real implementation, you would call an API to switch modes
      // This would require backend support for mode switching
      console.log(`Switching to ${newMode} mode...`);
      
      // For now, just update the UI
      setCurrentMode(newMode);
      
      // You could also trigger a page refresh to reinitialize with new mode
      // window.location.reload();
    } catch (error) {
      console.error('Failed to switch mode:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className={`flex items-center space-x-2 ${className}`}>
      <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
        Mode:
      </span>
      <div className="flex bg-gray-200 dark:bg-gray-700 rounded-lg p-1">
        <button
          onClick={() => handleModeChange('demo')}
          disabled={isLoading}
          className={`px-3 py-1 text-sm font-medium rounded-md transition-colors ${
            currentMode === 'demo'
              ? 'bg-blue-600 text-white shadow-sm'
              : 'text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200'
          } ${isLoading ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
        >
          Demo
        </button>
        <button
          onClick={() => handleModeChange('live')}
          disabled={isLoading}
          className={`px-3 py-1 text-sm font-medium rounded-md transition-colors ${
            currentMode === 'live'
              ? 'bg-green-600 text-white shadow-sm'
              : 'text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200'
          } ${isLoading ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
        >
          Live
        </button>
      </div>
      <div className="flex items-center space-x-1">
        <div
          className={`w-2 h-2 rounded-full ${
            currentMode === 'live' ? 'bg-green-500' : 'bg-blue-500'
          }`}
        />
        <span className="text-xs text-gray-500 dark:text-gray-400">
          {currentMode === 'live' ? 'Live Data' : 'Demo Data'}
        </span>
      </div>
    </div>
  );
}
