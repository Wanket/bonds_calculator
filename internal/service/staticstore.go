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

type StaticStoreService struct {
	client api.IMoexClient

	bonds        box.ConcurrentCacheBox[[]moex.Bond]
	bondsMap     box.ConcurrentBox[map[string]moex.Bond]
	bondizations box.ConcurrentCacheBox[map[string]moex.Bondization]

	clock clock.Clock
}

func NewStaticStoreService(client api.IMoexClient, timer ITimerService, clock clock.Clock) *StaticStoreService {
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
	result := staticStore.bonds.LockAndRead()
	staticStore.bonds.UnlockRead()

	return result
}

func (staticStore *StaticStoreService) GetBondById(id string) (moex.Bond, error) {
	result, exist := staticStore.bondsMap.LockAndRead()[id]
	staticStore.bondsMap.UnlockRead()

	if !exist {
		return moex.Bond{}, fmt.Errorf("bond with id: %s not found", id)
	}

	return result, nil
}

func (staticStore *StaticStoreService) GetBondization(id string) (moex.Bondization, error) {
	result, exist := staticStore.bondizations.LockAndRead()[id]
	staticStore.bondizations.UnlockRead()

	if !exist {
		return moex.Bondization{}, fmt.Errorf("bondization with id: %s not found", id)
	}

	return result, nil
}

func (staticStore *StaticStoreService) GetBondsChangedTime() time.Time {
	_, bondsTime := staticStore.bonds.LockAndReadWithTime()
	staticStore.bonds.UnlockRead()

	return bondsTime
}

func (staticStore *StaticStoreService) GetBondizationsChangedTime() time.Time {
	_, bondizationsTime := staticStore.bondizations.LockAndReadWithTime()
	staticStore.bondizations.UnlockRead()

	return bondizationsTime
}

func (staticStore *StaticStoreService) reloadBond() {
	log.Info("StaticStoreService: Bonds updating started")

	bonds, err := staticStore.client.GetBonds()

	for err != nil {
		log.Errorf("StaticStoreService: Error while updating bonds: %v, retrying...", err)

		bonds, err = staticStore.client.GetBonds()
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

		bondizations[bond.Id] = bondization
	}

	staticStore.bondizations.Set(bondizations)

	log.Info("StaticStoreService: Bondizations updating ended")
}
