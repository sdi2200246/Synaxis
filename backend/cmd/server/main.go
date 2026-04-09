package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sdi2200246/synaxis/internal/controllers"
	"github.com/sdi2200246/synaxis/internal/middleware"
	"github.com/sdi2200246/synaxis/internal/repos"
	"github.com/sdi2200246/synaxis/internal/services"
)

func main() {
    // load .env
    godotenv.Load()

    // connect to DB
    pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    categoryRepo:= repos.NewCategoryRepo(pool)
    
    userRepo    := repos.NewUserRepo(pool)
    userService := services.NewUserService(userRepo)
    userHandler := controllers.NewUserHandler(userService)

    authService := services.NewAuthService(userRepo , "jason_derullo") //TDO realsecret .
    authHandler := middleware.NewAuthHandler(authService)

    eventRepo   := repos.NewEventRepo(pool)
    eventsService := services.NewEventService(eventRepo)
    eventsHandler := controllers.NewEventsHandler(eventsService)

    venueRepo := repos.NewVenueRepo(pool)
    venueService:=services.NewVenueService(venueRepo)
    venueHandlers:=controllers.NewVenueHandler(venueService)


    r := gin.Default()

    r.POST("/users", userHandler.Register)
    r.POST("/auth/login" , authHandler.Login)

    r.GET("/categories", func(c *gin.Context) {
        categories, err := categoryRepo.GetAll(c.Request.Context())
        if err != nil {
            c.JSON(500, gin.H{"error": "failed to fetch categories"})
            return
        }
        c.JSON(200, categories)
    })


    auth := r.Group("/")
    auth.Use(authHandler.AuthMiddleware())
    {

        admin := auth.Group("/admin")
        admin.Use(authHandler.AdminOnly())
        {
            admin.GET("/users" , userHandler.GetUsers)
            admin.POST("/users/:id/approve" , userHandler.ApproveUser)
            admin.POST("/users/:id/reject" , userHandler.RejectUser)
        }

        auth.POST("/events", eventsHandler.Create)
        auth.PATCH("/events/:id" , eventsHandler.UpdateEvent)
        auth.GET("/events", eventsHandler.GetOrganizerEvents)
        auth.GET("/venues" , venueHandlers.GetVenues)
    }
    
    // start server
    r.Run(":8080")
}