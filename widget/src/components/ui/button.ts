import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';

@customElement('ui-button')
export class Button extends LitElement {
  @property({ type: String }) variant: 'default' | 'ghost' | 'link' = 'default';
  @property({ type: String }) size: 'default' | 'sm' | 'lg' | 'icon' = 'default';

  static styles = css`
    :host {
      display: inline-flex;
      align-items: center;
      justify-content: center;
    }

    button {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      white-space: nowrap;
      border-radius: 0.375rem;
      font-size: 0.875rem;
      font-weight: 500;
      transition: colors 0.2s;
      cursor: pointer;
      border: 1px solid transparent;
      background: none;
    }

    button:focus-visible {
      outline: none;
      ring: 2px solid var(--announcable-border);
      ring-offset: 2px;
    }

    button:disabled {
      pointer-events: none;
      opacity: 0.5;
    }

    /* Variants */
    .variant-default {
      background-color: #0a0a0b;
      color: #ffffff;
      border-color: var(--announcable-border);
    }

    .variant-default:hover {
      background-color: rgba(10, 10, 11, 0.9);
    }

    .variant-ghost:hover {
      background-color: rgba(0, 0, 0, 0.05);
    }

    .variant-link {
      text-decoration: underline;
      text-underline-offset: 4px;
    }

    .variant-link:hover {
      text-decoration: none;
    }

    /* Sizes */
    .size-default {
      height: 2.5rem;
      padding: 0.5rem 1rem;
    }

    .size-sm {
      height: 2.25rem;
      padding: 0.5rem 0.75rem;
    }

    .size-lg {
      height: 2.75rem;
      padding: 0.5rem 2rem;
    }

    .size-icon {
      height: 2.5rem;
      width: 2.5rem;
      padding: 0;
    }
  `;

  render() {
    return html`
      <button class="variant-${this.variant} size-${this.size}">
        <slot></slot>
      </button>
    `;
  }
}
