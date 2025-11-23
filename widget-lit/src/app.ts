import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { TaskStatus } from '@lit/task';
import type { WidgetInit } from './lib/types';
import { WidgetToggleController } from './controllers/widget-toggle';
import { AnchorsController } from './controllers/anchors';
import { ReleaseNotesTask } from './tasks/release-notes';
import { WidgetConfigTask } from './tasks/widget-config';
import { ReleaseNoteStatusTask } from './tasks/release-note-status';
import './components/ui/button';
import './components/ui/indicator';
import './components/ui/error-panel';
import './components/ui/release-notes-list';
import './components/widget';
import './components/icons/gift';

/**
 * Main Announcable App Component
 * Orchestrates all controllers, tasks, and widget rendering
 */
@customElement('announcable-app')
export class AnnouncableApp extends LitElement {
  @property({ type: Object }) init!: WidgetInit;

  private toggleController!: WidgetToggleController;
  private anchorsController!: AnchorsController;
  private statusTask!: ReleaseNoteStatusTask;

  static styles = css`
    :host {
      display: block;
    }

    .floating-button {
      position: fixed;
      z-index: 50;
      bottom: 1rem;
      right: 1rem;
    }

    .indicator-wrapper {
      position: absolute;
      top: 0;
      right: 0;
    }

    .icon-medium {
      width: 1rem;
      height: 1rem;
    }
  `;

  connectedCallback() {
    super.connectedCallback();

    // Initialize controllers
    this.toggleController = new WidgetToggleController(
      this,
      this.init.anchor_query_selector
    );

    this.anchorsController = new AnchorsController(
      this,
      this.init.anchor_query_selector
    );

    this.statusTask = new ReleaseNoteStatusTask(this, this.init.org_id);
  }

  updated() {
    // Update anchor element datasets
    const hasUnseen = this.hasUnseenReleaseNotes();

    if (this.anchorsController.anchors) {
      this.anchorsController.anchors.forEach((anchor) => {
        this.updateIndicatorDataset(anchor, hasUnseen);
      });
    }

    // Handle instant open
    if (
      this.shouldInstantOpen() &&
      !this.toggleController.isOpen &&
      this.anchorsController.anchors
    ) {
      this.toggleController.setIsOpen(true);
      this.anchorsController.anchors.forEach((anchor) => {
        this.updateInstantOpenDataset(anchor, true);
      });
    }
  }

  private hasUnseenReleaseNotes(): boolean {
    if (this.statusTask.task.status !== TaskStatus.COMPLETE) return false;

    const releaseNoteStatus = this.statusTask.task.value;
    const lastOpened = this.toggleController.lastOpened;

    if (!releaseNoteStatus) return false;
    if (!lastOpened) return true;

    const lastOpenedTimestamp = parseInt(lastOpened);
    return releaseNoteStatus.some((note) => {
      if (!note.last_update_on) return false;
      const noteTimestamp = new Date(note.last_update_on).getTime();
      return noteTimestamp > lastOpenedTimestamp;
    });
  }

  private shouldInstantOpen(): boolean {
    if (this.statusTask.task.status !== TaskStatus.COMPLETE) return false;

    const releaseNoteStatus = this.statusTask.task.value;
    const lastOpened = this.toggleController.lastOpened;

    if (!releaseNoteStatus) return false;
    if (!lastOpened) return true;

    const lastOpenedTimestamp = parseInt(lastOpened);
    return releaseNoteStatus.some((note) => {
      if (!note.last_update_on || note.attention_mechanism !== 'instant_open') return false;
      const noteTimestamp = new Date(note.last_update_on).getTime();
      return noteTimestamp > lastOpenedTimestamp;
    });
  }

  private updateIndicatorDataset(
    anchorElement: HTMLElement | null,
    shouldDisplay: boolean
  ) {
    if (!anchorElement) return;
    const newValue = shouldDisplay ? 'true' : 'false';
    anchorElement.dataset.newReleaseNotes = newValue;
  }

