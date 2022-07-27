package controller_test

import (
	"bonds_calculator/internal/controller"
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/service"
	mockservice "bonds_calculator/internal/service/mock"
	"bonds_calculator/test"
	testcontroller "bonds_calculator/test/controller"
	testservice "bonds_calculator/test/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"io"
	"net/http/httptest"
	"testing"
)

func TestSearchSuccess(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	app, searchService := createAndRegisterNewSearchController(mockController)

	searchService.EXPECT().Search("test_q").Return(service.SearchResults{
		service.SearchResult{
			Bond: moex.Bond{
				ID: "test",
			},
			MaturityIncome: datastruct.NewOptional(3.5),
			CurrentIncome:  datastruct.Optional[float64]{},
		},
	})

	req := httptest.NewRequest("GET", "/api/static/search?query=test_q", nil)

	resp, err := app.Test(req)
	assert.NoError(err)

	assert.Equal(fiber.StatusOK, resp.StatusCode)
	assert.Equal("application/json", resp.Header.Get("Content-Type"))

	expected, err := testservice.LoadExpectedJSON[service.SearchResults]("test/data/marshaling/search_success.json")
	assert.NoError(err)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Equal(expected, string(body))
}

func TestSearchBadRequest(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	app, _ := createAndRegisterNewSearchController(mockController)

	urls := []string{
		"/api/static/search",
		"/api/static/search?query=",
		"/api/static/search?query=te",
	}

	for _, url := range urls {
		req := httptest.NewRequest("GET", url, nil)

		resp, err := app.Test(req)

		assert.NoError(err)
		assert.Equal(fiber.StatusBadRequest, resp.StatusCode)
	}
}

func createAndRegisterNewSearchController(
	mockController *gomock.Controller,
) (*fiber.App, *mockservice.MockISearchService) {
	searchService := mockservice.NewMockISearchService(mockController)

	ctr := controller.NewSearchController(searchService)

	app := testcontroller.CreateAppAndRegistryController("api/static", ctr)

	return app, searchService
}
