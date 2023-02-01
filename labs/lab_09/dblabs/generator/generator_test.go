package generator

import (
	"flag"
	"fmt"
	"runtime"
	"testing"
)

var size int

func init() {
	flag.IntVar(&size, "size", 0, "Data Size")
	fmt.Printf("Default value runtime.NumCPU = %d\n", runtime.NumCPU())
}

func BenchmarkGenerateDataSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateDataSequential("data", size)
	}
}

func BenchmarkGenerateDataParallel1(b *testing.B) {
	runtime.GOMAXPROCS(1)
	for i := 0; i < b.N; i++ {
		_ = GenerateDataParallel("data", size)
	}
}

func BenchmarkGenerateDataParallel2(b *testing.B) {
	runtime.GOMAXPROCS(2)
	for i := 0; i < b.N; i++ {
		_ = GenerateDataParallel("data", size)
	}
}

func BenchmarkGenerateDataParallel4(b *testing.B) {
	runtime.GOMAXPROCS(4)
	for i := 0; i < b.N; i++ {
		_ = GenerateDataParallel("data", size)
	}
}

func BenchmarkGenerateDataParallel8(b *testing.B) {
	runtime.GOMAXPROCS(8)
	for i := 0; i < b.N; i++ {
		_ = GenerateDataParallel("data", size)
	}
}

func BenchmarkGenerateDataParallel16(b *testing.B) {
	runtime.GOMAXPROCS(16)
	for i := 0; i < b.N; i++ {
		_ = GenerateDataParallel("data", size)
	}
}

func BenchmarkGenerateDataParallel32(b *testing.B) {
	runtime.GOMAXPROCS(32)
	for i := 0; i < b.N; i++ {
		_ = GenerateDataParallel("data", size)
	}
}

func BenchmarkGenerateDataParallel64(b *testing.B) {
	runtime.GOMAXPROCS(64)
	for i := 0; i < b.N; i++ {
		_ = GenerateDataParallel("data", size)
	}
}

func BenchmarkGenerateDataParallel128(b *testing.B) {
	runtime.GOMAXPROCS(128)
	for i := 0; i < b.N; i++ {
		_ = GenerateDataParallel("data", size)
	}
}
