import { LitElement, html, css } from "lit";
import { customElement, property } from "lit/decorators.js";

@customElement("ui-dialog")
export class Dialog extends LitElement {
  @property({ type: Boolean }) isOpen = false;
  @property({ type: String }) title = "";
  @property({ type: String }) description = "";

  static styles = css`
    :host {
      display: contents;
    }

    .backdrop {
      position: fixed;
      inset: 0;
      background-color: rgba(0, 0, 0, 0.7);
      backdrop-filter: blur(4px);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 9999;
    }

    dialog {
      border: none;
      background: transparent;
      padding: 1rem;
      max-width: 40rem;
      width: 100%;
      box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
    }

    .content {
      display: grid;
      gap: 1rem;
      background-color: var(--dialog-bg-color, #ffffff);
      color: var(--dialog-font-color, #0a0a0b);
      border-radius: var(--dialog-border-radius, 0.5rem);
      border-width: var(--dialog-border-width, 1px);
      border-style: solid;
      border-color: var(--dialog-border-color, #e5e5e5);
      padding: 1.5rem;
      font-family: var(--dialog-font-family, inherit);
    }

    .header {
      position: relative;
    }

    .actions {
      position: absolute;
      top: 0;
      right: 0;
      transform: translate(12px, -12px);
      display: flex;
      align-items: center;
      gap: 0.25rem;
    }

    h2 {
      font-size: 1.125rem;
      font-weight: 600;
      letter-spacing: -0.025em;
      margin: 0;
    }

    .description {
      margin-top: 0.375rem;
      font-size: 0.875rem;
      color: rgba(10, 10, 11, 0.6);
    }
  `;

  private handleBackdropClick(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      this.dispatchEvent(
        new CustomEvent("close", { bubbles: true, composed: true }),
      );
    }
  }

  private handleKeyDown(e: KeyboardEvent) {
    if (e.key === "Escape") {
      this.dispatchEvent(
        new CustomEvent("close", { bubbles: true, composed: true }),
      );
    }
  }

  connectedCallback() {
    super.connectedCallback();
    document.addEventListener("keydown", this.handleKeyDown.bind(this));
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    document.removeEventListener("keydown", this.handleKeyDown.bind(this));
  }

  render() {
    if (!this.isOpen) return html``;

    return html`
      <div class="backdrop" @click=${this.handleBackdropClick}>
        <dialog
          open
          role="dialog"
          aria-modal="true"
          aria-labelledby="dialog-title"
          tabindex="-1"
        >
          <div class="content">
            <div class="header">
              <div class="actions">
                <slot name="actions"></slot>
              </div>
              <h2 id="dialog-title">${this.title}</h2>
              ${this.description
                ? html` <p class="description">${this.description}</p> `
                : ""}
            </div>
            <slot></slot>
          </div>
        </dialog>
      </div>
    `;
  }
}
