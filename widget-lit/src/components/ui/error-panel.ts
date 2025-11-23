import { LitElement, html } from 'lit';
import { customElement } from 'lit/decorators.js';
import './card';

@customElement('ui-error-panel')
export class ErrorPanel extends LitElement {
  render() {
    return html`
      <ui-card>
        <ui-card-header>
          <ui-card-title>Error</ui-card-title>
        </ui-card-header>
        <ui-card-content>
          <p>There was an error loading the release notes.</p>
        </ui-card-content>
      </ui-card>
    `;
  }
}
