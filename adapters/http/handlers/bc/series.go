package bc

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

const baseURL = "http://api.bcb.gov.br/dados/serie/bcdata.sgs."

type seriesQuery struct {
	From string `query:"from"`
	To   string `query:"to"`
}

func Series() fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		if c.Params("series_number") == "" {
			c.Status(http.StatusBadRequest).SendString("Invalid parameter!")
		}

		q := new(seriesQuery)

		if err := c.QueryParser(q); err != nil {
			return err
		}

		now := time.Now()

		if q.From == "" {
			q.From = fmt.Sprintf("%v/%v/%v", "01", int(now.Month()), now.Year())
		}

		if q.To == "" {
			q.To = fmt.Sprintf("%v/%v/%v", now.Day(), int(now.Month()), now.Year())
		}

		finalURL := fmt.Sprintf("%s%s/dados?formato=json&dataInicial=%s&dataFinal=%s", baseURL, c.Params("series_number"), q.From, q.To)

		response, err := http.Get(finalURL)
		if err != nil {
			return c.Status(http.StatusServiceUnavailable).SendString("Banco Central service is unavailable. Try again later!")
		}

		return c.SendStream(response.Body)
	})
}

func SeriesUltimos() fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		if c.Params("series_number") == "" {
			c.Status(http.StatusBadRequest).SendString("Invalid parameters!")
		}
		quantity := c.Params("quantity")
		_, err := strconv.Atoi(quantity)
		if err != nil {
			quantity = "1"
		}

		finalURL := fmt.Sprintf("%s%s/dados/ultimos/%s?formato=json", baseURL, c.Params("series_number"), quantity)

		response, err := http.Get(finalURL)
		if err != nil {
			return c.Status(http.StatusServiceUnavailable).SendString("Banco Central service is unavailable. Try again later!")
		}

		return c.SendStream(response.Body)
	})
}
