package cmd

import (
	"context"
	"fmt"
	"net/http"
	"order-services/clients"
	response "order-services/common/reponses"
	"order-services/config"
	"order-services/constants"
	controllers "order-services/controllers/http"
	kafkaController "order-services/controllers/kafka"
	kafkaConfig "order-services/controllers/kafka/config"
	"os"
	"os/signal"
	"syscall"

	"order-services/domain/models"
	"order-services/middlewares"
	"order-services/repositories"
	"order-services/routes"
	"order-services/services"
	"time"

	"github.com/IBM/sarama"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load()
		config.Init()
		db, err := config.InitDatabase()
		if err != nil {
			panic(err)
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}
		time.Local = loc

		err = db.AutoMigrate(
			&models.Order{},
			&models.OrderField{},
			&models.OrderHistory{},
		)

		if err != nil {
			panic(err)
		}

		client := clients.NewClientRegistry()
		repository := repositories.NewRepositoryRegistry(db)
		service := services.NewServiceRegistry(repository, client)
		controller := controllers.NewContollerRegistry(service)

		serverHttp(controller, client)
		serverKafkaConsumer(service)
	},
}

func Run() {
	err := command.Execute()

	if err != nil {
		panic(err)
	}
}

func serverHttp(
	controller controllers.IControllerRegistry,
	client clients.IClientRegistry,
) {
	router := *gin.Default()
	router.Use(middlewares.HandlePanic())
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, response.Response{
			Status:  constants.Error,
			Message: fmt.Sprintf("Path %s", http.StatusNotFound),
		})
	})

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, response.Response{
			Status:  constants.Success,
			Message: "Welcome to payment service",
		})
	})
	router.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-service-name, x-api-key, x-request-at")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	})
	lmt := tollbooth.NewLimiter(
		config.Config.RateLimiterMaxRequest,
		&limiter.ExpirableOptions{
			DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSecond) * time.Second,
		},
	)
	router.Use(middlewares.RateLimiter(lmt))

	group := router.Group("/api/v1")
	route := routes.NewRouteRegistry(controller, client, group)
	route.Serve()

	// Go routine
	go func() {
		port := fmt.Sprintf(":%d", config.Config.Port)
		router.Run(port)
	}()
}

func serverKafkaConsumer(service services.IServiceRegistry) {
	kafkaConsumerConfig := sarama.NewConfig()
	kafkaConsumerConfig.Consumer.MaxWaitTime = time.Duration(config.Config.Kafka.MaxWaitTimeInMs) * time.Millisecond
	kafkaConsumerConfig.Consumer.MaxProcessingTime = time.Duration(config.Config.Kafka.MaxProcessingTimeInMs) * time.Millisecond
	kafkaConsumerConfig.Consumer.Retry.Backoff = time.Duration(config.Config.Kafka.BackOffTimeInMs)
	kafkaConsumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	kafkaConsumerConfig.Consumer.Offsets.AutoCommit.Enable = true
	kafkaConsumerConfig.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	kafkaConsumerConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	brokers := config.Config.Kafka.Brokers
	groupID := config.Config.Kafka.GroupID
	topics := config.Config.Kafka.Topic
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, kafkaConsumerConfig)
	if err != nil {
		logrus.Errorf("failed to create consumer group: %v", err)
		return
	}
	defer consumerGroup.Close()

	consumer := kafkaConfig.NewConsumerGroup()
	kafkaRegistry := kafkaController.NewKafkaRegistry(service)
	kafkaConsumer := kafkaConfig.NewKafkaConsumer(consumer, kafkaRegistry)
	kafkaConsumer.Register()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			err = consumerGroup.Consume(ctx, topics, consumer)
			if err != nil {
				logrus.Errorf("failed to consume: %v", err)
				panic(err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	logrus.Info("Kafka consumer started")
	<-signals
	logrus.Info("Kafka consumer stopped")

}
