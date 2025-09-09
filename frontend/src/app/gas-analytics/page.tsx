'use client';

import React, { useState, useEffect } from 'react';
import GasPriceChart from '@/components/GasAnalytics/GasPriceChart';
import GasPriceCard from '@/components/GasAnalytics/GasPriceCard';
import GasFeeCalculator from '@/components/GasAnalytics/GasFeeCalculator';

interface GasPrice {
  slow: number;
  standard: number;
  fast: number;
}

interface GasPriceStats {
  current: GasPrice;
  average_24h: GasPrice;
  median_24h: GasPrice;
  min_24h: GasPrice;
  max_24h: GasPrice;
  trend: string;
  network_utilization: number;
}

interface GasPriceHistoryPoint {
  timestamp: string;
  slow: number;
  standard: number;
  fast: number;
}

interface GasPriceRecommendation {
  transaction_type: string;
  gas_price: number;
  estimated_time_seconds: number;
  description: string;
}

export default function GasAnalyticsPage() {
  const [currentPrices, setCurrentPrices] = useState<GasPrice | null>(null);
  const [stats, setStats] = useState<GasPriceStats | null>(null);
  const [history, setHistory] = useState<GasPriceHistoryPoint[]>([]);
  const [recommendations, setRecommendations] = useState<GasPriceRecommendation[]>([]);
  const [timeframe, setTimeframe] = useState('24h');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchGasPrices = async () => {
    try {
      const response = await fetch('/api/v1/gas/prices');
      if (!response.ok) throw new Error('Failed to fetch gas prices');
      const data = await response.json();
      setCurrentPrices(data);
    } catch (err) {
      console.error('Error fetching gas prices:', err);
      setError('Failed to load gas prices');
    }
  };

  const fetchGasStats = async () => {
    try {
      const response = await fetch('/api/v1/gas/stats');
      if (!response.ok) throw new Error('Failed to fetch gas stats');
      const data = await response.json();
      setStats(data);
    } catch (err) {
      console.error('Error fetching gas stats:', err);
    }
  };

  const fetchGasHistory = async (tf: string) => {
    try {
      const response = await fetch(`/api/v1/gas/history?timeframe=${tf}`);
      if (!response.ok) throw new Error('Failed to fetch gas history');
      const data = await response.json();
      setHistory(data.data || []);
    } catch (err) {
      console.error('Error fetching gas history:', err);
    }
  };

  const fetchRecommendations = async () => {
    try {
      const response = await fetch('/api/v1/gas/recommendations');
      if (!response.ok) throw new Error('Failed to fetch recommendations');
      const data = await response.json();
      setRecommendations(data.recommendations || []);
    } catch (err) {
      console.error('Error fetching recommendations:', err);
    }
  };

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      await Promise.all([
        fetchGasPrices(),
        fetchGasStats(),
        fetchGasHistory(timeframe),
        fetchRecommendations(),
      ]);
      setLoading(false);
    };

    loadData();
  }, [timeframe]);

  const handleTimeframeChange = (newTimeframe: string) => {
    setTimeframe(newTimeframe);
    fetchGasHistory(newTimeframe);
  };

  const getTrendColor = (trend: string) => {
    switch (trend) {
      case 'rising':
        return 'text-red-600';
      case 'falling':
        return 'text-green-600';
      default:
        return 'text-gray-600';
    }
  };

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'rising':
        return '↗️';
      case 'falling':
        return '↘️';
      default:
        return '→';
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Error Loading Gas Analytics</h1>
          <p className="text-gray-600 mb-4">{error}</p>
          <button
            onClick={() => window.location.reload()}
            className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Gas Analytics Dashboard</h1>
          <p className="text-gray-600">
            Real-time Ethereum gas price tracking, trends, and fee optimization tools
          </p>
        </div>

        {/* Current Gas Prices */}
        {currentPrices && (
          <div className="mb-8">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Current Gas Prices</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <GasPriceCard
                title="Slow"
                price={currentPrices.slow}
                description="~5 minutes • Cost-effective for non-urgent transactions"
                priority="low"
              />
              <GasPriceCard
                title="Standard"
                price={currentPrices.standard}
                description="~3 minutes • Balanced speed and cost"
                priority="standard"
              />
              <GasPriceCard
                title="Fast"
                price={currentPrices.fast}
                description="~1 minute • Priority for time-sensitive transactions"
                priority="fast"
              />
            </div>
          </div>
        )}

        {/* Network Statistics */}
        {stats && (
          <div className="mb-8">
            <div className="bg-white rounded-lg shadow-sm p-6">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">Network Statistics</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <div className="text-center">
                  <div className="text-2xl font-bold text-gray-900">
                    {stats.network_utilization.toFixed(1)}%
                  </div>
                  <div className="text-sm text-gray-600">Network Utilization</div>
                </div>
                <div className="text-center">
                  <div className={`text-2xl font-bold ${getTrendColor(stats.trend)}`}>
                    {getTrendIcon(stats.trend)} {stats.trend}
                  </div>
                  <div className="text-sm text-gray-600">Price Trend</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-gray-900">
                    {stats.average_24h.standard} gwei
                  </div>
                  <div className="text-sm text-gray-600">24h Average</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-gray-900">
                    {stats.min_24h.standard} - {stats.max_24h.standard} gwei
                  </div>
                  <div className="text-sm text-gray-600">24h Range</div>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Gas Price Chart */}
        <div className="mb-8">
          <div className="bg-white rounded-lg shadow-sm p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold text-gray-900">Gas Price History</h2>
              <div className="flex space-x-2">
                {['1h', '24h', '7d', '30d'].map((tf) => (
                  <button
                    key={tf}
                    onClick={() => handleTimeframeChange(tf)}
                    className={`px-3 py-1 rounded-md text-sm font-medium ${
                      timeframe === tf
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                  >
                    {tf}
                  </button>
                ))}
              </div>
            </div>
            <GasPriceChart data={history} timeframe={timeframe} loading={false} />
          </div>
        </div>

        {/* Gas Fee Calculator */}
        {currentPrices && (
          <div className="mb-8">
            <div className="bg-white rounded-lg shadow-sm p-6">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">Gas Fee Calculator</h2>
              <GasFeeCalculator currentGasPrices={currentPrices} />
            </div>
          </div>
        )}

        {/* Transaction Type Recommendations */}
        {recommendations.length > 0 && (
          <div className="mb-8">
            <div className="bg-white rounded-lg shadow-sm p-6">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">
                Gas Price Recommendations by Transaction Type
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {recommendations.map((rec, index) => (
                  <div key={index} className="border border-gray-200 rounded-lg p-4">
                    <h3 className="font-semibold text-gray-900 mb-2">
                      {rec.transaction_type.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                    </h3>
                    <div className="text-2xl font-bold text-blue-600 mb-1">
                      {rec.gas_price} gwei
                    </div>
                    <div className="text-sm text-gray-600 mb-2">
                      ~{Math.round(rec.estimated_time_seconds / 60)} minutes
                    </div>
                    <p className="text-sm text-gray-700">{rec.description}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* Tips and Information */}
        <div className="bg-white rounded-lg shadow-sm p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Gas Optimization Tips</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 className="font-semibold text-gray-900 mb-2">💡 Best Practices</h3>
              <ul className="text-sm text-gray-700 space-y-1">
                <li>• Use slow gas prices for non-urgent transactions</li>
                <li>• Monitor network congestion before making transactions</li>
                <li>• Batch multiple operations when possible</li>
                <li>• Consider transaction timing (weekends often cheaper)</li>
              </ul>
            </div>
            <div>
              <h3 className="font-semibold text-gray-900 mb-2">⚡ When to Use Fast Gas</h3>
              <ul className="text-sm text-gray-700 space-y-1">
                <li>• DeFi arbitrage opportunities</li>
                <li>• NFT minting during high demand</li>
                <li>• Time-sensitive trading</li>
                <li>• Avoiding MEV attacks</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
