package service

import (
	"bonds_calculator/internal/api"
	"bonds_calculator/internal/model/datastruct/box"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/util"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	errBondByIDNotFound        = fmt.Errorf("bond by id not found")
	errBondizationByIDNotFound = fmt.Errorf("bondization by id not found")
)

//go:generate go run github.com/golang/mock/mockgen -destination=mock/staticstore_gen.go . IStaticStoreService
type IStaticStoreService interface {
	GetBonds() []moex.Bond
	GetBondsWithUpdateTime() ([]moex.Bond, time.Time)
	GetBondByID(id string) (moex.Bond, error)
	GetBondsChangedTime() time.Time

	GetBondization(id string) (moex.Bondization, error)
	GetBondizationsChangedTime() time.Time
}

type StaticStoreService struct {
	client api.IMoexClient

	bonds        box.ConcurrentCacheBox[[]moex.Bond]
	bondsMap     box.ConcurrentBox[map[string]moex.Bond]
	bondizations box.ConcurrentCacheBox[map[string]moex.Bondization]

	timeHelper util.ITimeHelper
}

func NewStaticStoreService(
	client api.IMoexClient,
	timer ITimerService,
	timeHelper util.ITimeHelper,
) *StaticStoreService {
	staticStore := StaticStoreService{
		client: client,

		bonds:        box.ConcurrentCacheBox[[]moex.Bond]{},
		bondsMap:     box.ConcurrentBox[map[string]moex.Bond]{},
		bondizations: box.ConcurrentCacheBox[map[string]moex.Bondization]{},

		timeHelper: timeHelper,
	}

	staticStore.reloadBond()
	staticStore.reloadBondization()

	timer.SubscribeEvery(time.Minute*5, staticStore.reloadBond) //nolint:gomnd
	timer.SubscribeEveryStartFrom(
		util.Day,
		staticStore.timeHelper.GetMoexMidnight().Add(util.Day), staticStore.reloadBondization,
	)

	return &staticStore
}

func (staticStore *StaticStoreService) GetBonds() []moex.Bond {
	result := staticStore.bonds.SafeRead()

	return result
}

func (staticStore *StaticStoreService) GetBondsWithUpdateTime() ([]moex.Bond, time.Time) {
	result, updateTime := staticStore.bonds.SafeReadWithTime()

	return result, updateTime
}

func (staticStore *StaticStoreService) GetBondByID(id string) (moex.Bond, error) {
	result, exist := staticStore.bondsMap.SafeRead()[id]

	if !exist {
		return moex.Bond{}, fmt.Errorf("GetBondByID: %w, id: %s", errBondByIDNotFound, id)
	}

	return result, nil
}

func (staticStore *StaticStoreService) GetBondization(id string) (moex.Bondization, error) {
	result, exist := staticStore.bondizations.SafeRead()[id]

	if !exist {
		return moex.Bondization{}, fmt.Errorf("GetBondization: %w, id: %s", errBondizationByIDNotFound, id)
	}

	return result, nil
}

func (staticStore *StaticStoreService) GetBondsChangedTime() time.Time {
	_, bondsTime := staticStore.bonds.SafeReadWithTime()

	return bondsTime
}

func (staticStore *StaticStoreService) GetBondizationsChangedTime() time.Time {
	_, bondizationsTime := staticStore.bondizations.SafeReadWithTime()

	return bondizationsTime
}

func (staticStore *StaticStoreService) reloadBond() {
	log.Info("StaticStoreService: bonds updating started")

	bonds, err := staticStore.client.GetBonds()

	for err != nil {
		log.WithError(err).Error("StaticStoreService: error while updating bonds, retrying...")

		bonds, err = staticStore.client.GetBonds()
	}

	for bondInx, end := 0, len(bonds); bondInx < end; bondInx++ {
		if err := bonds[bondInx].IsValid(); err != nil { // impossible cause of tests but just in case
			log.WithFields(log.Fields{
				"bond":       bonds[bondInx],
				log.ErrorKey: err,
			}).Errorf("StaticStoreService: got invalid bond")

			bonds[bondInx] = bonds[end-1]
			bonds = bonds[:end-1]
			bondInx--
			end--
		}
	}

	staticStore.bondsMap.Set(util.SliceToMapBy(bonds, func(bond moex.Bond) string { return bond.ID }))
	staticStore.bonds.Set(bonds)

	log.WithField("count", len(bonds)).Info("StaticStoreService: bonds updated")
}

func (staticStore *StaticStoreService) reloadBondization() {
	log.Info("StaticStoreService: bondizations updating started")

	bonds := staticStore.GetBonds()

	bondizations := make(map[string]moex.Bondization, len(bonds))

	bondInx := -1
	for bondiztionResult := range staticStore.client.GetBondizationsAsync(bonds) {
		bondInx++

		for tryCount := 0; bondiztionResult.Error != nil && tryCount < 5; tryCount++ {
			log.WithFields(log.Fields{
				"bond":       bondiztionResult.Bond,
				log.ErrorKey: bondiztionResult.Error,
				"tryCount":   tryCount,
			}).Errorf("StaticStoreService: error while updating bondization, retrying...")

			bondiztionResult.Bondization, bondiztionResult.Error =
				staticStore.client.GetBondization(bondiztionResult.Bond.ID)
		}

		if bondiztionResult.Error != nil {
			log.WithFields(log.Fields{
				"bond":       bondiztionResult.Bond,
				log.ErrorKey: bondiztionResult.Error,
			}).Errorf("StaticStoreService: error while updating bondization, skipping")

			continue
		}

		// impossible cause of tests but just in case
		if err := bondiztionResult.Bondization.IsValid(bondiztionResult.Bond.EndDate); err != nil {
			log.WithFields(log.Fields{
				"bond":        bondiztionResult.Bond,
				"bondization": bondiztionResult.Bondization,
				log.ErrorKey:  err,
			}).Errorf("StaticStoreService: got invalid bondization")

			continue
		}

		bondizations[bondiztionResult.Bond.ID] = bondiztionResult.Bondization

		if bondInx%250 == 0 {
			log.WithField("bondInx", bondInx).Info("StaticStoreService: updating bondizations...")
		}
	}

	staticStore.bondizations.Set(bondizations)

	log.WithField("count", len(bondizations)).Info("StaticStoreService: bondizations updating ended")
}
