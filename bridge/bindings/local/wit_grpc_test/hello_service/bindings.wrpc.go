// Generated by `wit-bindgen-wrpc-go` 0.11.0. DO NOT EDIT!
package hello_service

import (
	bytes "bytes"
	context "context"
	binary "encoding/binary"
	errors "errors"
	fmt "fmt"
	io "io"
	slog "log/slog"
	math "math"
	sync "sync"
	atomic "sync/atomic"
	utf8 "unicode/utf8"
	wrpc "wrpc.io/go"
)

// message HelloRequest
type HelloRequest struct {
	Message string
}

func (v *HelloRequest) String() string { return "HelloRequest" }

func (v *HelloRequest) WriteToIndex(w wrpc.ByteWriter) (func(wrpc.IndexWriter) error, error) {
	writes := make(map[uint32]func(wrpc.IndexWriter) error, 1)
	slog.Debug("writing field", "name", "message")
	write0, err := (func(wrpc.IndexWriter) error)(nil), func(v string, w io.Writer) (err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return fmt.Errorf("string byte length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing string byte length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return fmt.Errorf("failed to write string byte length of %d: %w", n, err)
		}
		slog.Debug("writing string bytes")
		_, err = w.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("failed to write string bytes: %w", err)
		}
		return nil
	}(v.Message, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `message` field: %w", err)
	}
	if write0 != nil {
		writes[0] = write0
	}

	if len(writes) > 0 {
		return func(w wrpc.IndexWriter) error {
			var wg sync.WaitGroup
			var wgErr atomic.Value
			for index, write := range writes {
				wg.Add(1)
				w, err := w.Index(index)
				if err != nil {
					return fmt.Errorf("failed to index nested record writer: %w", err)
				}
				write := write
				go func() {
					defer wg.Done()
					if err := write(w); err != nil {
						wgErr.Store(err)
					}
				}()
			}
			wg.Wait()
			err := wgErr.Load()
			if err == nil {
				return nil
			}
			return err.(error)
		}, nil
	}
	return nil, nil
}

// message HelloResponse
type HelloResponse struct {
	Message string
}

func (v *HelloResponse) String() string { return "HelloResponse" }

func (v *HelloResponse) WriteToIndex(w wrpc.ByteWriter) (func(wrpc.IndexWriter) error, error) {
	writes := make(map[uint32]func(wrpc.IndexWriter) error, 1)
	slog.Debug("writing field", "name", "message")
	write0, err := (func(wrpc.IndexWriter) error)(nil), func(v string, w io.Writer) (err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return fmt.Errorf("string byte length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing string byte length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return fmt.Errorf("failed to write string byte length of %d: %w", n, err)
		}
		slog.Debug("writing string bytes")
		_, err = w.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("failed to write string bytes: %w", err)
		}
		return nil
	}(v.Message, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `message` field: %w", err)
	}
	if write0 != nil {
		writes[0] = write0
	}

	if len(writes) > 0 {
		return func(w wrpc.IndexWriter) error {
			var wg sync.WaitGroup
			var wgErr atomic.Value
			for index, write := range writes {
				wg.Add(1)
				w, err := w.Index(index)
				if err != nil {
					return fmt.Errorf("failed to index nested record writer: %w", err)
				}
				write := write
				go func() {
					defer wg.Done()
					if err := write(w); err != nil {
						wgErr.Store(err)
					}
				}()
			}
			wg.Wait()
			err := wgErr.Load()
			if err == nil {
				return nil
			}
			return err.(error)
		}, nil
	}
	return nil, nil
}

// rpc Hello(HelloRequest) returns (HelloResponse) {};
func Hello(ctx__ context.Context, wrpc__ wrpc.Invoker, req *HelloRequest) (r0__ *HelloResponse, err__ error) {
	var buf__ bytes.Buffer
	write0__, err__ := (req).WriteToIndex(&buf__)
	if err__ != nil {
		err__ = fmt.Errorf("failed to write `req` parameter: %w", err__)
		return
	}
	if write0__ != nil {
		err__ = errors.New("unexpected deferred write for synchronous `req` parameter")
		return
	}
	var w__ wrpc.IndexWriteCloser
	var r__ wrpc.IndexReadCloser
	w__, r__, err__ = wrpc__.Invoke(ctx__, "local:wit-grpc-test/hello-service@0.1.0", "hello", buf__.Bytes())
	if err__ != nil {
		err__ = fmt.Errorf("failed to invoke `hello`: %w", err__)
		return
	}
	defer func() {
		if err := r__.Close(); err != nil {
			slog.ErrorContext(ctx__, "failed to close reader", "instance", "local:wit-grpc-test/hello-service@0.1.0", "name", "hello", "err", err)
		}
	}()
	if cErr__ := w__.Close(); cErr__ != nil {
		slog.DebugContext(ctx__, "failed to close outgoing stream", "instance", "local:wit-grpc-test/hello-service@0.1.0", "name", "hello", "err", cErr__)
	}
	r0__, err__ = func(r wrpc.IndexReadCloser, path ...uint32) (*HelloResponse, error) {
		v := &HelloResponse{}
		var err error
		slog.Debug("reading field", "name", "message")
		v.Message, err = func(r interface {
			io.ByteReader
			io.Reader
		}) (string, error) {
			var x uint32
			var s uint8
			for i := 0; i < 5; i++ {
				slog.Debug("reading string length byte", "i", i)
				b, err := r.ReadByte()
				if err != nil {
					if i > 0 && err == io.EOF {
						err = io.ErrUnexpectedEOF
					}
					return "", fmt.Errorf("failed to read string length byte: %w", err)
				}
				if s == 28 && b > 0x0f {
					return "", errors.New("string length overflows a 32-bit integer")
				}
				if b < 0x80 {
					x = x | uint32(b)<<s
					if x == 0 {
						return "", nil
					}
					buf := make([]byte, x)
					slog.Debug("reading string bytes", "len", x)
					_, err = r.Read(buf)
					if err != nil {
						return "", fmt.Errorf("failed to read string bytes: %w", err)
					}
					if !utf8.Valid(buf) {
						return string(buf), errors.New("string is not valid UTF-8")
					}
					return string(buf), nil
				}
				x |= uint32(b&0x7f) << s
				s += 7
			}
			return "", errors.New("string length overflows a 32-bit integer")
		}(r)
		if err != nil {
			return nil, fmt.Errorf("failed to read `message` field: %w", err)
		}
		return v, nil
	}(r__, []uint32{0}...)
	if err__ != nil {
		err__ = fmt.Errorf("failed to read result 0: %w", err__)
		return
	}
	return
}
