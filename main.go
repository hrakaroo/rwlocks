package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// queue to show what has a lock and what is waiting on a lock
type queue struct {
	items []string
	mutex sync.Mutex
}

func (q *queue) add(s string) {
	q.mutex.Lock()
	q.items = append(q.items, s)
	q.mutex.Unlock()
}

func (q *queue) remove(s string) {
	q.mutex.Lock()
	var cp []string
	for _, item := range q.items {
		if item != s {
			cp = append(cp, item)
		}
	}
	q.items = cp
	q.mutex.Unlock()
}

func (q *queue) toString() string {
	q.mutex.Lock()
	s := strings.Join(q.items, ", ")
	q.mutex.Unlock()
	return s
}

type rwLock struct {
	wait queue
	held queue

	sync.RWMutex
}

func (r *rwLock) RLock(name string) {
	// Add it onto the waiting list
	r.wait.add(name)

	// Request the real lock
	r.RWMutex.RLock()

	// Add it to the holding list and remove it from the waiting list
	r.held.add(name)
	r.wait.remove(name)
}

func (r *rwLock) RUnlock(name string) {
	// Release the mutex
	r.RWMutex.RUnlock()

	// Remove it from the held list
	r.held.remove(name)
}

func (r *rwLock) Lock(name string) {
	// Add it onto the waiting list
	r.wait.add(name)

	// Request the real lock
	r.RWMutex.Lock()

	// Add it to the holding list and remove it from the waiting list
	r.held.add(name)
	r.wait.remove(name)
}

func (r *rwLock) Unlock(name string) {
	// Release the mutex
	r.RWMutex.Unlock()

	// Remove from held
	r.held.remove(name)
}

func (r *rwLock) holding() string {
	return r.held.toString()
}

func (r *rwLock) waiting() string {
	return r.wait.toString()
}

var reader = bufio.NewReader(os.Stdin)

func read(name string, mutex *rwLock, releaseCh <-chan bool) {
	fmt.Printf("\t%s -> requesting read lock\n", name)

	mutex.RLock(name)

	// Sleep just a moment so release messages show up first
	time.Sleep(1 * time.Second)
	fmt.Printf("\t%s -> obtained read lock\n", name)

	<-releaseCh

	mutex.RUnlock(name)
	fmt.Printf("\t%s -> released read lock\n", name)
}

func write(name string, mutex *rwLock, releaseCh <-chan bool) {
	fmt.Printf("\t%s -> requesting write lock\n", name)

	mutex.Lock(name)

	// Sleep just a moment so release messages show up first
	time.Sleep(1 * time.Second)
	fmt.Printf("\t%s -> obtained write lock\n", name)

	<-releaseCh

	mutex.Unlock(name)
	fmt.Printf("\t%s -> released write lock\n", name)
}

func prompt(s string, mutex *rwLock) {
	time.Sleep(2 * time.Second)
	fmt.Println()

	fmt.Printf("holding: [%s]\n", mutex.holding())

	fmt.Printf("waiting: [%s]\n", mutex.waiting())

	fmt.Printf("%s <return> ", s)
	reader.ReadString('\n')
}

func ex5() {
	println("Running example 5")
	mutex := &rwLock{}

	read1Chan := make(chan bool)
	read2Chan := make(chan bool)
	read3Chan := make(chan bool)
	read4Chan := make(chan bool)
	read5Chan := make(chan bool)
	read6Chan := make(chan bool)

	write1Chan := make(chan bool)
	write2Chan := make(chan bool)

	prompt("read1 request lock", mutex)
	go read("read1", mutex, read1Chan)

	prompt("read2 request lock", mutex)
	go read("read2", mutex, read2Chan)

	prompt("write1 request lock", mutex)
	go write("write1", mutex, write1Chan)

	prompt("read3 request lock", mutex)
	go read("read3", mutex, read3Chan)

	prompt("write2 request lock", mutex)
	go write("write2", mutex, write2Chan)

	prompt("read4 request lock", mutex)
	go read("read4", mutex, read4Chan)

	prompt("read1 release", mutex)
	read1Chan <- true

	prompt("read2 release", mutex)
	read2Chan <- true

	prompt("read5 request lock", mutex)
	go read("read5", mutex, read5Chan)

	prompt("write1 release", mutex)
	write1Chan <- true

	prompt("read6 request lock", mutex)
	go read("read6", mutex, read6Chan)

	prompt("read3 release", mutex)
	read3Chan <- true

	prompt("read4 release", mutex)
	read4Chan <- true

	prompt("read5 release", mutex)
	read5Chan <- true

	prompt("write2 release", mutex)
	write2Chan <- true

	prompt("read6 release", mutex)
	read6Chan <- true
}

func ex4() {
	println("Running example 4")
	mutex := &rwLock{}

	read1Chan := make(chan bool)
	read2Chan := make(chan bool)
	read3Chan := make(chan bool)

	write1Chan := make(chan bool)

	prompt("read1 request lock", mutex)
	go read("read1", mutex, read1Chan)

	prompt("read2 request lock", mutex)
	go read("read2", mutex, read2Chan)

	prompt("write1 request lock", mutex)
	go write("write1", mutex, write1Chan)

	prompt("read3 request lock", mutex)
	go read("read3", mutex, read3Chan)

	prompt("read1 release", mutex)
	read1Chan <- true

	prompt("read2 release", mutex)
	read2Chan <- true

	prompt("write1 release", mutex)
	write1Chan <- true

	prompt("read3 release", mutex)
	read3Chan <- true
}

func ex3() {
	println("Running example 3")
	mutex := &rwLock{}

	read1Chan := make(chan bool)
	read2Chan := make(chan bool)

	write1Chan := make(chan bool)

	prompt("read1 request lock", mutex)
	go read("read1", mutex, read1Chan)

	prompt("read2 request lock", mutex)
	go read("read2", mutex, read2Chan)

	prompt("write1 request lock", mutex)
	go write("write1", mutex, write1Chan)

	prompt("read1 release", mutex)
	read1Chan <- true

	prompt("read2 release", mutex)
	read2Chan <- true

	prompt("write1 release", mutex)
	write1Chan <- true
}

func ex2() {
	println("Running example 2")
	mutex := &rwLock{}

	read1Chan := make(chan bool)
	read2Chan := make(chan bool)

	prompt("read1 request lock", mutex)
	go read("read1", mutex, read1Chan)

	prompt("read2 request lock", mutex)
	go read("read2", mutex, read2Chan)

	prompt("read1 release", mutex)
	read1Chan <- true

	prompt("read2 release", mutex)
	read2Chan <- true
}

func ex1() {
	println("Running example 1")
	mutex := &rwLock{}

	read1Chan := make(chan bool)

	prompt("read1 request lock", mutex)
	go read("read1", mutex, read1Chan)

	prompt("read1 release", mutex)
	read1Chan <- true
}

func main() {
	println()

	//ex1()
	//ex2()
	//ex3()
	//ex4()
	//ex5()

	time.Sleep(1 * time.Second)
}
