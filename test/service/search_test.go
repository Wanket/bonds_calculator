package service

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/service"
	mockservice "bonds_calculator/internal/service/mock"
	"errors"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	asserts "github.com/stretchr/testify/assert"
	"io/ioutil"
	"runtime"
	"testing"
	"time"
)

func TestReloadSearcher(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	assert := asserts.New(t)

	mockController := gomock.NewController(t)

	staticCalculator := mockservice.NewMockIStaticCalculatorService(mockController)
	staticStore := mockservice.NewMockIStaticStoreService(mockController)

	updatedTime := time.Now()
	staticStore.EXPECT().GetBondsWithUpdateTime().Return([]moex.Bond{}, updatedTime)
	staticStore.EXPECT().GetBondsChangedTime().Return(updatedTime)

	search := service.NewSearchService(staticCalculator, staticStore)

	runtime.Gosched()

	assert.Empty(search.Search("1"))

	useBonds := []moex.Bond{
		{
			Id: "1",
		},
	}

	updatedTime = updatedTime.Add(time.Second)
	staticStore.EXPECT().GetBondsChangedTime().Return(updatedTime)
	staticStore.EXPECT().GetBondsWithUpdateTime().Return(useBonds, updatedTime)

	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBonds[0], calculator.Maturity).Return(10.0, nil)
	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBonds[0], calculator.Current).Return(20.0, nil)

	result := search.Search("1")
	for i := 0; len(result) == 0; i++ {
		staticStore.EXPECT().GetBondsChangedTime().Return(updatedTime)

		time.Sleep(time.Millisecond)
		runtime.Gosched()

		result = search.Search("1")
	}

	assert.Equal(1, len(result))
	assert.Equal("1", result[0].Bond.Id)
	assert.Equal(datastuct.NewOptional(10.0), result[0].MaturityIncome)
	assert.Equal(datastuct.NewOptional(20.0), result[0].CurrentIncome)
}

func TestSearchErrors(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	assert := asserts.New(t)

	mockController := gomock.NewController(t)

	staticCalculator := mockservice.NewMockIStaticCalculatorService(mockController)
	staticStore := mockservice.NewMockIStaticStoreService(mockController)

	updatedTime := time.Now()

	useBonds := []moex.Bond{
		{
			Id: "1",
		},
	}

	staticStore.EXPECT().GetBondsWithUpdateTime().Return(useBonds, updatedTime)
	staticStore.EXPECT().GetBondsChangedTime().Return(updatedTime)

	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBonds[0], calculator.Maturity).Return(0.0, errors.New("mat"))
	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBonds[0], calculator.Current).Return(0.0, errors.New("cur"))

	search := service.NewSearchService(staticCalculator, staticStore)

	runtime.Gosched()

	result := search.Search("1")

	assert.Equal(1, len(result))
	assert.Equal("1", result[0].Bond.Id)
	assert.Equal(datastuct.Optional[float64]{}, result[0].MaturityIncome)
	assert.Equal(datastuct.Optional[float64]{}, result[0].CurrentIncome)
}
