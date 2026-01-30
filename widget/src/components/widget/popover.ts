import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import type { WidgetConfig, WidgetInit } from '@/lib/types';
import '../ui/card';
import '../ui/button';
import '../ui/scroll-area';
import '../icons/external-link';
import '../icons/x';

@customElement('widget-popover')
export class PopoverWidget extends LitElement {
  @property({ type: Object }) config!: WidgetConfig;
  @property({ type: Object }) init!: WidgetInit;

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

    .popover {
      width: 32rem;
      position: fixed;
      bottom: 5rem;
      right: 1rem;
      z-index: 9999;
      border-radius: var(--widget-border-radius, 0.5rem);
      border-width: var(--widget-border-width, 1px);
      border-style: solid;
      border-color: var(--widget-border-color, #e5e5e5);
      background-color: var(--widget-bg-color, #ffffff);
      color: var(--widget-font-color, #0a0a0b);
      font-family: var(--widget-font-family, inherit);
    }

    .actions {
      position: absolute;
      top: 0;
      right: 0;
      padding: 0.5rem;
      display: flex;
      align-items: center;
      gap: 0.25rem;
    }

    .icon-medium {
      width: 1rem;
      height: 1rem;
    }

    .title {
      font-size: 1.125rem;
      display: inline-flex;
      align-items: center;
    }

    .scroll-content {
      height: 32rem;
    }
  `;

  private handleClose() {
    this.dispatchEvent(new CustomEvent('close', { bubbles: true, composed: true }));
  }

  render() {
    const { title, description } = this.config;

    return html`
      <div
        class="popover"
        style="
          --widget-border-radius: ${this.config.widget_border_radius}px;
          --widget-border-width: ${this.config.widget_border_width}px;
          --widget-border-color: ${this.config.widget_border_color};
          --widget-bg-color: ${this.config.widget_bg_color};
          --widget-font-color: ${this.config.widget_font_color};
          --widget-font-family: ${this.init.font_family?.join(',') || 'inherit'};
        "
      >
        <div class="actions">
          ${!this.config.disable_release_page ? html`
            <ui-button size="icon" variant="ghost">
              <icon-external-link class="icon-medium" .size=${16}></icon-external-link>
            </ui-button>
          ` : ''}
          <ui-button size="icon" variant="ghost" @click=${this.handleClose}>
            <icon-x class="icon-medium" .size=${16}></icon-x>
          </ui-button>
        </div>
        <ui-card-header>
          <ui-card-title class="title">
            ${title}
          </ui-card-title>
          ${description ? html`
            <ui-card-description style="color: ${this.config.widget_font_color}">
              ${description}
            </ui-card-description>
          ` : ''}
        </ui-card-header>
        <ui-card-content>
          <ui-scroll-area
            class="scroll-content"
            style="--scroll-area-height: 32rem;"
          >
            <slot></slot>
          </ui-scroll-area>
        </ui-card-content>
      </div>
    `;
  }
}
