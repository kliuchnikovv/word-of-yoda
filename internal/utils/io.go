package utils

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"golang.org/x/net/context"
)

const headerSize = 4

type Reader[T any] struct {
	ctx context.Context
	r   *bufio.Reader

	ch  chan T
	err chan error
}

func NewReader[T any](ctx context.Context, r *bufio.Reader, capacity int) *Reader[T] {
	if capacity < 0 {
		capacity = 0
	}

	return &Reader[T]{
		ctx: ctx,
		r:   r,
		ch:  make(chan T, capacity),
		err: make(chan error, capacity),
	}
}

func (reader *Reader[T]) Start() (<-chan T, <-chan error) {
	go reader.read()

	return reader.ch, reader.err
}

func (reader *Reader[T]) read() {
	for {
		select {
		case <-reader.ctx.Done():
			return
		default:
			msg, err := ReadMessage[T](reader.r)
			if err != nil {
				reader.err <- err
			} else {
				reader.ch <- *msg
			}
		}
	}
}

// readMessage читает одно сообщение: сначала 4-байтовый заголовок длины, затем JSON-данные
func ReadMessage[T any](r *bufio.Reader) (*T, error) {
	// читаем заголовок длины
	header := make([]byte, headerSize)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, fmt.Errorf("can't read header: %w", err)
	}

	length := binary.BigEndian.Uint32(header) // длина тела
	if length == 0 {
		return nil, errors.New("message length is zero")
	}

	// читаем тело
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("can't read message: %w", err)
	}

	var msg T
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("can't unmarshal message: %w", err)
	}

	return &msg, nil
}

// writeMessage сериализует Message в JSON, посылает с 4-байтовым заголовком длины
func WriteMessage(w *bufio.Writer, msg any) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	var header = make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))

	if _, err := w.Write(header); err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	return w.Flush()
}
