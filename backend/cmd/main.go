package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"example.com/defect-control-system/internal/db"
	"example.com/defect-control-system/internal/handler"
	"example.com/defect-control-system/internal/middleware"
	"example.com/defect-control-system/internal/repository"
	"example.com/defect-control-system/internal/service"
)

func main() {
	viper.SetConfigFile("configs/config.yml")
	// allow overriding via environment variables
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()
	// common env bindings
	_ = viper.BindEnv("database.url", "DATABASE_URL")
	_ = viper.BindEnv("uploads.path", "UPLOADS_PATH")
	_ = viper.BindEnv("jwt.secret", "JWT_SECRET")

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	gdb, err := db.Connect()
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	// repositories
	userRepo := repository.NewUserRepository(gdb)
	projectRepo := repository.NewProjectRepository(gdb)
	defectRepo := repository.NewDefectRepository(gdb)
	attachRepo := repository.NewAttachmentRepository(gdb)

	// services & handlers
	jwtSecret := viper.GetString("jwt.secret")
	authSvc := service.NewAuthServiceWithSecret(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authSvc)

	// project/defect services & handlers
	projectSvc := service.NewProjectService(projectRepo)
	defectSvc := service.NewDefectService(defectRepo, projectRepo, userRepo)
	projectHandler := handler.NewProjectHandler(projectSvc, defectSvc)
	// attachments
	storageSvc := service.NewLocalStorage()
	attachHandler := handler.NewAttachmentHandler(storageSvc, attachRepo, defectSvc)
	userHandler := handler.NewUserHandler(userRepo, authSvc)
	// comments
	commentRepo := repository.NewCommentRepository(gdb)
	commentSvc := service.NewCommentService(commentRepo)
	commentHandler := handler.NewCommentHandler(commentSvc)

	r := gin.Default()

	// CORS for development (Vite dev server)
	// in production configure allowed origins properly
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// trust local proxy only (adjust in production)
	_ = r.SetTrustedProxies([]string{"127.0.0.1"})
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", middleware.JWTAuthMiddleware(), authHandler.Me)
		// projects
		projects := api.Group("/projects")
		projects.POST("/", middleware.JWTAuthMiddleware(), middleware.RequireRole("manager", "admin"), projectHandler.Create)
		projects.GET("/", projectHandler.List)
		projects.PATCH(":id", middleware.JWTAuthMiddleware(), middleware.RequireRole("manager", "admin"), projectHandler.UpdateProject)
		projects.POST(":id/defects", middleware.JWTAuthMiddleware(), middleware.RequireRole("engineer", "manager", "admin"), projectHandler.CreateDefect)
		projects.GET("/:id/defects", projectHandler.ListDefects)
		projects.GET("/:id", projectHandler.GetProject)
		projects.GET(":id/defects/:defectId", projectHandler.GetDefect)
		projects.PATCH(":id/defects/:defectId", middleware.JWTAuthMiddleware(), middleware.RequireRole("engineer", "manager", "admin"), projectHandler.UpdateDefect)
		// attachments (upload under defects)
		projects.POST(":id/attachments", middleware.JWTAuthMiddleware(), middleware.RequireRole("engineer", "manager", "admin"), attachHandler.Upload)
		api.GET("/attachments/:id", middleware.JWTAuthMiddleware(), attachHandler.Download)
		// listing attachments by defect
		api.GET("/attachments", middleware.JWTAuthMiddleware(), attachHandler.List)
		projects.GET(":id/defects/:defectId/attachments", middleware.JWTAuthMiddleware(), attachHandler.List)
		// users list for autocomplete
		api.GET("/users", userHandler.ListUsers)
		api.GET("/users/me", middleware.JWTAuthMiddleware(), userHandler.Me)
		api.PATCH("/users/me", middleware.JWTAuthMiddleware(), userHandler.UpdateMe)
		// admin: update arbitrary user
		api.PATCH("/users/:id", middleware.JWTAuthMiddleware(), middleware.RequireRole("admin"), userHandler.UpdateUser)
		// comments under defects
		projects.POST(":id/defects/:defectId/comments", middleware.JWTAuthMiddleware(), commentHandler.Create)
		projects.GET(":id/defects/:defectId/comments", commentHandler.List)
		// also expose a global comments list endpoint that accepts ?defect_id= for flexibility
		api.GET("/comments", commentHandler.List)
	}

	// Serve generated swagger files and Swagger UI
	r.Static("/docs", "./docs")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

	addr := viper.GetString("server.addr")
	if addr == "" {
		addr = ":8080"
	}
	if err := r.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
