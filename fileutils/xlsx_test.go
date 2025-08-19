package fileutils

import (
	"fmt"
	"testing"
)

func TestSheet(t *testing.T) {
	xlsx, err := OpenXlsx("xxx.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	defer xlsx.Close()

	xlsx.NewSheet("domain").
		NewSheet("ips").
		NewSheet("urls").
		SetActiveSheetByName("ips").
		DelSheet("Sheet1").
		DelSheet("ips")

	t.Log(xlsx.ActiveSheet())
	t.Log(xlsx.GetSheetList())
}

func TestCell(t *testing.T) {
	xlsx, err := OpenXlsx("xxx.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	defer xlsx.Close()

	xlsx.NewSheet("urls").DelSheet("Sheet1").
		NewSheet("domain")

	// 设置单个单元格的值
	xlsx.SetCell("A1", "hello world")
	xlsx.SetCell("B1", "hello world")
	xlsx.SetCell("A2", "hello world")
	xlsx.SetCell("A2", "hello world -- 2")
	xlsx.SetCellWithSheet("domain", "A2", "hello world")

	// 从指定位置处增加一行
	xlsx.SetRow("A5", []interface{}{
		"hello world",
		"hello world",
		"hello world",
		"hello world",
		"hello world",
	})

	// 末尾追加一行
	xlsx.AddRowWithSheet(xlsx.ActiveSheet(), []interface{}{
		"hello world",
		"hello world",
		123123,
	})

	// 追加一列值
	xlsx.SetCol("F1", []interface{}{
		"add-col",
		2, 3, 4, 5, 6, 7, 8, 9,
	})

	// 获取指定的值
	value, _ := xlsx.GetCell("A2")
	fmt.Println(value)

	// 获取指定行数
	values, _ := xlsx.GetRow(6)
	fmt.Println(values)

	// 获取指定列的值
	cols, _ := xlsx.GetCol("C")
	fmt.Println(cols)

	// 获取全部值
	rows, _ := xlsx.GetRows()
	fmt.Println(rows)
}
