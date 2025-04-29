package handlers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gocolly/colly"
)

// Global variable to store the scraped version (optional for caching)
var cachedVersion string = ""

// ScrapeVersion fetches the latest version from the website
func ScrapeVersion() (string, error) {
	if cachedVersion != "" {
		return cachedVersion, nil
	}

	var arrVersao []string
	var classbutton string
	classbutton = ButtonDownloadClass()
	// Criando um coletor normal
	c := colly.NewCollector(
		colly.AllowedDomains("sisaps.saude.gov.br"),
	)

	// Modificando diretamente o Transport do coletor
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Ignora SSL
	})

	//downloadButton_GNxa

	c.OnHTML(classbutton, func(e *colly.HTMLElement) {
		versaoElement := e.Text
		arr := strings.Split(versaoElement, " ")

		// Check if splitting resulted in at least 2 elements
		if len(arr) < 2 {
			fmt.Println("Warning: Version element format unexpected")
			return
		}

		arrVersao = append(arrVersao, arr[2])
	})

	err := c.Visit("https://sisaps.saude.gov.br/sistemas/esusaps/")
	if err != nil {
		return "", fmt.Errorf("Error visiting website: %w", err)
	}

	if len(arrVersao) == 0 {
		return "", fmt.Errorf("No version found on the page")
	}

	cachedVersion = arrVersao[0]
	return cachedVersion, nil
}

func testeClass() (string, error) {

	var arrClass []string
	var classAttr string
	// Criando um coletor normal
	c := colly.NewCollector(
		colly.AllowedDomains("sisaps.saude.gov.br"),
	)

	// Modificando diretamente o Transport do coletor
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Ignora SSL
	})

	c.OnHTML("button", func(e *colly.HTMLElement) {
		classAttr = e.Attr("class")
		arrClass = append(arrClass, classAttr)
	})

	err := c.Visit("https://sisaps.saude.gov.br/sistemas/esusaps/")
	if err != nil {
		return "", fmt.Errorf("Error visiting website: %w", err)
	}

	classAttr = "." + arrClass[3]
	return classAttr, nil
}

func ButtonDownloadClass () string {
	classbutton, err := testeClass()
	if err != nil {
		return ""
	}
	return classbutton
}

func Teste(c *fiber.Ctx) error {
	classbutton, err := testeClass()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error fetching version: %v", err))
	}
	return c.SendString(classbutton)
} 

// Home handler displays available endpoints
func Home(c *fiber.Ctx) error {
	return c.SendString("Available endpoints: /versaoesus, /verificaversao?v=X.X.X")
}

// VersaoEsus handler fetches and returns the current version
func VersaoEsus(c *fiber.Ctx) error {
	version, err := ScrapeVersion()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error fetching version: %v", err))
	}
	return c.SendString(version)
}

// VerificaVersao handler compares provided version with the scraped version
func VerificaVersao(c *fiber.Ctx) error {
	requestedVersion := c.Query("v")
	if requestedVersion == "" {
		return c.Status(fiber.StatusBadRequest).SendString("-1")
	}

	currentVersion, err := ScrapeVersion()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error fetching version: %v", err))
	}

	var comparisonResult string
	switch {
	case requestedVersion < currentVersion:
		comparisonResult = "0" // Update required
	case requestedVersion == currentVersion:
		comparisonResult = "1" // Versions match
	default:
		comparisonResult = "0" // Downgrade not recommended
	}

	return c.SendString(comparisonResult)
}