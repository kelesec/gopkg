package fileutils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
)

// Ref: https://github.com/qax-os/excelize?tab=readme-ov-file

type XlsxFile struct {
	filename    string
	activeSheet string
	file        *excelize.File
}

// OpenXlsx 打开/创建一个XLSX文件，如果文件不存在则会自动创建
// 新创建的文件默认有一个 Sheet1，如果需要删除的话，需要先创建一个新的sheet，再进行删除
func OpenXlsx(filename string) (*XlsxFile, error) {
	// 文件不存在，则创建新文件
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return &XlsxFile{
			filename:    filename,
			activeSheet: "Sheet1",
			file:        excelize.NewFile(),
		}, nil
	} else if err != nil {
		return nil, fmt.Errorf("fileutils::OpenXlsx open xlsx %s error: %w", filename, err)
	}

	// 文件已存在，则通过读写方式打开
	perm, _ := GetFilePerm(filename)
	file, err := os.OpenFile(filename, os.O_RDWR, perm)
	if err != nil {
		return nil, fmt.Errorf("fileutils::OpenXlsx open xlsx %s error: %w", filename, err)
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("fileutils::OpenXlsx create reader error: %w", err)
	}

	return &XlsxFile{filename: filename, activeSheet: f.GetSheetName(f.GetActiveSheetIndex()), file: f}, nil
}

// GetXlsxFile 获取原始的 excelize.File 对象，便于扩展别的操作
func (x *XlsxFile) GetXlsxFile() *excelize.File {
	return x.file
}

// NewSheet 创建信息 sheet
func (x *XlsxFile) NewSheet(sheet string) *XlsxFile {
	x.file.NewSheet(sheet)
	return x
}

// DelSheet 删除 sheet
func (x *XlsxFile) DelSheet(sheet string) *XlsxFile {
	x.file.DeleteSheet(sheet)
	if x.activeSheet == sheet {
		x.activeSheet = x.file.GetSheetName(x.file.GetActiveSheetIndex())
	}
	return x
}

// GetSheetList 获取全部 sheet 名称
func (x *XlsxFile) GetSheetList() []string {
	return x.file.GetSheetList()
}

// SetActiveSheet 将指定的 sheet 设置为活动状态
func (x *XlsxFile) SetActiveSheet(index int) *XlsxFile {
	x.file.SetActiveSheet(index)
	x.activeSheet = x.file.GetSheetName(index)
	return x
}

// SetActiveSheetByName 将指定的 sheet 设置为活动状态
func (x *XlsxFile) SetActiveSheetByName(sheet string) *XlsxFile {
	index, err := x.file.GetSheetIndex(sheet)
	if err != nil {
		return x
	}

	x.file.SetActiveSheet(index)
	x.activeSheet = sheet
	return x
}

// ActiveSheet 默认活动状态 sheet
func (x *XlsxFile) ActiveSheet() string {
	return x.activeSheet
}

// SetCellWithSheet 设置表格的值
// @example: SetCellWithSheet("Sheet1", "A1", "name")
func (x *XlsxFile) SetCellWithSheet(sheet, cell string, value interface{}) *XlsxFile {
	x.file.SetCellValue(sheet, cell, value)
	return x
}

// SetCell 设置表格的值，默认使用 active 状态的 sheet
// @example: SetCell("A1", "name")
func (x *XlsxFile) SetCell(cell string, value interface{}) *XlsxFile {
	return x.SetCellWithSheet(x.activeSheet, cell, value)
}

// SetRowWithSheet 给指定的 sheet 设置一行值
// @param startCell: 起始的单元格位置
// @example: SetRowWithSheet("Sheet1", "A5", []interface{"hello", 2, 3, 4, 5})
func (x *XlsxFile) SetRowWithSheet(sheet, startCell string, values []interface{}) *XlsxFile {
	x.file.SetSheetRow(sheet, startCell, &values)
	return x
}

// SetRow 给活动状态的sheet设置一行值
// @param startCell: 起始的单元格位置
// @example: SetRow("A5", []interface{"hello", 2, 3, 4, 5})
func (x *XlsxFile) SetRow(startCell string, values []interface{}) *XlsxFile {
	return x.SetRowWithSheet(x.activeSheet, startCell, values)
}

