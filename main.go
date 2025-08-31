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
	app.Delete("/:id", deleteTodo)
	app.Put("/:id", updateTodoStatus)

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

func deleteTodo(ctx *fiber.Ctx) error {
	// get param from url

	todoId, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, " Invalid todo id "+err.Error())
	}

	var todo model.Todo
	if err := database.DB.First(&todo, todoId).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, " There is no todo with this ID")
	}

	// delete record

	if err := database.DB.Delete(&todo).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, " Delete error! ")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo is deleted successfully",
	})

}

func updateTodoStatus(ctx *fiber.Ctx) error {
	todoId, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, " Invalid todo id "+err.Error())
	}

	type StatusUpdate struct {
		Status bool
	}
	var req StatusUpdate
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid payload !",
		})
	}

	result := database.DB.Model(&model.Todo{}).Where("id = ?", todoId).Update("Status", req.Status)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Todo Record is not found !",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Todo status is updated successfully",
	})
}
