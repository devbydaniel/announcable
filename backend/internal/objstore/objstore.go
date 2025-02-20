package objstore

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/devbydaniel/release-notes-go/config"
	"github.com/devbydaniel/release-notes-go/internal/logger"
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

type ObjStore struct {
	Client *minio.Client
}

type Bucket string

func (b Bucket) String() string {
	return string(b)
}

const (
	ReleaseNotesBucket Bucket = "release-notes"
	LandingPageBucket  Bucket = "landing-page"
)

var buckets = []Bucket{"release-notes", "landing-page"}

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

func (o *ObjStore) GetImageUrl(bucket, path string) (string, error) {
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
	if cfg.Env == "production" {
		// in production, we can directly use the Hetzner ObjStore URL
		return url.String(), nil
	}
	// in development, we need to proxy the URL through the API
	// because Minio doesn't support different URLs for signing and accessing
	urlProxy := strings.Replace(url.String(), "http://objstorage:9000", "/img", 1)
	return urlProxy, nil
}

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
