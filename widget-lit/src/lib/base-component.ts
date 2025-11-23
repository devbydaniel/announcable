import { LitElement } from 'lit';

/**
 * Base component class for all widget components
 * Provides common utilities and setup for Lit elements
 */
export class BaseComponent extends LitElement {
  /**
   * Disable Shadow DOM for this component
   * Override in subclasses if needed
   */
  createRenderRoot() {
    // Most components will use Shadow DOM from main.ts
    // This can be overridden in specific components if needed
    return this;
  }
}
