import { LitElement, html, css } from 'lit';
import { customElement } from 'lit/decorators.js';

@customElement('ui-skeleton')
export class Skeleton extends LitElement {
  static styles = css`
    :host {
      display: block;
    }

    .skeleton {
      animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
      border-radius: 0.375rem;
      background-color: rgba(10, 10, 11, 0.1);
    }

    @keyframes pulse {
      0%, 100% {
        opacity: 1;
      }
      50% {
        opacity: 0.5;
      }
    }
  `;

  render() {
    return html`
      <div class="skeleton">
        <slot></slot>
      </div>
    `;
  }
}
