package binsemaphore

// #include <linux/futex.h>
import "C"

import (
	"fmt"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

// S (semaphore) is just int32 value somewhere in (shared) memory
type S int32

// Post (signal) operation of semaphore
func (semaphore *S) Post() {
	if atomic.CompareAndSwapInt32((*int32)(semaphore), 0, 1) {
		_, _, err := syscall.Syscall6(
			syscall.SYS_FUTEX, uintptr(unsafe.Pointer(semaphore)), uintptr(C.FUTEX_WAKE), 1,
			0, 0, 0,
		)
		if err != 0 {
			panic(fmt.Sprintf("FUTEX_WAKE error: %v", err))
		}
	}
}

// Wait operation of semaphore
func (semaphore *S) Wait() {
	for {
		if atomic.CompareAndSwapInt32((*int32)(semaphore), 1, 0) {
			return
		}
		_, _, err := syscall.Syscall6(
			syscall.SYS_FUTEX, uintptr(unsafe.Pointer(semaphore)), uintptr(C.FUTEX_WAIT), 0,
			0, 0, 0,
		)
		if err == syscall.EAGAIN || err == syscall.EINTR {
			continue
		}
		if err != 0 {
			panic(fmt.Sprintf("FUTEX_WAIT error: %v", err))
		}
	}
}

// WaitWithTimeout is wait operation of semaphore with timeout
func (semaphore *S) WaitWithTimeout(timeout time.Duration) (timedout bool) {
	seconds := timeout.Truncate(time.Second)
	timespec := syscall.Timespec{int64(seconds.Seconds()), int64((timeout - seconds).Nanoseconds())}
	for {
		if atomic.CompareAndSwapInt32((*int32)(semaphore), 1, 0) {
			return
		}
		_, _, err := syscall.Syscall6(
			syscall.SYS_FUTEX, uintptr(unsafe.Pointer(semaphore)), uintptr(C.FUTEX_WAIT), 0,
			uintptr(unsafe.Pointer(&timespec)), 0, 0,
		)
		if err == syscall.ETIMEDOUT {
			return true
		}
		if err == syscall.EAGAIN || err == syscall.EINTR {
			continue
		}
		if err != 0 {
			panic(fmt.Sprintf("FUTEX_WAIT with timeout error: %v", err))
		}
	}
}
