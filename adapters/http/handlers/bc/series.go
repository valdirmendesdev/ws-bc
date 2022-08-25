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

type serie struct {
	Date  string `json:"data"`
	Value string `json:"valor"`
}

type HTTPError struct {
	Code    int    `json:"code" example:"999"`
	Message string `json:"message" example:"Mensagem de resposta"`
}

// @Summary     Retorna uma lista de valores de uma série histórica do banco central em um período específico
// @Description Executa a consulta dos valores da série informada no parâmetro no período informado como parâmetro
// @Tags        Séries
// @Param       series_number path  int    true  "Código da série do banco central. Pode encontrada aqui: https://www3.bcb.gov.br/sgspub/localizarseries/localizarSeries.do?method=prepararTelaLocalizarSeries"
// @Param       from          query string false "Data inicial no formato dd/MM/yyyy. Se não informado, será considerado o primeiro dia do mês atual"
// @Param       to            query string false "Data final no formato dd/MM/yyyy. Se não informado, será considerada a data atual"
// @Produce     json
// @Success     200 {array}  serie
// @Failure     400 {object} HTTPError
// @Failure     503 {object} HTTPError
// @Router      /series/{series_number} [get]
func Series() fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		if c.Params("series_number") == "" {
			return c.Status(http.StatusBadRequest).JSON(HTTPError{Code: http.StatusBadRequest, Message: "Código da série é obrigatório!"})
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
			return c.Status(http.StatusServiceUnavailable).JSON(HTTPError{Code: http.StatusServiceUnavailable, Message: "Serviço do Banco Central service está indisponível no momento. Tente novamente mais tarde!"})
		}

		return c.SendStream(response.Body)
	})
}

// @Summary     Retorna o último ou os últimos valores de uma série histórica do banco central
// @Description Executa a consulta dos últimos valores da série informada no parâmetro
// @Tags        Séries
// @Param       series_number path int true "Código da série do banco central. Pode encontrada aqui: https://www3.bcb.gov.br/sgspub/localizarseries/localizarSeries.do?method=prepararTelaLocalizarSeries"
// @Param       quantity      path int true "Quantidade desejada de últimos valores da série informada. Se informado um valor inválido, será considerado valor 1, ou seja, apenas o último valor."
// @Produce     json
// @Success     200 {array}  serie
// @Failure     400 {object} HTTPError
// @Failure     503 {object} HTTPError
// @Router      /series/{series_number}/latest/{quantity} [get]
func SeriesUltimos() fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		if c.Params("series_number") == "" {
			return c.Status(http.StatusBadRequest).JSON(HTTPError{Code: http.StatusBadRequest, Message: "Código da série é obrigatório!"})
		}

		quantity, err := strconv.Atoi(c.Params("quantity"))
		if err != nil || quantity == 0 {
			quantity = 1
		}

		finalURL := fmt.Sprintf("%s%s/dados/ultimos/%v?formato=json", baseURL, c.Params("series_number"), quantity)

		response, err := http.Get(finalURL)
		if err != nil {
			return c.Status(http.StatusServiceUnavailable).JSON(HTTPError{Code: http.StatusServiceUnavailable, Message: "Serviço do Banco Central service está indisponível no momento. Tente novamente mais tarde!"})
		}

		return c.SendStream(response.Body)
	})
}
