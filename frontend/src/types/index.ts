// Shared TypeScript interfaces and types

export interface Block {
  number: number;
  hash: string;
  parent_hash: string;
  timestamp: string;
  gas_limit: number;
  gas_used: number;
  difficulty: string;
  total_difficulty: string;
  size: number;
  transaction_count: number;
  miner: string;
  extra_data: string;
  base_fee_per_gas?: string;
  created_at: string;
  updated_at: string;
}

export interface Transaction {
  hash: string;
  block_number: number;
  transaction_index: number;
  from_address: string;
  to_address: string | null;
  value: string;
  gas_limit: number;
  gas_used: number | null;
  gas_price: string;
  max_fee_per_gas?: string;
  max_priority_fee_per_gas?: string;
  nonce: number;
  input_data: string;
  status: number | null;
  contract_address: string | null;
  logs_bloom: string | null;
  created_at: string;
  updated_at: string;
}

export interface Address {
  address: string;
  balance: string;
  nonce: number;
  is_contract: boolean;
  contract_creator?: string;
  creation_transaction?: string;
  first_seen_block?: number;
  last_seen_block?: number;
  transaction_count: number;
  label?: string;
  tags: string[];
  created_at: string;
  updated_at: string;
}

export interface Pagination {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface BlocksResponse {
  blocks: Block[];
  pagination: Pagination;
}

export interface TransactionsResponse {
  transactions: Transaction[];
  pagination: Pagination;
}

export interface TransactionNode {
  id: string;
  address: string;
  label?: string;
  value: number;
  type: 'address' | 'contract';
}

export interface TransactionLink {
  source: string;
  target: string;
  value: number;
  hash: string;
  timestamp: string;
}

export interface TransactionFlowData {
  nodes: TransactionNode[];
  links: TransactionLink[];
}

export interface ApiError {
  error: string;
  details?: string;
}

export interface LoadingState {
  isLoading: boolean;
  error: string | null;
}

export type ApiStatus = 'loading' | 'healthy' | 'error';

export interface NetworkStats {
  latest_block?: number;
  total_blocks?: number;
  total_transactions?: number;
  avg_gas_price?: string;
}

export interface SearchResult {
  type: 'block' | 'transaction' | 'address';
  hash?: string;
  number?: number;
  block_number?: number;
  address?: string;
  url: string;
}

export interface SearchResponse {
  results: {
    block?: SearchResult;
    transaction?: SearchResult;
    address?: SearchResult;
  };
}
