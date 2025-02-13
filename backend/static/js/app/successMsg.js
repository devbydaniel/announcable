// show success message when redirected with success url param
document.addEventListener("DOMContentLoaded", () => {
  const urlParams = new URLSearchParams(window.location.search);
  const success = urlParams.get("success");
  if (success) {
    toastSuccess(success);
    return;
  }
  const info = urlParams.get("info");
  if (info) {
    toastInfo(info);
    return;
  }
  const error = urlParams.get("error");
  if (error) {
    toastError(error);
    return;
  }
});
