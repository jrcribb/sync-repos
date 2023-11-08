package cui

import (
	"fmt"
	"strings"
)

// Funciones de actualización y posicionamiento de pantalla
const CLEAR_SCREEN = "\033[2J"
const BOLD = "\033[1m"
const UNDERLINE = "\033[4m"
const BLINK = "\033[5m"
const INVERTED = "\033[7m"
const ITALIC = "\033[3m"
const NORMAL = "\033[0m"
const CONSOLE_WIDTH = 200
const HIDE_CURSOR = "\033[?25l"
const SHOW_CURSOR = "\033[?25h"

// ┌───┬──┐
// │   │  │
// ├───┼──┤
// │   │  │
// └───┴──┘
const ESI = "┌"
const ESD = "┐"
const EII = "└"
const EID = "┘"
const LAT = "│"
const TRA = "─"
const INT = "┼"
const ILI = "├"
const ILD = "┤"
const ISU = "┬"
const IIN = "┴"

var lastRow int = 0
var lastCol int = 0

func Box(_x1, _y1, _x2, _y2 int) {
	ls := ESI + strings.Repeat(TRA, _y2-_y1-1) + ESD
	lm := LAT + strings.Repeat(" ", _y2-_y1-1) + LAT
	li := EII + strings.Repeat(TRA, _y2-_y1-1) + EID
	XyPrintf(_x1, _y1, 0, "%s", ls)
	for i := _x1 + 1; i <= _x2-1; i++ {
		XyPrintf(i, _y1, 0, "%s", lm)
	}
	XyPrintf(_x2, _y1, 0, "%s", li)
}
func ClearScreen() {
	// Borra la pantalla (simulación)
	fmt.Print(CLEAR_SCREEN)
}
func ClearLines(_initial, _end int) {
	// Borra lineas de la pantalla desde _initial hasta _end (desde columna 1 a CONSOLE_WIDTH)
	fmt.Print(HIDE_CURSOR) // oculta Cursor
	for i := _initial; i <= _end; i++ {
		XyPrintf(i, 1, CONSOLE_WIDTH, "")
	}
	fmt.Print(SHOW_CURSOR) // revela Cursor
}
func XyPrintfBold(_row int, _col int, _cleanSpaceReservation int, _format string, _val ...interface{}) {
	formattedString := fmt.Sprintf(_format, _val...)
	fmt.Printf(BOLD)
	XyPrintf(_row, _col, _cleanSpaceReservation, formattedString)
	fmt.Printf(NORMAL)
}
func XyPrintfUnderline(_row int, _col int, _cleanSpaceReservation int, _format string, _val ...interface{}) {
	formattedString := fmt.Sprintf(_format, _val...)
	fmt.Printf(UNDERLINE)
	XyPrintf(_row, _col, _cleanSpaceReservation, formattedString)
	fmt.Printf(NORMAL)
}
func XyPrintfInverted(_row int, _col int, _cleanSpaceReservation int, _format string, _val ...interface{}) {
	formattedString := fmt.Sprintf(_format, _val...)
	fmt.Printf(INVERTED)
	XyPrintf(_row, _col, _cleanSpaceReservation, formattedString)
	fmt.Printf(NORMAL)
}

func XyPrintf(_row int, _col int, _cleanSpaceReservation int, _format string, _val ...interface{}) {
	formattedString := fmt.Sprintf(_format, _val...)
	cleanSpaceReservation := strings.Repeat(" ", _cleanSpaceReservation)
	csrLength := len(cleanSpaceReservation)
	row := 0
	col := 0
	if _row == 0 {
		lastRow++
		row = lastRow
	} else {
		row = _row
		lastRow = _row
	}
	if _col == 0 {
		lastCol++
		col = lastCol
	} else {
		col = _col
		lastCol = _col
	}
	if csrLength == 0 {
		fmt.Printf("\033[%d;%dH%s", row, col, formattedString)
	} else {
		if len(formattedString) > csrLength {
			fmt.Printf("\033[%d;%dH%s", row, col, formattedString[0:csrLength])
		} else {
			fmt.Printf("\033[%d;%dH%s\033[%d;%dH%s", row, col, cleanSpaceReservation, row, col, formattedString)
		}
	}
}
