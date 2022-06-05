package env

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	KafkaTopic = os.Getenv("KAFKA_TOPIC")
	KafkaUser = os.Getenv("KAFKA_USER")
	KafkaPass = os.Getenv("KAFKA_PASS")
	KafkaBrokers = os.Getenv("KAFKA_BROKERS")
	KafkaVersion = "2.7.0"
	ClusterName = os.Getenv("K8S_CLUSTER_NAME")
)

var Log *logrus.Logger