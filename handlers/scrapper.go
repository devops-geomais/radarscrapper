package handlers

import (
	"github.com/gofiber/fiber/v2"
	"strings"
     "fmt"
	"github.com/gocolly/colly"
)

func Home(c *fiber.Ctx) error {
		return c.SendString("/versaoesus ou /verificaversao?v=5.x.xx")
}

func VersaoEsus(c *fiber.Ctx) error {
	return c.SendString(scrapVersao())
}

func VerificaVersao(c *fiber.Ctx) error {
	var versao string = c.Query("v")
	fmt.Println(versao)
	return c.SendString(comparaVersao(versao))
}

func scrapVersao() string {
	var arrVersao []string
	// instantiate a new collector object
	c := colly.NewCollector(
		colly.AllowedDomains("sisaps.saude.gov.br"),
	)
	c.OnHTML(".pricing-table li h3", func(e *colly.HTMLElement) {
		var versaoElement = e.Text
		arr := strings.Split(versaoElement, " ")
		arrVersao = append(arrVersao, arr[1])
	})
	// open the target URL
	c.Visit("https://sisaps.saude.gov.br/esus/")
	return arrVersao[0]
}

func comparaVersao(versao string) string {
	var versaoAtual string
	versaoAtual = scrapVersao()
	var retorno = "0"
	if(versaoAtual == versao){
		retorno = "1"
	}
	return retorno
}