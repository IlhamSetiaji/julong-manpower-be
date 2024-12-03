package response

import "github.com/IlhamSetiaji/julong-manpower-be/internal/entity"

type FindAllPaginatedMPPPeriodResponse struct {
	MPPPeriods *[]entity.MPPPeriod `json:"mppperiods"`
	Total      int64               `json:"total"`
}

type FindByIdMPPPeriodResponse struct {
	MPPPeriod *entity.MPPPeriod `json:"mppperiod"`
}
