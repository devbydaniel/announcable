package objstore

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/devbydaniel/announcable/config"
	"github.com/devbydaniel/announcable/internal/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	log = logger.Get()
	cfg = config.New()
)

var (
	endpoint = cfg.ObjStorage.Endpoint
	minioCfg = minio.Options{
		Secure: cfg.ObjStorage.UseSSL,
		Creds: credentials.NewStaticV4(
			cfg.ObjStorage.AccessKey,
			cfg.ObjStorage.SecretKey,
			"",
		),
	}
	bucketOptions = minio.MakeBucketOptions{Region: cfg.ObjStorage.Region}
)

// ObjStore wraps a Minio client for object storage operations.
type ObjStore struct {
	Client *minio.Client
}

// Bucket represents a named storage bucket.
type Bucket string

func (b Bucket) String() string {
	return string(b)
}

// Available storage buckets.
const (
	ReleaseNotesBucket Bucket = "release-notes"
	LandingPageBucket  Bucket = "landing-page"
)

var buckets = []Bucket{"release-notes", "landing-page"}

// Init initializes the object store client and creates required buckets.
func Init(ctx context.Context) (*ObjStore, error) {
	log.Trace().Msg("Init")
	client, err := createClient(endpoint, &minioCfg)
	if err != nil {
		log.Error().Err(err).Msg("Error creating client")
		return nil, err
	}
	err = createBuckets(client, ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating buckets")
		return nil, err
	}
	return &ObjStore{Client: client}, nil
}

func createClient(endpoint string, config *minio.Options) (*minio.Client, error) {
	log.Trace().Msg("createClient")
	client, err := minio.New(endpoint, config)
	if err != nil {
		log.Error().Err(err).Msg("Error creating client")
		return nil, err
	}
	return client, nil
}

func createBuckets(client *minio.Client, ctx context.Context) error {
	log.Trace().Msg("createBuckets")
	for _, bucket := range buckets {
		exists, err := client.BucketExists(ctx, bucket.String())
		if err != nil {
			log.Error().Err(err).Str("bucket", bucket.String()).Msg("Error checking bucket exists")
			return err
		}
		if exists {
			log.Debug().Str("bucket", bucket.String()).Msg("Bucket already exists, skipping")
			continue
		}
		err = client.MakeBucket(ctx, bucket.String(), bucketOptions)
		if err != nil {
			log.Error().Err(err).Str("bucket", bucket.String()).Msg("Error creating bucket")
			return err
		}
		log.Info().Str("bucket", bucket.String()).Msg("Bucket created")

	}
	return nil
}

// GetImageURL returns a presigned URL for accessing the image at the given bucket and path.
func (o *ObjStore) GetImageURL(bucket, path string) (string, error) {
	cfg := config.New()
	ctx := context.Background()
	// check if object exists
	_, err := o.Client.StatObject(ctx, bucket, path, minio.StatObjectOptions{})
	if err != nil {
		log.Debug().Err(err).Msg("No image found")
		return "", nil
	}
	expiresIn := time.Second * 60 * 60 * 24
	url, err := o.Client.PresignedGetObject(ctx, bucket, path, expiresIn, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error getting image url")
		return "", err
	}

	// Minio doesn't support different URLs for signing and accessing
	// so we proxy through the API using our public base URL
	internalScheme := "http://"
	if cfg.ObjStorage.UseSSL {
		internalScheme = "https://"
	}
	internalURL := internalScheme + cfg.ObjStorage.Endpoint

	// Build public URL, handling case where BaseURL may already include scheme
	var publicURL string
	if strings.HasPrefix(cfg.BaseURL, "http://") || strings.HasPrefix(cfg.BaseURL, "https://") {
		publicURL = cfg.BaseURL + "/api/img"
	} else {
		publicScheme := "http://"
		if cfg.Env == "production" {
			publicScheme = "https://"
		}
		publicURL = publicScheme + cfg.BaseURL + "/api/img"
	}

	urlProxy := strings.Replace(url.String(), internalURL, publicURL, 1)
	return urlProxy, nil
}

// UpdateImage uploads or replaces an image in the specified bucket and path.
func (o *ObjStore) UpdateImage(bucket, path string, img *io.Reader) error {
	log.Trace().Msg("UpdateImage")
	ctx := context.Background()
	info, err := o.Client.PutObject(ctx, bucket, path, *img, -1, minio.PutObjectOptions{})
	log.Debug().Interface("info", info).Msg("PutObject")
	if err != nil {
		log.Error().Err(err).Msg("Error uploading image")
		return err
	}
	return nil
}

// DeleteImage removes an image from the specified bucket and path.
func (o *ObjStore) DeleteImage(bucket, path string) error {
	log.Trace().Str("path", path).Str("bucket", bucket).Msg("DeleteImage")
	ctx := context.Background()
	err := o.Client.RemoveObject(ctx, bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Error deleting image")
		return err
	}
	return nil
}
