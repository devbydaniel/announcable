import React from "react";
import ReactDOM from "react-dom/client";
import styles from "./index.css?inline";
import type { WidgetInit } from "./lib/types";
import Providers from "./components/provider";
import App from "./App";

declare global {
  interface Window {
    release_beacon_widget_init: WidgetInit;
    ReleaseBeaconWidget: {
      init: (config: WidgetInit) => void;
    };
  }
}

function initialize(init: WidgetInit) {
  const widgetRoot = document.createElement("div");
  widgetRoot.id = "release-beacon-widget-root";
  document.body.appendChild(widgetRoot);

  // Create a closed Shadow DOM
  const shadowRoot = widgetRoot.attachShadow({ mode: "closed" });

  // Inject styles into Shadow DOM
  const styleElement = document.createElement("style");
  styleElement.textContent = styles;
  shadowRoot.appendChild(styleElement);

  // Create a container for the React app inside the Shadow DOM
  const reactContainer = document.createElement("div");
  reactContainer.className = "release-beacon-widget";
  shadowRoot.appendChild(reactContainer);

  // Extract the backend url from the src attribute of the script tag
  const script = document.currentScript as HTMLScriptElement;
  const scriptSrc = script.src;
  const scriptSrcParts = scriptSrc?.split("/");
  const backendUrl = scriptSrcParts?.slice(0, -1).join("/");

  const root = ReactDOM.createRoot(reactContainer);
  root.render(
    <React.StrictMode>
      <Providers>
        <App init={init} backendUrl={backendUrl} />
      </Providers>
    </React.StrictMode>,
  );
}

// Expose initialization function globally
window.ReleaseBeaconWidget = {
  init: (config: WidgetInit) => {
    initialize(config);
  },
};

// Automatically initialize if config is present
if (window.release_beacon_widget_init) {
  window.ReleaseBeaconWidget.init(window.release_beacon_widget_init);
}
