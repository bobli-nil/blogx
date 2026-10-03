package main

import (
	"bytes"
	"errors"
	"fmt"
)

// CloseableBuffer 是一个可关闭的内存缓冲区
type CloseableBuffer struct {
	buf    *bytes.Buffer
	closed bool
}

// NewCloseableBuffer 创建一个可关闭的缓冲区
func NewCloseableBuffer() *CloseableBuffer {
	return &CloseableBuffer{
		buf: &bytes.Buffer{},
	}
}

// Write 实现 io.Writer 接口
func (c *CloseableBuffer) Write(p []byte) (int, error) {
	if c.closed {
		return 0, errors.New("buffer is closed")
	}
	return c.buf.Write(p)
}

// Read 实现 io.Reader 接口
func (c *CloseableBuffer) Read(p []byte) (int, error) {
	if c.closed {
		return 0, errors.New("buffer is closed")
	}
	return c.buf.Read(p)
}

// Close 实现 io.Closer 接口
func (c *CloseableBuffer) Close() error {
	if c.closed {
		return nil // 幂等：重复关闭不报错
	}
	c.closed = true
	c.buf = nil // 释放对底层字节切片的引用，让 GC 回收
	return nil
}

// String 返回缓冲区内容（仅在未关闭时可用）
func (c *CloseableBuffer) String() string {
	if c.closed || c.buf == nil {
		return ""
	}
	return c.buf.String()
}

func main() {
	buf := NewCloseableBuffer()
	buf.Write([]byte("hello, world"))
	fmt.Println("关闭前:", buf.String()) // 输出: hello, world

	buf.Close()
	fmt.Println("关闭后:", buf.String()) // 输出: （空）

	_, err := buf.Write([]byte("test"))
	fmt.Println("关闭后写入:", err) // 输出: buffer is closed
}
