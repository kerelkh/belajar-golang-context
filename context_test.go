package belajargolangcontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

// membuat context kosong
func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

// context menganut ke konsep parent child, dimana setiap context itu immutable, sehingga jika menambah sesuatu ke context
// maka akan dibuat child context
// child context mewarisi semua context parent
func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "B")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
}

// ketika akses value gunakan Value()
// jika value tidak ditemukan maka akan dicari value tersebut ke parent diatas atasnya
// jika parent paling tinggi tidak ada maka yang hasilnya nil
func TestContextWithReturnValue(t *testing.T) {
	contextA := context.Background()
	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextB, "c", "C")
	contextD := context.WithValue(contextC, "d", "D")
	contextE := context.WithValue(contextC, "e", "E")
	contextF := context.WithValue(contextB, "f", "F")

	fmt.Println(contextE.Value("b"))
	fmt.Println(contextE.Value("c"))
	fmt.Println(contextE.Value("d"))
	fmt.Println(contextE.Value("e"))
	fmt.Println(contextE.Value("f"))
	fmt.Println(contextD.Value("b"))
	fmt.Println(contextF.Value("c"))
}

// context with cancel
// melakukan cancel jika terdapat signal cancel Done() yang dikirimkan di context
func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)
	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
			}
		}
	}()

	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	for data := range destination {
		fmt.Println("Counter", data)
		if data > 3 {
			break
		}
	}

	cancel()

	time.Sleep(5 * time.Second)
	fmt.Println(runtime.NumGoroutine())
}

// context with timeout
// membatalkan context jika melebihi waktu timeout
func CreateCounterSlowProccess(ctx context.Context) chan int {
	destination := make(chan int)
	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) //simulasi slow proccess
			}
		}
	}()

	return destination
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	destination := CreateCounterSlowProccess(ctx)
	for counter := range destination {
		fmt.Println("Counter", counter)
	}

	fmt.Println(runtime.NumGoroutine())

}

// context with deadline
// berbeda dengan timeout, kl timeout waktu dari sekarang + waktu
// sedangkan deadline seperti jam 12 siang hari ini
// context.WithDeadline(parent, time)
// bedanya di parameternya saja
func TestContextWithDeadline(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(10*time.Second))
	defer cancel()

	destination := CreateCounterSlowProccess(ctx)
	for counter := range destination {
		fmt.Println("Counter", counter)
	}

	fmt.Println(runtime.NumGoroutine())

}
