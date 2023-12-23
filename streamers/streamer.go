package streamers

import (
	"context"
	"fmt"
	"kafka-connector/configs"
	"kafka-connector/loggers"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Streamer struct {
	kinesisClient *kinesis.Client
	conf          *configs.ViperConfig
	mutex         sync.Mutex
}

func NewStreamer() *Streamer {
	conf := configs.GetConfig()
	session, err := getOriginalSession()
	if err != nil {
		loggers.GlobalLogger.Fatal("fail to create service session:", err)
	}
	crossAccount := createCrossSession(session)
	kinesisClient, err := createCrossKinesisSession(crossAccount)
	if err != nil {
		loggers.GlobalLogger.Fatal("fail to create kinesis client:,", err)
	}
	return &Streamer{
		kinesisClient: kinesisClient,
		conf:          conf,
	}
}

func getOriginalSession() (aws.Config, error) {
	session, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return session, err
	}
	return session, nil
}

func createCrossSession(session aws.Config) aws.CredentialsProviderFunc {
	conf := configs.GetConfig()
	stsClient := sts.NewFromConfig(session)
	roleSessionName := fmt.Sprintf("doka-%s-%s", viper.GetString("ENV"), uuid.NewString())
	stsInput := &sts.AssumeRoleInput{
		RoleArn:         aws.String(conf.GetString("aws.crossAccountArn")),
		RoleSessionName: aws.String(roleSessionName),
	}
	stsOutput, err := stsClient.AssumeRole(context.TODO(), stsInput)
	if err != nil {
		loggers.GlobalLogger.Fatal(err)
	}
	//
	crossAccountSts := aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     *stsOutput.Credentials.AccessKeyId,
			SecretAccessKey: *stsOutput.Credentials.SecretAccessKey,
			SessionToken:    *stsOutput.Credentials.SessionToken,
			Source:          conf.GetString("aws.crossAccountArn"),
		}, nil
	})
	return crossAccountSts
}

func createCrossKinesisSession(crossAccountSts aws.CredentialsProviderFunc) (*kinesis.Client, error) {
	conf := configs.GetConfig()
	kinesisConfig := aws.Config{
		Region:      conf.GetString("aws.region"),
		Credentials: aws.NewCredentialsCache(crossAccountSts),
	}
	kinesisClient := kinesis.NewFromConfig(kinesisConfig)
	return kinesisClient, nil
}

func (s *Streamer) PutStreamData(message string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	streamName := s.conf.GetString("aws.kinesis.streamName")
	partitionKey := fmt.Sprintf("%s-%s", s.conf.GetString("aws.kinesis.partitionKey"), uuid.NewString())
	input := &kinesis.PutRecordInput{
		Data:         []byte(message),
		StreamName:   &streamName,
		PartitionKey: aws.String(partitionKey),
	}

	_, err := s.kinesisClient.PutRecord(context.TODO(), input)
	if err != nil {
		// TODO: Create metric to monitoring fail kinesis put record
		loggers.GlobalLogger.Printf("fail to put record: %v", err)
	}
	return nil
}

func (s *Streamer) RenewClient() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	for {
		loggers.GlobalLogger.Println("crossaccount session renew after 30 minute")
		select {
		case <-ticker.C:
			s.mutex.Lock()
			session, _ := getOriginalSession()
			newCrossAccount := createCrossSession(session)
			newKinesisClient, _ := createCrossKinesisSession(newCrossAccount)
			s.kinesisClient = newKinesisClient
			loggers.GlobalLogger.Println("renew crossaccount session")
			s.mutex.Unlock()
		}
	}
}
