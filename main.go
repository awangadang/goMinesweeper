package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Cell represents an individual square which has a mine and a counter of the mines around it
type Cell struct {
	minesAround int
	hasMine     bool
	isRevealed  bool
	isFlagged   bool
}

// Grid holds the cells in the grid in a slice of slices.
// The cellArray is organized by rows.
type Grid struct {
	totalMines    int
	width, height int
	totalRevealed int
	cellArray     [][]Cell
}

// Prints the coordinates as well as the grid
func (grid Grid) print(revealAll bool) {
	fmt.Println("Revealed", grid.totalRevealed)
	output := ""
	for rowIndex, row := range grid.cellArray {
		rowStr := fmt.Sprint(rowIndex, " ")
		if len(rowStr) < 3 {
			rowStr += " "
		}
		for _, cell := range row {
			if grid.width > 10 {
				rowStr += " "
			}
			if !cell.isRevealed && !revealAll {
				if cell.isFlagged {
					rowStr = fmt.Sprint(rowStr, " ", "F")
				} else {
					rowStr = fmt.Sprint(rowStr, " ", "?")
				}
			} else if cell.hasMine {
				rowStr = fmt.Sprint(rowStr, " ", "*")
			} else {
				rowStr = fmt.Sprint(rowStr, " ", cell.minesAround)
			}
		}
		output = fmt.Sprint(output, rowStr, "\n")
	}
	// Horizontal coordinates
	widthCoordinates := "   " // two spaces to account for height coordinates
	for i := 0; i < grid.width; i++ {
		if grid.width > 10 {
			if i < 10 {
				widthCoordinates = fmt.Sprintf("%s  %d", widthCoordinates, i)
			} else {
				widthCoordinates = fmt.Sprintf("%s %d", widthCoordinates, i)
			}
		} else {
			widthCoordinates = fmt.Sprintf("%s %d", widthCoordinates, i)
		}
	}
	output += fmt.Sprint(widthCoordinates, "\n")
	fmt.Println(output)
}

// Difficulty indicates the difficulty level chosen by player
type Difficulty int

// Difficulty levels
const (
	Beginner Difficulty = iota
	Intermediate
	Advanced
)

// Coordinate locates a point on the 2D Grid
type Coordinate struct {
	x, y int
}

// Uses breadth first search to expand 0-cells and returns the coordinates of all connected cells that are also 0
func searchEmptyCells(grid Grid, coord Coordinate) (ret []Coordinate) {
	neighborsAt := func(grid Grid, c Coordinate) (ret []Coordinate) {
		validCoord := func(grid Grid, c Coordinate) bool {
			return c.x >= 0 && c.y >= 0 && c.x < grid.height && c.y < grid.width
		}
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}
				newCoord := Coordinate{c.x + dx, c.y + dy}
				if validCoord(grid, newCoord) {
					ret = append(ret, newCoord)
				}
			}
		}
		return
	}
	seenCoordinates := make(map[Coordinate]bool)
	seenCoordinates[coord] = true
	frontier := []Coordinate{coord}
	for len(frontier) > 0 {
		coordinate := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		for _, n := range neighborsAt(grid, coordinate) {
			cell := grid.cellArray[n.x][n.y]
			if _, ok := seenCoordinates[n]; !cell.isRevealed && !ok {
				seenCoordinates[n] = true
				if cell.minesAround == 0 {
					frontier = append(frontier, n)
				}
			}
		}
	}
	for k := range seenCoordinates {
		ret = append(ret, k)
	}
	return
}

