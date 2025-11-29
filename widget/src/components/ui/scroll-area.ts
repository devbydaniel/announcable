import { LitElement, html, css } from 'lit';
import { customElement } from 'lit/decorators.js';

/**
 * ScrollArea component - simplified version without Radix
 * Uses native CSS overflow with custom scrollbar styling
 */
@customElement('ui-scroll-area')
export class ScrollArea extends LitElement {
  static styles = css`
    :host {
      display: block;
      height: var(--scroll-area-height, auto);
    }

    .scroll-container {
      position: relative;
      overflow-y: auto;
      overflow-x: hidden;
      height: 100%;
      max-height: var(--scroll-area-max-height, none);
    }

    /* Custom scrollbar styling */
    .scroll-container::-webkit-scrollbar {
      width: 10px;
    }

    .scroll-container::-webkit-scrollbar-track {
      background: transparent;
    }

    .scroll-container::-webkit-scrollbar-thumb {
      background: var(--announcable-border);
      border-radius: 5px;
    }

    .scroll-container::-webkit-scrollbar-thumb:hover {
      background: rgba(229, 229, 229, 0.8);
    }
  `;

  render() {
    return html`
      <div
        class="scroll-container"
        data-scroll-area-viewport
      >
        <slot></slot>
      </div>
    `;
  }
}
