package stocks

import (
	"github.com/tealeg/xlsx"
	"log"
	"strconv"
)

func ParseCCCFile(file string) (stocks []Stock) {
	log.Printf("Using CCC file at %s\n", file)

	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		log.Println(err)
		return
	}

	var s [200]Stock

	log.Printf("\nGetting Info From Sheet %s\n", xlFile.Sheets[0].Name)

	count := 0
	for j, row := range xlFile.Sheets[0].Rows {
		if j < 6 {
			continue
		}

		s[count].Ticker = row.Cells[1].String()
		s[count].Name = row.Cells[0].String()
		if f, err := strconv.ParseFloat(row.Cells[8].String(), 32); err == nil {
			s[count].Dividend = float32(f)
		}
		count++
	}

	return s[:count-1]
}
