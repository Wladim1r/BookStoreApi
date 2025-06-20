package perskafka

import (
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	producer     *kafka.Producer
	producerOnce sync.Once
	producerErr  error
)

func GetProducer(address string) (*kafka.Producer, error) {
	producerOnce.Do(func() {
		config := &kafka.ConfigMap{
			// Основные настройки
			"bootstrap.servers": address, // Список брокеров
			"acks":              "all",   // Количество подтверждений: all (ждем все реплики)
			"retries":           3,       // Количество попыток повтора при ошибке
			"retry.backoff.ms":  1000,    // Задержка между повторами (1 сек)

			// Гарантии доставки
			"enable.idempotence": true,   // Идемпотентность (предотвращение дублей)
			"message.timeout.ms": 300000, // Макс. время доставки сообщения (5 мин)

			// Настройки сети
			"socket.keepalive.enable":            true,  // TCP keepalive
			"socket.timeout.ms":                  30000, // Таймаут сокета (30 сек)
			"socket.connection.setup.timeout.ms": 30000, // Таймаут установки соединения (30 сек)

			// Настройки переподключения
			"reconnect.backoff.ms":     1000,  // Начальная задержка переподключения
			"reconnect.backoff.max.ms": 10000, // Макс. задержка переподключения
		}

		producer, producerErr = kafka.NewProducer(config)
	})

	return producer, producerErr
}
