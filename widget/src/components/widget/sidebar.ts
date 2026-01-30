import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import type { WidgetConfig, WidgetInit } from '@/lib/types';
import '../ui/card';
import '../ui/button';
import '../icons/external-link';
import '../icons/x';

@customElement('widget-sidebar')
export class SidebarWidget extends LitElement {
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

    .sidebar {
      width: 32rem;
      position: fixed;
      top: 0;
      right: 0;
      height: 100vh;
      transition: transform 0.3s ease-in-out;
      z-index: 9999;
      border-radius: 0;
      background-color: var(--widget-bg-color, #ffffff);
      color: var(--widget-font-color, #0a0a0b);
      font-family: var(--widget-font-family, inherit);
    }

    .sidebar.open {
      transform: translateX(0);
    }

    .sidebar.closed {
      transform: translateX(100%);
    }

    .actions {
      position: fixed;
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

    .content {
      height: calc(100vh - 120px);
      overflow-y: auto;
    }
  `;

  private handleClose() {
    this.dispatchEvent(new CustomEvent('close', { bubbles: true, composed: true }));
  }

  render() {
    const { title, description } = this.config;

    return html`
      <div
        class="sidebar ${this.isOpen ? 'open' : 'closed'}"
        style="
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
        <ui-card-content class="content">
          <slot></slot>
        </ui-card-content>
      </div>
    `;
  }
}
