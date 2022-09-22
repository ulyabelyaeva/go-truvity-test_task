package main

import (
	"fmt"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	answers := []int{23, 34, 45, 7535327, 7784597, 271895, 0, 53533889, 242627852, 373006156, 272237365}
	for i := 0; i < len(answers); i++ {
		before := time.Now()
		input(fmt.Sprintf("testcases\\testcase%v.txt", i+1))
		result := dinic()
		after := time.Now()
		if result != answers[i] {
			t.Errorf("Wrong answer on test %v, expected %v, found %v.", i+1, answers[i], result)
		}
		// fmt.Println(after.Sub(before).Seconds())
		if after.Sub(before).Seconds() > 1.0 {
			t.Errorf("Time limit exceeded on test %v.", i+1)
		}
	}
}
