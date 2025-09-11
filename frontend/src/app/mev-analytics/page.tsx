'use client';

import { useState, useEffect } from 'react';
import { formatValue, formatGasPrice } from '../../utils/formatting';

interface MEVTransaction {
  hash: string;
  block_number: number;
  transaction_index: number;
  from_address: string;
  to_address?: string;
  value: string;
  gas_price: string;
  gas_used?: number;
  gas_limit: number;
  timestamp: string;
  mev_type: string;
  mev_score: number;
  potential_profit?: string;
  related_tx_hashes: string[];
}

interface MEVBot {
  address: string;
  transaction_count: number;
  high_gas_transactions: number;
  average_gas_price: string;
  mev_score: number;
  first_seen_block: number;
  last_seen_block: number;
  suspicious_patterns: string[];
}

interface MEVStats {
  time_period: {
    hours: number;
    start_time: string;
    end_time: string;
  };
  total_transactions: number;
  mev_transactions: number;
  mev_percentage: number;
  average_gas_price: string;
  top_mev_bots_count: number;
  network_health: {
    mev_activity_level: string;
  };
}

export default function MEVAnalyticsPage() {
  const [activeTab, setActiveTab] = useState<'overview' | 'transactions' | 'bots' | 'block-analysis'>('overview');
  const [stats, setStats] = useState<MEVStats | null>(null);
  const [mevBots, setMevBots] = useState<MEVBot[]>([]);
  const [suspiciousTransactions, setSuspiciousTransactions] = useState<MEVTransaction[]>([]);
  const [blockNumber, setBlockNumber] = useState<string>('');
  const [blockAnalysis, setBlockAnalysis] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [timeRange, setTimeRange] = useState<number>(24);

  useEffect(() => {
    loadMEVStats();
    loadMEVBots();
  }, [timeRange]);

  const loadMEVStats = async () => {
    try {
      const response = await fetch(`/api/v1/mev-analytics/stats?hours=${timeRange}`);
      const data = await response.json();
      if (data.success) {
        setStats(data.data);
      }
    } catch (error) {
      console.error('Error loading MEV stats:', error);
    }
  };

  const loadMEVBots = async () => {
    try {
      const response = await fetch(`/api/v1/mev-analytics/mev-bots?hours=${timeRange}&min_transactions=10`);
      const data = await response.json();
      if (data.success) {
        setMevBots(data.data.mev_bots || []);
      }
    } catch (error) {
      console.error('Error loading MEV bots:', error);
    }
  };

  const loadSuspiciousTransactions = async (blockNum: string) => {
    if (!blockNum) return;
    
    setLoading(true);
    try {
      const response = await fetch(`/api/v1/mev-analytics/suspicious-transactions?block_number=${blockNum}&threshold=2.0`);
      const data = await response.json();
      if (data.success) {
        setSuspiciousTransactions(data.data.suspicious_transactions || []);
      }
    } catch (error) {
      console.error('Error loading suspicious transactions:', error);
    } finally {
      setLoading(false);
    }
  };

  const analyzeBlock = async (blockNum: string) => {
    if (!blockNum) return;
    
    setLoading(true);
    try {
      const response = await fetch(`/api/v1/mev-analytics/block/${blockNum}`);
      const data = await response.json();
      if (data.success) {
        setBlockAnalysis(data.data);
      }
    } catch (error) {
      console.error('Error analyzing block:', error);
    } finally {
      setLoading(false);
    }
  };

  const getMEVActivityColor = (level: string) => {
    switch (level) {
      case 'high': return 'text-red-600 bg-red-100';
      case 'moderate': return 'text-yellow-600 bg-yellow-100';
      case 'low': return 'text-green-600 bg-green-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const getMEVTypeColor = (type: string) => {
    switch (type) {
      case 'sandwich_attack': return 'bg-red-100 text-red-800';
      case 'high_gas_price': return 'bg-orange-100 text-orange-800';
      case 'arbitrage': return 'bg-blue-100 text-blue-800';
      case 'front_running': return 'bg-purple-100 text-purple-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">MEV Analytics</h1>
          <p className="mt-2 text-gray-600">
            Monitor Maximal Extractable Value (MEV) activity, detect suspicious patterns, and identify MEV bots on the network.
          </p>
        </div>

        {/* Navigation Tabs */}
        <div className="mb-8">
          <nav className="flex space-x-8">
            {[
              { id: 'overview', label: 'Overview' },
              { id: 'transactions', label: 'Suspicious Transactions' },
              { id: 'bots', label: 'MEV Bots' },
              { id: 'block-analysis', label: 'Block Analysis' },
            ].map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id as any)}
                className={`py-2 px-1 border-b-2 font-medium text-sm ${
                  activeTab === tab.id
                    ? 'border-indigo-500 text-indigo-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </nav>
        </div>

        {/* Time Range Selector */}
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">Time Range</label>
          <select
            value={timeRange}
            onChange={(e) => setTimeRange(Number(e.target.value))}
            className="border border-gray-300 rounded-md px-3 py-2 bg-white"
          >
            <option value={1}>Last 1 hour</option>
            <option value={6}>Last 6 hours</option>
            <option value={24}>Last 24 hours</option>
            <option value={168}>Last 7 days</option>
          </select>
        </div>

        {/* Overview Tab */}
        {activeTab === 'overview' && stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            <div className="bg-white rounded-lg p-6 shadow">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                    <span className="text-blue-600 font-bold text-sm">TX</span>
                  </div>
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">Total Transactions</dt>
                    <dd className="text-lg font-medium text-gray-900">{stats.total_transactions.toLocaleString()}</dd>
                  </dl>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg p-6 shadow">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <div className="w-8 h-8 bg-red-100 rounded-full flex items-center justify-center">
                    <span className="text-red-600 font-bold text-sm">MEV</span>
                  </div>
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">MEV Transactions</dt>
                    <dd className="text-lg font-medium text-gray-900">{stats.mev_transactions.toLocaleString()}</dd>
                  </dl>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg p-6 shadow">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <div className={`w-8 h-8 rounded-full flex items-center justify-center ${getMEVActivityColor(stats.network_health.mev_activity_level)}`}>
                    <span className="font-bold text-sm">%</span>
                  </div>
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">MEV Percentage</dt>
                    <dd className="text-lg font-medium text-gray-900">{stats.mev_percentage.toFixed(2)}%</dd>
                  </dl>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg p-6 shadow">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <div className="w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center">
                    <span className="text-purple-600 font-bold text-sm">🤖</span>
                  </div>
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">MEV Bots</dt>
                    <dd className="text-lg font-medium text-gray-900">{stats.top_mev_bots_count}</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Network Health Status */}
        {activeTab === 'overview' && stats && (
          <div className="bg-white rounded-lg p-6 shadow mb-8">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Network Health</h3>
            <div className="flex items-center space-x-4">
              <div className="flex items-center">
                <span className="text-sm font-medium text-gray-500 mr-2">MEV Activity Level:</span>
                <span className={`px-2 py-1 rounded-full text-xs font-medium ${getMEVActivityColor(stats.network_health.mev_activity_level)}`}>
                  {stats.network_health.mev_activity_level.toUpperCase()}
                </span>
              </div>
              <div className="flex items-center">
                <span className="text-sm font-medium text-gray-500 mr-2">Average Gas Price:</span>
                <span className="text-sm text-gray-900">{formatGasPrice(stats.average_gas_price)}</span>
              </div>
            </div>
          </div>
        )}

        {/* Suspicious Transactions Tab */}
        {activeTab === 'transactions' && (
          <div className="bg-white rounded-lg shadow">
            <div className="p-6 border-b border-gray-200">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Search Suspicious Transactions</h3>
              <div className="flex space-x-4">
                <input
                  type="text"
                  placeholder="Enter block number"
                  value={blockNumber}
                  onChange={(e) => setBlockNumber(e.target.value)}
                  className="flex-1 border border-gray-300 rounded-md px-3 py-2"
                />
                <button
                  onClick={() => loadSuspiciousTransactions(blockNumber)}
                  disabled={loading}
                  className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 disabled:opacity-50"
                >
                  {loading ? 'Loading...' : 'Search'}
                </button>
              </div>
            </div>

            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Transaction</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">MEV Type</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Gas Price</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">MEV Score</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">From</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {suspiciousTransactions.map((tx) => (
                    <tr key={tx.hash} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm font-mono text-gray-900">
                          {tx.hash.slice(0, 10)}...{tx.hash.slice(-8)}
                        </div>
                        <div className="text-sm text-gray-500">Block #{tx.block_number}</div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`px-2 py-1 rounded-full text-xs font-medium ${getMEVTypeColor(tx.mev_type)}`}>
                          {tx.mev_type.replace('_', ' ')}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {formatGasPrice(tx.gas_price)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center">
                          <div className={`h-2 w-16 bg-gray-200 rounded-full mr-2`}>
                            <div 
                              className={`h-2 rounded-full ${tx.mev_score > 7 ? 'bg-red-500' : tx.mev_score > 4 ? 'bg-yellow-500' : 'bg-green-500'}`}
                              style={{ width: `${Math.min(tx.mev_score * 10, 100)}%` }}
                            ></div>
                          </div>
                          <span className="text-sm text-gray-900">{tx.mev_score.toFixed(1)}</span>
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">
                        {tx.from_address.slice(0, 6)}...{tx.from_address.slice(-4)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {suspiciousTransactions.length === 0 && !loading && (
                <div className="text-center py-8 text-gray-500">
                  No suspicious transactions found. Try searching for a specific block number.
                </div>
              )}
            </div>
          </div>
        )}

        {/* MEV Bots Tab */}
        {activeTab === 'bots' && (
          <div className="bg-white rounded-lg shadow">
            <div className="p-6 border-b border-gray-200">
              <h3 className="text-lg font-medium text-gray-900">Identified MEV Bots</h3>
              <p className="mt-1 text-sm text-gray-600">
                Addresses showing suspicious MEV-like behavior patterns
              </p>
            </div>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Address</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total TXs</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">High Gas TXs</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">MEV Score</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Patterns</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {mevBots.map((bot) => (
                    <tr key={bot.address} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">
                        {bot.address.slice(0, 8)}...{bot.address.slice(-6)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {bot.transaction_count.toLocaleString()}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {bot.high_gas_transactions.toLocaleString()}
                        <span className="text-gray-500 ml-1">
                          ({((bot.high_gas_transactions / bot.transaction_count) * 100).toFixed(1)}%)
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center">
                          <div className="h-2 w-16 bg-gray-200 rounded-full mr-2">
                            <div 
                              className={`h-2 rounded-full ${bot.mev_score > 7 ? 'bg-red-500' : bot.mev_score > 4 ? 'bg-yellow-500' : 'bg-green-500'}`}
                              style={{ width: `${Math.min(bot.mev_score * 10, 100)}%` }}
                            ></div>
                          </div>
                          <span className="text-sm text-gray-900">{bot.mev_score.toFixed(1)}</span>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex flex-wrap gap-1">
                          {bot.suspicious_patterns.map((pattern) => (
                            <span
                              key={pattern}
                              className="px-2 py-1 bg-red-100 text-red-800 text-xs rounded-full"
                            >
                              {pattern.replace('_', ' ')}
                            </span>
                          ))}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {mevBots.length === 0 && (
                <div className="text-center py-8 text-gray-500">
                  No MEV bots identified in the current time range.
                </div>
              )}
            </div>
          </div>
        )}

        {/* Block Analysis Tab */}
        {activeTab === 'block-analysis' && (
          <div className="space-y-6">
            <div className="bg-white rounded-lg p-6 shadow">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Analyze Block for MEV Activity</h3>
              <div className="flex space-x-4">
                <input
                  type="text"
                  placeholder="Enter block number"
                  value={blockNumber}
                  onChange={(e) => setBlockNumber(e.target.value)}
                  className="flex-1 border border-gray-300 rounded-md px-3 py-2"
                />
                <button
                  onClick={() => analyzeBlock(blockNumber)}
                  disabled={loading}
                  className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 disabled:opacity-50"
                >
                  {loading ? 'Analyzing...' : 'Analyze Block'}
                </button>
              </div>
            </div>

            {blockAnalysis && (
              <div className="bg-white rounded-lg p-6 shadow">
                <h4 className="text-lg font-medium text-gray-900 mb-4">
                  Block #{blockAnalysis.block_number} Analysis
                </h4>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div className="text-center">
                    <div className="text-2xl font-bold text-gray-900">
                      {blockAnalysis.total_transactions}
                    </div>
                    <div className="text-sm text-gray-500">Total Transactions</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-orange-600">
                      {blockAnalysis.mev_transactions}
                    </div>
                    <div className="text-sm text-gray-500">MEV Transactions</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-red-600">
                      {blockAnalysis.mev_percentage?.toFixed(2) || 0}%
                    </div>
                    <div className="text-sm text-gray-500">MEV Percentage</div>
                  </div>
                </div>

                <div className="mt-6 grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div className="text-center">
                    <div className="text-lg font-semibold text-gray-900">
                      {blockAnalysis.high_gas_tx_count || 0}
                    </div>
                    <div className="text-sm text-gray-500">High Gas Transactions</div>
                  </div>
                  <div className="text-center">
                    <div className="text-lg font-semibold text-gray-900">
                      {blockAnalysis.sandwich_attacks || 0}
                    </div>
                    <div className="text-sm text-gray-500">Sandwich Attacks</div>
                  </div>
                  <div className="text-center">
                    <div className="text-lg font-semibold text-gray-900">
                      {blockAnalysis.average_gas_price ? formatGasPrice(blockAnalysis.average_gas_price) : 'N/A'}
                    </div>
                    <div className="text-sm text-gray-500">Average Gas Price</div>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
