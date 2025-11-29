import { LitElement, html, css } from 'lit';
import { customElement } from 'lit/decorators.js';

@customElement('ui-card')
export class Card extends LitElement {
  static styles = css`
    :host {
      display: block;
      border-radius: 0.5rem;
      border: 1px solid var(--announcable-border);
      background-color: var(--announcable-card);
      color: var(--announcable-card-foreground);
      box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
    }
  `;

  render() {
    return html`<slot></slot>`;
  }
}

@customElement('ui-card-header')
export class CardHeader extends LitElement {
  static styles = css`
    :host {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;
      padding: 1.5rem;
    }
  `;

  render() {
    return html`<slot></slot>`;
  }
}

@customElement('ui-card-title')
export class CardTitle extends LitElement {
  static styles = css`
    :host {
      display: block;
    }

    h3 {
      font-size: 1.5rem;
      font-weight: 600;
      line-height: 1;
      letter-spacing: -0.025em;
      margin: 0;
    }
  `;

  render() {
    return html`<h3><slot></slot></h3>`;
  }
}

@customElement('ui-card-description')
export class CardDescription extends LitElement {
  static styles = css`
    :host {
      display: block;
    }

    p {
      font-size: 0.875rem;
      color: rgba(10, 10, 11, 0.6);
      margin: 0;
    }
  `;

  render() {
    return html`<p><slot></slot></p>`;
  }
}

@customElement('ui-card-content')
export class CardContent extends LitElement {
  static styles = css`
    :host {
      display: block;
      padding: 1.5rem;
      padding-top: 0;
    }
  `;

  render() {
    return html`<slot></slot>`;
  }
}

@customElement('ui-card-footer')
export class CardFooter extends LitElement {
  static styles = css`
    :host {
      display: flex;
      align-items: center;
      padding: 1.5rem;
      padding-top: 0;
    }
  `;

  render() {
    return html`<slot></slot>`;
  }
}
