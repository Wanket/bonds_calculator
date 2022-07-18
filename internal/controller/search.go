package controller

import (
	"bonds_calculator/internal/service"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type SearchController struct {
	searchService service.ISearchService
}

func NewSearchController(searchService service.ISearchService) *SearchController {
	controller := SearchController{
		searchService: searchService,
	}

	return &controller
}

func (controller *SearchController) Configure(router fiber.Router) {
	router.Get("/search", controller.Search)

	log.Info("SearchController: configured")
}

func (controller *SearchController) Search(ctx *fiber.Ctx) error {
	query := ctx.Query("query")
	if query == "" {
		return fiber.NewError(fiber.StatusBadRequest)
	}

	searchResults := controller.searchService.Search(query)

	return ctx.JSON(searchResults)
}
