import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { ref, createRef, Ref } from 'lit/directives/ref.js';
import type { ReleaseNote, WidgetConfig } from '@/lib/types';
import { ReleaseNoteMetricsController } from '@/tasks/release-note-metrics';
import { ReleaseNoteLikesController } from '@/tasks/release-note-likes';
import { getOrCreateClientId } from '@/lib/clientId';
import './card';
import './skeleton';
import '../icons/thumbs-up';
import '../icons/external-link';

/**
 * Release Notes List container
 */
@customElement('release-notes-list')
export class ReleaseNotesList extends LitElement {
  static styles = css`
    :host {
      display: block;
    }

    .list {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }
  `;

  render() {
    return html`
      <div class="list">
        <slot></slot>
      </div>
    `;
  }
}

/**
 * Individual Release Note Entry
 */
@customElement('release-note-entry')
export class ReleaseNoteEntry extends LitElement {
  @property({ type: Object }) config!: WidgetConfig;
  @property({ type: Object }) releaseNote!: ReleaseNote;

  private metricsController!: ReleaseNoteMetricsController;
  private likesController!: ReleaseNoteLikesController;
  private cardRef: Ref<HTMLElement> = createRef();

  static styles = css`
    :host {
      display: block;
    }

    .content {
      width: 100%;
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .media-container {
      position: relative;
      width: 100%;
      aspect-ratio: 16 / 9;
    }

    .media-container iframe {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
    }

    .media-container img {
      width: 100%;
      height: auto;
    }

    .text {
      white-space: pre-wrap;
    }

    .actions {
      width: 100%;
      display: flex;
      justify-content: space-around;
      margin-top: 0.5rem;
    }

    .action-wrapper {
      width: 100%;
      display: flex;
      justify-content: center;
      margin: 0 auto;
    }

    .like-button {
      display: inline-flex;
      align-items: center;
      gap: 0.25rem;
      background: none;
      border: none;
      cursor: pointer;
      padding: 0;
      font: inherit;
      color: inherit;
    }

    .like-button:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .like-button span {
      font-size: 0.875rem;
    }

    .cta-link {
      text-decoration: none;
      color: inherit;
    }

    .cta-content {
      display: flex;
      align-items: center;
      gap: 0.25rem;
    }

    .icon-small {
      width: 0.75rem;
      height: 0.75rem;
    }

    .icon-filled {
      fill: currentColor;
    }
  `;

  connectedCallback() {
    super.connectedCallback();

    // Setup metrics controller
    this.metricsController = new ReleaseNoteMetricsController(this, {
      releaseNoteId: this.releaseNote.id,
      orgId: this.config.org_id,
    });

    // Setup likes controller
    const clientId = getOrCreateClientId();
    this.likesController = new ReleaseNoteLikesController(this, {
      releaseNoteId: this.releaseNote.id,
      orgId: this.config.org_id,
      clientId,
    });
  }

  updated() {
    // Set element for metrics tracking after render
    if (this.cardRef.value) {
      this.metricsController.setElement(this.cardRef.value as HTMLElement);
    }
  }

  private handleImageError(e: Event) {
    const img = e.target as HTMLImageElement;
    console.error(
      `Image failed to load for ${this.releaseNote.title}`,
      this.releaseNote.imageSrc,
      e
    );
    img.style.display = 'none';
  }

