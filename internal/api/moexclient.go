package api

import (
	"bonds_calculator/internal/model/moex"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

//nolint:lll
const (
	allBondsURL = "https://iss.moex.com/iss/engines/stock/markets/bonds/securities.csv?iss.meta=off&iss.only=marketdata,securities&securities.columns=SECID,SHORTNAME,COUPONVALUE,NEXTCOUPON,ACCRUEDINT,PREVPRICE,FACEVALUE,COUPONPERIOD,MINSTEP,COUPONPERCENT,MATDATE,FACEUNIT&marketdata.columns=SECID,LCURRENTPRICE"

	allBondizationURL = "https://iss.moex.com/iss/securities/${bond}/bondization.csv?limit=unlimited&iss.meta=off&iss.only=amortizations,coupons&amortizations.columns=amortdate,value&coupons.columns=coupondate,value"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mock/moexclient_gen.go . IMoexClient
type IMoexClient interface {
	Close() error

	GetBonds() ([]moex.Bond, error)
	GetBondization(ID string) (moex.Bondization, error)
}

type MoexClient struct {
	innerClient http.Client

	workQueue chan worker

	context context.Context //nolint:containedctx
	cancel  context.CancelFunc
}

func NewMoexClient(queueSize int) *MoexClient {
	return NewMoexClientWithContext(context.Background(), queueSize)
}

func NewMoexClientWithContext(ctx context.Context, queueSize int) *MoexClient {
	ctx, cancel := context.WithCancel(ctx)
	client := MoexClient{
		innerClient: http.Client{},
		workQueue:   make(chan worker, queueSize),
		context:     ctx,
		cancel:      cancel,
	}

	for i := 0; i < queueSize; i++ {
		go queueListener(client)
	}

	return &client
}

func (client *MoexClient) Close() error {
	client.cancel()

	return nil
}

func (client *MoexClient) GetBonds() ([]moex.Bond, error) {
	request := bondsRequest{newCommonRequest[[]moex.Bond]()}

	return request.commonReturn(client, &request)
}

func (client *MoexClient) GetBondization(id string) (moex.Bondization, error) {
	request := bondizationRequest{
		bondID:        id,
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

func (request *commonRequest[R]) commonReturn(client *MoexClient, worker worker) (R, error) { //nolint:ireturn
	client.workQueue <- worker

	<-request.doneChan

	return request.result, request.err
}

type bondsRequest struct {
	commonRequest[[]moex.Bond]
}

func (bondsRequest *bondsRequest) work(client MoexClient) {
	defer close(bondsRequest.doneChan)

	resp, err := client.innerClient.Get(allBondsURL)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		bondsRequest.err = err

		return
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		bondsRequest.err = fmt.Errorf("cannot read body %w", err)

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
	bondID string
	commonRequest[moex.Bondization]
}

func (bondizationRequest *bondizationRequest) work(client MoexClient) {
	defer close(bondizationRequest.doneChan)

	resp, err := client.innerClient.Get(strings.Replace(allBondizationURL, "${bond}", bondizationRequest.bondID, 1))

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		bondizationRequest.err = err

		return
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		bondizationRequest.err = fmt.Errorf("cannot read body %w", err)

		return
	}

	bondization, err := moex.ParseBondization(bondizationRequest.bondID, buf)
	if err != nil {
		bondizationRequest.err = err

		return
	}

	bondizationRequest.result = bondization
}
