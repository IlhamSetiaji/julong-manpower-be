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
	err := db.AutoMigrate(&entity.JobPlafon{}, &entity.MPPPeriod{}, &entity.MPPlanningHeader{}, &entity.MPPlanningLine{}, &entity.Major{}, &entity.ManpowerAttachment{},
		&entity.RequestCategory{}, &entity.MPRequestHeader{}, &entity.RequestMajor{}, &entity.MPRequestApprovalHistory{})
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

	requestCategories := []entity.RequestCategory{
		{
			Name:          "Undur Diri",
			IsReplacement: true,
		},
		{
			Name:          "Dimutasikan",
			IsReplacement: true,
		},
		{
			Name:          "Pensiun",
			IsReplacement: true,
		},
		{
			Name:          "Diberhentikan",
			IsReplacement: true,
		},
		{
			Name:          "Promosi",
			IsReplacement: true,
		},
		{
			Name:          "Meninggal Dunia",
			IsReplacement: true,
		},
		{
			Name:          "Posisi Baru",
			IsReplacement: false,
		},
		{
			Name:          "Pegawai Baru",
			IsReplacement: false,
		},
	}

	for _, requestCategory := range requestCategories {
		if err := db.Create(&requestCategory).Error; err != nil {
			log.Fatalf("Error when creating requestCategory: " + err.Error())
		}
	}

	majors := []entity.Major{
		{
			Major:          "Teknik Informatika",
			EducationLevel: entity.EducationLevelEnumBachelor,
		},
		{
			Major:          "Teknik Elektro",
			EducationLevel: entity.EducationLevelEnumBachelor,
		},
		{
			Major:          "Teknik Mesin",
			EducationLevel: entity.EducationLevelEnumBachelor,
		},
	}

	for _, major := range majors {
		if err := db.Create(&major).Error; err != nil {
			log.Fatalf("Error when creating major: " + err.Error())
		}
	}

	log.Info("Seed success")
}
