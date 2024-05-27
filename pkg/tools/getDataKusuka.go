package tools

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type RequestToken struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

func GetDataKusuka(c *fiber.Ctx) error {
	baseUrl := "https://statistik.kkp.go.id/api-statistik/index.php/"

	// Prepare form data
	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args) // Make sure to release args after use
	args.Set("username", "bppsdmkp")
	args.Set("password", "EF700E623E10945FA9B55EBE42139D96")

	// Create POST request for token
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(baseUrl + "Token/create")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetBody(args.QueryString())
	req.Header.SetContentType("application/x-www-form-urlencoded")

	// Create response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Send request
	if err := fasthttp.Do(req, resp); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Pesan": "gagal mengirim request",
		})
	}

	// Read response from server
	responseBody := resp.Body()

	// Parse JSON response to get the token
	var apiResp ApiResponse
	if err := json.Unmarshal(responseBody, &apiResp); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Pesan": "gagal mengurai response JSON",
		})
	}
	fmt.Println(apiResp)
	// Check if token is retrieved successfully
	if apiResp.Status != "success" {
		return c.Status(500).JSON(fiber.Map{
			"Pesan": "gagal mendapatkan token",
		})
	}

	// Use the token to get data from Kusuka API
	nomor := c.Query("nomor_kusuka")
	apiURL := baseUrl + "Kusuka?nomor_kusuka=" + nomor

	// Create GET request for Kusuka data
	req.SetRequestURI(apiURL)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("Token", apiResp.Data.Token)

	// Send request
	if err := fasthttp.Do(req, resp); err != nil {
		log.Println("Error making request to API:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Check response status code
	if resp.StatusCode() != fasthttp.StatusOK {
		log.Println("Non-OK status code from API:", resp.StatusCode())
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Send response to client as JSON
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Send(resp.Body())
}
