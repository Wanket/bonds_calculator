//go:generate go run github.com/mailru/easyjson/easyjson -no_std_marshalers -lower_camel_case $GOFILE
package controller

import (
	controllerutil "bonds_calculator/internal/controller/util"
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/service"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type BondInfoController struct {
	bondInfoService service.IBondInfoService
}

//easyjson:json
type BondInfoResult struct {
	controllerutil.BaseResult
	Result datastruct.Optional[service.BondInfoResult]
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
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	bondInfo, err := controller.bondInfoService.GetBondInfo(id)
	if err != nil {
		return ctx.JSON(BondInfoResult{
			BaseResult: controllerutil.BaseResult{
				Status: controllerutil.ResultBondNotFound,
			},
			Result: datastruct.Optional[service.BondInfoResult]{},
		})
	}

	return ctx.JSON(BondInfoResult{
		BaseResult: controllerutil.BaseResult{
			Status: controllerutil.ResultOK,
		},
		Result: datastruct.NewOptional(bondInfo),
	})
}
