// toggle the visibility of fields
document.addEventListener("alpine:init", () => {
  Alpine.data(
    "form",
    (
      textWebsiteOverrideIsChecked = false,
      hideCtaIsChecked = false,
      ctaLabelOverrideIsChecked = false,
      ctaUrlOverrideIsChecked = false,
      mediaType,
    ) => ({
      textWebsiteOverrideIsChecked,
      hideCtaIsChecked,
      ctaLabelOverrideIsChecked,
      ctaUrlOverrideIsChecked,
      mediaType,
      onSubmitError: function (event) {
        toastError(event.detail.xhr.response);
      },
      onSubmitSuccess: function () {
        toastSuccess("Release note updated");
      },
    }),
  );
});

// submit form
document.getElementById("submit-button").addEventListener("click", () => {
  document.getElementById("form").requestSubmit();
});