  render() {
    const ctaLabel = this.releaseNote.cta_label_override || this.config.cta_text;
    const baseUrl = this.config.release_page_baseurl;
    const ctaHref = this.releaseNote.cta_href_override || `${baseUrl}#${this.releaseNote.id}`;

    const clientId = getOrCreateClientId();

    return html`
      <ui-card
        ${ref(this.cardRef)}
        style="
          border-radius: ${this.config.release_note_border_radius}px;
          border-color: ${this.config.release_note_border_color};
          border-width: ${this.config.release_note_border_width}px;
          color: ${this.config.release_note_font_color};
          background-color: ${this.config.release_note_bg_color};
        "
      >
        <ui-card-header style="padding-bottom: 1rem;">
          <ui-card-title>${this.releaseNote.title}</ui-card-title>
          <ui-card-description style="color: ${this.config.release_note_font_color}">
            ${this.releaseNote.date || ''}
          </ui-card-description>
        </ui-card-header>
        <ui-card-content>
          <div class="content">
            ${this.releaseNote.media_link ? html`
              <div class="media-container">
                <iframe
                  src="${this.releaseNote.media_link}"
                  allow="fullscreen"
                  allowfullscreen
                  loading="lazy"
                  referrerpolicy="no-referrer"
                  sandbox="allow-scripts allow-presentation allow-same-origin"
                  title="${this.releaseNote.title}"
                ></iframe>
              </div>
            ` : this.releaseNote.imageSrc ? html`
              <div class="media-container">
                <img
                  src="${this.releaseNote.imageSrc}"
                  alt="${this.releaseNote.title}"
                  @error=${this.handleImageError}
                />
              </div>
            ` : ''}

            ${this.releaseNote.text ? html`
              <div class="text">${this.releaseNote.text}</div>
            ` : ''}

            ${(this.config.enable_likes || !this.releaseNote.hide_cta) ? html`
              <div class="actions">
                ${this.config.enable_likes ? html`
                  <div class="action-wrapper">
                    <button
                      class="like-button"
                      @click=${() => this.likesController.toggleLike()}
                      ?disabled=${!clientId}
                    >
                      <span>
                        ${this.likesController.isLiked
                          ? this.config.unlike_button_text
                          : this.config.like_button_text}
                      </span>
                      <icon-thumbs-up
                        class="icon-small ${this.likesController.isLiked ? 'icon-filled' : ''}"
                        .size=${12}
                      ></icon-thumbs-up>
                    </button>
                  </div>
                ` : ''}

                ${!this.releaseNote.hide_cta ? html`
                  <div class="action-wrapper">
                    <a
                      class="cta-link"
                      href="${ctaHref}"
                      target="_blank"
                      @click=${this.metricsController.trackCtaClick}
                    >
                      <span class="cta-content">
                        ${ctaLabel}
                        <icon-external-link class="icon-small" .size=${12}></icon-external-link>
                      </span>
                    </a>
                  </div>
                ` : ''}
              </div>
            ` : ''}
          </div>
        </ui-card-content>
      </ui-card>
    `;
  }
}

/**
 * Release Note Skeleton (loading state)
 */
@customElement('release-note-skeleton')
export class ReleaseNoteSkeleton extends LitElement {
  @property({ type: Object }) config!: WidgetConfig;

  static styles = css`
    :host {
      display: block;
    }

    .header-skeletons {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .content {
      width: 100%;
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .text-skeletons {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .center {
      width: 100%;
      display: flex;
      justify-content: center;
    }

    .sk-title {
      height: 1.75rem;
      width: 75%;
    }

    .sk-date {
      height: 1rem;
      width: 25%;
    }

    .sk-image {
      height: 12rem;
      width: 100%;
    }

    .sk-text {
      height: 1rem;
      width: 100%;
    }

    .sk-text-short {
      height: 1rem;
      width: 75%;
    }

    .sk-cta {
      height: 1rem;
      width: 6rem;
    }
  `;

  render() {
    const skeletonBgColor = this.config.widget_bg_color;
    const skeletonBorderRadius = this.config.release_note_border_radius;

    return html`
      <ui-card
        style="
          border-radius: ${this.config.release_note_border_radius}px;
          border-color: ${this.config.release_note_border_color};
          border-width: ${this.config.release_note_border_width}px;
          color: ${this.config.release_note_font_color};
          background-color: ${this.config.release_note_bg_color};
        "
      >
        <ui-card-header style="padding-bottom: 1rem;">
          <div class="header-skeletons">
            <ui-skeleton
              class="sk-title"
              style="
                background-color: ${skeletonBgColor};
                border-radius: ${skeletonBorderRadius}px;
              "
            ></ui-skeleton>
            <ui-skeleton
              class="sk-date"
              style="
                background-color: ${skeletonBgColor};
                border-radius: ${skeletonBorderRadius}px;
              "
            ></ui-skeleton>
          </div>
        </ui-card-header>
        <ui-card-content>
          <div class="content">
            <ui-skeleton
              class="sk-image"
              style="
                background-color: ${skeletonBgColor};
                border-radius: ${skeletonBorderRadius}px;
              "
            ></ui-skeleton>
            <div class="text-skeletons">
              <ui-skeleton
                class="sk-text"
                style="
                  background-color: ${skeletonBgColor};
                  border-radius: ${skeletonBorderRadius}px;
                "
              ></ui-skeleton>
              <ui-skeleton
                class="sk-text"
                style="
                  background-color: ${skeletonBgColor};
                  border-radius: ${skeletonBorderRadius}px;
                "
              ></ui-skeleton>
              <ui-skeleton
                class="sk-text"
                style="
                  background-color: ${skeletonBgColor};
                  border-radius: ${skeletonBorderRadius}px;
                "
              ></ui-skeleton>
              <ui-skeleton
                class="sk-text-short"
                style="
                  background-color: ${skeletonBgColor};
                  border-radius: ${skeletonBorderRadius}px;
                "
              ></ui-skeleton>
            </div>
            <div class="center">
              <ui-skeleton
                class="sk-cta"
                style="
                  background-color: ${skeletonBgColor};
                  border-radius: ${skeletonBorderRadius}px;
                "
              ></ui-skeleton>
            </div>
          </div>
        </ui-card-content>
      </ui-card>
    `;
  }
}
