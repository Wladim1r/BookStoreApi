package service

import (
	"bookstore-api/internal/lib/errs"
	"bookstore-api/internal/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	getAllBooksMethod  = "GetAllBooksMethod"
	getUserBooksMethod = "GetUserBooks"
	postBookMethod     = "PostBookMethod"
	updateBookMethod   = "UpdateBookMethod"
	deleteBookMethod   = "DeleteBookMethod"

	requestType  = "request"
	responseType = "response"
)

func (s *bookService) sendKafkaRequest(req []byte) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &s.topic,
			Partition: kafka.PartitionAny,
		},
		Value: req,
		Key:   nil,
	}

	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	if err := s.producer.Produce(kafkaMsg, deliveryChan); err != nil {
		return fmt.Errorf("%w: %v", errs.ErrKafkaProducer, err)
	}

	e := <-deliveryChan
	switch event := e.(type) {
	case *kafka.Message:
		if event.TopicPartition.Error != nil {
			return fmt.Errorf("%w: %v", errs.ErrKafkaProducer, event.TopicPartition.Error)
		}
		return nil
	case *kafka.Error:
		return fmt.Errorf("%w: %v", errs.ErrKafkaProducer, event)
	}

	return nil
}

func (s *bookService) consumptionMessage() {
	if err := s.consumer.Subscribe(s.topic, nil); err != nil {
		log.Println("service/kafka.go -> 44 line", err)
	}
	defer s.consumer.Close()

	semaphore := make(chan struct{}, 10)

	for {
		msg, err := s.consumer.ReadMessage(-1)
		if err != nil {
			log.Println("service/kafka.go -> 53 line |", err)
		}
		if msg == nil {
			log.Println("Received nil message, skip")
			continue
		}

		var msgCommon struct {
			Method     string `json:"method"`
			Type       string `json:"type"`
			RelationID string `json:"relation_id"`
		}

		if err := json.Unmarshal(msg.Value, &msgCommon); err != nil {
			log.Println("service/kafka.go -> 66 line |", err)
		}

		switch msgCommon.Type {
		case requestType:
			go func() {
				semaphore <- struct{}{}
				defer func() {
					<-semaphore
				}()

				if err := s.proccessRequest(msg, msgCommon.Method); err != nil {
					log.Println("service/kafka.go -> 88 line |", err)
				}
			}()
		case responseType:
			var res models.KafkaBookResponse
			if err := json.Unmarshal(msg.Value, &res); err != nil {
				log.Println("service/kafka.go -> 100 |", err)
			}

			if ch, ok := s.responses.Load(res.RelationID); ok {
				select {
				case ch.(chan models.KafkaBookResponse) <- res:
				default:
					log.Println("service/kafka.go -> 108 | Could not sent chan")
				}
			} else {
				log.Println("Channel not found for RelationID:", res.RelationID)
			}

		}
	}
}

func (s *bookService) proccessRequest(msg *kafka.Message, method string) error {
	var kafkaReq struct {
		Method     string          `json:"method"`
		RelationID string          `json:"relation_id"`
		Payload    json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(msg.Value, &kafkaReq); err != nil {
		return fmt.Errorf("%w: %v", errs.ErrInternal, err)
	}

	switch method {
	case getUserBooksMethod:
		var req models.GetUserBooksRequest
		if err := json.Unmarshal(kafkaReq.Payload, &req); err != nil {
			return fmt.Errorf("%w: %v", errs.ErrInternal, err)
		}

		books, errKafka := s.repo.GetUserBooks(req.UserID, req.Author, req.Title, req.Limit)

		booksRes := models.GetUserBooksResponse{
			Books: books,
		}

		rawMes, err := json.Marshal(booksRes)
		if err != nil {
			return fmt.Errorf("%w: %v", errs.ErrInternal, err)
		}

		res := models.KafkaBookResponse{
			Method:     kafkaReq.Method,
			Type:       responseType,
			RelationID: kafkaReq.RelationID,
			Result:     json.RawMessage(rawMes),
			Error:      errKafka,
		}

		resBytes, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("%w: %v", errs.ErrInternal, err)
		}

		if err := s.proccessResponse(resBytes); err != nil {
			return err
		}
	case postBookMethod:
		var req models.Book
		if err := json.Unmarshal(kafkaReq.Payload, &req); err != nil {
			return fmt.Errorf("%w: %v", errs.ErrInternal, err)
		}

		errKafka := s.repo.PostBook(req)

		res := models.KafkaBookResponse{
			Method:     kafkaReq.Method,
			Type:       responseType,
			RelationID: kafkaReq.RelationID,
			Error:      errKafka,
		}

		resBytes, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("%w: %v", errs.ErrInternal, err)
		}

		if err := s.proccessResponse(resBytes); err != nil {
			return err
		}
	case updateBookMethod:
		var req models.Book
		if err := json.Unmarshal(kafkaReq.Payload, &req); err != nil {
			return fmt.Errorf("%w: %v", errs.ErrInternal, err)
		}

		errKafka := s.repo.UpdateBook(req.UserID, req.ID, req)

		res := models.KafkaBookResponse{
			Method:     kafkaReq.Method,
			Type:       responseType,
			RelationID: kafkaReq.RelationID,
			Error:      errKafka,
		}
		resBytes, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("%w: %v", errs.ErrInternal, err)
		}

		if err := s.proccessResponse(resBytes); err != nil {
			return err
		}
	case deleteBookMethod:
		var req models.DeleteBook
		if err := json.Unmarshal(kafkaReq.Payload, &req); err != nil {
			return fmt.Errorf("%w: %v", errs.ErrInternal, err)
		}

		errKafka := s.repo.DeleteBook(req.UserID, req.ID)

		res := models.KafkaBookResponse{
			Method:     kafkaReq.Method,
			Type:       responseType,
			RelationID: kafkaReq.RelationID,
			Error:      errKafka,
		}
		resBytes, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("%w: %v", errs.ErrInternal, err)
		}

		if err := s.proccessResponse(resBytes); err != nil {
			return err
		}
	}

	return nil
}

func (s *bookService) proccessResponse(res []byte) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &s.topic,
			Partition: kafka.PartitionAny,
		},
		Value: res,
		Key:   nil,
	}

	deliveryChan := make(chan kafka.Event)

	if err := s.producer.Produce(msg, deliveryChan); err != nil {
		return fmt.Errorf("%w: %v", errs.ErrInternal, err)
	}

	e := <-deliveryChan
	switch event := e.(type) {
	case *kafka.Message:
		if event.TopicPartition.Error != nil {
			return fmt.Errorf("%w: %v", errs.ErrKafkaProducer, event.TopicPartition.Error)
		}
		return nil
	case kafka.Error:
		return fmt.Errorf("%w: %v", errs.ErrKafkaProducer, event)
	}

	return nil
}