  private updateInstantOpenDataset(
    anchorElement: HTMLElement | null,
    shouldDisplay: boolean
  ) {
    if (!anchorElement) return;
    const newValue = shouldDisplay ? 'true' : 'false';
    anchorElement.dataset.instantOpen = newValue;
  }

  render() {
    const hasUnseen = this.hasUnseenReleaseNotes();
    const shouldDisplayIndicator =
      this.anchorsController.anchors && !this.init.hide_indicator && hasUnseen;

    return html`
      ${shouldDisplayIndicator
        ? Array.from(this.anchorsController.anchors!).map(
            (anchor) => html`
              <ui-anchor-indicator .anchorElement=${anchor}></ui-anchor-indicator>
            `
          )
        : ''}

      ${!this.init.anchor_query_selector
        ? html`
            <ui-button
              class="floating-button"
              size="icon"
              style="border-radius: 0.5rem"
              @click=${() => this.toggleController.setIsOpen(!this.toggleController.isOpen)}
            >
              ${hasUnseen
                ? html`<ui-indicator class="indicator-wrapper"></ui-indicator>`
                : ''}
              <icon-gift class="icon-medium" .size=${16}></icon-gift>
            </ui-button>
          `
        : ''}

      ${this.toggleController.isOpen
        ? html`
            <widget-content
              .init=${this.init}
              .isOpen=${this.toggleController.isOpen}
              @close=${() => this.toggleController.setIsOpen(false)}
            ></widget-content>
          `
        : ''}
    `;
  }
}

/**
 * Widget Content Component
 * Handles data fetching and rendering of widget content
 */
@customElement('widget-content')
export class WidgetContent extends LitElement {
  @property({ type: Object }) init!: WidgetInit;
  @property({ type: Boolean }) isOpen = false;

  private notesTask!: ReleaseNotesTask;
  private configTask!: WidgetConfigTask;

  static styles = css`
    :host {
      display: block;
    }

    .container {
      position: relative;
    }
  `;

  connectedCallback() {
    super.connectedCallback();

    this.notesTask = new ReleaseNotesTask(this, this.init.org_id);
    this.configTask = new WidgetConfigTask(this, this.init.org_id);
  }

  private handleClose() {
    this.dispatchEvent(new CustomEvent('close', { bubbles: true, composed: true }));
  }

  render() {
    const configIsLoading = this.configTask.task.status === TaskStatus.PENDING;
    const notesError = this.notesTask.task.error;
    const configError = this.configTask.task.error;

    // Log errors
    if (notesError) console.error(notesError);
    if (configError) console.error(configError);

    const isReadyToMount = !configIsLoading && !notesError && !configError;

    if (!isReadyToMount) return html``;

    const widgetConfig = this.configTask.task.value;
    if (!widgetConfig) return html``;

    return html`
      <div class="container">
        <widget-container
          .config=${widgetConfig}
          .init=${this.init}
          .isOpen=${this.isOpen}
          @close=${this.handleClose}
        >
          ${notesError || configError
            ? html`<ui-error-panel></ui-error-panel>`
            : this.notesTask.task.render({
                pending: () => html`
                  <release-notes-list>
                    <release-note-skeleton .config=${widgetConfig}></release-note-skeleton>
                  </release-notes-list>
                `,
                complete: (releaseNotes) => html`
                  <release-notes-list>
                    ${releaseNotes.map(
                      (releaseNote) => html`
                        <release-note-entry
                          .config=${widgetConfig}
                          .releaseNote=${releaseNote}
                        ></release-note-entry>
                      `
                    )}
                  </release-notes-list>
                `,
                error: (e) => {
                  console.error('Error loading release notes:', e);
                  return html`<ui-error-panel></ui-error-panel>`;
                },
              })}
        </widget-container>
      </div>
    `;
  }
}
