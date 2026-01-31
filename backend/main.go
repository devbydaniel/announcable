package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/devbydaniel/announcable/config"
	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/domain/rbac"
	apiShared "github.com/devbydaniel/announcable/internal/handler/api/shared"
	apiWidget "github.com/devbydaniel/announcable/internal/handler/api/widget"
	"github.com/devbydaniel/announcable/internal/handler/pages/admin/dashboard"
	"github.com/devbydaniel/announcable/internal/handler/pages/admin/organisation"
	"github.com/devbydaniel/announcable/internal/handler/pages/auth/invite_accept"
	"github.com/devbydaniel/announcable/internal/handler/pages/auth/login"
	"github.com/devbydaniel/announcable/internal/handler/pages/auth/logout"
	"github.com/devbydaniel/announcable/internal/handler/pages/auth/password_forgot"
	"github.com/devbydaniel/announcable/internal/handler/pages/auth/password_reset"
	"github.com/devbydaniel/announcable/internal/handler/pages/auth/register"
	"github.com/devbydaniel/announcable/internal/handler/pages/auth/verify_email"
	"github.com/devbydaniel/announcable/internal/handler/pages/public/home"
	"github.com/devbydaniel/announcable/internal/handler/pages/public/release_page"
	"github.com/devbydaniel/announcable/internal/handler/pages/public/widget_script"
	rnCreateHandler "github.com/devbydaniel/announcable/internal/handler/pages/release_notes/create"
	rnDetailHandler "github.com/devbydaniel/announcable/internal/handler/pages/release_notes/detail"
	rnListHandler "github.com/devbydaniel/announcable/internal/handler/pages/release_notes/list"
	releasePageConfigHandler "github.com/devbydaniel/announcable/internal/handler/pages/release_page/config"
	"github.com/devbydaniel/announcable/internal/handler/pages/settings/account"
	"github.com/devbydaniel/announcable/internal/handler/pages/users"
	widgetConfigHandler "github.com/devbydaniel/announcable/internal/handler/pages/widget/config"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	"github.com/devbydaniel/announcable/internal/logger"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/devbydaniel/announcable/internal/objstore"
	"github.com/devbydaniel/announcable/static"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	log = logger.Get()
	cfg = config.New()
)

