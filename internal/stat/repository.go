package stat

import (
	"linkshorter/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (s *StatRepository) AddClick(linkId uint) {
	var stat Stat

	currentDate := datatypes.Date(time.Now())
	s.Db.Find(&stat, "link_id = ? and date = ? ", linkId, currentDate)
	if stat.ID == 0 {
		s.Db.Create(&Stat{
			LinkId: linkId,
			Clicks: 1,
			Date:   currentDate,
		})
	} else {
		stat.Clicks += 1
		s.Db.Save(&stat)
	}
}

func (s *StatRepository) GetStats(by string, from, to time.Time) []GetStatResponse {
	var stats []GetStatResponse
	var selectQuery string

	switch by {
	case GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case GroupByMounth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}

	s.DB.Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)

	return stats
}
