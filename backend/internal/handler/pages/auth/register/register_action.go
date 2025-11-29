package register

import (
	"errors"
	"net/http"

	"github.com/devbydaniel/announcable/config"
	"github.com/devbydaniel/announcable/internal/domain/organisation"
	releasepageconfig "github.com/devbydaniel/announcable/internal/domain/release-page-configs"
	"github.com/devbydaniel/announcable/internal/domain/session"
	"github.com/devbydaniel/announcable/internal/domain/user"
	widgetconfigs "github.com/devbydaniel/announcable/internal/domain/widget-configs"
	"github.com/devbydaniel/announcable/internal/password"
	"github.com/devbydaniel/announcable/internal/ratelimit"
)

type registerForm struct {
	OrgName         string
	Email           string
	Password        string
	ConfirmPassword string
}

var registerRateLimiter = ratelimit.New(60, 5)

// HandleRegister handles POST /register/
func (h *Handlers) HandleRegister(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleRegister")
	userService := user.NewService(*user.NewRepository(h.deps.DB))
	orgService := organisation.NewService(*organisation.NewRepository(h.deps.DB))
	sessionService := session.NewService(*session.NewRepository(h.deps.DB))
	releasePageConfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.deps.DB, h.deps.ObjStore))
	widgetConfigService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.deps.DB))

	if err := r.ParseForm(); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	req := registerForm{
		OrgName:         r.FormValue("org"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm"),
	}
	h.deps.Log.Debug().Interface("req", req).Msg("Register request")

	// check rate limit
	if err := registerRateLimiter.Deduct(req.Email, 1); err != nil {
		h.deps.Log.Warn().Str("email", req.Email).Msg("Rate limit exceeded for register requests")
		http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
		return
	}
	if err := registerRateLimiter.Deduct(req.OrgName, 1); err != nil {
		h.deps.Log.Warn().Str("org_name", req.OrgName).Msg("Rate limit exceeded for register requests")
		http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
		return
	}

	if err := orgService.IsValidOrgName(req.OrgName); err != nil {
		h.deps.Log.Debug().Err(err).Msg("Invalid org name")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := password.IsValidPassword(req.Password); err != nil {
		h.deps.Log.Debug().Err(err).Msg("Invalid password")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Password != req.ConfirmPassword {
		h.deps.Log.Debug().Msg("Passwords do not match")
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// check if user already exists
	existing, err := userService.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, h.deps.DB.ErrRecordNotFound) {
		h.deps.Log.Debug().Err(err).Msg("Error creating user")
		http.Error(w, "Error creating user", http.StatusBadRequest)
		return
	}
	if existing != nil {
		h.deps.Log.Debug().Msg("User already exists")
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// create user
	user, err := userService.Create(req.Email, req.Password, false)
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error creating user")
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// create org with user as admin
	ou, err := orgService.CreateOrgWithAdmin(req.OrgName, user)
	if err != nil {
		userService.Delete(user.ID)
		http.Error(w, "Error creating organisation", http.StatusInternalServerError)
		return
	}

	// create release page config
	if _, err := releasePageConfigService.Init(ou.Organisation.ID, ou.Organisation.Name); err != nil {
		h.deps.Log.Warn().Err(err).Msg("Error creating release page config")
	}

	// create widget config
	if _, err := widgetConfigService.Init(ou.Organisation.ID); err != nil {
		h.deps.Log.Warn().Err(err).Msg("Error creating widget config")
	}

	cfg := config.New()

	if cfg.IsEmailEnabled() {
		// Email enabled: send verification email
		token := sessionService.CreateToken()
		if err := sessionService.Create(token, user.ID); err != nil {
			http.Error(w, "Error creating session", http.StatusInternalServerError)
			return
		}

		if err := userService.SendVerifcationEmail(user, token); err != nil {
			http.Error(w, "Error sending verification email", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/verify-email")
	} else {
		// Email disabled: auto-verify and redirect to login
		if err := userService.VerifyEmail(user.ID); err != nil {
			h.deps.Log.Error().Err(err).Msg("Error verifying email")
			http.Error(w, "Error completing registration", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/login?success=Registration+successful")
	}

	w.WriteHeader(http.StatusCreated)
}
