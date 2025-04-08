package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/devbydaniel/release-notes-go/config"
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/domain/rbac"
	"github.com/devbydaniel/release-notes-go/internal/handler"
	"github.com/devbydaniel/release-notes-go/internal/logger"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/internal/objstore"
	"github.com/devbydaniel/release-notes-go/internal/stripeUtil"
	"github.com/devbydaniel/release-notes-go/static"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var (
	log = logger.Get()
	cfg = config.New()
)

func main() {
	log.Info().Msg("Starting application")
	initEnv()
	log.Info().Msg("Environment loaded")
	db := initDb()
	defer database.Close(db)
	defer logger.Cleanup() // Ensure logger is cleaned up on shutdown
	
	initStripe()
	objStore := initObjStore()
	mwHandler := mw.NewHandler(db)
	handler := handler.NewHandler(db, objStore)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handler.HandleHomePage)

	// ONBOARDING AND AUTH

	r.Route("/login", func(r chi.Router) {
		r.Get("/", handler.HandleLoginPage)
		r.Post("/", handler.HandleLogin)
	})

	r.Route("/register", func(r chi.Router) {
		r.Get("/", handler.HandleRegisterPage)
		r.Post("/", handler.HandleRegister)
	})

	r.Route("/verify-email", func(r chi.Router) {
		r.Get("/", handler.HandleEmailVerifyPage)
		r.Post("/", handler.HandleEmailVerifyResend)
	})

	r.Route("/invite-accept/{token}", func(r chi.Router) {
		r.Get("/", handler.HandleInviteAcceptPage)
		r.Post("/", handler.HandleInviteAccept)
	})

	// USER AND PASSWORDS

	r.With(
		mwHandler.Authenticate,
		mwHandler.Authorize(rbac.PermissionManageAccess),
	).Route("/users", func(r chi.Router) {
		r.Get("/", handler.HandleUsersPage)
		r.Delete("/{id}", handler.HandleUserDelete)
	})

	r.With(mwHandler.Authenticate).Route("/profile", func(r chi.Router) {
		r.Patch("/password", handler.HandlePwUpdate)
	})

	r.With(mwHandler.Authenticate, mwHandler.Authorize(rbac.PermissionManageAccess)).Route("/invites", func(r chi.Router) {
		r.Post("/", handler.HandleInvite)
		r.Delete("/{id}", handler.HandleInviteDelete)
	})

	r.Route("/forgot-pw", func(r chi.Router) {
		r.Get("/", handler.HandlePwForgotPage)
		r.Post("/", handler.HandlePwForgot)
	})

	r.Route("/reset-pw/{token}", func(r chi.Router) {
		r.Get("/", handler.HandlePwResetPage)
		r.Post("/", handler.HandlePwReset)
	})

	r.With(mwHandler.Authenticate).Route("/logout", func(r chi.Router) {
		r.Get("/", handler.HandleLogout)
	})

	// RELEASE NOTES

	r.With(
		mwHandler.Authenticate,
		mwHandler.Authorize(rbac.PermissionManageReleaseNote),
		mwHandler.WithSubscriptionStatus,
	).Route("/release-notes", func(r chi.Router) {
		r.Get("/", handler.HandleReleaseNotesPage)
		r.Post("/", handler.HandleReleaseNoteCreate)
		r.Get("/new", handler.HandleReleaseNoteCreatePage)
		r.Get("/{id}", handler.HandleReleaseNotePage)
		r.Patch("/{id}", handler.HandleReleaseNoteUpdate)
		r.Delete("/{id}", handler.HandleReleaseNoteDelete)
		r.Patch("/{id}/publish", handler.HandleReleaseNotePublish)
	})

	// WIDGET

	r.With(
		mwHandler.Authenticate,
		mwHandler.Authorize(rbac.PermissionManageReleaseNote),
		mwHandler.WithSubscriptionStatus,
	).Route("/widget-config", func(r chi.Router) {
		r.Get("/", handler.HandleWidgetPage)
		r.Patch("/", handler.HandleWidgetUpdate)
		r.Patch("/external-id", handler.HandleOrgExternalIdRegenerate)
		r.Patch("/base-url", handler.HandleReleasePageBaseUrlUpdate)
	})

	// RELEASE PAGE CONFIG

	r.With(
		mwHandler.Authenticate,
		mwHandler.Authorize(rbac.PermissionManageReleaseNote),
		mwHandler.WithSubscriptionStatus,
	).Route("/release-page-config", func(r chi.Router) {
		r.Get("/", handler.HandleReleasePageConfigPage)
		r.Patch("/", handler.HandleReleasePageConfigUpdate)
	})

	// SETTINGS

	r.With(
		mwHandler.Authenticate,
		mwHandler.Authorize(rbac.PermissionManageAccess),
		mwHandler.WithSubscriptionStatus,
	).Route("/settings", func(r chi.Router) {
		r.Get("/", handler.HandleSettingsPage)
	})

	// ADMIN DASHBOARD
	
	r.With(
		mwHandler.Authenticate,
		mwHandler.AuthorizeSuperAdmin,
	).Route("/admin", func(r chi.Router) {
		r.Get("/", handler.HandleAdminDashboard)
		r.Get("/organisations/{orgId}", handler.HandleAdminOrgDetails)
		r.Patch("/organisations/{orgId}", handler.HandleAdminOrgUpdate)
		r.Patch("/organisations/{orgId}/release-page", handler.HandleAdminOrgReleasePageUpdate)
		r.Post("/organisations/{orgId}/subscriptions", handler.HandleAdminOrgSubscriptionCreate)
		r.Delete("/organisations/{orgId}/subscriptions/{id}", handler.HandleAdminOrgSubscriptionDelete)
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
		r.Get("/release-notes/{orgId}", handler.HandleReleaseNotesServe)
		r.Get("/release-notes/{orgId}/status", handler.HandleReleaseNotesStatusServe)
		r.Post("/release-notes/{orgId}/metrics", handler.HandleReleaseNoteMetricCreate)
		r.Get("/release-notes/{orgId}/{releaseNoteId}/like", handler.HandleGetReleaseNoteLikeState)
		r.Post("/release-notes/{orgId}/{releaseNoteId}/like", handler.HandleReleaseNoteToggleLike)
		r.Get("/widget-config/{orgId}", handler.HandleWidgetConfigServe)
		r.Get("/img/*", handler.HandleObjStore)
	})

	// WIDGET SCRIPT

	r.Route("/widget", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		r.Get("/", handler.HandleWidgetjsServe)
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
		r.Get("/{orgSlug}", handler.HandleReleasePage)
	})

	// STRIPE
	r.With(mwHandler.Authenticate, mwHandler.Authorize(rbac.PermissionManageAccess)).Route("/payment", func(r chi.Router) {
		r.Post("/create-checkout-session", handler.HandleCheckoutSession)
		r.Post("/create-portal-session", handler.HandlePortalSession)
	})

	// Subscription confirmation page (no auth required)
	r.Get("/subscription/confirm", handler.HandleSubscriptionConfirm)
	r.Get("/subscription/cancel", handler.HandleSubscriptionCancel)

	r.Route("/stripe", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		r.Post("/webhook", handler.HandleWebhook)
	})

	// STATIC
	fs := http.FileServer(http.FS(static.Assets))
	r.Get("/static/*", http.StripPrefix("/static/", fs).ServeHTTP)
	if cfg.Env != "production" {
		// proxy to obj storage since Minio doesn't support different
		// urls for signing and accessing
		r.Get("/img/*", handler.HandleObjStore)
	}

	// OTHER

	r.NotFound(handler.HandleNotFound)

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

func initEnv() {
	log.Trace().Msg("initEnv")
	if err := godotenv.Load(".env"); err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
		os.Exit(1)
	}
	log.Info().Msg("Environment loaded")
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

func initStripe() {
	log.Trace().Msg("initStripe")
	stripeUtil.Setup()
}