func initializeGrid(diff Difficulty) (grid Grid) {
	switch diff {
	case Beginner:
		grid = Grid{width: 9, height: 9, totalMines: 10}
	case Intermediate:
		grid = Grid{width: 16, height: 16, totalMines: 40}
	case Advanced:
		grid = Grid{width: 24, height: 24, totalMines: 99}
	}
	cellArray := [][]Cell{}
	for i := 0; i < grid.height; i++ {
		row := make([]Cell, grid.width)
		cellArray = append(cellArray, row)
	}
	rand.Seed(time.Now().UnixNano())
	max := grid.width * grid.height
	mineMap := make(map[Coordinate]bool)
	for i := grid.totalMines; i > 0; {
		offset := rand.Intn(max)
		rowIndex, columnIndex := 0, 0
		if offset > grid.width-1 {
			quotient := offset / grid.width
			rem := offset - quotient*grid.width
			rowIndex = quotient
			columnIndex = rem
		} else {
			columnIndex = offset
		}
		// fmt.Println("coordinates", rowIndex, columnIndex)
		coordinate := Coordinate{rowIndex, columnIndex}
		if _, ok := mineMap[coordinate]; ok {
			continue
		}
		mineMap[coordinate] = true
		cellArray[rowIndex][columnIndex].hasMine = true
		// update minesAround for upper row around mine
		if rowIndex-1 >= 0 {
			cellArray[rowIndex-1][columnIndex].minesAround++
			if columnIndex-1 >= 0 {
				cellArray[rowIndex-1][columnIndex-1].minesAround++
			}
			if columnIndex+1 < grid.width {
				cellArray[rowIndex-1][columnIndex+1].minesAround++
			}
		}
		if columnIndex-1 >= 0 {
			cellArray[rowIndex][columnIndex-1].minesAround++
		}
		if columnIndex+1 < grid.width {
			cellArray[rowIndex][columnIndex+1].minesAround++
		}
		// update minesAround for lower row around mine
		if rowIndex+1 < grid.height {
			cellArray[rowIndex+1][columnIndex].minesAround++
			if columnIndex-1 >= 0 {
				cellArray[rowIndex+1][columnIndex-1].minesAround++
			}
			if columnIndex+1 < grid.width {
				cellArray[rowIndex+1][columnIndex+1].minesAround++
			}
		}
		i--
	}
	grid.cellArray = cellArray
	return
}

func playGame(grid Grid) {
	var lose bool
	scanner := bufio.NewScanner(os.Stdin)
	for !lose {
		grid.print(false)
		fmt.Print("Mark or reveal a coordinate Ex. \"R 1 2\", \"F 1 2\" : ")
		scanner.Scan()
		text := scanner.Text()
		var x, y int
		var reveal bool
		split := strings.Split(text, " ")
		if len(split) == 3 {
			if strings.ToLower(split[0]) == "r" {
				reveal = true
			} else if strings.ToLower(split[0]) == "f" {
				reveal = false
			} else {
				fmt.Println("Invalid input")
				continue
			}
			row, err := strconv.Atoi(split[1])
			if err != nil || !(row >= 0 && row < grid.height) {
				fmt.Println("Invalid input")
				continue
			}
			col, err := strconv.Atoi(split[2])
			if err != nil || !(col >= 0 && col < grid.width) {
				fmt.Println("Invalid input")
				continue
			}
			x, y = row, col
		} else {
			fmt.Println("Invalid input")
			continue
		}
		if grid.cellArray[x][y].isRevealed {
			fmt.Println("Cell already revealed.")
			continue
		}
		if reveal {
			if grid.cellArray[x][y].hasMine {
				grid.cellArray[x][y].isRevealed = true
				lose = true
				break
			}
			grid.cellArray[x][y].isRevealed = true
			grid.totalRevealed++
			if grid.cellArray[x][y].minesAround == 0 {
				coordToReveal := searchEmptyCells(grid, Coordinate{x, y})
				for _, c := range coordToReveal {
					grid.cellArray[c.x][c.y].isRevealed = true
					grid.totalRevealed++
				}
			}
			if grid.totalRevealed+grid.totalMines == grid.width*grid.height {
				break
			}
		} else { // flag
			if grid.cellArray[x][y].isFlagged {
				fmt.Println("Cell already flagged")
			} else {
				grid.cellArray[x][y].isFlagged = true
			}
		}
	}
	if lose {
		grid.print(true)
		fmt.Println("Game Over")
	} else {
		grid.print(true)
		fmt.Println("Good Game!")
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Choose difficulty level (1, 2, 3): ")
	if scanner.Scan() {
	}
	text := scanner.Text()
	number, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println("Invalid entry")
		return
	}
	difficulty := Difficulty(number - 1)
	if difficulty < Beginner || difficulty > Advanced {
		fmt.Println("Invalid entry")
		return
	}
	grid := initializeGrid(Difficulty(difficulty))
	playGame(grid)
}
