package writer

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/olekukonko/tablewriter"
)

const (
	CLEAR     = "\033[H\033[2J"
	FirstCell = "\033[0;0H"
)

var (
	RED = color.New(color.FgRed).SprintFunc()
)

type TableWriter struct {
	size  int
	data  [][]string
	table *tablewriter.Table
}

func New(size int) *TableWriter {
	t := TableWriter{size: size}
	t.prepare()
	return &t
}

func (t *TableWriter) AddItem(cell item) {
	x, y := cell.Positions()
	t.data[y][x] = RED("X")
}

func (t *TableWriter) prepareData() {
	data := make([][]string, t.size)
	for i := 0; i < t.size; i++ {
		data[i] = make([]string, t.size)
	}
	t.data = data
}

func (t *TableWriter) prepare() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetBorder(true)
	table.SetRowLine(true)
	t.table = table
	t.prepareData()
}

func (t *TableWriter) Draw() {
	t.table.AppendBulk(t.data)
	fmt.Print(CLEAR)
	fmt.Print(FirstCell)
	t.table.Render()
	t.prepareData()
	t.table.ClearRows()
}
