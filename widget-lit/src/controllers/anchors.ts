import { ReactiveController, ReactiveControllerHost } from 'lit';

/**
 * Controller for managing anchor element references
 * 
 * Features:
 * - Queries and maintains references to anchor elements
 * - Watches for dynamically added/removed anchors
 * - Handles edge cases (invalid selectors, missing elements)
 * - Automatic cleanup on disconnect
 * 
 * Replaces useAnchorsRef hook from React version
 */
export class AnchorsController implements ReactiveController {
  host: ReactiveControllerHost;
  
  anchors: NodeListOf<HTMLElement> | null = null;
  
  private querySelector?: string;
  private mutationObserver: MutationObserver | null = null;
  private retryTimeout: number | null = null;
  private readonly RETRY_DELAY = 500; // ms to wait before retrying query

  constructor(host: ReactiveControllerHost, querySelector?: string) {
    this.host = host;
    this.querySelector = querySelector;
    host.addController(this);
  }

  hostConnected() {
    this.queryAnchors();
    
    // Setup observer for dynamic anchor elements
    if (this.querySelector && typeof MutationObserver !== 'undefined') {
      this.setupMutationObserver();
    } else if (this.querySelector && !this.anchors) {
      // Retry query after a delay if no elements found initially
      // Useful for SPAs that haven't finished initial render
      this.scheduleRetry();
    }
  }

  hostDisconnected() {
    this.cleanup();
  }

  /**
   * Query for anchor elements using the provided selector
   * Includes validation and error handling
   */
  private queryAnchors() {
    if (!this.querySelector) {
      return;
    }

    try {
      // Validate selector syntax
      document.querySelector(this.querySelector);
      
      const elements = document.querySelectorAll(
        this.querySelector
      ) as NodeListOf<HTMLElement>;

      if (!elements || elements.length === 0) {
        console.debug(
          `[Announcable] No anchor elements found for selector: ${this.querySelector}`
        );
        this.anchors = null;
        return;
      }

      // Only update if anchors actually changed
      const anchorsChanged = 
        !this.anchors || 
        this.anchors.length !== elements.length ||
        Array.from(this.anchors).some((anchor, i) => anchor !== elements[i]);

      if (anchorsChanged) {
        this.anchors = elements;
        console.debug(`[Announcable] Found ${elements.length} anchor element(s)`);
        this.host.requestUpdate();
      }
    } catch (error) {
      if (error instanceof DOMException && error.name === 'SyntaxError') {
        console.error(
          `[Announcable] Invalid CSS selector: ${this.querySelector}`,
          error
        );
      } else {
        console.error('[Announcable] Error querying anchor elements:', error);
      }
      this.anchors = null;
    }
  }

  /**
   * Setup MutationObserver to watch for DOM changes
   * Re-queries anchors when DOM structure changes
   */
  private setupMutationObserver() {
    if (!this.querySelector) return;

    try {
      this.mutationObserver = new MutationObserver(() => {
        this.queryAnchors();
      });

      this.mutationObserver.observe(document.body, {
        childList: true,
        subtree: true,
      });
    } catch (error) {
      console.error('[Announcable] Error setting up MutationObserver:', error);
    }
  }

  /**
   * Schedule a retry for querying anchors
   * Used when initial query finds no elements
   */
  private scheduleRetry() {
    if (this.retryTimeout) {
      clearTimeout(this.retryTimeout);
    }

    this.retryTimeout = window.setTimeout(() => {
      console.debug('[Announcable] Retrying anchor query...');
      this.queryAnchors();
    }, this.RETRY_DELAY);
  }

  /**
   * Clean up all resources
   */
  private cleanup() {
    if (this.mutationObserver) {
      try {
        this.mutationObserver.disconnect();
      } catch (error) {
        console.error('[Announcable] Error disconnecting MutationObserver:', error);
      }
      this.mutationObserver = null;
    }

    if (this.retryTimeout) {
      clearTimeout(this.retryTimeout);
      this.retryTimeout = null;
    }

    this.anchors = null;
  }
}
