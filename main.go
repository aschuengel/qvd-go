package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

func readSymbols(field QvdFieldHeader, buffer []byte) []string {
	subBuffer := buffer[field.Offset : field.Offset+field.Length]
	i := 0
	stringStart := 0
	var symbols []string
	file, _ := os.Create(field.Name + ".dat")
	_, _ = file.Write(subBuffer)
	_ = file.Close()
	fmt.Printf("Field %s, offset: %d, length: %d\n", field.Name, field.Offset, field.Length)
	for i < field.Length {
		currentByte := subBuffer[i]
		switch currentByte {
		case 0:
			// End of string
			symbol := string(subBuffer[stringStart:i])
			symbols = append(symbols, symbol)
			i += 1
			break
		case 1:
			// Integer
			intValue := decodeInt4(buffer[i+1 : i+5])
			symbols = append(symbols, fmt.Sprintf("%d", intValue))
			i += 5
		case 2:
			// Double
			symbols = append(symbols, "8-byte double")
			i += 9
		case 4:
			i += 1
			stringStart = i
		case 5:
			i += 5
			stringStart = i
		case 6:
			i += 5
			stringStart = i
		default:
			i += 1
		}
	}
	return symbols
}

func decodeInt4(bytes []byte) int {
	return int(bytes[0])<<24 + int(bytes[1])<<16 + int(bytes[2])<<8 + int(bytes[3])
}

func main() {

	// Open our xmlFile
	xmlFile, err := os.Open("Z_U1_PROJECTS_AC.qvd")
	// if we os.Open returns an error then handle it
	if err != nil {
		panic(err)
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer func(xmlFile *os.File) {
		_ = xmlFile.Close()
	}(xmlFile)

	data := readXml(xmlFile)

	var header QvdHeader
	err = xml.Unmarshal(data, &header)
	if err != nil {
		panic(err)
	}
	fmt.Println("Name", header.Name)
	fmt.Println("Comment", header.Comment)
	fmt.Println("Length", header.Length)
	fmt.Println("Offset", header.Offset)

	bytes, err := io.ReadAll(xmlFile)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(bytes))

	for i := 0; i < len(header.Fields); i++ {
		field := header.Fields[i]
		field.Symbols = readSymbols(field, bytes)
		fmt.Printf("Number of symbols for %s: %d\n", field.Name, len(field.Symbols))
		fmt.Printf("Number of symbols for %s: %d\n", field.Name, field.SymbolCount)
	}
}

func readXml(file *os.File) []byte {
	token := "</QvdTableHeader>"
	buffer := make([]byte, len(token))
	var i int64 = 0
	for {
		read, err := file.ReadAt(buffer, i)
		if err != nil {
			return nil
		}
		if read < len(buffer) {
			return nil
		}
		if string(buffer) == token {
			buffer := make([]byte, int64(len(token))+i+2)
			_, err := io.ReadFull(file, buffer)
			if err != nil {
				return nil
			}
			return buffer
		}
		i++
	}
}
