package main

import (
	"log"
	"todo_backend/database"
	"todo_backend/model"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// db connection
	database.ConnectDB()

	if err := database.DB.AutoMigrate(&model.Todo{}); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/", getAllTodos)
	app.Post("/", createTodo)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	app.Listen(":3000")

}

func getAllTodos(ctx *fiber.Ctx) error {
	var todos []model.Todo

	if err := database.DB.Find(&todos).Error; err != nil {
		return fiber.NewError(500, err.Error())
	}
	return ctx.JSON(todos)
}

func createTodo(ctx *fiber.Ctx) error {
	var todo model.Todo

	// bind req body to model
	if err := ctx.BodyParser(&todo); err != nil {
		return fiber.NewError(400, "Invalid Json : "+err.Error())
	}

	// save to db
	if err := database.DB.Create(&todo).Error; err != nil {
		return fiber.NewError(500, "Db Error "+err.Error())
	}
	// return saved todo as json

	return ctx.Status(fiber.StatusCreated).JSON(todo)

}
