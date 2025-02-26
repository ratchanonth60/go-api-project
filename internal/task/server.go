// internal/task/server.go
package task

import (
	"fmt"
	"net/http"
	"time"

	sesConfig "project-api/internal/infra/config"
	"project-api/internal/infra/logger"

	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/backends/eager"
	broker "github.com/RichardKnop/machinery/v2/brokers/sqs"
	"github.com/RichardKnop/machinery/v2/config"
	lock "github.com/RichardKnop/machinery/v2/locks/eager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/zap"
)

func NewMachineryServer() (*machinery.Server, error) {
	// โหลด AWS credentials และ region จาก config
	awsCredentials := sesConfig.Config.GetCredentialSQS()
	awsSQSConfig := sesConfig.Config.GetSQSConfig()
	brokerSQS := *awsSQSConfig.Endpoint
	// สร้าง AWS session และ SQS client
	sess, err := session.NewSession(&aws.Config{
		Region:      awsSQSConfig.Region,
		Credentials: awsCredentials,
		HTTPClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	})
	if err != nil {
		logger.Error("Failed to create AWS session", zap.Error(err))
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	sqsClient := sqs.New(sess)

	// SQS configuration
	visibilityTimeout := 20
	cfg := &config.Config{
		Broker:        brokerSQS,
		DefaultQueue:  "golang-queue",
		ResultBackend: "", // ไม่ใช้ result backend
		SQS: &config.SQSConfig{
			Client:            sqsClient, // ใช้ SQS client ที่เราสร้าง
			VisibilityTimeout: &visibilityTimeout,
			WaitTimeSeconds:   20,
		},
	}

	// สร้าง Machinery server
	server := machinery.NewServer(
		cfg,
		broker.New(cfg), // SQS broker รับ cfg ไม่ใช่ sess
		eager.New(),     // Eager backend
		lock.New(),      // Eager lock
	)
	logger.Info("Machinery server initialized with SQS",
		zap.String("region", *awsSQSConfig.Region),
		zap.String("queue", cfg.DefaultQueue))
	return server, nil
}
