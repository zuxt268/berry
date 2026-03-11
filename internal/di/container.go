package di

import (
	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/interface/adapter"
	"github.com/zuxt268/berry/internal/interface/handlers"
	xmiddleware "github.com/zuxt268/berry/internal/interface/middleware"
	"github.com/zuxt268/berry/internal/repository"
	"github.com/zuxt268/berry/internal/usecase"
)

type Container struct {
	AuthMiddleware       *xmiddleware.AuthMiddleware
	UserHandler          *handlers.UserHandler
	OperatorHandler      *handlers.OperatorHandler
	UserAuthHandler      *handlers.UserAuthHandler
	OperatorAuthHandler  *handlers.OperatorAuthHandler
	GA4AuthHandler       *handlers.GA4AuthHandler
	GSCAuthHandler       *handlers.GSCAuthHandler
	GBPAuthHandler       *handlers.GBPAuthHandler
	InstagramAuthHandler *handlers.InstagramAuthHandler
	LineAuthHandler      *handlers.LineAuthHandler
	ReportHandler        *handlers.ReportHandler
	DBClose              func()
}

func NewContainer() (*Container, error) {

	// Infrastructure
	db, err := infrastructure.NewMySQLConnection()
	if err != nil {
		return nil, err
	}

	dbClose := func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	dbDriver := infrastructure.NewDBDriver(db, db)

	// Adapters
	userOAuthAdapter := adapter.NewOAuthClient(config.Env.UserGoogleRedirectURL)
	operatorOAuthAdapter := adapter.NewOAuthClient(config.Env.OperatorGoogleRedirectURL)
	userSessionAdapter := adapter.NewSessionStore("user_session")
	operatorSessionAdapter := adapter.NewSessionStore("operator_session")

	// Repositories
	baseRepository := repository.NewBaseRepository(dbDriver)
	userRepository := repository.NewUserRepository(dbDriver)
	userSessionRepository := repository.NewUserSessionRepository(dbDriver)
	operatorRepository := repository.NewOperatorRepository(dbDriver)
	operatorSessionRepository := repository.NewOperatorSessionRepository(dbDriver)

	// GA4
	ga4OAuthAdapter := adapter.NewGA4OAuthClient(config.Env.GA4GoogleRedirectURL)
	ga4ConnectionRepository := repository.NewGA4ConnectionRepository(dbDriver)

	// GSC
	gscOAuthAdapter := adapter.NewGSCOAuthClient(config.Env.GSCGoogleRedirectURL)
	gscConnectionRepository := repository.NewGSCConnectionRepository(dbDriver)

	// GBP
	gbpOAuthAdapter := adapter.NewGBPOAuthClient(config.Env.GBPGoogleRedirectURL)
	gbpConnectionRepository := repository.NewGBPConnectionRepository(dbDriver)

	// Instagram
	instagramOAuthAdapter := adapter.NewInstagramOAuthClient(config.Env.InstagramRedirectURL)
	instagramConnectionRepository := repository.NewInstagramConnectionRepository(dbDriver)

	// LINE
	lineTokenAdapter := adapter.NewLineTokenAdapter()
	lineConnectionRepository := repository.NewLineConnectionRepository(dbDriver)

	// Daily Report Repositories
	ga4DailyReportRepository := repository.NewGA4DailyReportRepository(dbDriver)
	gscDailyReportRepository := repository.NewGSCDailyReportRepository(dbDriver)
	gbpDailyReportRepository := repository.NewGBPDailyReportRepository(dbDriver)
	instagramDailyReportRepository := repository.NewInstagramDailyReportRepository(dbDriver)
	lineDailyReportRepository := repository.NewLineDailyReportRepository(dbDriver)

	// Usecases
	userUsecase := usecase.NewUserUsecase(baseRepository, userRepository)
	operatorUsecase := usecase.NewOperatorUsecase(baseRepository, operatorRepository)
	userAuthUseCase := usecase.NewAuthUseCase(userOAuthAdapter, userRepository, userSessionRepository)
	operatorAuthUseCase := usecase.NewOperatorAuthUseCase(operatorOAuthAdapter, operatorRepository, operatorSessionRepository)
	ga4AuthUseCase := usecase.NewGA4AuthUseCase(ga4OAuthAdapter, ga4ConnectionRepository)
	gscAuthUseCase := usecase.NewGSCAuthUseCase(gscOAuthAdapter, gscConnectionRepository)
	gbpAuthUseCase := usecase.NewGBPAuthUseCase(gbpOAuthAdapter, gbpConnectionRepository)
	instagramAuthUseCase := usecase.NewInstagramAuthUseCase(instagramOAuthAdapter, instagramConnectionRepository)
	lineAuthUseCase := usecase.NewLineAuthUseCase(lineTokenAdapter, lineConnectionRepository)

	// Report Usecases
	ga4ReportUseCase := usecase.NewGA4ReportUseCase(ga4ConnectionRepository, ga4DailyReportRepository)
	gscReportUseCase := usecase.NewGSCReportUseCase(gscConnectionRepository, gscDailyReportRepository)
	gbpReportUseCase := usecase.NewGBPReportUseCase(gbpConnectionRepository, gbpDailyReportRepository)
	instagramReportUseCase := usecase.NewInstagramReportUseCase(instagramConnectionRepository, instagramDailyReportRepository)
	lineReportUseCase := usecase.NewLineReportUseCase(lineConnectionRepository, lineDailyReportRepository)

	// Handlers
	userHandler := handlers.NewUserHandler(userUsecase)
	operatorHandler := handlers.NewOperatorHandler(operatorUsecase)
	userAuthHandler := handlers.NewUserAuthHandler(userAuthUseCase, userSessionAdapter)
	operatorAuthHandler := handlers.NewOperatorAuthHandler(operatorAuthUseCase, operatorSessionAdapter)
	ga4AuthHandler := handlers.NewGA4AuthHandler(ga4AuthUseCase)
	gscAuthHandler := handlers.NewGSCAuthHandler(gscAuthUseCase)
	gbpAuthHandler := handlers.NewGBPAuthHandler(gbpAuthUseCase)
	instagramAuthHandler := handlers.NewInstagramAuthHandler(instagramAuthUseCase)
	lineAuthHandler := handlers.NewLineAuthHandler(lineAuthUseCase)
	reportHandler := handlers.NewReportHandler(ga4ReportUseCase, gscReportUseCase, gbpReportUseCase, instagramReportUseCase, lineReportUseCase)

	// Middleware
	authMiddleware := xmiddleware.NewAuthMiddleware(userAuthUseCase, operatorAuthUseCase, userSessionAdapter, operatorSessionAdapter)

	return &Container{
		AuthMiddleware:       authMiddleware,
		UserHandler:          userHandler,
		OperatorHandler:      operatorHandler,
		UserAuthHandler:      userAuthHandler,
		OperatorAuthHandler:  operatorAuthHandler,
		GA4AuthHandler:       ga4AuthHandler,
		GSCAuthHandler:       gscAuthHandler,
		GBPAuthHandler:       gbpAuthHandler,
		InstagramAuthHandler: instagramAuthHandler,
		LineAuthHandler:      lineAuthHandler,
		ReportHandler:        reportHandler,
		DBClose:              dbClose,
	}, nil
}