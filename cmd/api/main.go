package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/All-Things-Muchiri/server/internal/auth"
	"github.com/All-Things-Muchiri/server/internal/clerk"
	"github.com/All-Things-Muchiri/server/internal/config"
	"github.com/All-Things-Muchiri/server/internal/database"
	"github.com/All-Things-Muchiri/server/internal/handler"
	"github.com/All-Things-Muchiri/server/internal/repository"
	"github.com/All-Things-Muchiri/server/internal/router"
	"github.com/All-Things-Muchiri/server/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type application struct {
	config       config.Config
	authProvider auth.AuthProvider
	db           *sqlx.DB
}

func main() {
	// TODO: Create a util function to load and parse env variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env variables: %v", err)
	}

	required := map[string]string{
		"PORT":                               os.Getenv("PORT"),
		"DATABASE_PORT":                      os.Getenv("DATABASE_PORT"),
		"DATABASE_HOST":                      os.Getenv("DATABASE_HOST"),
		"DATABASE_USER":                      os.Getenv("DATABASE_USER"),
		"DATABASE_PASSWORD":                  os.Getenv("DATABASE_PASSWORD"),
		"DATABASE_NAME":                      os.Getenv("DATABASE_NAME"),
		"SSL_MODE":                           os.Getenv("SSL_MODE"),
		"AUTH_SECRET_KEY":                    os.Getenv("AUTH_SECRET_KEY"),
		"CLERK_USERS_WEBHOOK_SECRET":         os.Getenv("CLERK_USERS_WEBHOOK_SECRET"),
		"CLERK_ORGANIZATIONS_WEBHOOK_SECRET": os.Getenv("CLERK_ORGANIZATIONS_WEBHOOK_SECRET"),
		"CLERK_MEMBERS_WEBHOOK_SECRET":       os.Getenv("CLERK_MEMBERS_WEBHOOK_SECRET"),
	}
	for k, v := range required {
		if v == "" {
			log.Fatalf("%s environment variable is required but not set", k)
		}
	}

	secret := required["AUTH_SECRET_KEY"]

	authConfig := &config.AuthConfig{
		SecretKey: secret,
	}

	dbPort, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		log.Fatalf("Failed to parse DATABASE_PORT: %v", err)
	}

	db, err := database.New(&database.DBConfig{
		Host:     required["DATABASE_HOST"],
		Port:     dbPort,
		User:     required["DATABASE_USER"],
		Password: required["DATABASE_PASSWORD"],
		Database: required["DATABASE_NAME"],
		SSLMode:  required["SSL_MODE"],
	})

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// REPOSITORIES
	userRepo := repository.NewUserRepository(db)
	organizationsRepo := repository.NewOrganizationsRepository(db)
	membersRepo := repository.NewMembershipRepository(db)

	// SERVICES
	userService := service.NewUserService(userRepo)
	organizationsService := service.NewOrganizationsService(organizationsRepo)
	membersService := service.NewMembershipService(membersRepo)

	// HANDLERS
	whHandler := handler.NewUsersWebhookHandler(required["CLERK_USERS_WEBHOOK_SECRET"], userService)
	OrgsWhHandler := handler.NewOrganizationsWebhookHandler(required["CLERK_ORGANIZATIONS_WEBHOOK_SECRET"], organizationsService)
	MembersWhHandler := handler.NewMembershipWebhookHandler(required["CLERK_MEMBERS_WEBHOOK_SECRET"], membersService)

	cfg := config.LoadConfig(authConfig)
	clerkProvider := clerk.NewProvider(*cfg.AuthConfig)
	apiRouter := &router.AppRouter{
		AuthProvider:           clerkProvider,
		WhHandler:              whHandler,
		OrganizationsWhHandler: OrgsWhHandler,
		MembershipWhHandler:    MembersWhHandler,
	}

	srv := &application{
		config: *cfg,
		db:     db,
	}

	mux := apiRouter.Mount()
	log.Fatal(srv.run(mux))
}

func (a *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:    a.config.Addr,
		Handler: mux,
	}
	log.Printf("Server started on port %s", a.config.Addr)

	return srv.ListenAndServe()
}
