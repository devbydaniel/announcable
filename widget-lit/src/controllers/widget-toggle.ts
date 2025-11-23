import { ReactiveController, ReactiveControllerHost } from 'lit';
import { getLastOpened, setLastOpened, getCurrentTimestamp } from '@/lib/storage';

/**
 * Controller for managing widget open/close state
 * 
 * Features:
 * - Manages widget open/close state
 * - Tracks last opened timestamp
 * - Attaches click listeners to anchor elements
 * - Handles edge cases (missing elements, duplicate listeners)
 * - Proper cleanup on disconnect
 * 
 * Replaces useWidgetToggle hook from React version
 */
export class WidgetToggleController implements ReactiveController {
  host: ReactiveControllerHost;
  
  isOpen = false;
  lastOpened: string | null = null;
  
  private querySelector?: string;
  private anchors: NodeListOf<HTMLElement> | null = null;
  private boundToggleWidget = this.toggleWidget.bind(this);
  private observerSetupAttempted = false;
  private mutationObserver: MutationObserver | null = null;

  constructor(host: ReactiveControllerHost, querySelector?: string) {
    this.host = host;
    this.querySelector = querySelector;
    host.addController(this);
    
    // Load last opened timestamp from localStorage
    this.lastOpened = getLastOpened();
  }

  hostConnected() {
    this.setupAnchorListeners();
    
    // Watch for dynamically added anchor elements
    if (this.querySelector && typeof MutationObserver !== 'undefined') {
      this.setupMutationObserver();
    }
  }

  hostDisconnected() {
    this.cleanup();
  }

  /**
   * Setup click listeners on anchor elements
   * Includes retry logic for dynamically loaded elements
   */
  private setupAnchorListeners() {
    if (!this.querySelector) {
      return;
    }

    try {
      const elements = document.querySelectorAll(this.querySelector) as NodeListOf<HTMLElement>;
      
      if (!elements || elements.length === 0) {
        if (!this.observerSetupAttempted) {
          console.warn(
            `[Announcable] No elements found for selector: ${this.querySelector}`
          );
          console.warn('[Announcable] Widget will use floating button instead');
        }
        return;
      }

      // Remove listeners from old anchors if they exist
      if (this.anchors) {
        this.removeAnchorListeners();
      }

      this.anchors = elements;
      
      this.anchors.forEach((anchor) => {
        // Prevent duplicate listeners
        anchor.removeEventListener('click', this.boundToggleWidget);
        anchor.addEventListener('click', this.boundToggleWidget);
        
        // Add visual indicator that element is clickable
        if (!anchor.style.cursor) {
          anchor.style.cursor = 'pointer';
        }
      });

      console.debug(`[Announcable] Attached listeners to ${this.anchors.length} anchor(s)`);
    } catch (error) {
      console.error('[Announcable] Error setting up anchor listeners:', error);
    }
  }

  /**
   * Setup MutationObserver to watch for dynamically added anchors
   * Useful for SPAs that render anchors after initial load
   */
  private setupMutationObserver() {
    if (!this.querySelector) return;

    try {
      this.mutationObserver = new MutationObserver(() => {
        // Re-setup listeners when DOM changes
        this.setupAnchorListeners();
      });

      this.mutationObserver.observe(document.body, {
        childList: true,
        subtree: true,
      });

      this.observerSetupAttempted = true;
    } catch (error) {
      console.error('[Announcable] Error setting up MutationObserver:', error);
    }
  }

  /**
   * Remove click listeners from anchor elements
   */
  private removeAnchorListeners() {
    if (this.anchors && this.anchors.length > 0) {
      this.anchors.forEach((anchor) => {
        try {
          anchor.removeEventListener('click', this.boundToggleWidget);
        } catch (error) {
          console.error('[Announcable] Error removing listener:', error);
        }
      });
    }
  }

  /**
   * Clean up all resources
   */
  private cleanup() {
    this.removeAnchorListeners();
    
    if (this.mutationObserver) {
      try {
        this.mutationObserver.disconnect();
      } catch (error) {
        console.error('[Announcable] Error disconnecting MutationObserver:', error);
      }
      this.mutationObserver = null;
    }
    
    this.anchors = null;
  }

  /**
   * Toggle widget open/close
   * Prevents default anchor behavior (e.g., navigation)
   */
  private toggleWidget(event?: Event) {
    if (event) {
      event.preventDefault();
      event.stopPropagation();
    }
    
    this.setIsOpen(!this.isOpen);
  }

  /**
   * Set widget open/close state
   * Updates last opened timestamp when opening
   */
  setIsOpen(value: boolean) {
    this.isOpen = value;
    
    // Update last opened timestamp when opening
    if (value) {
      const now = getCurrentTimestamp();
      this.lastOpened = now;
      const success = setLastOpened(now);
      
      if (!success) {
        console.warn('[Announcable] Failed to save last opened timestamp');
      }
    }
    
    // Request update to re-render host component
    this.host.requestUpdate();
  }
}
