package table

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func Print(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetHeader(header)
	table.SetColWidth(100)

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}
