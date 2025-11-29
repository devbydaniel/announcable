import { LitElement, svg, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';

@customElement('icon-x')
export class XIcon extends LitElement {
  @property({ type: String }) class = '';
  @property({ type: Number }) size = 24;

  static styles = css`
    :host {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      vertical-align: middle;
    }

    svg {
      display: block;
    }
  `;

  render() {
    return svg`
      <svg 
        class=${this.class}
        width=${this.size} 
        height=${this.size} 
        viewBox="0 0 24 24" 
        fill="none" 
        stroke="currentColor" 
        stroke-width="2" 
        stroke-linecap="round" 
        stroke-linejoin="round"
      >
        <path d="M18 6 6 18"></path>
        <path d="m6 6 12 12"></path>
      </svg>
    `;
  }
}
