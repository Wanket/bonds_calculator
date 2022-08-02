package controller_test

import (
	"bonds_calculator/internal/controller"
	controllerutil "bonds_calculator/internal/controller/util"
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/service"
	mockservice "bonds_calculator/internal/service/mock"
	"bonds_calculator/test"
	testcontroller "bonds_calculator/test/controller"
	testservice "bonds_calculator/test/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/mailru/easyjson"
	"io"
	"net/http/httptest"
	"testing"
)

func TestBondInfoSuccess(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	app, bondInfoService := createAndRegisterNewBondInfoController(mockController)

	bondInfoService.EXPECT().GetBondInfo("test_id").Return(service.BondInfoResult{
		Bond: moex.Bond{
			ID: "test_id",
		},
		Bondization: moex.Bondization{
			ID: "test_id",
			Amortizations: []moex.Amortization{
				{
					Value: 1000.0,
				},
			},
			Coupons: []moex.Coupon{},
		},
		MaturityIncome: datastruct.NewOptional(3.5),
		CurrentIncome:  datastruct.Optional[float64]{},
	}, nil)

	req := httptest.NewRequest("GET", "/api/static/bond_info?id=test_id", nil)

	resp, err := app.Test(req)
	assert.NoError(err)

	assert.Equal(fiber.StatusOK, resp.StatusCode)
	assert.Equal("application/json", resp.Header.Get("Content-Type"))

	expected, err := testservice.LoadExpectedJSON[controller.BondInfoResult]("test/data/marshaling/bond_info_success.json")
	assert.NoError(err)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Equal(expected, string(body))
}

func TestBondInfoBadRequest(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	app, _ := createAndRegisterNewBondInfoController(mockController)

	req := httptest.NewRequest("GET", "/api/static/bond_info", nil)

	resp, err := app.Test(req)
	assert.NoError(err)
	assert.Equal(fiber.StatusBadRequest, resp.StatusCode)
}

func TestBondInfoErrors(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	app, bondInfoService := createAndRegisterNewBondInfoController(mockController)

	bondInfoService.EXPECT().GetBondInfo("test_id").Return(service.BondInfoResult{}, test.ErrTest) //nolint:exhaustruct

	expected, err := easyjson.Marshal(controller.BondInfoResult{
		BaseResult: controllerutil.BaseResult{
			Status: controllerutil.StatusBondNotFound,
		},
		Result: datastruct.Optional[service.BondInfoResult]{},
	})
	assert.NoError(err)

	req := httptest.NewRequest("GET", "/api/static/bond_info?id=test_id", nil)

	resp, err := app.Test(req)
	assert.NoError(err)
	assert.Equal(fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Equal(string(expected), string(body))
}

func createAndRegisterNewBondInfoController(
	mockController *gomock.Controller,
) (*fiber.App, *mockservice.MockIBondInfoService) {
	bondInfoService := mockservice.NewMockIBondInfoService(mockController)

	ctr := controller.NewBondInfoController(bondInfoService)

	app := testcontroller.CreateAppAndRegistryController("api/static", ctr)

	return app, bondInfoService
}
