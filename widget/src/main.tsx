import React from "react";
import ReactDOM from "react-dom/client";
import styles from "./index.css?inline";
import type { WidgetInit } from "./lib/types";
import Providers from "./components/provider";
import App from "./App";

declare global {
  interface Window {
    release_beacon_widget_init: WidgetInit;
    announcable_widget_init: WidgetInit;
    AnnouncableWidget: {
      init: (config: WidgetInit) => void;
    };
  }
}

function initialize(init: WidgetInit) {
  if (!init.org_id) {
    console.error("org_id is required to initialize release notes widget");
    return;
  }
  const widgetRoot = document.createElement("div");
  widgetRoot.id = "announcable-widget-root";
  document.body.appendChild(widgetRoot);

  // Create a closed Shadow DOM
  const shadowRoot = widgetRoot.attachShadow({ mode: "closed" });

  // Inject styles into Shadow DOM
  const styleElement = document.createElement("style");
  styleElement.textContent = styles;
  shadowRoot.appendChild(styleElement);

  // Create a container for the React app inside the Shadow DOM
  const reactContainer = document.createElement("div");
  reactContainer.className = "announcable-widget";
  shadowRoot.appendChild(reactContainer);

  const root = ReactDOM.createRoot(reactContainer);
  root.render(
    <React.StrictMode>
      <Providers>
        <App init={init} />
      </Providers>
    </React.StrictMode>,
  );
}

// Expose initialization function globally
window.AnnouncableWidget = {
  init: (config: WidgetInit) => {
    initialize(config);
  },
};

// Automatically initialize if config is present
if (window.announcable_widget_init) {
  window.AnnouncableWidget.init(window.announcable_widget_init);
}

// Use release beacon legacy init for backwards compatibility
if (window.release_beacon_widget_init) {
  window.AnnouncableWidget.init(window.release_beacon_widget_init);
}
