'use client';

import React, { useState, useEffect } from 'react';
import GasPriceCard from '@/components/GasAnalytics/GasPriceCard';
import GasPriceChart from '@/components/GasAnalytics/GasPriceChart';
import GasFeeCalculator from '@/components/GasAnalytics/GasFeeCalculator';

interface GasPriceData {
  timestamp: string;
  slow: number;
  standard: number;
  fast: number;
}

interface GasPriceStats {
  average: number;
  median: number;
  min: number;
  max: number;
  percentile_25: number;
  percentile_75: number;
}

interface CurrentGasPrices {
  slow: number;
  standard: number;
  fast: number;
  timestamp: string;
}

export default function GasAnalyticsPage() {
  const [currentPrices, setCurrentPrices] = useState<CurrentGasPrices | null>(null);
  const [historicalData, setHistoricalData] = useState<GasPriceData[]>([]);
  const [stats, setStats] = useState<GasPriceStats | null>(null);
  const [selectedTimeframe, setSelectedTimeframe] = useState<string>('24h');
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const timeframes = [
    { value: '1h', label: '1 Hour' },
    { value: '24h', label: '24 Hours' },
    { value: '7d', label: '7 Days' },
    { value: '30d', label: '30 Days' },
  ];

  const fetchCurrentPrices = async () => {
    try {
      const response = await fetch('/api/v1/gas/prices');
      if (!response.ok) throw new Error('Failed to fetch current prices');
      const data = await response.json();
      setCurrentPrices(data);
    } catch (err) {
      console.error('Error fetching current prices:', err);
      setError('Failed to load current gas prices');
    }
  };

  const fetchHistoricalData = async (timeframe: string) => {
    try {
      const response = await fetch(`/api/v1/gas/history?period=${timeframe}&limit=100`);
      if (!response.ok) throw new Error('Failed to fetch historical data');
      const data = await response.json();
      setHistoricalData(data.data || []);
    } catch (err) {
      console.error('Error fetching historical data:', err);
      setError('Failed to load historical data');
    }
  };

  const fetchStats = async (timeframe: string) => {
    try {
      const response = await fetch(`/api/v1/gas/stats?period=${timeframe}`);
      if (!response.ok) throw new Error('Failed to fetch stats');
      const data = await response.json();
      setStats(data);
    } catch (err) {
      console.error('Error fetching stats:', err);
    }
  };

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      setError(null);
      
      await Promise.all([
        fetchCurrentPrices(),
        fetchHistoricalData(selectedTimeframe),
        fetchStats(selectedTimeframe)
      ]);
      
      setLoading(false);
    };

    loadData();
  }, [selectedTimeframe]);

  // Auto-refresh current prices every 30 seconds
  useEffect(() => {
    const interval = setInterval(fetchCurrentPrices, 30000);
    return () => clearInterval(interval);
  }, []);

  const calculateChange = (current: number, historical: GasPriceData[]) => {
    if (!historical.length) return 0;
    const previous = historical[0];
    const previousValue = previous.standard; // Use standard as reference
    return ((current - previousValue) / previousValue) * 100;
  };

  if (loading && !currentPrices) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading gas analytics...</p>
        </div>
      </div>
    );
  }

  if (error && !currentPrices) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-red-500 text-xl mb-4">⚠️</div>
          <p className="text-gray-600 mb-4">{error}</p>
          <button
            onClick={() => window.location.reload()}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Gas Price Analytics</h1>
          <p className="text-gray-600">
            Real-time Ethereum gas prices, historical trends, and fee calculations
          </p>
        </div>

        {/* Current Gas Prices */}
        {currentPrices && (
          <div className="mb-8">
            <h2 className="text-xl font-semibold text-gray-800 mb-4">Current Gas Prices</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <GasPriceCard
                title="Slow"
                price={currentPrices.slow}
                change={calculateChange(currentPrices.slow, historicalData)}
                description="~10+ minutes"
                priority="low"
              />
              <GasPriceCard
                title="Standard"
                price={currentPrices.standard}
                change={calculateChange(currentPrices.standard, historicalData)}
                description="~3-5 minutes"
                priority="standard"
              />
              <GasPriceCard
                title="Fast"
                price={currentPrices.fast}
                change={calculateChange(currentPrices.fast, historicalData)}
                description="~1-2 minutes"
                priority="fast"
              />
            </div>
            <p className="text-sm text-gray-500 mt-4">
              Last updated: {new Date(currentPrices.timestamp).toLocaleString()}
            </p>
          </div>
        )}

        {/* Historical Chart */}
        <div className="mb-8">
          <div className="bg-white rounded-lg border border-gray-200 p-6">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-6">
              <h2 className="text-xl font-semibold text-gray-800 mb-4 sm:mb-0">
                Gas Price History
              </h2>
              <div className="flex space-x-2">
                {timeframes.map((timeframe) => (
                  <button
                    key={timeframe.value}
                    onClick={() => setSelectedTimeframe(timeframe.value)}
                    className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                      selectedTimeframe === timeframe.value
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                  >
                    {timeframe.label}
                  </button>
                ))}
              </div>
            </div>
            
            <GasPriceChart
              data={historicalData}
              timeframe={selectedTimeframe}
              loading={loading}
            />
          </div>
        </div>

        {/* Statistics */}
        {stats && (
          <div className="mb-8">
            <h2 className="text-xl font-semibold text-gray-800 mb-4">
              Statistics ({timeframes.find(t => t.value === selectedTimeframe)?.label})
            </h2>
            <div className="grid grid-cols-2 md:grid-cols-6 gap-4">
              <div className="bg-white rounded-lg border border-gray-200 p-4 text-center">
                <p className="text-sm text-gray-600">Average</p>
                <p className="text-lg font-semibold text-gray-900">{stats.average?.toFixed(1) || '0.0'} gwei</p>
              </div>
              <div className="bg-white rounded-lg border border-gray-200 p-4 text-center">
                <p className="text-sm text-gray-600">Median</p>
                <p className="text-lg font-semibold text-gray-900">{stats.median?.toFixed(1) || '0.0'} gwei</p>
              </div>
              <div className="bg-white rounded-lg border border-gray-200 p-4 text-center">
                <p className="text-sm text-gray-600">Min</p>
                <p className="text-lg font-semibold text-green-600">{stats.min?.toFixed(1) || '0.0'} gwei</p>
              </div>
              <div className="bg-white rounded-lg border border-gray-200 p-4 text-center">
                <p className="text-sm text-gray-600">Max</p>
                <p className="text-lg font-semibold text-red-600">{stats.max?.toFixed(1) || '0.0'} gwei</p>
              </div>
              <div className="bg-white rounded-lg border border-gray-200 p-4 text-center">
                <p className="text-sm text-gray-600">25th %ile</p>
                <p className="text-lg font-semibold text-gray-900">{stats.percentile_25?.toFixed(1) || '0.0'} gwei</p>
              </div>
              <div className="bg-white rounded-lg border border-gray-200 p-4 text-center">
                <p className="text-sm text-gray-600">75th %ile</p>
                <p className="text-lg font-semibold text-gray-900">{stats.percentile_75?.toFixed(1) || '0.0'} gwei</p>
              </div>
            </div>
          </div>
        )}

        {/* Fee Calculator */}
        {currentPrices && (
          <GasFeeCalculator
            currentGasPrices={{
              slow: currentPrices.slow,
              standard: currentPrices.standard,
              fast: currentPrices.fast
            }}
          />
        )}
      </div>
    </div>
  );
}
