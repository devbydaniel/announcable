// Import HTMX
import htmx from 'htmx.org';

// Import Toastify
import Toastify from 'toastify-js';

// Import SweetAlert
import swal from 'sweetalert';

// Import Feather Icons
import feather from 'feather-icons';

// Expose libraries as globals
// Note: Alpine.js is loaded separately from CDN in root.html
window.htmx = htmx;
window.Toastify = Toastify;
window.swal = swal;
window.feather = feather;

// Initialize HTMX (it auto-initializes when imported)
// No additional initialization needed

// Export for potential future ES module usage
export { htmx, Toastify, swal, feather };
