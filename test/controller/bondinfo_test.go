package controller

import (
	"bonds_calculator/internal/controller"
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/service"
	mockservice "bonds_calculator/internal/service/mock"
	"bonds_calculator/test"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"io"
	"net/http/httptest"
	"testing"
)

func TestBondInfoSuccess(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	app, bondInfoService := createAndRegisterNewBondInfoController(mockController)

	bondInfoService.EXPECT().GetBondInfo("test_id").Return(service.BondInfoResult{
		Bond: moex.Bond{
			Id: "test_id",
		},
		Bondization: moex.Bondization{
			Id: "test_id",
			Amortizations: []moex.Amortization{
				{
					Value: 1000.0,
				},
			},
			Coupons: []moex.Coupon{},
		},
		MaturityIncome: datastuct.NewOptional(3.5),
	}, nil)

	req := httptest.NewRequest("GET", "/api/static/bond_info?id=test_id", nil)

	resp, err := app.Test(req)
	assert.NoError(err)

	assert.Equal(fiber.StatusOK, resp.StatusCode)
	assert.Equal("application/json", resp.Header.Get("Content-Type"))

	expected := `{"Bond":{"Id":"test_id","CurrentPricePercent":0,"ShortName":"","Coupon":0,"NextCoupon":"0001-01-01T00:00:00Z","AccCoupon":0,"PrevPricePercent":0,"Value":0,"CouponPeriod":0,"PriceStep":0,"CouponPercent":null,"EndDate":"0001-01-01T00:00:00Z","Currency":""},"Bondization":{"Id":"test_id","Amortizations":[{"Date":"0001-01-01T00:00:00Z","Value":1000}],"Coupons":[]},"MaturityIncome":3.5,"CurrentIncome":null}`
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

	bondInfoService.EXPECT().GetBondInfo("test_id").Return(service.BondInfoResult{}, errors.New("test error"))

	req := httptest.NewRequest("GET", "/api/static/bond_info?id=test_id", nil)

	resp, err := app.Test(req)
	assert.NoError(err)
	assert.Equal(fiber.StatusNotFound, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Equal("test error", string(body))
}

func createAndRegisterNewBondInfoController(mockController *gomock.Controller) (*fiber.App, *mockservice.MockIBondInfoService) {
	bondInfoService := mockservice.NewMockIBondInfoService(mockController)

	ctr := controller.NewBondInfoController(bondInfoService)

	app := createAppAndRegistryController("api/static", ctr)

	return app, bondInfoService
}
