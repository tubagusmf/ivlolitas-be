package storage

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryStorage struct {
	client *cloudinary.Cloudinary
}

func NewCloudinaryStorage(cloudName string, apiKey string, apiSecret string) (Storage, error) {
	cld, err := cloudinary.NewFromParams(
		cloudName,
		apiKey,
		apiSecret,
	)
	if err != nil {
		return nil, err
	}

	return &CloudinaryStorage{
		client: cld,
	}, nil
}

func (c *CloudinaryStorage) Upload(ctx context.Context, file multipart.File, folder string) (*UploadResult, error) {
	result, err := c.client.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{
			Folder: folder,
		},
	)

	if err != nil {
		return nil, err
	}

	return &UploadResult{
		URL:      result.SecureURL,
		PublicID: result.PublicID,
	}, nil
}

func (c *CloudinaryStorage) Delete(ctx context.Context, publicID string) error {
	_, err := c.client.Upload.Destroy(
		ctx,
		uploader.DestroyParams{
			PublicID: publicID,
		},
	)

	return err
}