// SetColWithSheet 给指定的 sheet 设置一列的值
// @param startCell: 起始的单元格位置
// @example: SetColWithSheet("Sheet1", "D1", []interface{1, 2, 3, 4, 5})
func (x *XlsxFile) SetColWithSheet(sheet, startCell string, values []interface{}) *XlsxFile {
	x.file.SetSheetCol(sheet, startCell, &values)
	return x
}

// SetCol 给活动状态的 sheet 设置一列的值
// @param startCell: 起始的单元格位置
// @example: SetCol("A5", []interface{1, 2, 3, 4, 5})
func (x *XlsxFile) SetCol(startCell string, values []interface{}) *XlsxFile {
	return x.SetColWithSheet(x.activeSheet, startCell, values)
}

// AddRowWithSheet 给指定的sheet从末尾追加一行值
// @example: AddRowWithSheet("Sheet1", []interface{"hello", 2, 3, 4, 5})
func (x *XlsxFile) AddRowWithSheet(sheet string, values []interface{}) *XlsxFile {
	rows, err := x.file.GetRows(sheet)
	if err != nil {
		return x
	}

	// 从最后一行的下一行位置的第一列开始追加
	err = x.file.SetSheetRow(sheet, fmt.Sprintf("A%d", len(rows)+1), &values)
	return x
}

// AddRow 给活动状态的sheet从末尾追加一行值
// @example: AddRow([]interface{"hello", 2, 3, 4, 5})
func (x *XlsxFile) AddRow(values []interface{}) *XlsxFile {
	return x.AddRowWithSheet(x.activeSheet, values)
}

// GetCellWithSheet 获取指定sheet的单元格的值
// @example: GetCellWithSheet("Sheet1", "A2)
func (x *XlsxFile) GetCellWithSheet(sheet, cell string) (string, error) {
	return x.file.GetCellValue(sheet, cell)
}

// GetCell 获取活动状态sheet的单元格的值
// @example: GetCell("A2)
func (x *XlsxFile) GetCell(cell string) (string, error) {
	return x.file.GetCellValue(x.activeSheet, cell)
}

// GetRowWithSheet 获取指定sheet的某一行的全部值
// @example: GetRowWithSheet("Sheet1", 5)
func (x *XlsxFile) GetRowWithSheet(sheet string, row int) ([]string, error) {
	rows, err := x.file.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("fileutils::GetRowWithSheet get rows error: %w", err)
	}

	return rows[row], nil
}

// GetRow 获取活动状态的某一行的全部值
// @example: GetRow(5)
func (x *XlsxFile) GetRow(row int) ([]string, error) {
	return x.GetRowWithSheet(x.activeSheet, row)
}

// GetColWithSheet 获取指定sheet的某一行的全部值
// @example: GetColWithSheet("Sheet1", "F")
func (x *XlsxFile) GetColWithSheet(sheet string, col string) ([]string, error) {
	colNumber, err := excelize.ColumnNameToNumber(col)
	if err != nil {
		return nil, fmt.Errorf("fileutils::GetColWithSheet get column number error: %w", err)
	}

	cols, err := x.file.GetCols(sheet)
	if err != nil {
		return nil, fmt.Errorf("fileutils::GetColWithSheet get cols error: %w", err)
	}

	if colNumber > len(cols) {
		return nil, fmt.Errorf("fileutils::GetColWithSheet column number out of range: %d", colNumber)
	}

	return cols[colNumber-1], nil
}

// GetCol 获取活动状态的某一行的全部值
// @example: GetCol("F")
func (x *XlsxFile) GetCol(col string) ([]string, error) {
	return x.GetColWithSheet(x.activeSheet, col)
}

// GetRowsWithSheet 获取指定sheet的全部值
func (x *XlsxFile) GetRowsWithSheet(sheet string) ([][]string, error) {
	return x.file.GetRows(sheet)
}

// GetRows 获取活动状态的全部值
func (x *XlsxFile) GetRows() ([][]string, error) {
	return x.GetRowsWithSheet(x.activeSheet)
}

// Close 保存并关闭资源文件
func (x *XlsxFile) Close() {
	if x.file != nil {
		x.file.SaveAs(x.filename)
		x.file.Close()
		x.file = nil
	}
}
