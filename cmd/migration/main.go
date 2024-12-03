package main

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
)

func main() {
	viper := config.NewViper()
	log := config.NewLogrus(viper)
	db := config.NewDatabase()

	// migrate the schema
	err := db.AutoMigrate(&entity.JobPlafon{}, &entity.MPPPeriod{}, &entity.MPPlanningHeader{}, &entity.MPPlanningLine{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Migration success")
	}
}
