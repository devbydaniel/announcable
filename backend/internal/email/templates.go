package email

import (
	"embed"
	"html/template"
)

//go:embed templates/*.html
var emailTemplates embed.FS

var (
	welcomeTmpl       *template.Template
	passwordResetTmpl *template.Template
	userInviteTmpl    *template.Template
)

func init() {
	base := "templates/base.html"

	welcomeTmpl = template.Must(
		template.ParseFS(emailTemplates, base, "templates/welcome.html"),
	)
	passwordResetTmpl = template.Must(
		template.ParseFS(emailTemplates, base, "templates/password-reset.html"),
	)
	userInviteTmpl = template.Must(
		template.ParseFS(emailTemplates, base, "templates/user-invitation.html"),
	)
}
