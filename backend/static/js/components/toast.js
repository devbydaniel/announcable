function toastSuccess(message) {
  Toastify({
    text: message || "success",
    className: "toast toast__success",
    gravity: "bottom",
    position: "center",
    style: {
      background: "#40a02b",
      borderRadius: "0.5rem",
    },
  }).showToast();
}

function toastInfo(message) {
  Toastify({
    text: message,
    className: "toast toast__info",
    gravity: "bottom",
    position: "center",
    style: {
      background: "#1e66f5",
      borderRadius: "0.5rem",
    },
  }).showToast();
}

function toastError(message) {
  Toastify({
    text: message || "an error occurred",
    className: "toast toast__error",
    gravity: "bottom",
    position: "center",
    style: {
      background: "#d20f39",
      borderRadius: "0.5rem",
    },
  }).showToast();
}
