// Validation utilities for the frontend

export const isValidEthereumAddress = (address: string): boolean => {
  if (address.length !== 42) return false;
  if (!address.startsWith('0x')) return false;
  // Check if the rest are valid hex characters
  const hexPart = address.slice(2);
  return /^[0-9a-fA-F]+$/.test(hexPart);
};

export const isValidTransactionHash = (hash: string): boolean => {
  return hash.startsWith('0x') && hash.length === 66;
};

export const isValidBlockHash = (hash: string): boolean => {
  return hash.startsWith('0x') && hash.length === 66;
};

export const isValidBlockNumber = (blockNumber: string): boolean => {
  const num = parseInt(blockNumber, 10);
  return !isNaN(num) && num >= 0;
};

export const validateSearchQuery = (query: string): {
  isValid: boolean;
  type: 'address' | 'transaction' | 'block' | 'unknown';
  error?: string;
} => {
  if (!query || query.trim().length === 0) {
    return { isValid: false, type: 'unknown', error: 'Search query cannot be empty' };
  }

  const trimmedQuery = query.trim();

  // Check if it's a block number
  if (/^\d+$/.test(trimmedQuery)) {
    const blockNum = parseInt(trimmedQuery, 10);
    if (blockNum >= 0) {
      return { isValid: true, type: 'block' };
    }
  }

  // Check if it's an address
  if (isValidEthereumAddress(trimmedQuery)) {
    return { isValid: true, type: 'address' };
  }

  // Check if it's a transaction hash or block hash
  if (isValidTransactionHash(trimmedQuery) || isValidBlockHash(trimmedQuery)) {
    return { isValid: true, type: 'transaction' };
  }

  return { 
    isValid: false, 
    type: 'unknown', 
    error: 'Invalid format. Please enter a valid address, transaction hash, or block number.' 
  };
};
