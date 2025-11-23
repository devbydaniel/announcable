// Import Alpine.js
import Alpine from 'alpinejs'

// Import HTMX
import 'htmx.org'

// Import Basecoat UI JavaScript components
import 'basecoat-css/all'

// Initialize Alpine.js
window.Alpine = Alpine
Alpine.start()

// Note: Toast notifications are handled automatically by Basecoat UI
// The basecoat-css package includes HTMX integration that:
// 1. Listens for HX-Trigger headers with "basecoat:toast" key
// 2. Automatically dispatches basecoat:toast custom events
// 3. Triggers toast notifications based on the config
// No additional event listeners needed!
