package top_test

import (
	"bytes"
	"io"
	"sync"
)

type syncedBuffer struct {
	b  bytes.Buffer
	mu sync.Mutex
}

func (b *syncedBuffer) Bytes() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Bytes()
}

func (b *syncedBuffer) Cap() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Cap()
}

func (b *syncedBuffer) Grow(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.b.Grow(n)
}

func (b *syncedBuffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Len()
}

func (b *syncedBuffer) Next(n int) []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Next(n)
}

func (b *syncedBuffer) Read(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Read(p)
}

func (b *syncedBuffer) ReadByte() (c byte, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.ReadByte()
}

func (b *syncedBuffer) ReadBytes(delim byte) (line []byte, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.ReadBytes(delim)
}

func (b *syncedBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.ReadFrom(r)
}

func (b *syncedBuffer) ReadRune() (r rune, size int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.ReadRune()
}

func (b *syncedBuffer) ReadString(delim byte) (line string, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.ReadString(delim)
}

func (b *syncedBuffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.b.Reset()
}

func (b *syncedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.String()
}

func (b *syncedBuffer) Truncate(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.b.Truncate(n)
}

func (b *syncedBuffer) UnreadByte() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.UnreadByte()
}

func (b *syncedBuffer) UnreadRune() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.UnreadRune()
}

func (b *syncedBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Write(p)
}

func (b *syncedBuffer) WriteByte(c byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.WriteByte(c)
}

func (b *syncedBuffer) WriteRune(r rune) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.WriteRune(r)
}

func (b *syncedBuffer) WriteString(s string) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.WriteString(s)
}

func (b *syncedBuffer) WriteTo(w io.Writer) (n int64, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.WriteTo(w)
}
