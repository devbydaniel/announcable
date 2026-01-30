import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import type { WidgetConfig, WidgetInit } from '@/lib/types';
import '../ui/dialog';
import '../ui/button';
import '../ui/scroll-area';
import '../icons/external-link';
import '../icons/x';

@customElement('widget-modal')
export class ModalWidget extends LitElement {
  @property({ type: Object }) config!: WidgetConfig;
  @property({ type: Object }) init!: WidgetInit;
  @property({ type: Boolean }) isOpen = false;

  static styles = css`
    :host {
      display: block;
    }

    a {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      text-decoration: none;
      color: inherit;
    }

    .icon-medium {
      width: 1rem;
      height: 1rem;
    }

    .scroll-content {
      height: 32rem;
    }
  `;

  private handleClose() {
    this.dispatchEvent(new CustomEvent('close', { bubbles: true, composed: true }));
  }

  render() {
    return html`
      <ui-dialog
        .isOpen=${this.isOpen}
        .title=${this.config.title}
        .description=${this.config.description}
        @close=${this.handleClose}
        style="
          --dialog-border-radius: ${this.config.widget_border_radius}px;
          --dialog-border-color: ${this.config.widget_border_color};
          --dialog-border-width: ${this.config.widget_border_width}px;
          --dialog-bg-color: ${this.config.widget_bg_color};
          --dialog-font-color: ${this.config.widget_font_color};
          --dialog-font-family: ${this.init.font_family?.join(',') || 'inherit'};
        "
      >
        <div slot="actions">
          ${!this.config.disable_release_page ? html`
            <ui-button size="icon" variant="ghost">
              <a href="${this.config.release_page_baseurl}" target="_blank">
                <icon-external-link class="icon-medium" .size=${16}></icon-external-link>
              </a>
            </ui-button>
          ` : ''}
          <ui-button size="icon" variant="ghost" @click=${this.handleClose}>
            <icon-x class="icon-medium" .size=${16}></icon-x>
          </ui-button>
        </div>
        <ui-scroll-area
          class="scroll-content"
          style="--scroll-area-height: 32rem;"
        >
          <slot></slot>
        </ui-scroll-area>
      </ui-dialog>
    `;
  }
}
