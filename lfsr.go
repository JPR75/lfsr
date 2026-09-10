// --------------------------------------------------------------------------------
// File        : lfsr.go
// Project     : LFSR
// Description : Calculate the LFSR for frequency division.
//               The main purpose is for implementation in an FPGA, based on
//               Xilinx application note XAPP052.
// --------------------------------------------------------------------------------
// Author      : JPR75 (https://github.com/JPR75/lfsr.git)
// --------------------------------------------------------------------------------
// Copyright (C) 2020 - 2026 JPR75
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>
// --------------------------------------------------------------------------------
package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
	"strings"
)

var lfsrTaps = [][]uint{
	{0}, {0}, {0}, {3, 2}, {4, 3}, {5, 3}, {6, 5}, {7, 6}, {8, 6, 5, 4}, {9, 5}, {10, 7}, {11, 9}, {12, 6, 4, 1}, {13, 4, 3, 1}, {14, 5, 3, 1}, {15, 14}, {16, 15, 13, 4}, {17, 14}, {18, 11}, {19, 6, 2, 1}, {20, 17}, {21, 19}, {22, 21}, {23, 18}, {24, 23, 22, 17}, {25, 22}, {26, 6, 2, 1}, {27, 5, 2, 1}, {28, 25}, {29, 27}, {30, 6, 4, 1}, {31, 28}, {32, 22, 2, 1}, {33, 20}, {34, 27, 2, 1}, {35, 33}, {36, 25}, {37, 5, 4, 3, 2, 1}, {38, 6, 5, 1}, {39, 35}, {40, 38, 21, 19}, {41, 38}, {42, 41, 20, 19}, {43, 42, 38, 37}, {44, 43, 18, 17}, {45, 44, 42, 41}, {46, 45, 26, 25}, {47, 42}, {48, 47, 21, 20}, {49, 40}, {50, 49, 24, 23}, {51, 50, 36, 35}, {52, 49}, {53, 52, 38, 37}, {54, 53, 18, 17}, {55, 31}, {56, 55, 35, 34}, {57, 50}, {58, 39}, {59, 58, 38, 37}, {60, 59}, {61, 60, 46, 45}, {62, 61, 6, 5}, {63, 62},
}

func mask(n int) uint64 {
	if n == 64 {
		return ^uint64(0)
	}
	return (uint64(1) << n) - 1
}

func main() {
	lfsr := uint64(0)
	lfsrFeedback := uint64(0)
	iteration := uint64(0)
	maxIteration := uint64(0xFFFFFFFFFFFFFFFE)
	var polynomial string
	
	// Check command-line argument
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <frequency> in Hz\n", os.Args[0])
		os.Exit(1)
	}

	// Convert argument from string to integer
	inputFrequency, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid argument: %s\n", os.Args[1])
		os.Exit(1)
	}

	// Check user input
	frequency := uint64(inputFrequency)
	n := int(math.Ceil(math.Log2(float64(frequency + 1))))
	//fmt.Printf("Polynomial size = %d\n", n)
	if n < 3 {
		fmt.Println("The divider ratio is too small — use counters or flip-flops instead!")
		os.Exit(1)
	}
	if n > 64 {
		fmt.Println("The divisor cannot exceed 64 bits (GO lang uint64 limitation).")
		os.Exit(1)
	}
	//fmt.Printf("mask = 0x%X\n", mask(n))

	// Starting iteration to find the target count
	fmt.Printf("Starting iteration ...\n")
	iterationDuration := time.Now()

	for iteration != frequency-1 {
		if len(lfsrTaps[n]) == 2 {
			lfsrFeedback = ^((lfsr >> (lfsrTaps[n][0] - 1)) ^ (lfsr >> (lfsrTaps[n][1] - 1))) & 1
			lfsr = ((lfsr << 1) | lfsrFeedback) & mask(n)
		} else {
			lfsrFeedback = ^((lfsr >> (lfsrTaps[n][0] - 1)) ^ (lfsr >> (lfsrTaps[n][1] - 1)) ^ (lfsr >> (lfsrTaps[n][2] - 1)) ^ (lfsr >> (lfsrTaps[n][3] - 1))) & 1
			lfsr = ((lfsr << 1) | lfsrFeedback) & mask(n)
		}
		iteration++

		if iteration > maxIteration {
			fmt.Println("Max iteration reached, target count not found!\n")
			fmt.Printf("Iteration duration: %v\n", time.Since(iterationDuration))
			os.Exit(1)
		}
	}
	fmt.Printf("Iteration duration: %v\n", time.Since(iterationDuration))

	// Results
	if len(lfsrTaps[n]) == 2 {
		fmt.Printf("\n=> P = X^%d + X^%d ; LFSR target count = 0x%X (%d)\n", lfsrTaps[n][0], lfsrTaps[n][1], lfsr, lfsr)
		polynomial = fmt.Sprintf("s_lfsr(%d) xnor s_lfsr(%d)", lfsrTaps[n][0] - 1, lfsrTaps[n][1] - 1)
	} else {
		fmt.Printf("\n=> P = X^%d + X^%d + X^%d + X^%d ; LFSR target count = 0x%X (%d)\n", lfsrTaps[n][0], lfsrTaps[n][1], lfsrTaps[n][2], lfsrTaps[n][3], lfsr, lfsr)
		polynomial = fmt.Sprintf("s_lfsr(%d) xnor s_lfsr(%d) xnor s_lfsr(%d) xnor s_lfsr(%d)", lfsrTaps[n][0] - 1, lfsrTaps[n][1] - 1, lfsrTaps[n][2] - 1, lfsrTaps[n][3] - 1)
	}

	// Write lfsr.vhdl file from the template file lfsr_template.vhdl
	err = writeLFSRFile(strconv.Itoa(n - 1), fmt.Sprintf("X\"%X\"", lfsr), polynomial)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nCreated lfsr.vhdl file\n")

}


// --------------------------------------------------------------------------------
// Read and replace the VHDL template
// --------------------------------------------------------------------------------
func writeLFSRFile(size, count, polynomial string) error {
	// Read the input file
	content, err := os.ReadFile("lfsr_template.vhdl")
	if err != nil {
		return fmt.Errorf("Error reading lfsr_template.vhdl: %w", err)
	}

	// Perform replacements
	result := string(content)
	result = strings.ReplaceAll(result, "{!size}", size)
	result = strings.ReplaceAll(result, "{!count}", count)
	result = strings.ReplaceAll(result, "{!polynomial}", polynomial)

	// Write the result to the output file
	err = os.WriteFile("lfsr.vhdl", []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("Error writing lfsr.vhdl: %w", err)
	}

	return nil
}

