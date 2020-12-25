package main

import (
	"fmt"
	"time"
)

type Scanner struct {
	N      int
	end    bool
	done   chan bool
	cancel chan bool
	msg    chan string
}

func CreateScanner() Scanner {
	done := make(chan bool, 1)
	cancel := make(chan bool, 1)
	msg := make(chan string, 5)
	s := Scanner{
		done:   done,
		cancel: cancel,
		msg:    msg,
	}
	return s
}

func (s *Scanner) Run() {
	s.end = false
	for i := 0; i < 10 && !s.end; i++ {

		select {
		case _ = <-s.cancel:
			println("Canceling")
			s.end = true
		default:
			s.N = i
			time.Sleep(1 * time.Second)
			s.msg <- fmt.Sprintf("VALUE %d", i)
		}
	}
	s.done <- true
}

func SCmain() {
	s := CreateScanner()
	go s.Run()

	end := false
	for !end {
		select {
		case _ = <-s.done:
			end = true
		case msg := <-s.msg:
			println(msg)
		default:
		}
	}

	//time.Sleep(5 * time.Second)
	//cancel <- true

	println("Done!")
}
