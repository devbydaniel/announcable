/**
 * LocalStorage utilities for the widget
 * 
 * Handles persistence of:
 * - Last opened timestamp
 * - Client ID for analytics
 * - User preferences
 * 
 * Includes edge case handling for:
 * - LocalStorage not available (Safari private mode)
 * - QuotaExceededError
 * - Invalid data format
 */

const LAST_OPENED_KEY = 'announcable_last_opened';
const CLIENT_ID_KEY = 'announcable_client_id';

/**
 * Check if localStorage is available
 * Handles Safari private mode and other edge cases
 */
function isLocalStorageAvailable(): boolean {
  try {
    const test = '__storage_test__';
    localStorage.setItem(test, test);
    localStorage.removeItem(test);
    return true;
  } catch (e) {
    return false;
  }
}

/**
 * Safely get item from localStorage with error handling
 */
function safeGetItem(key: string): string | null {
  if (!isLocalStorageAvailable()) {
    console.warn('[Announcable] localStorage not available');
    return null;
  }
  
  try {
    return localStorage.getItem(key);
  } catch (e) {
    console.error('[Announcable] Error reading from localStorage:', e);
    return null;
  }
}

/**
 * Safely set item in localStorage with error handling
 */
function safeSetItem(key: string, value: string): boolean {
  if (!isLocalStorageAvailable()) {
    console.warn('[Announcable] localStorage not available');
    return false;
  }
  
  try {
    localStorage.setItem(key, value);
    return true;
  } catch (e) {
    // Handle quota exceeded or other errors
    if (e instanceof Error && e.name === 'QuotaExceededError') {
      console.error('[Announcable] localStorage quota exceeded');
    } else {
      console.error('[Announcable] Error writing to localStorage:', e);
    }
    return false;
  }
}

/**
 * Get the timestamp of when the widget was last opened
 * Returns null if not set or invalid
 */
export function getLastOpened(): string | null {
  const value = safeGetItem(LAST_OPENED_KEY);
  
  // Validate it's a valid timestamp
  if (value && !isNaN(parseInt(value))) {
    return value;
  }
  
  return null;
}

/**
 * Set the timestamp of when the widget was last opened
 * Returns true if successful
 */
export function setLastOpened(timestamp: string): boolean {
  // Validate timestamp format
  if (!timestamp || isNaN(parseInt(timestamp))) {
    console.error('[Announcable] Invalid timestamp format:', timestamp);
    return false;
  }
  
  return safeSetItem(LAST_OPENED_KEY, timestamp);
}

/**
 * Get current timestamp as string
 * Always returns a valid timestamp
 */
export function getCurrentTimestamp(): string {
  return Date.now().toString();
}

/**
 * Check if a date is after the last opened timestamp
 * Safely handles null values and invalid dates
 */
export function isAfterLastOpened(date: Date, lastOpened: string | null): boolean {
  if (!lastOpened) return true;
  
  try {
    const lastOpenedTimestamp = parseInt(lastOpened);
    if (isNaN(lastOpenedTimestamp)) {
      console.warn('[Announcable] Invalid last opened timestamp');
      return true;
    }
    
    const lastOpenedDate = new Date(lastOpenedTimestamp);
    
    // Check if date is valid
    if (isNaN(lastOpenedDate.getTime())) {
      console.warn('[Announcable] Invalid last opened date');
      return true;
    }
    
    return date > lastOpenedDate;
  } catch (e) {
    console.error('[Announcable] Error comparing dates:', e);
    return true; // Default to showing as new on error
  }
}

/**
 * Get or create a persistent client ID for analytics
 * Uses random UUID v4 format
 */
export function getOrCreateClientId(): string {
  const existing = safeGetItem(CLIENT_ID_KEY);
  
  if (existing) {
    return existing;
  }
  
  // Generate a simple UUID v4
  const uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
  
  safeSetItem(CLIENT_ID_KEY, uuid);
  return uuid;
}

/**
 * Clear all widget storage (for testing/debugging)
 * Not exposed in production builds
 */
export function clearAllStorage(): void {
  if (!isLocalStorageAvailable()) return;
  
  try {
    localStorage.removeItem(LAST_OPENED_KEY);
    localStorage.removeItem(CLIENT_ID_KEY);
  } catch (e) {
    console.error('[Announcable] Error clearing storage:', e);
  }
}
