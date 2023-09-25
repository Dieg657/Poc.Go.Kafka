package base

import (
	"errors"

	"github.com/google/uuid"
)

type IMessage[TData any] interface {
	New(data TData) error
	NewWithSameCorrelationId(data TData, correlationId uuid.UUID) error
}

type Message[TData any] struct {
	CorrelationId uuid.UUID `json:"correlationId"`
	Data          TData     `json:"data"`
}

func (message *Message[TData]) New(data *TData) error {
	if data == nil {
		return errors.New("data cannot be null")
	}

	message.Data = *data
	return nil
}

func (message *Message[TData]) NewWithSameCorrelationId(data *TData, correlationId uuid.UUID) error {
	if data == nil {
		return errors.New("data cannot be null")
	}

	message.Data = *data
	message.CorrelationId = correlationId

	return nil
}
