package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	address := ":8081"
	e := echo.New()
	var h handlers

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	// Добавление User:
	e.POST("/api/users", h.addUser)

	// Получение списка User:
	e.GET("/api/users", h.getAllUsers)

	// Получение User по id:
	e.GET("/api/users/:id", h.getUser)

	// Редактирование User по id:
	e.PATCH("/api/users/:id", h.editUser)

	// Удаление user по id:
	e.DELETE("/api/users/:id", h.deleteUser)

	fmt.Println(e.Start(address))
}
