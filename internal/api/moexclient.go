package api

import (
	"bonds_calculator/internal/model/moex"
	"fmt"
	"github.com/valyala/fasthttp"
	"strings"
	"sync"
	"time"
)

//nolint:lll
const (
	allBondsURL = "https://iss.moex.com/iss/engines/stock/markets/bonds/securities.csv?iss.meta=off&iss.only=marketdata,securities&securities.columns=SECID,SHORTNAME,COUPONVALUE,NEXTCOUPON,ACCRUEDINT,PREVPRICE,FACEVALUE,COUPONPERIOD,MINSTEP,COUPONPERCENT,MATDATE,FACEUNIT&marketdata.columns=SECID,LCURRENTPRICE"

	allBondizationURL = "https://iss.moex.com/iss/securities/${bond}/bondization.csv?limit=unlimited&iss.meta=off&iss.only=amortizations,coupons&amortizations.columns=amortdate,value&coupons.columns=coupondate,value"
)

type GetBondizationsResult struct {
	Bondization moex.Bondization
	Bond        moex.Bond

	Error error
}

//go:generate go run github.com/golang/mock/mockgen -destination=mock/moexclient_gen.go . IMoexClient
type IMoexClient interface {
	Close() error

	GetBonds() ([]moex.Bond, error)
	GetBondization(ID string) (moex.Bondization, error)

	GetBondizationsAsync([]moex.Bond) <-chan GetBondizationsResult
}

type MoexClient struct {
	innerClient fasthttp.Client
}

func NewMoexClient(queueSize int) *MoexClient {
	return &MoexClient{
		innerClient: fasthttp.Client{
			MaxConnsPerHost:    queueSize,
			MaxConnWaitTimeout: time.Second,
		},
	}
}

func (client *MoexClient) Close() error {
	client.innerClient.CloseIdleConnections()

	return nil
}

func (client *MoexClient) GetBonds() ([]moex.Bond, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(allBondsURL)

	if err := client.innerClient.Do(req, resp); err != nil {
		return nil, fmt.Errorf("failed to get bonds: %w", err)
	}

	bonds, err := moex.ParseBondsCp1251(resp.Body())
	if err != nil {
		return nil, fmt.Errorf("failed to parse bonds: %w", err)
	}

	return bonds, nil
}

func (client *MoexClient) GetBondization(bondID string) (moex.Bondization, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(strings.Replace(allBondizationURL, "${bond}", bondID, 1))

	if err := client.innerClient.Do(req, resp); err != nil {
		return moex.Bondization{}, fmt.Errorf("failed to get bondization: %w", err)
	}

	bonds, err := moex.ParseBondization(bondID, resp.Body())
	if err != nil {
		return moex.Bondization{}, fmt.Errorf("failed to parse bondization: %w", err)
	}

	return bonds, nil
}

func (client *MoexClient) GetBondizationsAsync(bonds []moex.Bond) <-chan GetBondizationsResult {
	workerCount := client.innerClient.MaxConnsPerHost

	resultChan := make(chan GetBondizationsResult, workerCount)

	go func() {
		waitGroup := sync.WaitGroup{}
		waitGroup.Add(workerCount - 1)

		for workerID := 0; workerID < workerCount-1; workerID++ {
			go func(workerID int) {
				defer waitGroup.Done()

				client.processBondizationWorker(workerID, bonds, workerCount, resultChan)
			}(workerID)
		}

		client.processBondizationWorker(workerCount-1, bonds, workerCount, resultChan)

		waitGroup.Wait()

		close(resultChan)
	}()

	return resultChan
}

func (client *MoexClient) processBondizationWorker(
	workerID int,
	bonds []moex.Bond,
	workerCount int,
	resultChan chan GetBondizationsResult,
) {
	for i := workerID; i < len(bonds); i += workerCount {
		bondization, err := client.GetBondization(bonds[i].ID)

		resultChan <- GetBondizationsResult{
			Bondization: bondization,
			Bond:        bonds[i],

			Error: err,
		}
	}
}
