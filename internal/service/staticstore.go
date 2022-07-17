package service

import (
	"bonds_calculator/internal/api"
	"bonds_calculator/internal/model/datastuct/box"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/util"
	"fmt"
	"github.com/benbjohnson/clock"
	log "github.com/sirupsen/logrus"
	"time"
)

//go:generate mockgen -destination=mock/staticstore_gen.go . IStaticStoreService
type IStaticStoreService interface {
	GetBonds() []moex.Bond
	GetBondsWithUpdateTime() ([]moex.Bond, time.Time)
	GetBondById(id string) (moex.Bond, error)
	GetBondsChangedTime() time.Time

	GetBondization(id string) (moex.Bondization, error)
	GetBondizationsChangedTime() time.Time
}

type StaticStoreService struct {
	client api.IMoexClient

	bonds        box.ConcurrentCacheBox[[]moex.Bond]
	bondsMap     box.ConcurrentBox[map[string]moex.Bond]
	bondizations box.ConcurrentCacheBox[map[string]moex.Bondization]

	clock clock.Clock
}

func NewStaticStoreService(client api.IMoexClient, timer ITimerService, clock clock.Clock) IStaticStoreService {
	staticStore := StaticStoreService{
		client: client,

		clock: clock,
	}

	staticStore.reloadBond()
	staticStore.reloadBondization()

	timer.SubscribeEvery(time.Minute*5, staticStore.reloadBond)
	timer.SubscribeEveryStartFrom(time.Hour*24, util.GetMoexMidnight(staticStore.clock), staticStore.reloadBondization)

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

func (staticStore *StaticStoreService) GetBondById(id string) (moex.Bond, error) {
	result, exist := staticStore.bondsMap.SafeRead()[id]

	if !exist {
		return moex.Bond{}, fmt.Errorf("bond with id: %s not found", id)
	}

	return result, nil
}

func (staticStore *StaticStoreService) GetBondization(id string) (moex.Bondization, error) {
	result, exist := staticStore.bondizations.SafeRead()[id]

	if !exist {
		return moex.Bondization{}, fmt.Errorf("bondization with id: %s not found", id)
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
	log.Info("StaticStoreService: Bonds updating started")

	bonds, err := staticStore.client.GetBonds()

	for err != nil {
		log.Errorf("StaticStoreService: Error while updating bonds: %v, retrying...", err)

		bonds, err = staticStore.client.GetBonds()
	}

	for i, end := 0, len(bonds); i < end; i++ {
		if err := bonds[i].IsValid(); err != nil { // impossible cause of tests but just in case
			log.Errorf("StaticStoreService: Got invalid bond: %v, error: %v", bonds[i], err)

			bonds[i] = bonds[end-1]
			bonds = bonds[:end-1]
			i--
			end--
		}
	}

	staticStore.bondsMap.Set(util.SliceToMapBy(bonds, func(bond moex.Bond) string { return bond.Id }))
	staticStore.bonds.Set(bonds)

	log.Info("StaticStoreService: Bonds updated successfully")
}

func (staticStore *StaticStoreService) reloadBondization() {
	log.Info("StaticStoreService: Bondizations updating started")

	bonds := staticStore.GetBonds()

	bondizations := make(map[string]moex.Bondization, len(bonds))

	for _, bond := range bonds {
		bondization, err := staticStore.client.GetBondization(bond.Id)

		for tryCount := 0; err != nil && tryCount < 5; tryCount++ {
			log.Errorf("StaticStoreService: Error while updating bondization for bond id: %s, error: %v", bond.Id, err)

			bondization, err = staticStore.client.GetBondization(bond.Id)
		}

		if err != nil {
			log.Errorf("StaticStoreService: Error while updating bondization for bond id: %s, error: %v, skipping", bond.Id, err)

			continue
		}

		if err := bondization.IsValid(bond.EndDate); err != nil { // impossible cause of tests but just in case
			log.Errorf("StaticStoreService: Got invalid bondization for bond id: %s, error: %v", bond.Id, err)

			continue
		}

		bondizations[bond.Id] = bondization
	}

	staticStore.bondizations.Set(bondizations)

	log.Info("StaticStoreService: Bondizations updating ended")
}