func main() {
	log.Info().Msg("Starting application")
	if cfg.Env == "production" {
		runMigrations()
	}
	db := initDb()
	defer database.Close(db)
	defer logger.Cleanup() // Ensure logger is cleaned up on shutdown

	objStore := initObjStore()
	mwHandler := mw.NewHandler(db)

	// All handlers now use shared dependencies
	deps := shared.New(db, objStore)

	// Auth handlers
	loginHandler := login.New(deps)
	registerHandler := register.New(deps)
	verifyEmailHandler := verify_email.New(deps)
	inviteAcceptHandler := invite_accept.New(deps)
	passwordForgotHandler := password_forgot.New(deps)
	passwordResetHandler := password_reset.New(deps)
	logoutHandler := logout.New(deps)

	// User & Settings handlers
	usersHandler := users.New(deps)
	settingsHandler := account.New(deps)

	// Release Notes handlers
	rnListHandler := rnListHandler.New(deps)
	rnCreateHandler := rnCreateHandler.New(deps)
	rnDetailHandler := rnDetailHandler.New(deps)

	// Config handlers
	widgetHandler := widgetConfigHandler.New(deps)
	releasePageHandler := releasePageConfigHandler.New(deps)

	// Admin handlers
	adminDashboardHandler := dashboard.New(deps)
	adminOrgHandler := organisation.New(deps)

	// Public handlers
	homeHandler := home.New(deps)
	releasePagePublicHandler := release_page.New(deps)
	widgetScriptHandler := widget_script.New(deps)

	// API handlers
	widgetAPIHandler := apiWidget.New(deps)
	sharedAPIHandler := apiShared.New(deps)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", homeHandler.ServeHomePage)

	// ONBOARDING AND AUTH

	r.Route("/login", func(r chi.Router) {
		r.Get("/", loginHandler.ServeLoginPage)
		r.Post("/", loginHandler.HandleLogin)
	})

	r.Route("/register", func(r chi.Router) {
		r.Get("/", registerHandler.ServeRegisterPage)
		r.Post("/", registerHandler.HandleRegister)
	})

	r.Route("/verify-email", func(r chi.Router) {
		r.Get("/", verifyEmailHandler.ServeVerifyEmailPage)
		r.Post("/", verifyEmailHandler.HandleResend)
	})

	r.Route("/invite-accept/{token}", func(r chi.Router) {
		r.Get("/", inviteAcceptHandler.ServeInviteAcceptPage)
		r.Post("/", inviteAcceptHandler.HandleAccept)
	})

	// USER AND PASSWORDS

	r.With(
		mwHandler.Authenticate,
		mwHandler.Authorize(rbac.PermissionManageAccess),
	).Route("/users", func(r chi.Router) {
		r.Get("/", usersHandler.ServeUsersPage)
		r.Delete("/{id}", usersHandler.HandleUserDelete)
		r.Post("/{id}/password-reset", usersHandler.HandlePasswordResetTrigger)
	})

	r.With(mwHandler.Authenticate, mwHandler.Authorize(rbac.PermissionManageAccess)).Route("/invites", func(r chi.Router) {
		r.Post("/", usersHandler.HandleInviteCreate)
		r.Delete("/{id}", usersHandler.HandleInviteDelete)
	})

	r.Route("/forgot-pw", func(r chi.Router) {
		r.Get("/", passwordForgotHandler.ServeForgotPasswordPage)
		r.Post("/", passwordForgotHandler.HandleForgotPassword)
	})

	r.Route("/reset-pw/{token}", func(r chi.Router) {
		r.Get("/", passwordResetHandler.ServeResetPasswordPage)
		r.Post("/", passwordResetHandler.HandleResetPassword)
	})

	r.With(mwHandler.Authenticate).Route("/logout", func(r chi.Router) {
		r.Get("/", logoutHandler.HandleLogout)
	})

	// RELEASE NOTES

	r.With(
		mwHandler.Authenticate,
		mwHandler.Authorize(rbac.PermissionManageReleaseNote),
	).Route("/release-notes", func(r chi.Router) {
		r.Get("/", rnListHandler.ServeReleaseNotesListPage)
		r.Post("/", rnCreateHandler.HandleReleaseNoteCreate)
		r.Get("/new", rnCreateHandler.ServeReleaseNoteCreatePage)
		r.Get("/{id}", rnDetailHandler.ServeReleaseNoteDetailPage)
		r.Patch("/{id}", rnDetailHandler.HandleReleaseNoteUpdate)
		r.Delete("/{id}", rnDetailHandler.HandleReleaseNoteDelete)
		r.Patch("/{id}/publish", rnDetailHandler.HandleReleaseNotePublish)
	})

	// WIDGET

	r.With(
		mwHandler.Authenticate,
		mwHandler.Authorize(rbac.PermissionManageReleaseNote),
	).Route("/widget-config", func(r chi.Router) {
		r.Get("/", widgetHandler.ServeWidgetConfigPage)
		r.Patch("/", widgetHandler.HandleConfigUpdate)
	})

	// RELEASE PAGE CONFIG

	r.With(
		mwHandler.Authenticate,
		mwHandler.Authorize(rbac.PermissionManageReleaseNote),
	).Route("/release-page-config", func(r chi.Router) {
		r.Get("/", releasePageHandler.ServeReleasePageConfigPage)
		r.Patch("/", releasePageHandler.HandleConfigUpdate)
	})

	// SETTINGS

	r.With(
		mwHandler.Authenticate,
		mwHandler.Authorize(rbac.PermissionManageAccess),
	).Route("/settings", func(r chi.Router) {
		r.Get("/", settingsHandler.ServeSettingsPage)
		r.Patch("/password", settingsHandler.HandlePasswordUpdate)
		r.Patch("/widget-id", settingsHandler.HandleWidgetIdRegenerate)
		r.Patch("/release-page-url", settingsHandler.HandleReleasePageUrlUpdate)
	})

	// ADMIN DASHBOARD

	r.With(
		mwHandler.Authenticate,
		mwHandler.AuthorizeSuperAdmin,
	).Route("/admin", func(r chi.Router) {
		r.Get("/", adminDashboardHandler.ServeDashboardPage)
		r.Get("/organisations/{orgId}", adminOrgHandler.ServeOrganisationDetailsPage)
		r.Patch("/organisations/{orgId}", adminOrgHandler.HandleOrgUpdate)
		r.Patch("/organisations/{orgId}/release-page", adminOrgHandler.HandleReleasePageUpdate)
	})

	// API

	// !! this route path is hardcoded in the widget script
	r.Route("/api", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		r.Use(mwHandler.RateLimitPublicByIP)
		r.Use(mwHandler.RateLimitPublicByOrg)
		r.Get("/release-notes/{orgId}", widgetAPIHandler.HandleReleaseNotesServe)
		r.Get("/release-notes/{orgId}/status", widgetAPIHandler.HandleReleaseNotesStatusServe)
		r.Post("/release-notes/{orgId}/metrics", widgetAPIHandler.HandleReleaseNoteMetricCreate)
		r.Get("/release-notes/{orgId}/{releaseNoteId}/like", widgetAPIHandler.HandleGetReleaseNoteLikeState)
		r.Post("/release-notes/{orgId}/{releaseNoteId}/like", widgetAPIHandler.HandleReleaseNoteToggleLike)
		r.Get("/widget-config/{orgId}", widgetAPIHandler.HandleWidgetConfigServe)
		r.Get("/img/*", sharedAPIHandler.HandleObjStore)
	})

	// WIDGET SCRIPT

	r.Route("/widget", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}))
		r.Get("/", widgetScriptHandler.ServeWidgetScript)
	})

	// RELEASE PAGE

	// !! this route path is hardcoded in the widget script
	r.Route("/s", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers

		}))
		r.Get("/{orgSlug}", releasePagePublicHandler.ServeReleasePage)
	})

	// STATIC
	fs := http.FileServer(http.FS(static.Assets))
	r.Get("/static/*", http.StripPrefix("/static/", fs).ServeHTTP)
	if cfg.Env != "production" {
		// proxy to obj storage since Minio doesn't support different
		// urls for signing and accessing
		r.Get("/img/*", sharedAPIHandler.HandleObjStore)
	}

	// OTHER

	r.NotFound(sharedAPIHandler.HandleNotFound)

	// Create server with timeout configurations
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: r,
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server in the background
	go func() {
		log.Info().Msgf("Server listening on port %d", cfg.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Error().Err(err).Msg("Server error")
	case sig := <-shutdown:
		log.Info().Msgf("Start shutdown due to %s signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Graceful shutdown did not complete")
			if err := srv.Close(); err != nil {
				log.Error().Err(err).Msg("Could not stop server")
			}
		}
	}
}

func initDb() *database.DB {
	log.Trace().Msg("initDb")
	db, err := database.Connect()
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to database")
		os.Exit(1)
	}
	log.Info().Msg("Database connected")
	return db
}

func initObjStore() *objstore.ObjStore {
	log.Trace().Msg("initObjStore")
	ctx := context.Background()
	store, err := objstore.Init(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Could not initialize object store")
		os.Exit(1)
	}
	log.Info().Msg("Object store initialized")
	return store
}

func runMigrations() {
	log.Info().Msg("Running database migrations")
	dsn := "postgres://" + cfg.Postgres.User + ":" + cfg.Postgres.Password +
		"@" + cfg.Postgres.Host + ":" + strconv.Itoa(cfg.Postgres.Port) +
		"/" + cfg.Postgres.Name + "?sslmode=disable"
	if err := database.RunMigrations(dsn); err != nil {
		log.Error().Err(err).Msg("Could not run migrations")
		os.Exit(1)
	}
}
