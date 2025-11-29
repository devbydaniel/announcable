import { LitElement, html, css } from "lit";
import { customElement, property } from "lit/decorators.js";

/**
 * Simple indicator dot (for showing unseen items)
 */
@customElement("ui-indicator")
export class Indicator extends LitElement {
  static styles = css`
    :host {
      display: inline-block;
    }

    .indicator {
      background-color: #ef4444;
      border-radius: 50%;
      width: 0.375rem;
      height: 0.375rem;
      transform: translate(0.25rem, -0.25rem);
    }
  `;

  render() {
    return html`<div class="indicator"></div>`;
  }
}

/**
 * Indicator that attaches to an anchor element in the host page
 * Creates an indicator outside the shadow DOM
 */
@customElement("ui-anchor-indicator")
export class AnchorIndicator extends LitElement {
  @property({ type: Object }) anchorElement!: HTMLElement;

  private indicatorElement: HTMLDivElement | null = null;

  connectedCallback() {
    super.connectedCallback();
    this.attachIndicator();
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    this.removeIndicator();
  }

  private attachIndicator() {
    if (!this.anchorElement) return;

    const indicator = document.createElement("div");
    // Set styles directly (can't use classes outside shadow DOM)
    indicator.style.position = "absolute";
    indicator.style.top = "0";
    indicator.style.right = "0";
    indicator.style.backgroundColor = "red";
    indicator.style.borderRadius = "50%";
    indicator.style.width = "8px";
    indicator.style.height = "8px";
    indicator.style.transform = "translate(50%, -50%)";
    indicator.style.zIndex = "9999";

    this.indicatorElement = indicator;
    this.anchorElement.appendChild(indicator);
  }

  private removeIndicator() {
    if (this.indicatorElement && this.anchorElement) {
      try {
        this.anchorElement.removeChild(this.indicatorElement);
      } catch {
        // Element may already be removed
      }
      this.indicatorElement = null;
    }
  }

  // This component doesn't render anything in the shadow DOM
  render() {
    return html``;
  }
}
