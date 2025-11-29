import { LitElement, html } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import type { WidgetConfig, WidgetInit } from '@/lib/types';
import './popover';
import './modal';
import './sidebar';

/**
 * Main widget component that renders the appropriate widget type
 */
@customElement('widget-container')
export class WidgetContainer extends LitElement {
  @property({ type: Object }) config!: WidgetConfig;
  @property({ type: Object }) init!: WidgetInit;
  @property({ type: Boolean }) isOpen = false;

  private handleClose() {
    this.dispatchEvent(new CustomEvent('close', { bubbles: true, composed: true }));
  }

  render() {
    const widgetType = this.config?.widget_type || 'popover';

    switch (widgetType) {
      case 'sidebar':
        return html`
          <widget-sidebar
            .config=${this.config}
            .init=${this.init}
            .isOpen=${this.isOpen}
            @close=${this.handleClose}
          >
            <slot></slot>
          </widget-sidebar>
        `;
      
      case 'modal':
        return html`
          <widget-modal
            .config=${this.config}
            .init=${this.init}
            .isOpen=${this.isOpen}
            @close=${this.handleClose}
          >
            <slot></slot>
          </widget-modal>
        `;
      
      default:
        return html`
          <widget-popover
            .config=${this.config}
            .init=${this.init}
            @close=${this.handleClose}
          >
            <slot></slot>
          </widget-popover>
        `;
    }
  }
}
