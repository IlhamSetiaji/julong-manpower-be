package main

import (
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
)

func main() {
	viper := config.NewViper()
	log := config.NewLogrus(viper)
	db := config.NewDatabase()

	// migrate the schema
	err := db.AutoMigrate(&entity.JobPlafon{}, &entity.MPPPeriod{}, &entity.MPPlanningHeader{}, &entity.MPPlanningLine{}, &entity.Major{}, &entity.MPPlanningHeaderAttachment{},
		&entity.RequestCategory{}, &entity.MPRequestHeader{}, &entity.RequestMajor{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Migration success")
	}

	mppPeriods := []entity.MPPPeriod{}

	startDate1, err := time.Parse("2006-01-02", "2024-06-01")
	if err != nil {
		log.Fatal(err)
	}
	endDate1, err := time.Parse("2006-01-02", "2025-07-01")
	if err != nil {
		log.Fatal(err)
	}
	startDate2, err := time.Parse("2006-01-02", "2023-06-01")
	if err != nil {
		log.Fatal(err)
	}
	endDate2, err := time.Parse("2006-01-02", "2024-07-01")
	if err != nil {
		log.Fatal(err)
	}

	mppPeriods = append(mppPeriods, entity.MPPPeriod{
		Title:     "MPP Period 1",
		StartDate: startDate1,
		EndDate:   endDate1,
		Status:    entity.MPPeriodStatusOpen,
	}, entity.MPPPeriod{
		Title:     "MPP Period 2",
		StartDate: startDate2,
		EndDate:   endDate2,
		Status:    entity.MPPeriodStatusComplete,
	})

	for _, mppPeriod := range mppPeriods {
		if err := db.Create(&mppPeriod).Error; err != nil {
			log.Fatalf("Error when creating mppPeriod: " + err.Error())
		}
	}

	log.Info("Seed success")
}
