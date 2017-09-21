package excelize

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// GetRows return all the rows in a sheet by given "sheet" + index. For now you
// should use sheet_name like "sheet3" where "sheet" is a constant part and "3"
// is a sheet number. For example, if sheet named as "SomeUniqueData" and it is
// second if spreadsheet program interface - you should use "sheet2" here. For
// example:
//
//    index := xlsx.GetSheetIndex("Sheet2")
//    rows := xlsx.GetRows("sheet" + strconv.Itoa(index))
//    for _, row := range rows {
//        for _, colCell := range row {
//            fmt.Print(colCell, "\t")
//        }
//        fmt.Println()
//    }
//
func (f *File) GetRows(sheet string) [][]string {
	xlsx := f.workSheetReader(sheet)
	rows := [][]string{}
	name := "xl/worksheets/" + strings.ToLower(sheet) + ".xml"
	if xlsx != nil {
		output, _ := xml.Marshal(f.Sheet[name])
		f.saveFileList(name, replaceWorkSheetsRelationshipsNameSpace(string(output)))
	}
	decoder := xml.NewDecoder(strings.NewReader(f.readXML(name)))
	d := f.sharedStringsReader()
	var inElement string
	var r xlsxRow
	var row []string
	tr, tc := f.getTotalRowsCols(sheet)
	for i := 0; i < tr; i++ {
		row = []string{}
		for j := 0; j <= tc; j++ {
			row = append(row, "")
		}
		rows = append(rows, row)
	}
	decoder = xml.NewDecoder(strings.NewReader(f.readXML(name)))
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.StartElement:
			inElement = startElement.Name.Local
			if inElement == "row" {
				r = xlsxRow{}
				decoder.DecodeElement(&r, &startElement)
				cr := r.R - 1
				for _, colCell := range r.C {
					c := TitleToNumber(strings.Map(letterOnlyMapF, colCell.R))
					val, _ := colCell.getValueFrom(f, d)
					rows[cr][c] = val
				}
			}
		default:
		}
	}
	return rows
}

// getTotalRowsCols provides a function to get total columns and rows in a
// sheet.
func (f *File) getTotalRowsCols(sheet string) (int, int) {
	name := "xl/worksheets/" + strings.ToLower(sheet) + ".xml"
	decoder := xml.NewDecoder(strings.NewReader(f.readXML(name)))
	var inElement string
	var r xlsxRow
	var tr, tc int
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.StartElement:
			inElement = startElement.Name.Local
			if inElement == "row" {
				r = xlsxRow{}
				decoder.DecodeElement(&r, &startElement)
				tr = r.R
				for _, colCell := range r.C {
					col := TitleToNumber(strings.Map(letterOnlyMapF, colCell.R))
					if col > tc {
						tc = col
					}
				}
			}
		default:
		}
	}
	return tr, tc
}

// SetRowHeight provides a function to set the height of a single row.
// For example:
//
//    xlsx := excelize.NewFile()
//    xlsx.SetRowHeight("Sheet1", 0, 50)
//    err := xlsx.Save()
//    if err != nil {
//        fmt.Println(err)
//        os.Exit(1)
//    }
//
func (f *File) SetRowHeight(sheet string, rowIndex int, height float64) {
	xlsx := f.workSheetReader(sheet)
	rows := rowIndex + 1
	cells := 0
	completeRow(xlsx, rows, cells)
	xlsx.SheetData.Row[rowIndex].Ht = height
	xlsx.SheetData.Row[rowIndex].CustomHeight = true
}

// getRowHeight provides function to get row height in pixels by given sheet
// name and row index.
func (f *File) getRowHeight(sheet string, row int) int {
	xlsx := f.workSheetReader(sheet)
	for _, v := range xlsx.SheetData.Row {
		if v.R == row+1 && v.Ht != 0 {
			return int(convertRowHeightToPixels(v.Ht))
		}
	}
	// Optimisation for when the row heights haven't changed.
	return int(defaultRowHeightPixels)
}

// GetRowHeight provides function to get row height by given worksheet name and
// row index.
func (f *File) GetRowHeight(sheet string, row int) float64 {
	xlsx := f.workSheetReader(sheet)
	for _, v := range xlsx.SheetData.Row {
		if v.R == row+1 && v.Ht != 0 {
			return v.Ht
		}
	}
	// Optimisation for when the row heights haven't changed.
	return defaultRowHeightPixels
}

// sharedStringsReader provides function to get the pointer to the structure
// after deserialization of xl/sharedStrings.xml.
func (f *File) sharedStringsReader() *xlsxSST {
	if f.SharedStrings == nil {
		var sharedStrings xlsxSST
		xml.Unmarshal([]byte(f.readXML("xl/sharedStrings.xml")), &sharedStrings)
		f.SharedStrings = &sharedStrings
	}
	return f.SharedStrings
}

// getValueFrom return a value from a column/row cell, this function is inteded
// to be used with for range on rows an argument with the xlsx opened file.
func (xlsx *xlsxC) getValueFrom(f *File, d *xlsxSST) (string, error) {
	switch xlsx.T {
	case "s":
		xlsxSI := 0
		xlsxSI, _ = strconv.Atoi(xlsx.V)
		if len(d.SI[xlsxSI].R) > 0 {
			value := ""
			for _, v := range d.SI[xlsxSI].R {
				value += v.T
			}
			return value, nil
		}
		return f.formattedValue(xlsx.S, d.SI[xlsxSI].T), nil
	case "str":
		return f.formattedValue(xlsx.S, xlsx.V), nil
	default:
		return f.formattedValue(xlsx.S, xlsx.V), nil
	}
}

// SetRowVisible provides a function to set visible of a single row by given
// worksheet index and row index. For example, hide row 3 in Sheet1:
//
//    xlsx.SetRowVisible("Sheet1", 2, false)
//
func (f *File) SetRowVisible(sheet string, rowIndex int, visible bool) {
	xlsx := f.workSheetReader(sheet)
	rows := rowIndex + 1
	cells := 0
	completeRow(xlsx, rows, cells)
	if visible {
		xlsx.SheetData.Row[rowIndex].Hidden = false
		return
	}
	xlsx.SheetData.Row[rowIndex].Hidden = true
}

// GetRowVisible provides a function to get visible of a single row by given
// worksheet index and row index. For example, get visible state of row 3 in
// Sheet1:
//
//    xlsx.GetRowVisible("Sheet1", 2)
//
func (f *File) GetRowVisible(sheet string, rowIndex int) bool {
	xlsx := f.workSheetReader(sheet)
	rows := rowIndex + 1
	cells := 0
	completeRow(xlsx, rows, cells)
	return !xlsx.SheetData.Row[rowIndex].Hidden
}
