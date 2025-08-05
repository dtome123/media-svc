package media

import (
	"media-svc/internal/models"

	mongodb "github.com/dtome123/go-mongo-generic"
)

type MediaRepository struct {
	mediaCol mongodb.Collection[models.Media]
}

func NewMediaRepository(db *mongodb.Database) *MediaRepository {

	mediaCol := mongodb.NewCollection[models.Media](db)
	mediaCol.EnsureIndexes(GetMediaIndexes())

	return &MediaRepository{
		mediaCol: mediaCol,
	}
}
