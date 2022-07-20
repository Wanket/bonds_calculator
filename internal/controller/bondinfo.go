package controller

import (
	"bonds_calculator/internal/service"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type BondInfoController struct {
	bondInfoService service.IBondInfoService
}

func NewBondInfoController(bondInfoService service.IBondInfoService) *BondInfoController {
	controller := BondInfoController{
		bondInfoService: bondInfoService,
	}

	return &controller
}

func (controller *BondInfoController) Configure(router fiber.Router) {
	router.Get("/bond_info", controller.GetBondInfo)

	log.Info("BondInfoController: configured")
}

func (controller *BondInfoController) GetBondInfo(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest)
	}

	bondInfo, err := controller.bondInfoService.GetBondInfo(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	return ctx.JSON(bondInfo)
}
