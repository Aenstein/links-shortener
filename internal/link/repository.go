package link

import (
	"linkshorter/pkg/db"

	"gorm.io/gorm/clause"
)

type LinkRepository struct {
	Database *db.Db
}

func NewLinkRepository(database *db.Db) *LinkRepository {
	return &LinkRepository{
		Database: database,
	}
}

func (repository *LinkRepository) Create(link *Link) (*Link, error) {
	result := repository.Database.DB.Create(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repository *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	
	result := repository.Database.DB.First(&link, "hash = ?", hash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repository *LinkRepository) Update(link *Link) (*Link, error) {
	result := repository.Database.DB.Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repository *LinkRepository) Delete(id uint64) (error) {
	result := repository.Database.DB.Delete(&Link{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repository *LinkRepository) GetById(id uint64) (error) {
	result := repository.Database.DB.First(&Link{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repository *LinkRepository) Count() int64 {
	var count int64

	repository.Database.
		Table("links").
		Where("deleted_at is null").
		Count(&count)
	
	return count
}

func (repository *LinkRepository) GetAll(limit, offset uint) ([]Link) {
	var links []Link

	repository.Database.
		Table("links").
		Where("deleted_at is null").
		Order("id asc").
		Limit(int(limit)).
		Offset(int(offset)).
		Scan(&links)
	
	return links
}