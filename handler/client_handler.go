package handler

import (
	"fullstack2024-test/model"
	"fullstack2024-test/service"
	"fullstack2024-test/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ClientHandler struct {
	clientUseCase *usecase.ClientUseCase
	s3Service     *service.S3Service
}

func NewClientHandler(clientUseCase *usecase.ClientUseCase, s3Service *service.S3Service) *ClientHandler {
	return &ClientHandler{
		clientUseCase: clientUseCase,
		s3Service:     s3Service,
	}
}

func (client *ClientHandler) ClientRoutes(e *echo.Echo) {
	e.POST("/clients", client.CreateClient)
	e.GET("/clients/:id", client.GetClientByID)
	e.PUT("/clients/:id", client.UpdateClientByID)
	e.DELETE("/clients/:id", client.DeleteClientByID)
}

func (client *ClientHandler) CreateClient(c echo.Context) error {
	if c.Request().Header.Get("Content-Type") != "" && len(c.Request().Header.Get("Content-Type")) >= 19 && c.Request().Header.Get("Content-Type")[:19] == "multipart/form-data" {
		if err := c.Request().ParseMultipartForm(10 << 20); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid form data"})
		}
	}

	var clientModel model.Client

	if err := c.Bind(&clientModel); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	createdClient, err := client.clientUseCase.CreateClient(&clientModel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create client"})
	}

	if client.s3Service != nil {
		file, err := c.FormFile("logo")
		if err == nil && file != nil {
			fileURL, err := client.s3Service.UploadFile(file, createdClient.ID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upload logo"})
			}

			createdClient.ClientLogo = fileURL
			updatedClient, err := client.clientUseCase.UpdateClientByID(createdClient.ID, createdClient)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update client with logo"})
			}
			createdClient = updatedClient
		}
	}

	return c.JSON(http.StatusCreated, createdClient)
}

func (client *ClientHandler) GetClientByID(c echo.Context) error {
	id := c.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	clientModel, err := client.clientUseCase.GetClientByID(clientID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
	}

	return c.JSON(http.StatusOK, clientModel)
}

func (client *ClientHandler) UpdateClientByID(c echo.Context) error {
	id := c.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	if c.Request().Header.Get("Content-Type") != "" && len(c.Request().Header.Get("Content-Type")) >= 19 && c.Request().Header.Get("Content-Type")[:19] == "multipart/form-data" {
		if err := c.Request().ParseMultipartForm(10 << 20); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid form data"})
		}
	}

	var clientModel model.Client

	if err := c.Bind(&clientModel); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	existingClient, err := client.clientUseCase.GetClientByID(clientID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
	}

	if client.s3Service != nil {
		file, err := c.FormFile("logo")
		if err == nil && file != nil {
			fileURL, err := client.s3Service.UploadFile(file, clientID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upload logo"})
			}

			clientModel.ClientLogo = fileURL
		} else {
			clientModel.ClientLogo = existingClient.ClientLogo
		}
	}

	updatedClient, err := client.clientUseCase.UpdateClientByID(clientID, &clientModel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update client"})
	}

	return c.JSON(http.StatusOK, updatedClient)
}

func (client *ClientHandler) DeleteClientByID(c echo.Context) error {
	id := c.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	err = client.clientUseCase.DeleteClientByID(clientID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete client"})
	}

	return c.NoContent(http.StatusNoContent)
}
