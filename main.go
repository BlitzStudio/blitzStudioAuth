package main

// import "github.com/gofiber/fiber/v3"

// func main() {
// 	app := fiber.New()

// 	app.Post("/auth/signup", func(c *fiber.Ctx) error {})
// 	app.Post("/auth/signin", func(c *fiber.Ctx) error {})
// 	app.Post("/auth/refresh-token", func(c *fiber.Ctx) error {})

// 	app.Get("*", func(c fiber.Ctx) error {
// 		return c.SendStatus(fiber.StatusNotFound)
// 	})

// 	app.Listen(":3000")
// }

import (
	"database/sql"

	"github.com/BlitzStudio/blitzStudioAuth/types"
	"github.com/BlitzStudio/blitzStudioAuth/utils"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v3"
)

func main() {
	log := utils.GetLogger()

	db, err := sql.Open("mysql", "root:pass@/test")
	// ctx := context.Background()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// i dont know what those 2 options do
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	createdUser, err := utils.CreateUser(types.User{
		Name:     "test",
		Password: "here",
		Email:    "test@il.com",
	}, db)

	if err != nil {
		log.Debug(err)
	}
	spew.Dump(createdUser)
	app := fiber.New()

	// app.Post("/auth/signup", func(c *fiber.Ctx) error {})
	// app.Post("/auth/signin", func(c *fiber.Ctx) error {})
	// app.Post("/auth/refresh-token", func(c *fiber.Ctx) error {})

	app.Get("*", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	app.Listen(":3000")
}
