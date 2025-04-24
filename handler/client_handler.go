package handler

import (
	"fullstack2024-test/model"
	"fullstack2024-test/usecase"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ClientHandler struct {
	clientUseCase *usecase.ClientUseCase
}

func NewClientHandler(clientUseCase *usecase.ClientUseCase) *ClientHandler {
	return &ClientHandler{clientUseCase: clientUseCase}
}

func (client *ClientHandler) ClientRoutes(e *echo.Echo) {
	e.POST("/clients", client.CreateClient)
	e.GET("/clients/:id", client.GetClientByID)
	e.PUT("/clients/:id", client.UpdateClientByID)
	e.DELETE("/clients/:id", client.DeleteClientByID)
}

func (client *ClientHandler) CreateClient(c echo.Context) error {
	var clientModel model.Client

	if err := c.Bind(&clientModel); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid input"})
	}

	createdClient, err := client.clientUseCase.CreateClient(&clientModel)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create client"})
	}

	return c.JSON(201, createdClient)
}
func (client *ClientHandler) GetClientByID(c echo.Context) error {
	id := c.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid ID format"})
	}

	clientModel, err := client.clientUseCase.GetClientByID(clientID)
	if err != nil {
		return c.JSON(404, map[string]string{"error": "Client not found"})
	}

	return c.JSON(200, clientModel)
}

func (client *ClientHandler) UpdateClientByID(c echo.Context) error {
	id := c.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid ID format"})
	}

	var clientModel model.Client

	if err := c.Bind(&clientModel); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid input"})
	}

	updatedClient, err := client.clientUseCase.UpdateClientByID(clientID, &clientModel)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to update client"})
	}

	return c.JSON(200, updatedClient)
}

func (client *ClientHandler) DeleteClientByID(c echo.Context) error {
	id := c.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid ID format"})
	}

	err = client.clientUseCase.DeleteClientByID(clientID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to delete client"})
	}

	return c.NoContent(204)
}
