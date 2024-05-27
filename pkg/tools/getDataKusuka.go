package tools

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func GetDataKusuka(c *fiber.Ctx) error {
	// Inisialisasi aplikasi Fiber
	// URL API eksternal
	nomor := c.Query("nomor_kusuka")
	apiURL := "https://statistik.kkp.go.id/api-statistik/index.php/Kusuka?nomor_kusuka=" + nomor

	// Membuat objek request FastHTTP
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	// Set method dan URL request
	req.SetRequestURI(apiURL)

	// Menambahkan header "Token"
	req.Header.Set("Token", "SjAmrqx5mtSf2RmFGIWUrDxQNyx2odjg")

	// Membuat objek response FastHTTP
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

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
