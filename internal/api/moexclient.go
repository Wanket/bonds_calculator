package api

import (
	"bonds_calculator/internal/model/moex"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	allBondsUrl = "https://iss.moex.com/iss/engines/stock/markets/bonds/securities.csv?iss.meta=off&iss.only=marketdata,securities&securities.columns=SECID,SHORTNAME,COUPONVALUE,NEXTCOUPON,ACCRUEDINT,PREVPRICE,FACEVALUE,COUPONPERIOD,MINSTEP,COUPONPERCENT,MATDATE&marketdata.columns=SECID,LCURRENTPRICE"

	allBondizationUrl = "https://iss.moex.com/iss/securities/${bond}/bondization.csv?limit=unlimited&iss.meta=off&iss.only=amortizations,coupons&amortizations.columns=amortdate,value&coupons.columns=coupondate,value"
)

type MoexClient struct {
	innerClient http.Client

	workQueue chan worker

	context context.Context
	cancel  context.CancelFunc
}

func NewMoexClient(queueSize int) MoexClient {
	return NewMoexClientWithContext(queueSize, context.Background())
}

func NewMoexClientWithContext(queueSize int, ctx context.Context) MoexClient {
	ctx, cancel := context.WithCancel(ctx)
	client := MoexClient{
		workQueue: make(chan worker, queueSize),
		context:   ctx,
		cancel:    cancel,
	}

	for i := 0; i < queueSize; i++ {
		go queueListener(client)
	}

	return client
}

func (client *MoexClient) Close() error {
	client.cancel()

	return nil
}

func (client *MoexClient) GetBonds() ([]moex.Bond, error) {
	request := bondsRequest{newCommonRequest[[]moex.Bond]()}

	return request.commonReturn(client, &request)
}

func (client *MoexClient) GetBondization(Id string) (moex.Bondization, error) {
	request := bondizationRequest{
		bondId:        Id,
		commonRequest: newCommonRequest[moex.Bondization](),
	}

	return request.commonReturn(client, &request)
}

type worker interface {
	work(client MoexClient)
}

func queueListener(client MoexClient) {
	for {
		select {
		case worker := <-client.workQueue:
			worker.work(client)
		case <-client.context.Done():
			return
		}
	}
}

type commonRequest[T any] struct {
	result   T
	err      error
	doneChan chan struct{}
}

func newCommonRequest[T any]() commonRequest[T] {
	return commonRequest[T]{
		doneChan: make(chan struct{}),
	}
}

func (request *commonRequest[R]) commonReturn(client *MoexClient, worker worker) (R, error) {
	client.workQueue <- worker

	<-request.doneChan

	return request.result, request.err
}

type bondsRequest struct {
	commonRequest[[]moex.Bond]
}

func (bondsRequest *bondsRequest) work(client MoexClient) {
	defer close(bondsRequest.doneChan)

	resp, err := client.innerClient.Get(allBondsUrl)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		bondsRequest.err = err

		return
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		bondsRequest.err = fmt.Errorf("cannot read body %v", err)

		return
	}

	bonds, err := moex.ParseBondsCp1251(buf)
	if err != nil {
		bondsRequest.err = err

		return
	}

	bondsRequest.result = bonds
}

type bondizationRequest struct {
	bondId string
	commonRequest[moex.Bondization]
}

func (bondizationRequest *bondizationRequest) work(client MoexClient) {
	defer close(bondizationRequest.doneChan)

	resp, err := client.innerClient.Get(strings.Replace(allBondizationUrl, "${bond}", bondizationRequest.bondId, 1))

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		bondizationRequest.err = err

		return
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		bondizationRequest.err = fmt.Errorf("cannot read body %v", err)

		return
	}

	bondization, err := moex.ParseBondization(bondizationRequest.bondId, buf)
	if err != nil {
		bondizationRequest.err = err

		return
	}

	bondizationRequest.result = bondization
}
