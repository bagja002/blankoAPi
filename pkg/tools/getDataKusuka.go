package tools

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type RequesToken struct {
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
	// Inisialisasi aplikasi Fiber

	baseUrl := "https://statistik.kkp.go.id/api-statistik/index.php/"

	// Membuat form data
	args := fasthttp.AcquireArgs()
	args.Set("username", "bppsdmkp")
	args.Set("password", "EF700E623E10945FA9B55EBE42139D96")

	// Membuat request POST
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(baseUrl + "Token/create")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetBodyString(args.String())
	req.Header.SetContentType("application/x-www-form-urlencoded")

	// Membuat response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Mengirim request
	if err := fasthttp.Do(req, resp); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Pesan": "gagal mengirim request",
		})
	}

	// Membaca response dari server
	responseBody := resp.Body()

	// Mengurai JSON response untuk mendapatkan token
	var apiResp ApiResponse
	if err := json.Unmarshal(responseBody, &apiResp); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Pesan": "gagal mengurai response JSON",
		})
	}

	// Data token apiResp.Data.Token

	//Hit Gerate Tokennay

	// URL API eksternal
	nomor := c.Query("nomor_kusuka")
	apiURL := baseUrl + "Kusuka?nomor_kusuka=" + nomor

	// Membuat objek request FastHTTP

	// Set method dan URL request
	req.SetRequestURI(apiURL)

	// Menambahkan header "Token"
	req.Header.Set("Token", apiResp.Data.Token)

	// Melakukan permintaan HTTP GET ke API menggunakan FastHTTP
	if err := fasthttp.Do(req, resp); err != nil {
		log.Println("Error making request to API:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Memeriksa kode status respons
	if resp.StatusCode() != fasthttp.StatusOK {
		log.Println("Non-OK status code from API:", resp.StatusCode())
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Mengirimkan data respons ke klien sebagai JSON
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	c.Send(resp.Body())

	return nil
}
