// Formatting utilities for the frontend

export const formatHash = (hash: string): string => {
  if (hash.length < 18) return hash;
  return `${hash.slice(0, 10)}...${hash.slice(-8)}`;
};

export const formatAddress = (address: string): string => {
  if (address.length < 10) return address;
  return `${address.slice(0, 6)}...${address.slice(-4)}`;
};

export const formatNumber = (num: number): string => {
  return num.toLocaleString();
};

export const formatValue = (value: string): string => {
  // Convert wei to ETH (1 ETH = 10^18 wei)
  const eth = parseFloat(value) / Math.pow(10, 18);
  if (eth === 0) return '0 ETH';
  if (eth < 0.001) return '<0.001 ETH';
  return `${eth.toFixed(6)} ETH`;
};

export const formatGasPrice = (gasPrice: string): string => {
  // Convert wei to gwei (1 gwei = 10^9 wei)
  const gwei = parseFloat(gasPrice) / Math.pow(10, 9);
  return `${gwei.toFixed(2)} gwei`;
};

export const formatTimestamp = (timestamp: string): string => {
  return new Date(timestamp).toLocaleString();
};

export const formatTimeAgo = (timestamp: string): string => {
  const now = new Date();
  const date = new Date(timestamp);
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (diffInSeconds < 60) {
    return `${diffInSeconds}s ago`;
  }
  
  const diffInMinutes = Math.floor(diffInSeconds / 60);
  if (diffInMinutes < 60) {
    return `${diffInMinutes}m ago`;
  }
  
  const diffInHours = Math.floor(diffInMinutes / 60);
  if (diffInHours < 24) {
    return `${diffInHours}h ago`;
  }
  
  const diffInDays = Math.floor(diffInHours / 24);
  return `${diffInDays}d ago`;
};

export const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`;
};

export const formatPercentage = (value: number, total: number): string => {
  if (total === 0) return '0%';
  return `${((value / total) * 100).toFixed(1)}%`;
};
