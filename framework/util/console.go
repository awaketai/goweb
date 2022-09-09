package util

import "fmt"

// PrettyPrint 美观输出数组
func PrettyPrint(arr [][]string) {
	if len(arr) == 0 {
		return
	}
	rows := len(arr)
	columns := len(arr[0])
	lens := make([][]int, rows)
	// 每行第列的字符数
	for i := 0; i < rows; i++ {
		lens[i] = make([]int, columns)
		for j := 0; j < columns; j++ {
			lens[i][j] = len(arr[i][j])
		}
	}

	colMaxs := make([]int, columns)
	for j := 0; j < columns; j++ {
		for i := 0; i < rows; i++ {
			if colMaxs[j] < lens[i][j] {
				colMaxs[j] = lens[i][j]
			}
		}
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			fmt.Print(arr[i][j])
			padding := colMaxs[j] - lens[i][j] + 2
			for p := 0; p < padding; p++ {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}
