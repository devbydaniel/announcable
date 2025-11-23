document.addEventListener("htmx:confirm", function (e) {
  if (!e.detail.target.hasAttribute("hx-confirm")) return;
  e.preventDefault();

  swal({
    title: "Are you sure?",
    text: e.detail.question,
    buttons: true,
    dangerMode: true,
  }).then((willDelete) => {
    if (willDelete) {
      if (willDelete) e.detail.issueRequest(true);
    }
  });
});
