package perskafka

import (
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	consumer     *kafka.Consumer
	consumerOnce sync.Once
	consumerErr  error
)

func GetConsumer(address, groupID string) (*kafka.Consumer, error) {
	consumerOnce.Do(func() {
		config := &kafka.ConfigMap{
			// Основные настройки подключения
			"bootstrap.servers": address, // Список Kafka-брокеров для подключения
			"group.id":          groupID, // Идентификатор группы потребителей

			// Управление оффсетами
			"auto.offset.reset":        "earliest", // Откуда начинать чтение: earliest (с начала), latest (с конца)
			"enable.auto.offset.store": false,      // Отключаем автоматическое сохранение оффсетов
			"enable.auto.commit":       true,       // Отключаем автоматический коммит оффсетов
			"auto.commit.interval.ms":  5000,       // Интервал коммита

			// Таймауты и надежность
			"socket.timeout.ms":     60000,  // Таймаут сетевого соединения
			"session.timeout.ms":    30000,  // Таймаут сессии потребителя
			"heartbeat.interval.ms": 10000,  // Интервал heartbeat-сообщений
			"max.poll.interval.ms":  300000, // Макс. время между вызовами poll()

			// Настройки переподключения
			"socket.keepalive.enable":  true,  // Включаем TCP keepalive
			"reconnect.backoff.ms":     1000,  // Начальная задержка перед переподключением (1 сек)
			"reconnect.backoff.max.ms": 10000, // Макс. задержка переподключения (10 сек)
		}

		consumer, consumerErr = kafka.NewConsumer(config)
	})

	return consumer, consumerErr
}
