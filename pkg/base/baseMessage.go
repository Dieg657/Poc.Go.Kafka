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
	Header Header
	data   TData
}

func (message *Message[TData]) New(data *TData) error {
	if data == nil {
		return errors.New("data cannot be null")
	}

	message.data = *data
	message.Header = Header{CorrelationId: uuid.New()}
	return nil
}

func (message *Message[TData]) GetData() *TData {
	return &message.data
}

type Header struct {
	CorrelationId uuid.UUID `json:"correlationId"`
}
