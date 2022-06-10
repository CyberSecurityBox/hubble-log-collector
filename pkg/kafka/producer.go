package kafka

import (
	"fmt"
	"github.com/cilium/hubble/pkg/env"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type ProducerInfor struct {
	KafkaBrokers       string
	KafkaUser          string
	KafkaPassword      string
	KafkaVersion       string
	KafkaProducerTopic string
}

type producer struct {
	syncProducer  sarama.SyncProducer
	producerTopic string
	log           logrus.FieldLogger
}

type Producer interface {
	SendMessage(key string, message *string) (int32, int64, error)
	ClearProducer()
}

func NewProducer(log logrus.FieldLogger, producerInfor *ProducerInfor) (Producer, error) {
	/*
		https://stackoverflow.com/questions/50251277/i-always-get-0-partitions-when-i-call-sendmessagemsg-i-am-specifying-12-via-c
	*/
	version, err := sarama.ParseKafkaVersion(producerInfor.KafkaVersion)
	if err != nil {
		log.WithError(err).Error("Can't parse kafka version")
		return nil, err
	}

	configProducer := sarama.NewConfig()

	configProducer.Version = version
	configProducer.Producer.Partitioner = sarama.NewRandomPartitioner
	configProducer.Producer.Return.Successes = true
	configProducer.Producer.RequiredAcks = sarama.WaitForLocal //default
	configProducer.Producer.MaxMessageBytes = 10000000         // default * 10

	configProducer.Net.SASL.Enable = true
	configProducer.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	configProducer.Net.SASL.User = producerInfor.KafkaUser
	configProducer.Net.SASL.Password = producerInfor.KafkaPassword
	configProducer.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
		return &XDGSCRAMClient{
			Client:             nil,
			ClientConversation: nil,
			HashGeneratorFcn:   SHA512,
		}
	}

	syncProducer, err := sarama.NewSyncProducer(strings.Split(producerInfor.KafkaBrokers, ","), configProducer)
	if err != nil {
		log.WithError(err).Error("Can't create sync producer")
		return nil, err
	}

	return &producer{syncProducer: syncProducer, producerTopic: producerInfor.KafkaProducerTopic, log: log}, nil
}

func (p *producer) SendMessage(key string, message *string) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: p.producerTopic,
		Value: sarama.StringEncoder(*message),
	}

	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}

	part, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		p.log.WithError(err).Error("Can't produce message to kafka")
		return 0, 0, err
	}

	return part, offset, nil
}

func (p *producer) ClearProducer() {
	if p.syncProducer != nil {
		p.syncProducer.Close()
	}
}

func ProducerInit(log logrus.FieldLogger) (Producer, error) {
	producerInfor, err := newProducerInfor()
	if err != nil {
		log.WithError(err).Error("Can't create producer infor")
		return nil, err
	}

	producer, err := NewProducer(log, producerInfor)
	if err != nil {
		log.WithError(err).Error("Can't create producer")
		return nil, err
	}

	return producer, nil
}

func newProducerInfor() (*ProducerInfor, error) {
	kafkaBrokers := env.KafkaBrokers
	if kafkaBrokers == "" {
		return nil, fmt.Errorf("Missing kafka brokers in env")
	}

	kafkaUser := env.KafkaUser
	if kafkaUser == "" {
		return nil, fmt.Errorf("Missing kafka user in env")
	}

	kafkaPassword := env.KafkaPass
	if kafkaPassword == "" {
		return nil, fmt.Errorf("Missing kafka password in env")
	}

	kafkaProducerTopic := env.KafkaTopic
	if kafkaProducerTopic == "" {
		return nil, fmt.Errorf("Missing kafka producer topic in env")
	}

	return &ProducerInfor{
		KafkaBrokers:       kafkaBrokers,
		KafkaUser:          kafkaUser,
		KafkaPassword:      kafkaPassword,
		KafkaProducerTopic: kafkaProducerTopic,
		KafkaVersion:       env.KafkaVersion,
	}, nil
}