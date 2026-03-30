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
    
    userRepo    := repos.NewUserRepo(pool)
    userService := services.NewUserService(userRepo)
    userHandler := controllers.NewUserHandler(userService)

    authService := services.NewAuthService(userRepo , "jason_derullo") //TDO realsecret .
    authHnadler := middleware.NewAuthHandler(authService)

    r := gin.Default()

    r.POST("/users", userHandler.Register)
    r.GET("/auth/login" , authHnadler.Login)

    // start server
    r.Run(":8080")
}