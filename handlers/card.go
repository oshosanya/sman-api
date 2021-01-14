package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/oshosanya/sman/definitions"
	"github.com/oshosanya/sman/utils"
	"io"
	"net/http"
	"os"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type OnlineCard struct {
	URL string `json:"url"`
}

func CreateCard(c echo.Context) error {
	id := new(definitions.IDCard)

	if err := c.Bind(id); err != nil {
		println(err)
		return err
	}

	if err := id.Validate(); err != nil {
		b, _ := json.Marshal(err)
		c.Response().Header().Set(echo.HeaderContentType, "application/json")
		c.Response().WriteHeader(201)
		return c.String(http.StatusOK, string(b))
	}

	file, err := c.FormFile("passport")
	if err != nil {
		errorResponse := ErrorResponse{
			Message: err.Error(),
		}
		return c.JSON(http.StatusOK, errorResponse)
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("media/" + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	idFile, err := utils.CreateIDCard(dst, *id)
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	url, err := utils.UploadFile(idFile.Name())

	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Status: "success",
		Data: OnlineCard{
			URL: url,
		},
	})
}
