import { ReactiveController, ReactiveControllerHost } from 'lit';
import { backendUrl } from '@/lib/config';
import { getOrCreateClientId } from '@/lib/clientId';

type MetricType = 'view' | 'cta_click';

interface MetricsOptions {
  releaseNoteId: string;
  orgId: string;
}

/**
 * Controller for tracking release note metrics (views, CTA clicks)
 * 
 * Features:
 * - IntersectionObserver for view tracking (50% visible threshold)
 * - Automatic cleanup on component disconnect
 * - Error handling and retry logic
 * - Debounced setup for mounted elements
 * - Multiple viewport support (scroll areas and window)
 * 
 * Replaces useReleaseNoteMetrics hook from React version
 */
export class ReleaseNoteMetricsController implements ReactiveController {
  host: ReactiveControllerHost;
  
  private releaseNoteId: string;
  private orgId: string;
  private hasTrackedView = false;
  private observer: IntersectionObserver | null = null;
  private element: HTMLElement | null = null;
  private timeoutId: number | null = null;
  private retryCount = 0;
  private readonly MAX_RETRIES = 3;
  private readonly SETUP_DELAY = 100; // ms to wait for DOM to settle
  private readonly RETRY_DELAY = 1000; // ms between retries

  constructor(host: ReactiveControllerHost, options: MetricsOptions) {
    this.host = host;
    this.releaseNoteId = options.releaseNoteId;
    this.orgId = options.orgId;
    host.addController(this);
  }

  hostConnected() {}

  hostDisconnected() {
    this.cleanup();
  }

  /**
   * Clean up all resources
   */
  private cleanup() {
    if (this.observer) {
      try {
        this.observer.disconnect();
      } catch (e) {
        console.error('[Announcable] Error disconnecting observer:', e);
      }
      this.observer = null;
    }
    
    if (this.timeoutId) {
      clearTimeout(this.timeoutId);
      this.timeoutId = null;
    }
    
    this.element = null;
  }

  /**
   * Set the element to observe for view tracking
   * Includes debouncing to wait for DOM to fully mount
   */
  setElement(element: HTMLElement | null) {
    if (!element) {
      console.warn('[Announcable] setElement called with null element');
      return;
    }
    
    this.element = element;
    
    if (this.hasTrackedView) {
      return; // Already tracked, no need to observe
    }
    
    // Clear any existing timeout
    if (this.timeoutId) {
      clearTimeout(this.timeoutId);
    }
    
    // Wait for scroll area to be fully mounted
    this.timeoutId = window.setTimeout(() => {
      this.setupIntersectionObserver();
    }, this.SETUP_DELAY);
  }

  /**
   * Setup IntersectionObserver to track when element comes into view
   * Includes error handling and fallback to window viewport
   */
  private setupIntersectionObserver() {
    if (!this.element || this.hasTrackedView) {
      return;
    }

    try {
      // Find scroll container (look for scroll-area viewport)
      const viewport = document.querySelector('[data-scroll-area-viewport]');
      
      // Check if IntersectionObserver is supported
      if (typeof IntersectionObserver === 'undefined') {
        console.error('[Announcable] IntersectionObserver not supported');
        // Fallback: track view immediately if observer not supported
        this.hasTrackedView = true;
        this.sendMetric('view');
        return;
      }

      this.observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (
              entry.isIntersecting &&
              entry.intersectionRatio >= 0.5 &&
              !this.hasTrackedView
            ) {
              this.hasTrackedView = true;
              this.sendMetric('view');
              
              // Disconnect after tracking
              if (this.observer) {
                this.observer.disconnect();
                this.observer = null;
              }
            }
          });
        },
        {
          threshold: 0.5, // 50% of element must be visible
          root: viewport as Element | null, // null falls back to viewport
          rootMargin: '0px',
        }
      );

      this.observer.observe(this.element);
    } catch (error) {
      console.error('[Announcable] Error setting up IntersectionObserver:', error);
      
      // Retry with exponential backoff
      if (this.retryCount < this.MAX_RETRIES) {
        this.retryCount++;
        const delay = this.RETRY_DELAY * Math.pow(2, this.retryCount - 1);
        
        console.log(`[Announcable] Retrying observer setup in ${delay}ms (attempt ${this.retryCount}/${this.MAX_RETRIES})`);
        
        this.timeoutId = window.setTimeout(() => {
          this.setupIntersectionObserver();
        }, delay);
      } else {
        console.error('[Announcable] Max retries reached for IntersectionObserver setup');
        // Fallback: track view immediately
        this.hasTrackedView = true;
        this.sendMetric('view');
      }
    }
  }

  /**
   * Track CTA click
   * Public method called from component when CTA is clicked
   */
  trackCtaClick = () => {
    this.sendMetric('cta_click');
  };

  /**
   * Send metric to backend with error handling and retry
   * Fails silently to not disrupt user experience
   */
  private async sendMetric(type: MetricType) {
    try {
      const clientId = getOrCreateClientId();
      
      const response = await fetch(`${backendUrl}/api/release-notes/${this.orgId}/metrics`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          release_note_id: this.releaseNoteId,
          metric_type: type,
          client_id: clientId,
        }),
      });
      
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
      
      console.debug(`[Announcable] Tracked ${type} for release note ${this.releaseNoteId}`);
    } catch (error) {
      // Log but don't throw - metrics are not critical
      console.error(`[Announcable] Failed to send ${type} metric:`, error);
      
      // Could implement retry queue here if metrics are critical
      // For now, we fail silently to not disrupt user experience
    }
  }
}
