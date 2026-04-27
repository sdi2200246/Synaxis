package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sdi2200246/synaxis/internal/controllers"
	"github.com/sdi2200246/synaxis/internal/infastructure"
	"github.com/sdi2200246/synaxis/internal/middleware"
	"github.com/sdi2200246/synaxis/internal/repos"
	"github.com/sdi2200246/synaxis/internal/services"
)

func main() {

    godotenv.Load()

    pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    eventBus := infastructure.NewEventBus();


    categoryRepo := repos.NewCategoryRepo(pool)
    userRepo     := repos.NewUserRepo(pool)
    eventRepo    := repos.NewEventRepo(pool)
    venueRepo    := repos.NewVenueRepo(pool)
    bookingRepo  := repos.NewBookingsRepo(pool)
    ticketsRepo  := repos.NewTicketTypeRepo(pool)
    messagesRepo := repos.NewMessagesRepo(pool)

    userService  := services.NewUserService(userRepo)
    authService  := services.NewAuthService(userRepo, "jason_derullo")
    venueService := services.NewVenueService(venueRepo)

    eventsService := services.NewEventService(eventRepo, categoryRepo, bookingRepo , ticketsRepo ,eventBus , venueRepo)
    bookingService := services.NewBookingService(ticketsRepo, bookingRepo , eventRepo)
    ticketTypeService := services.NewTicketTypeService(ticketsRepo, eventRepo)
    messagesService := services.NewMessageService(messagesRepo , bookingRepo , eventRepo)
    eventCancelationService := services.NewCancelEventService(eventRepo , bookingRepo , messagesRepo , eventBus)
    eventCancelationService.Subscribe()


    baseHandler        := &controllers.BaseHandler{}
    userHandler        := controllers.NewUserHandler(userService)
    authHandler        := middleware.NewAuthHandler(authService)
    venueHandlers      := controllers.NewVenueHandler(venueService)
    eventsHandler      := controllers.NewEventsHandler(eventsService , baseHandler)
    ticketsHandler     := controllers.NewTicketTypeHandler(ticketTypeService , baseHandler)
    bookingHandler     := controllers.NewBookingHandler(bookingService , baseHandler)
    messagesHandler    := controllers.NewMessagesHandler(messagesService , baseHandler)
    // adminExportHandler := controllers.NewAdminExportHandler(eventsService, bookingService)


    r := gin.Default()

    r.POST("/users", userHandler.Register)
    r.POST("/auth/login", authHandler.Login)

    r.GET("/categories", func(c *gin.Context) {
        categories, err := categoryRepo.GetAll(c.Request.Context())
        if err != nil {
            c.JSON(500, gin.H{"error": "failed to fetch categories"})
            return
        }
        c.JSON(200, categories)
    })
    r.GET("/venues", venueHandlers.GetVenues)
    r.GET("/venues/:id", venueHandlers.GetVenue)
    r.GET("/events", authHandler.OptionalAuth(), eventsHandler.List)
    r.GET("/events/:id", eventsHandler.GetByID)
    r.GET("/events/:id/categories", eventsHandler.GetEventCategories)

    auth := r.Group("/")
    auth.Use(authHandler.AuthMiddleware())
    {
        admin := auth.Group("/admin")
        admin.Use(authHandler.AdminOnly())
        {
            admin.GET("/users", userHandler.GetUsers)
            admin.POST("/users/:id/approve", userHandler.ApproveUser)
            admin.POST("/users/:id/reject", userHandler.RejectUser)
            // admin.GET("/events" , adminExportHandler.Export)
        }
        auth.GET("/users/:id", userHandler.GetByID)

        auth.POST("/events", eventsHandler.Create)
        auth.PATCH("/events/:id", eventsHandler.UpdateEvent)
        auth.DELETE("/events/:id" , eventsHandler.Delete)
       
        auth.POST("/events/:id/tickets", ticketsHandler.Create)
        auth.GET("/events/:id/tickets", ticketsHandler.GetByEventID)
        auth.PATCH("/events/:id/tickets/:ticket_id", ticketsHandler.Update)
        auth.GET("/tickets/:id", ticketsHandler.GetByID)

        
        auth.GET("/events/:id/bookings", bookingHandler.GetEventBookings)
        auth.POST("/events/:id/bookings", bookingHandler.Create)
        auth.GET("/bookings", bookingHandler.GetUserBookings)


        auth.POST("/conversations" , messagesHandler.CreateConversation)
        auth.GET("/conversations" , messagesHandler.ListUserConversations)
        auth.PATCH("/conversations/:id/read", messagesHandler.MarkConversationAsRead)
        auth.POST("/conversations/:id/messages" , messagesHandler.CreateMessage)
        auth.GET("/conversations/:id/messages" , messagesHandler.GetConversationMessages)
        auth.PATCH("/messages/:id", messagesHandler.UpdateMessage)
    }
    
    // start server
    r.Run(":8080")
}