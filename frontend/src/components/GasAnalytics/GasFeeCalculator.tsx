'use client';

import React, { useState, useEffect } from 'react';

interface GasFeeCalculatorProps {
  currentGasPrices: {
    slow: number;
    standard: number;
    fast: number;
  };
}

interface TransactionType {
  value: string;
  label: string;
  gasLimit: number;
  description: string;
}

const transactionTypes: TransactionType[] = [
  { value: 'transfer', label: 'ETH Transfer', gasLimit: 21000, description: 'Simple ETH transfer' },
  { value: 'erc20', label: 'ERC-20 Transfer', gasLimit: 65000, description: 'Token transfer' },
  { value: 'uniswap', label: 'Uniswap Swap', gasLimit: 150000, description: 'DEX token swap' },
  { value: 'nft', label: 'NFT Transfer', gasLimit: 85000, description: 'NFT transfer' },
  { value: 'contract', label: 'Contract Interaction', gasLimit: 200000, description: 'Smart contract call' },
];

export default function GasFeeCalculator({ currentGasPrices }: GasFeeCalculatorProps) {
  const [selectedType, setSelectedType] = useState<string>('transfer');
  const [customGasLimit, setCustomGasLimit] = useState<string>('');
  const [ethPrice, setEthPrice] = useState<number>(2500); // Mock ETH price
  const [useCustomGas, setUseCustomGas] = useState<boolean>(false);

  const selectedTransaction = transactionTypes.find(t => t.value === selectedType);
  const gasLimit = useCustomGas && customGasLimit ? parseInt(customGasLimit) : selectedTransaction?.gasLimit || 21000;

  const calculateFee = (gasPrice: number) => {
    const feeInWei = gasPrice * gasLimit * 1e9; // Convert gwei to wei
    const feeInEth = feeInWei / 1e18;
    const feeInUsd = feeInEth * ethPrice;
    
    return {
      eth: feeInEth,
      usd: feeInUsd,
      gwei: gasPrice * gasLimit
    };
  };

  const slowFee = calculateFee(currentGasPrices.slow);
  const standardFee = calculateFee(currentGasPrices.standard);
  const fastFee = calculateFee(currentGasPrices.fast);

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6">
      <h3 className="text-xl font-semibold text-gray-800 mb-6">Gas Fee Calculator</h3>
      
      {/* Transaction Type Selection */}
      <div className="mb-6">
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Transaction Type
        </label>
        <select
          value={selectedType}
          onChange={(e) => setSelectedType(e.target.value)}
          className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        >
          {transactionTypes.map((type) => (
            <option key={type.value} value={type.value}>
              {type.label} - {type.gasLimit.toLocaleString()} gas
            </option>
          ))}
        </select>
        {selectedTransaction && (
          <p className="text-sm text-gray-600 mt-1">{selectedTransaction.description}</p>
        )}
      </div>

      {/* Custom Gas Limit */}
      <div className="mb-6">
        <div className="flex items-center mb-2">
          <input
            type="checkbox"
            id="customGas"
            checked={useCustomGas}
            onChange={(e) => setUseCustomGas(e.target.checked)}
            className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
          />
          <label htmlFor="customGas" className="ml-2 text-sm font-medium text-gray-700">
            Use custom gas limit
          </label>
        </div>
        {useCustomGas && (
          <input
            type="number"
            value={customGasLimit}
            onChange={(e) => setCustomGasLimit(e.target.value)}
            placeholder="Enter gas limit"
            className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        )}
      </div>

      {/* Fee Calculations */}
      <div className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {/* Slow */}
          <div className="bg-green-50 border border-green-200 rounded-lg p-4">
            <div className="flex items-center mb-2">
              <div className="w-3 h-3 bg-green-500 rounded-full mr-2"></div>
              <h4 className="font-medium text-green-800">Slow</h4>
            </div>
            <div className="text-sm text-green-700">
              <p>{currentGasPrices.slow} gwei</p>
              <p className="font-semibold">{slowFee.eth.toFixed(6)} ETH</p>
              <p className="text-green-600">${slowFee.usd.toFixed(2)}</p>
            </div>
          </div>

          {/* Standard */}
          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
            <div className="flex items-center mb-2">
              <div className="w-3 h-3 bg-yellow-500 rounded-full mr-2"></div>
              <h4 className="font-medium text-yellow-800">Standard</h4>
            </div>
            <div className="text-sm text-yellow-700">
              <p>{currentGasPrices.standard} gwei</p>
              <p className="font-semibold">{standardFee.eth.toFixed(6)} ETH</p>
              <p className="text-yellow-600">${standardFee.usd.toFixed(2)}</p>
            </div>
          </div>

          {/* Fast */}
          <div className="bg-red-50 border border-red-200 rounded-lg p-4">
            <div className="flex items-center mb-2">
              <div className="w-3 h-3 bg-red-500 rounded-full mr-2"></div>
              <h4 className="font-medium text-red-800">Fast</h4>
            </div>
            <div className="text-sm text-red-700">
              <p>{currentGasPrices.fast} gwei</p>
              <p className="font-semibold">{fastFee.eth.toFixed(6)} ETH</p>
              <p className="text-red-600">${fastFee.usd.toFixed(2)}</p>
            </div>
          </div>
        </div>

        {/* Summary */}
        <div className="bg-gray-50 rounded-lg p-4 mt-4">
          <h5 className="font-medium text-gray-800 mb-2">Calculation Details</h5>
          <div className="text-sm text-gray-600 space-y-1">
            <p>Gas Limit: {gasLimit.toLocaleString()}</p>
            <p>ETH Price: ${ethPrice.toLocaleString()}</p>
            <p>Total Gas Cost Range: {slowFee.gwei.toLocaleString()} - {fastFee.gwei.toLocaleString()} gwei</p>
          </div>
        </div>
      </div>

      {/* Tips */}
      <div className="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
        <h5 className="font-medium text-blue-800 mb-2">💡 Tips</h5>
        <ul className="text-sm text-blue-700 space-y-1">
          <li>• Use "Slow" for non-urgent transactions to save on fees</li>
          <li>• "Standard" is recommended for most transactions</li>
          <li>• "Fast" is best for time-sensitive operations</li>
          <li>• Gas prices fluctuate based on network congestion</li>
        </ul>
      </div>
    </div>
  );
}
