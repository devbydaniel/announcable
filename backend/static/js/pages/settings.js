// toggle the visibility of fields
document.addEventListener("alpine:init", () => {
  Alpine.data("landingPageSettings", (useCustomURL) => ({
    init() {
      const searchParams = new URLSearchParams(window.location.search);
      const focus = searchParams.get("focus");
      if (focus === "lp") {
        this.$refs.landingPageCard.classList.add("has-focus");
      }
    },
    useCustomURL,
    onSubmitError: function (event) {
      toastError(event.detail.xhr.response);
    },
    onSubmitSuccess: function () {
      toastSuccess("URL updated");
    },
  }));

  Alpine.data("widgetSettings", () => ({
    init() {
      const searchParams = new URLSearchParams(window.location.search);
      const focus = searchParams.get("focus");
      if (focus === "widget") {
        this.$refs.widgetCard.classList.add("has-focus");
      }
    },
    onSubmitError: function (event) {
      toastError(event.detail.xhr.response);
    },
    onSubmitSuccess: function () {
      toastSuccess("Widget ID regenerated");
    },
  }));

  Alpine.data("pwUpdate", () => ({
    onSubmitError: function (event) {
      toastError(event.detail.xhr.response);
    },
    onSubmitSuccess: function () {
      toastSuccess("Password updated");
    },
  }));
});

document
  .getElementById("custom-url-submit-button")
  .addEventListener("click", () => {
    document.getElementById("custom-url-form").requestSubmit();
  });

document
  .getElementById("pw-update-submit-button")
  .addEventListener("click", () => {
    document.getElementById("pw-update-form").requestSubmit();
  });
