/**
 * Client ID management for analytics and tracking
 * 
 * Re-exports the getOrCreateClientId function from storage.ts
 * for backward compatibility with existing code.
 * 
 * The client ID is a persistent UUID stored in localStorage
 * that identifies the user across sessions for analytics purposes.
 */

import { getOrCreateClientId as getOrCreateClientIdFromStorage } from './storage';

/**
 * Get or create a persistent client ID
 * Uses localStorage with error handling for edge cases
 */
export function getOrCreateClientId(): string {
  return getOrCreateClientIdFromStorage();
}
