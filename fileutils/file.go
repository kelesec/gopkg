package fileutils

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"io/fs"
	"os"
)

const (
	SeekStart   = io.SeekStart   // 	起始位置
	SeekCurrent = io.SeekCurrent // 当前位置
	SeekEnd     = io.SeekEnd     // 文件末尾
)

// GetFilePerm 获取文件的权限
func GetFilePerm(filename string) (fs.FileMode, error) {
	if info, err := os.Stat(filename); err == nil {
		return info.Mode().Perm(), nil
	} else {
		return 0644, fmt.Errorf("fileutils::GetFilePerm: %v", err)
	}
}

// Read 读取全部文件内容，并返回字节数组
func Read(filename string) ([]byte, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("fileutils::ReadFile: read file %s error: %v", filename, err)
	}

	return bytes, nil
}

// ReadString 读取全部文件内容，并返回 String 类型
func ReadString(filename string) (string, error) {
	bytes, err := Read(filename)
	if err != nil {
		return "", fmt.Errorf("fileutils::ReadString: %v", err)
	}

	return string(bytes), nil
}

// ReadLines 读取全部文件内容，并返回字符串数组（一行一个数组元素）
func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("fileutils::ReadLines: opening file %s error: %v", filename, err)
	}

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("fileutils::ReadLines: error reading file: %v", err)
	}
	return lines, nil
}

// ReadByChan 实时读取文件
// @param ctx: 上下文，监听退出文件实时读取信号
// @param filename: 待读取的文件名
// @param offset: 文件指针的相对偏移量
// @param whence: 文件指针开始计算的位置
// @return: 返回一个数组 chanel，用于获取实时文件内容
func ReadByChan(ctx context.Context, filename string, offset int64, whence int) (<-chan []byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("fileutils::ReadByChan: opening file %s error: %v", filename, err)
	}

	// 初始化 fsnotify 监听
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("fileutils::ReadByChan: failed to create watcher: %v", err)
	}

	if err := watcher.Add(filename); err != nil {
		return nil, fmt.Errorf("fileutils::ReadByChan: failed to watch file: %v", err)
	}

	// 通过 Chanel 传递文件内容
	outChan := make(chan []byte)

	go func() {
		// 退出时关闭资源
		defer func() {
			file.Close()
			watcher.Close()
			close(outChan)
		}()

		// 根据 start 将文件指针移动到指定位置
		if _, err := file.Seek(offset, whence); err != nil {
			panic(err)
			return
		}

		// 首次读取文件内容
		bytes, err := io.ReadAll(file)
		if err != nil {
			return
		}
		outChan <- bytes

		// 监听文件变化，实时读取新增内容
		reader := bufio.NewReader(file)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// 处理读取新的写入内容
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					// 如果文件被重新创建
					if event.Op&fsnotify.Create == fsnotify.Create {
						file, err = os.Open(filename)
						if err != nil {
							return
						}

						reader = bufio.NewReader(file)
					}

					// 读取新增内容
					for {
						line, err := reader.ReadBytes('\n')
						if err != nil && err == io.EOF {
							outChan <- line
							break
						}

						select {
						case outChan <- line:
						case <-ctx.Done():
							return
						}
					}
				}

			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	return outChan, nil
}

// Write 写入一个字节数组
func Write(filename string, data []byte) error {
	// 获取原文件的权限，如果没有则使用默认的 644 权限
	perm, _ := GetFilePerm(filename)
	return os.WriteFile(filename, data, perm)
}

// WriteString 写入一个字符串
func WriteString(filename, data string) error {
	return Write(filename, []byte(data))
}

// Append 从文件末尾追加写入一个字节数组
func Append(filename string, data []byte) error {
	perm, _ := GetFilePerm(filename)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, perm)
	if err != nil {
		return fmt.Errorf("fileutils::Append: opening file %s error: %v", filename, err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("fileutils::Append: writing file %s error: %v", filename, err)
	}

	return nil
}

// AppendString 从文件末尾追加写入一个字符串
func AppendString(filename, data string) error {
	return Append(filename, []byte(data))
}

// AppendByChan 通过接收 Chanel 的值写入文件
// @param ctx: 上下文，监听退出文件写入信号
// @param filename: 待读取的文件名
// @param inChan: 接收待写入文件的内容
func AppendByChan(ctx context.Context, filename string, inChan chan []byte) error {
	perm, _ := GetFilePerm(filename)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, perm)
	if err != nil {
		return fmt.Errorf("fileutils::AppendByChan: opening file %s error: %v", filename, err)
	}

	go func() {
		defer func() {
			file.Close()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case bytes, ok := <-inChan:
				if !ok {
					return
				}

				// 有数据接收时才写入
				if len(bytes) == 0 {
					continue
				}

				_, err := file.Write(bytes)
				if err != nil {
					return
				}
			}
		}

	}()

	return nil
}
