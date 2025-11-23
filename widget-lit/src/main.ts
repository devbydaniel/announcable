import { html, render } from 'lit';
import './index.css'; // CSS will be injected by vite-plugin-css-injected-by-js
import type { WidgetInit } from './lib/types';
import './app';

declare global {
  interface Window {
    announcable_init?: WidgetInit;
    AnnouncableWidget?: {
      init?: (config: WidgetInit) => void;
    };
    release_beacon_widget_init?: WidgetInit;
    ReleaseBeaconWidget?: {
      init?: (config: WidgetInit) => void;
    };
  }
}

/**
 * Generate custom font CSS based on font_family config
 * Applies to all text within the widget for consistent typography
 */
function generateFontCSS(fontFamily?: string[]): string {
  if (!fontFamily || fontFamily.length === 0) {
    return '';
  }

  // Format font family names, adding quotes if needed
  const formattedFonts = fontFamily.map((font) => {
    // Add quotes if font name contains spaces
    if (font.includes(' ') && !font.startsWith('"') && !font.startsWith("'")) {
      return `"${font}"`;
    }
    return font;
  });

  const fontFamilyString = formattedFonts.join(', ');

  return `
    .announcable-widget {
      font-family: ${fontFamilyString} !important;
    }
    
    .announcable-widget * {
      font-family: inherit !important;
    }
  `;
}

/**
 * Initialize the Announcable widget
 * 
 * Features:
 * - Creates Shadow DOM for style isolation
 * - Injects Tailwind styles
 * - Applies custom font family if provided
 * - Renders Lit app component
 * - Validates required config (org_id)
 */
function initialize(init: WidgetInit) {
  // Validate required config
  if (!init.org_id) {
    console.error('[Announcable] org_id is required to initialize widget');
    return;
  }

  try {
    // Create widget root element
    const widgetRoot = document.createElement('div');
    widgetRoot.id = 'announcable-widget-root';
    widgetRoot.className = 'announcable-widget';
    document.body.appendChild(widgetRoot);

    // Inject custom font styles if provided
    if (init.font_family && init.font_family.length > 0) {
      const fontStyleElement = document.createElement('style');
      fontStyleElement.textContent = generateFontCSS(init.font_family);
      document.head.appendChild(fontStyleElement);
      
      console.debug(
        `[Announcable] Applied custom fonts: ${init.font_family.join(', ')}`
      );
    }

    // Render Lit app directly (CSS is auto-injected by Vite plugin)
    render(
      html`<announcable-app .init=${init}></announcable-app>`,
      widgetRoot
    );

    console.debug('[Announcable] Widget initialized successfully');
  } catch (error) {
    console.error('[Announcable] Failed to initialize widget:', error);
  }
}

// Expose initialization function globally
window.AnnouncableWidget = {
  init: (config: WidgetInit) => {
    initialize(config);
  },
};

// Legacy support for release beacon branding
window.ReleaseBeaconWidget = {
  init: (config: WidgetInit) => {
    console.warn(
      '[Announcable] ReleaseBeaconWidget is deprecated, use AnnouncableWidget instead'
    );
    initialize(config);
  },
};

// Automatically initialize if config is present
if (window.announcable_init && window.AnnouncableWidget?.init) {
  console.log('[Announcable] Auto-initializing with announcable_init config');
  window.AnnouncableWidget.init(window.announcable_init);
} else if (
  window.release_beacon_widget_init &&
  window.ReleaseBeaconWidget?.init
) {
  console.log(
    '[Announcable] Auto-initializing with legacy release_beacon_widget_init config'
  );
  window.ReleaseBeaconWidget.init(window.release_beacon_widget_init);
} else {
  console.warn(
    '[Announcable] No auto-init config found. Call AnnouncableWidget.init() manually.'
  );
}
