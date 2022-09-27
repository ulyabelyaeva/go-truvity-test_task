package main

import (
	"fmt"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	answers := []int{23, 34, 45, 7535327, 7784597, 271895, 0, 53533889, 242627852, 373006156, 272237365}
	for i := 0; i < len(answers); i++ {
		fmt.Printf("Testing: %v\n", i+1)
		before := time.Now()
		input(fmt.Sprintf("testcases\\testcase%v.txt", i+1))
		result := dinic()
		after := time.Now()
		if result != answers[i] {
			t.Errorf("Wrong answer on test %v, expected %v, found %v.", i+1, answers[i], result)
		}
		fmt.Println("Time:", after.Sub(before).Seconds())
	}
}
