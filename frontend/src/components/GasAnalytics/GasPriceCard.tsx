'use client';

import React from 'react';

interface GasPriceCardProps {
  title: string;
  price: number;
  change?: number;
  description: string;
  priority: 'low' | 'standard' | 'fast';
}

export default function GasPriceCard({ title, price, change, description, priority }: GasPriceCardProps) {
  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'low':
        return 'border-green-500 bg-green-50';
      case 'standard':
        return 'border-yellow-500 bg-yellow-50';
      case 'fast':
        return 'border-red-500 bg-red-50';
      default:
        return 'border-gray-300 bg-gray-50';
    }
  };

  const getChangeColor = (change?: number) => {
    if (!change) return 'text-gray-500';
    return change > 0 ? 'text-red-500' : 'text-green-500';
  };

  return (
    <div className={`p-6 rounded-lg border-2 ${getPriorityColor(priority)} transition-all hover:shadow-md`}>
      <div className="flex justify-between items-start mb-2">
        <h3 className="text-lg font-semibold text-gray-800">{title}</h3>
        {change !== undefined && (
          <span className={`text-sm font-medium ${getChangeColor(change)}`}>
            {change > 0 ? '+' : ''}{change.toFixed(1)}%
          </span>
        )}
      </div>
      
      <div className="mb-2">
        <span className="text-3xl font-bold text-gray-900">{price}</span>
        <span className="text-lg text-gray-600 ml-1">gwei</span>
      </div>
      
      <p className="text-sm text-gray-600">{description}</p>
      
      <div className="mt-4 flex items-center">
        <div className={`w-3 h-3 rounded-full mr-2 ${
          priority === 'low' ? 'bg-green-500' : 
          priority === 'standard' ? 'bg-yellow-500' : 'bg-red-500'
        }`}></div>
        <span className="text-sm font-medium text-gray-700 capitalize">{priority} Priority</span>
      </div>
    </div>
  );
}
