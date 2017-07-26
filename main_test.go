package main_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func getItem() string {
	str := strconv.Itoa(rand.Int())
	return "Hello there!" + time.Now().String() + str + "\r\n"
}

func BenchmarkRawConsole(b *testing.B) {
	var _ = os.Stdout.Sync()
	for i := 0; i < b.N; i++ {
		fmt.Print(getItem())
	}
	var _ = os.Stdout.Sync()
}

func BenchmarkConsole(b *testing.B) {
	var _ = os.Stdout.Sync()
	for i := 0; i < b.N; i++ {
		logger := log.New(os.Stdout, "", log.LstdFlags)
		logger.Print(getItem())
	}
	var _ = os.Stdout.Sync()
}

func BenchmarkIoDiscard(b *testing.B) {
	for i := 0; i < b.N; i++ {
		logger := log.New(ioutil.Discard, "", log.LstdFlags)
		logger.Print(getItem())
	}
}

func BenchmarkFile(b *testing.B) {
	file, _ := os.OpenFile("test_log.log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0)
	for i := 0; i < b.N; i++ {
		logger := log.New(file, "", log.LstdFlags)
		logger.Print(getItem())
	}
	file.Close()
}
