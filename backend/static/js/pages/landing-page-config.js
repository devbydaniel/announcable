const titleChangeEventName = "title-change";
const descriptionChangeEventName = "description-change";
const bgColorChangeEventName = "bg-color-change";
const textColorChangeEventName = "text-color-change";
const textColorMutedChangeEventName = "text-color-muted-change";
const brandPositionChangeEventName = "brand-position-change";

document.addEventListener("alpine:init", () => {
  Alpine.data("form", () => ({
    onSubmitError: function (event) {
      toastError(event.detail.xhr.response);
    },
    onSubmitSuccess: function () {
      toastSuccess("Landing page config updated");
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
    brandPosition: {
      ["@change"]() {
        this.$dispatch(brandPositionChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    bgColor: {
      ["@input"]() {
        this.$dispatch(bgColorChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    textColor: {
      ["@input"]() {
        this.$dispatch(textColorChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
    textColorMuted: {
      ["@input"]() {
        this.$dispatch(textColorMutedChangeEventName, {
          value: this.$event.target.value,
        });
      },
    },
  }));

  Alpine.data(
    "lpContainer",
    (
      title,
      description,
      bgColor,
      textColor,
      textColorMuted,
      brandPosition,
    ) => ({
      title: {
        [`@${titleChangeEventName}.window`]() {
          this.$el.innerText = this.$event.detail.value;
        },
        [`@${textColorChangeEventName}.window`]() {
          this.$el.style.color = this.$event.detail.value;
        },
        ["x-init"]() {
          this.$el.innerText = title;
          this.$el.style.color = textColor;
        },
      },
      description: {
        [`@${descriptionChangeEventName}.window`]() {
          this.$el.innerText = this.$event.detail.value;
        },
        [`@${textColorMutedChangeEventName}.window`]() {
          this.$el.style.color = this.$event.detail.value;
        },
        ["x-init"]() {
          this.$el.innerText = description || "";
          this.$el.style.color = textColorMuted;
        },
      },
      lp: {
        [`@${bgColorChangeEventName}.window`]() {
          this.$el.style.backgroundColor = this.$event.detail.value;
        },
        ["x-init"]() {
          this.$el.style.backgroundColor = bgColor;
        },
      },
      brand: {
        [`@${brandPositionChangeEventName}.window`]() {
          switch (this.$event.detail.value) {
            case "left":
              this.$refs.left.style.display = "block";
              this.$refs.top.style.display = "none";
              break;
            case "top":
              this.$refs.top.style.display = "block";
              this.$refs.left.style.display = "none";
              break;
          }
        },
        [`@${textColorMutedChangeEventName}.window`]() {
          this.$el.style.backgroundColor = this.$event.detail.value;
        },
        ["x-init"]() {
          switch (brandPosition) {
            case "left":
              this.$refs.left.style.display = "block";
              this.$refs.top.style.display = "none";
              break;
            case "top":
              this.$refs.top.style.display = "block";
              this.$refs.left.style.display = "none";
              break;
          }
          this.$el.style.backgroundColor = textColorMuted;
        },
      },
      text: {
        [`@${textColorChangeEventName}.window`]() {
          this.$el.style.color = this.$event.detail.value;
        },
        ["x-init"]() {
          this.$el.style.color = textColor;
        },
      },
      skeleton: {
        [`@${textColorMutedChangeEventName}.window`]() {
          this.$el.style.backgroundColor = this.$event.detail.value;
        },
        ["x-init"]() {
          this.$el.style.backgroundColor = textColorMuted;
        },
      },
    }),
  );
});
