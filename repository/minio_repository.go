package repository

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/mxbikes/mxbikesclient.service.modImage/models"
)

type ModImageMinioRepository interface {
	GetModImagesByModID(ctx context.Context, modID string) ([]*models.ModImage, error)
}

type minioRepository struct {
	dbPool *minio.Client
}

func NewMinioRepository(c *minio.Client) *minioRepository {
	return &minioRepository{dbPool: c}
}

func (m *minioRepository) GetModImagesByModID(ctx context.Context, modID string) ([]*models.ModImage, error) {
	var bucket = "mod-images"
	objectCh := m.dbPool.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
		Prefix:       modID,
		Recursive:    true,
		WithMetadata: false,
	})

	var modImages []*models.ModImage
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			continue
		}
		modImages = append(modImages, &models.ModImage{
			Name:   object.Key,
			Bucket: bucket,
			Url:    "/" + bucket + "/" + object.Key,
		})
	}

	return modImages, nil
}
