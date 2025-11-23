const titleChangeEventName = "title-change";
const descriptionChangeEventName = "description-change";
const ctaChangeEventName = "cta-change";
const widgetBorderRadiusChangeEventName = "widget-border-radius-change";
const widgetBorderColorChangeEventName = "widget-border-color-change";
const widgetBorderWidthChangeEventName = "widget-border-width-change";
const widgetBackgroundColorChangeEventName = "widget-background-color-change";
const widgetTextColorChangeEventName = "widget-text-color-change";
const releaseNoteBorderRadiusChangeEventName =
  "release-note-border-radius-change";
const releaseNoteBorderColorChangeEventName =
  "release-note-border-color-change";
const releaseNoteBorderWidthChangeEventName =
  "release-note-border-width-change";
const releaseNoteBackgroundColorChangeEventName =
  "release-notes-background-color-change";
const releaseNoteTextColorChangeEventName = "release-notes-text-color-change";
const enableLikesChangeEventName = "enable-likes-change";
const likeButtonTextChangeEventName = "like-button-text-change";

document.addEventListener("alpine:init", () => {
  Alpine.data("form", () => ({
    onSubmitError: function (event) {
      toastError(event.detail.xhr.response);
    },
    onSubmitSuccess: function () {
      toastSuccess("Release note updated");
    },

    title: {
      ["@input"]() {
        this.$dispatch(titleChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    description: {
      ["@input"]() {
        this.$dispatch(descriptionChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    widgetBorderRadius: {
      ["@input"]() {
        this.$dispatch(widgetBorderRadiusChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    widgetBorderColor: {
      ["@input"]() {
        this.$dispatch(widgetBorderColorChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    widgetBorderWidth: {
      ["@input"]() {
        this.$dispatch(widgetBorderWidthChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    widgetBackgroundColor: {
      ["@input"]() {
        this.$dispatch(widgetBackgroundColorChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    widgetTextColor: {
      ["@input"]() {
        this.$dispatch(widgetTextColorChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    release_note_cta: {
      ["@input"]() {
        this.$dispatch(ctaChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    releaseNoteBorderRadius: {
      ["@input"]() {
        this.$dispatch(releaseNoteBorderRadiusChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    releaseNoteBorderColor: {
      ["@input"]() {
        this.$dispatch(releaseNoteBorderColorChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    releaseNoteBorderWidth: {
      ["@input"]() {
        this.$dispatch(releaseNoteBorderWidthChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },

    releaseNoteBackgroundColor: {
      ["@input"]() {
        this.$dispatch(releaseNoteBackgroundColorChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    releaseNoteTextColor: {
      ["@input"]() {
        this.$dispatch(releaseNoteTextColorChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    enableLikes: {
      ["@change"]() {
        this.$el.checked = this.$event.target.checked;
        this.$dispatch(enableLikesChangeEventName, {
          value: this.$event.target.checked,
        });
      },
      ["x-ref"]: "enableLikes",
    },
    likeButtonContainer: {
      [`@${enableLikesChangeEventName}.window`]() {
        this.$el.style.display = this.$event.detail.value ? "block" : "none";
      },
      ["x-init"]() {
        this.$el.style.display = this.$refs.enableLikes.checked
          ? "block"
          : "none";
      },
    },
    likeButtonText: {
      ["@input"]() {
        this.$dispatch(likeButtonTextChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
  }));

  Alpine.data(
    "widget",
    (
      title,
      description,
      cta,
      widgetBorderRadius,
      widgetBorderColor,
      widgetBorderWidth,
      widgetBgColor,
      widgetTextColor,
      releaseNoteBorderRadius,
      releaseNoteBorderColor,
      releaseNoteBorderWidth,
      releaseNoteBgColor,
      releaseNoteTextColor,
      enableLikes,
      likeButtonText,
    ) => ({
      title: {
        [`@${titleChangeEventName}.window`]() {
          this.$el.innerText = this.$event.detail.value;
        },
        [`@${widgetTextColorChangeEventName}.window`]() {
          this.$el.style.color = this.$event.detail.value;
        },
        ["x-init"]() {
          this.$el.innerText = decodeHtml(title);
          this.$el.style.color = widgetTextColor;
        },
      },
      description: {
        [`@${descriptionChangeEventName}.window`]() {
          this.$el.innerText = this.$event.detail.value;
        },
        [`@${widgetTextColorChangeEventName}.window`]() {
          this.$el.style.color = this.$event.detail.value;
        },
        ["x-init"]() {
          this.$el.innerText = decodeHtml(description) || "";
          this.$el.style.color = widgetTextColor;
        },
      },
      widget: {
        [`@${widgetBorderRadiusChangeEventName}.window`]() {
          this.$el.style.borderRadius = this.$event.detail.value + "px";
        },
        [`@${widgetBorderColorChangeEventName}.window`]() {
          this.$el.style.borderColor = this.$event.detail.value;
        },
        [`@${widgetBorderWidthChangeEventName}.window`]() {
          this.$el.style.borderWidth = this.$event.detail.value + "px";
        },
        [`@${widgetBackgroundColorChangeEventName}.window`]() {
          this.$el.style.backgroundColor = this.$event.detail.value;
        },
        ["x-init"]() {
          this.$el.style.borderRadius = widgetBorderRadius + "px";
          this.$el.style.borderColor = widgetBorderColor;
          this.$el.style.borderStyle = "solid";
          this.$el.style.borderWidth = widgetBorderWidth + "px";
          this.$el.style.backgroundColor = widgetBgColor;
        },
      },
      releaseNote: {
        [`@${releaseNoteBorderRadiusChangeEventName}.window`]() {
          this.$el.style.borderRadius = this.$event.detail.value + "px";
        },
        [`@${releaseNoteBorderColorChangeEventName}.window`]() {
          this.$el.style.borderColor = this.$event.detail.value;
        },
        [`@${releaseNoteBorderWidthChangeEventName}.window`]() {
          this.$el.style.borderWidth = this.$event.detail.value + "px";
        },
        [`@${releaseNoteBackgroundColorChangeEventName}.window`]() {
          this.$el.style.backgroundColor = this.$event.detail.value;
        },
        ["x-init"]() {
          this.$el.style.borderRadius = releaseNoteBorderRadius + "px";
          this.$el.style.borderColor = releaseNoteBorderColor;
          this.$el.style.borderStyle = "solid";
          this.$el.style.borderWidth = releaseNoteBorderWidth + "px";
          this.$el.style.backgroundColor = releaseNoteBgColor;
        },
      },
      releaseNoteText: {
        [`@${releaseNoteTextColorChangeEventName}.window`]() {
          this.$el.style.color = this.$event.detail.value;
        },
        ["x-init"]() {
          this.$el.style.color = releaseNoteTextColor;
        },
      },
      cta: {
        [`@${ctaChangeEventName}.window`]() {
          this.$el.innerText = this.$event.detail.value;
        },
        [`@${releaseNoteTextColorChangeEventName}.window`]() {
          this.$el.style.color = this.$event.detail.value;
        },
        ["x-init"]() {
          this.$el.innerText = decodeHtml(cta);
          this.$el.style.color = releaseNoteTextColor;
        },
      },
      likeButton: {
        [`@${enableLikesChangeEventName}.window`]() {
          console.log("widget -> likeButton");
          console.log(this.$event.detail.value);
          this.$el.style.display = this.$event.detail.value ? "block" : "none";
        },
        [`@${likeButtonTextChangeEventName}.window`]() {
          this.$el.innerText = this.$event.detail.value;
        },
        ["x-init"]() {
          this.$el.style.display =
            enableLikes === "true" || enableLikes === true ? "block" : "none";
          this.$el.innerText = likeButtonText || "Like";
          this.$el.style.color = releaseNoteTextColor;
        },
      },
    }),
  );
});

function decodeHtml(html) {
  const txt = document.createElement("textarea");
  txt.innerHTML = html;
  return txt.value;
}
