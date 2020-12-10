//Sale Order 数据整理

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// Data 数据内容
type Data struct {
	SON       string
	Customer  string
	Sodate    string
	Currency  string
	Shipdata  string
	Ponumber  string
	Itenumber string
	Orderqty  string
	Shipqty   string
	Unitprice string
}

func main() {
	fileName := "OEORLST1.csv"
	csvFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Could not loading the file, %s", err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	// test if read the data correctly
	fmt.Printf("Reader的信息,%T", reader)
	var data []Data
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(line)
		data = append(data, Data{
			SON:       line[12],
			Customer:  line[13],
			Sodate:    line[14],
			Currency:  line[17],
			Shipdata:  line[29],
			Ponumber:  line[45],
			Itenumber: line[110],
			Orderqty:  line[113],
			Shipqty:   line[114],
			Unitprice: line[117],
		})

	}
	newFile := "data.csv"
	checkFile := exists(newFile)
	fmt.Printf("checkfile value %v", checkFile)

	if checkFile {
		File, err := os.Create(newFile)
		if err != nil {
			fmt.Println(err)
		}
		defer File.Close()
		writer := csv.NewWriter(File)
		fmt.Println("Start....")

		for i := 0; i < len(data); i++ {
			var row []string
			row = append(row, data[i].SON)
			row = append(row, data[i].Customer)
			row = append(row, data[i].Sodate)
			row = append(row, data[i].Currency)
			row = append(row, data[i].Shipdata)
			row = append(row, data[i].Ponumber)
			row = append(row, data[i].Itenumber)
			row = append(row, data[i].Orderqty)
			row = append(row, data[i].Shipqty)
			row = append(row, data[i].Unitprice)
			writer.Write(row)

		}

		fmt.Println("Finished")
		writer.Flush()
	} else {
		fmt.Println("Already have file.")
	}

}

//判断是否已经存在原有数据。
func exists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsExist(err) {
			return false
		}
		return true
	}

	return true

}

