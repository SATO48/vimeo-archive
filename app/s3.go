package app

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/defval/di"
	"github.com/spf13/viper"
)

func S3(c *AppContainer) error {
	return c.Apply(
		di.Provide(func(ctx context.Context) (aws.Config, error) {
			return config.LoadDefaultConfig(
				ctx,
				config.WithRegion("us-east-1"),
				config.WithCredentialsProvider(
					credentials.NewStaticCredentialsProvider(
						viper.GetString("AWS_ACCESS_KEY_ID"),
						viper.GetString("AWS_SECRET_ACCESS_KEY"),
						"",
					),
				),
			)
		}),
		di.Provide(func(cfg aws.Config) *s3.Client {
			return s3.NewFromConfig(cfg, func(o *s3.Options) {
				o.BaseEndpoint = aws.String(viper.GetString("AWS_S3_ENDPOINT"))
			})
		}),
	)
}
