package main

import (
	"log"

	"github.com/chatbot-saas/backend/internal/config"
	"github.com/chatbot-saas/backend/internal/handlers"
	"github.com/chatbot-saas/backend/internal/middleware"
	"github.com/chatbot-saas/backend/internal/repository"
	"github.com/chatbot-saas/backend/internal/services"
	"github.com/chatbot-saas/backend/pkg/openai"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize OpenAI client
	openaiClient := openai.NewClient(cfg.OpenAIAPIKey)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	chatbotRepo := repository.NewChatbotRepository(db)
	knowledgeRepo := repository.NewKnowledgeRepository(db)
	conversationRepo := repository.NewConversationRepository(db)

	// New advanced KB repositories
	folderRepo := repository.NewKBFolderRepository(db)
	tagRepo := repository.NewKBTagRepository(db)
	versionRepo := repository.NewKBVersionRepository(db)
	chunkRepo := repository.NewKBChunkRepository(db)
	syncRepo := repository.NewKBSyncRepository(db)
	qualityRepo := repository.NewKBQualityRepository(db)

	// Initialize services
	knowledgeService := services.NewKnowledgeService(knowledgeRepo, openaiClient)
	chatService := services.NewChatService(conversationRepo, knowledgeService, openaiClient)

	// New advanced KB services
	folderService := services.NewKBFolderService(folderRepo)
	tagService := services.NewKBTagService(tagRepo)
	versionService := services.NewKBVersionService(versionRepo, knowledgeRepo)
	chunkingService := services.NewKBChunkingService(chunkRepo, openaiClient)
	qualityService := services.NewKBQualityService(qualityRepo, knowledgeService, openaiClient)
	syncService := services.NewKBSyncService(syncRepo, knowledgeRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, orgRepo, cfg.JWTSecret)
	chatbotHandler := handlers.NewChatbotHandler(chatbotRepo)
	knowledgeHandler := handlers.NewKnowledgeHandler(knowledgeService)
	chatHandler := handlers.NewChatHandler(chatService)
	widgetHandler := handlers.NewWidgetHandler("../widget")

	// New advanced KB handlers
	folderHandler := handlers.NewKBFolderHandler(folderService)
	tagHandler := handlers.NewKBTagHandler(tagService)
	versionHandler := handlers.NewKBVersionHandler(versionService)
	chunkHandler := handlers.NewKBChunkHandler(chunkingService, knowledgeRepo)
	qualityHandler := handlers.NewKBQualityHandler(qualityService)
	syncHandler := handlers.NewKBSyncHandler(syncService)

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORSMiddleware(cfg.AllowedOrigins))

	// Widget route (public)
	router.GET("/widget.js", widgetHandler.ServeWidget)

	// Public routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public chat endpoint (for widget)
		api.POST("/chat/:chatbot_id", chatHandler.HandleChat)
		api.GET("/chatbots/:id/settings", chatbotHandler.GetSettings)
	}

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Chatbot routes
		chatbots := protected.Group("/chatbots")
		{
			chatbots.POST("", chatbotHandler.Create)
			chatbots.GET("", chatbotHandler.GetAll)
			chatbots.GET("/:id", chatbotHandler.GetByID)
			chatbots.PUT("/:id", chatbotHandler.Update)
			chatbots.DELETE("/:id", chatbotHandler.Delete)
			chatbots.PUT("/:id/settings", chatbotHandler.UpdateSettings)
		}

		// Knowledge base routes
		knowledge := protected.Group("/knowledge")
		{
			knowledge.POST("", knowledgeHandler.Add)
			knowledge.POST("/upload", knowledgeHandler.UploadFile)
			knowledge.POST("/scrape", knowledgeHandler.ScrapeURL)
			knowledge.POST("/preview-chunks", chunkHandler.PreviewChunking)
			knowledge.GET("/chatbot/:chatbot_id", knowledgeHandler.GetByChatbot)
			knowledge.DELETE("/:id", knowledgeHandler.Delete)
		}

		// Knowledge base item-specific routes (separate group to avoid conflicts)
		kb := protected.Group("/kb")
		{
			// Tag routes
			kb.POST("/:id/tags", tagHandler.AssignTags)
			kb.GET("/:id/tags", tagHandler.GetKBTags)
			kb.POST("/:id/tags/:tag_id", tagHandler.AddTagToKB)
			kb.DELETE("/:id/tags/:tag_id", tagHandler.RemoveTagFromKB)

			// Version routes
			kb.GET("/:id/versions", versionHandler.GetVersionHistory)
			kb.GET("/:id/versions/:version", versionHandler.GetVersion)
			kb.POST("/:id/versions/:version/restore", versionHandler.RestoreVersion)
			kb.GET("/:id/versions/compare", versionHandler.CompareVersions)

			// Chunk routes
			kb.GET("/:id/chunks", chunkHandler.GetChunks)
			kb.POST("/:id/rechunk", chunkHandler.RechunkKB)
		}

		// Folder routes
		folders := protected.Group("/folders")
		{
			folders.POST("", folderHandler.CreateFolder)
			folders.GET("/:id", folderHandler.GetFolder)
			folders.PUT("/:id", folderHandler.UpdateFolder)
			folders.DELETE("/:id", folderHandler.DeleteFolder)
			folders.GET("/tree/:chatbot_id", folderHandler.GetFolderTree)
			folders.GET("/roots/:chatbot_id", folderHandler.GetRootFolders)
			folders.POST("/:id/move", folderHandler.MoveFolder)
			folders.GET("/:id/path", folderHandler.GetFolderPath)
		}

		// Tag routes (general)
		tags := protected.Group("/tags")
		{
			tags.POST("", tagHandler.CreateTag)
			tags.GET("/:id", tagHandler.GetTag)
			tags.PUT("/:id", tagHandler.UpdateTag)
			tags.DELETE("/:id", tagHandler.DeleteTag)
			tags.GET("/chatbot/:chatbot_id", tagHandler.GetAllTags)
		}

		// Quality testing routes
		quality := protected.Group("/quality")
		{
			quality.POST("/test", qualityHandler.RunTest)
			quality.POST("/batch-test", qualityHandler.RunBatchTests)
			quality.GET("/tests/:chatbot_id", qualityHandler.GetTestHistory)
			quality.GET("/stats/:chatbot_id", qualityHandler.GetTestStats)
			quality.GET("/failed/:chatbot_id", qualityHandler.GetFailedTests)
			quality.GET("/gaps/:chatbot_id", qualityHandler.IdentifyGaps)
			quality.GET("/score/:chatbot_id", qualityHandler.GetQualityScore)
			quality.GET("/suggestions/:chatbot_id", qualityHandler.GetSuggestions)
		}

		// Sync source routes
		syncSources := protected.Group("/sync-sources")
		{
			syncSources.POST("", syncHandler.CreateSyncSource)
			syncSources.GET("/:id", syncHandler.GetSyncSource)
			syncSources.PUT("/:id", syncHandler.UpdateSyncSource)
			syncSources.DELETE("/:id", syncHandler.DeleteSyncSource)
			syncSources.GET("/chatbot/:chatbot_id", syncHandler.GetSyncSources)
			syncSources.POST("/:id/toggle", syncHandler.ToggleSyncSource)
			syncSources.POST("/:id/trigger", syncHandler.TriggerSync)
			syncSources.GET("/:id/status", syncHandler.GetSyncStatus)
		}

		// Analytics routes
		analytics := protected.Group("/analytics")
		{
			analytics.GET("/conversations/:chatbot_id", chatHandler.GetConversations)
			analytics.GET("/messages/:conversation_id", chatHandler.GetMessages)
		}
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
