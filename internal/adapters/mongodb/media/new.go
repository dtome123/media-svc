package media

import (
	"media-svc/internal/models"

	mongodb "github.com/dtome123/go-mongo-generic"
)

type MediaRepository struct {
	mediaCol        mongodb.Collection[models.Media]
	transcodeJobCol mongodb.Collection[models.TranscodeJob]
}

func NewMediaRepository(db *mongodb.Database) *MediaRepository {

	mediaCol := mongodb.NewCollection[models.Media](db)
	mediaCol.EnsureIndexes(GetMediaIndexes())

	transcodeJobCol := mongodb.NewCollection[models.TranscodeJob](db)
	transcodeJobCol.EnsureIndexes(GetTranscodeJobIndexes())

	return &MediaRepository{
		mediaCol:        mediaCol,
		transcodeJobCol: transcodeJobCol,
	}
}
