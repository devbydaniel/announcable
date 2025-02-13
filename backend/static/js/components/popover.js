document.addEventListener("alpine:init", () => {
  Alpine.data("popover", () => ({
    isVisible: false,
    togglePopover() {
      this.isVisible = !this.isVisible;
    },
    showPopover() {
      this.showPopover = true;
    },
    hideIfVisible() {
      if (this.isVisible) {
        this.isVisible = false;
      }
    },
  }));
});
