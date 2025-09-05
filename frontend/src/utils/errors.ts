// Standardized error handling utilities for frontend

export interface ApiError {
  error: string;
  details?: string;
  code?: number;
}

export interface ErrorState {
  isError: boolean;
  message: string;
  details?: string;
  code?: number;
}

// Common error messages
export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'Network error: Unable to connect to the server. Please try again.',
  INVALID_ADDRESS: 'Invalid Ethereum address format. Address must be 42 characters long and start with 0x followed by hexadecimal characters.',
  INVALID_HASH: 'Invalid transaction hash format. Hash must be 66 characters long and start with 0x followed by hexadecimal characters.',
  INVALID_BLOCK: 'Invalid block number. Please enter a valid positive integer.',
  NOT_FOUND: 'Resource not found.',
  SERVER_ERROR: 'Server error occurred. Please try again later.',
  VALIDATION_ERROR: 'Validation failed. Please check your input.',
  UNKNOWN_ERROR: 'An unexpected error occurred.',
} as const;

// Error types for better categorization
export enum ErrorType {
  NETWORK = 'network',
  VALIDATION = 'validation',
  NOT_FOUND = 'not_found',
  SERVER = 'server',
  UNKNOWN = 'unknown',
}

// Create standardized error state
export const createErrorState = (
  message: string,
  details?: string,
  code?: number
): ErrorState => ({
  isError: true,
  message,
  details,
  code,
});

// Create success state (no error)
export const createSuccessState = (): ErrorState => ({
  isError: false,
  message: '',
});

// Parse API error response
export const parseApiError = async (response: Response): Promise<ApiError> => {
  try {
    const data = await response.json();
    return {
      error: data.error || ERROR_MESSAGES.SERVER_ERROR,
      details: data.details,
      code: response.status,
    };
  } catch {
    return {
      error: ERROR_MESSAGES.SERVER_ERROR,
      code: response.status,
    };
  }
};

// Handle fetch errors with standardized responses
export const handleFetchError = async (
  response: Response,
  operation: string
): Promise<ErrorState> => {
  if (response.status === 404) {
    return createErrorState(
      ERROR_MESSAGES.NOT_FOUND,
      `${operation} not found`
    );
  }

  if (response.status >= 500) {
    return createErrorState(
      ERROR_MESSAGES.SERVER_ERROR,
      `Server error during ${operation}`,
      response.status
    );
  }

  if (response.status >= 400) {
    const apiError = await parseApiError(response);
    return createErrorState(
      apiError.error,
      apiError.details,
      response.status
    );
  }

  return createErrorState(
    ERROR_MESSAGES.UNKNOWN_ERROR,
    `Unexpected error during ${operation}`,
    response.status
  );
};

// Handle network/connection errors
export const handleNetworkError = (operation: string): ErrorState => {
  return createErrorState(
    ERROR_MESSAGES.NETWORK_ERROR,
    `Failed to connect during ${operation}`
  );
};

// Handle validation errors
export const handleValidationError = (field: string, message: string): ErrorState => {
  return createErrorState(
    ERROR_MESSAGES.VALIDATION_ERROR,
    `${field}: ${message}`
  );
};

// Generic error handler for try-catch blocks
export const handleGenericError = (
  error: unknown,
  operation: string
): ErrorState => {
  if (error instanceof Error) {
    // Check if it's a network error
    if (error.message.includes('fetch') || error.message.includes('network')) {
      return handleNetworkError(operation);
    }
    
    return createErrorState(
      error.message,
      `Error during ${operation}`
    );
  }

  return createErrorState(
    ERROR_MESSAGES.UNKNOWN_ERROR,
    `Unknown error during ${operation}`
  );
};

// Utility to check if an error state represents a specific error type
export const isErrorType = (errorState: ErrorState, type: ErrorType): boolean => {
  switch (type) {
    case ErrorType.NETWORK:
      return errorState.message.includes('Network error') || 
             errorState.message.includes('Unable to connect');
    case ErrorType.VALIDATION:
      return errorState.message.includes('Validation failed') ||
             errorState.message.includes('Invalid');
    case ErrorType.NOT_FOUND:
      return errorState.message.includes('not found') ||
             errorState.code === 404;
    case ErrorType.SERVER:
      return errorState.message.includes('Server error') ||
             (errorState.code !== undefined && errorState.code >= 500);
    case ErrorType.UNKNOWN:
      return errorState.message.includes('unexpected') ||
             errorState.message.includes('Unknown');
    default:
      return false;
  }
};

// Format error for display
export const formatErrorForDisplay = (errorState: ErrorState): string => {
  if (!errorState.isError) return '';
  
  let displayMessage = errorState.message;
  
  if (errorState.details) {
    displayMessage += ` ${errorState.details}`;
  }
  
  return displayMessage;
};

// Retry utility for failed operations
export const withRetry = async <T>(
  operation: () => Promise<T>,
  maxRetries: number = 3,
  delay: number = 1000
): Promise<T> => {
  let lastError: Error;
  
  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      return await operation();
    } catch (error) {
      lastError = error instanceof Error ? error : new Error('Unknown error');
      
      if (attempt === maxRetries) {
        throw lastError;
      }
      
      // Wait before retrying
      await new Promise(resolve => setTimeout(resolve, delay * attempt));
    }
  }
  
  throw lastError!;
};
