package fileutils

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGetFilePerm(t *testing.T) {
	perm, err := GetFilePerm("/etc/passwd")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(perm)
}

func TestReadString(t *testing.T) {
	pd, err := ReadString("/etc/passwd")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pd)
}

func TestReadLines(t *testing.T) {
	lines, err := ReadLines("/etc/passwd")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(lines))
	t.Log(lines)
}

func TestReadByChan(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	outCh, err := ReadByChan(ctx, "xxx.log", 0, SeekStart)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for text := range outCh {
			fmt.Printf("%s", text)
		}
	}()

	// 10s 后退出
	time.Sleep(10 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
}

func TestWriteString(t *testing.T) {
	err := WriteString("./xx.log", "hello world --- 1\n")
	if err != nil {
		t.Fatal(err)
	}

	err = WriteString("./xx.log", "hello world --- 2\n")
	if err != nil {
		t.Fatal(err)
	}

	err = WriteString("./xx.log", "hello world --- 3\n")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAppendString(t *testing.T) {
	err := AppendString("./xx.log", "hello world --- 7\n")
	if err != nil {
		t.Fatal(err)
	}

	err = AppendString("./xx.log", "hello world --- 8\n")
	if err != nil {
		t.Fatal(err)
	}

	err = AppendString("./xx.log", "hello world --- 9\n")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAppendByChan(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	inCh := make(chan []byte)

	err := AppendByChan(ctx, "xxx.log", inCh)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		count := 0
		for {
			inCh <- []byte(fmt.Sprintf("write count -- %d\n", count))
			count++
			time.Sleep(500 * time.Millisecond)
		}
	}()

	time.Sleep(10 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
}
