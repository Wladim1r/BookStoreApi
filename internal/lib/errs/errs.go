package errs

import "errors"

var (
	ErrInvalidParam  = errors.New("invalid parameter value")
	ErrInvalidID     = errors.New("invalid ID format")
	ErrNotFound      = errors.New("record not found")
	ErrDBOperation   = errors.New("database operation failed")
	ErrInternal      = errors.New("internal server error")
	ErrNotRegistred  = errors.New("you have not registered yet")
	ErrNotAuthorized = errors.New("you are not authorized")
	ErrTimeout       = errors.New("request time expired")
	ErrKafkaProducer = errors.New("failed to produce message")
	ErrKafkaConsumer = errors.New("failed to consume message")
)
