document.addEventListener("alpine:init", () => {
  Alpine.data("fileInput", (imgUrl) => ({
    imgUrl: imgUrl || null,
    deleteImg: false,
    onImgChange(event) {
      this.imgUrl = URL.createObjectURL(event.target.files[0]);
      this.deleteImg = false;
    },
    onImgDelete(event) {
      event.stopPropagation();
      this.imgUrl = null;
      this.deleteImg = true;
    },
  }));
});
