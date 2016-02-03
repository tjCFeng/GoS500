/*
	Author: tjCFeng(LiuYang)
	EMail: tjCFeng@163.com
	[2016.02.02]
*/

package S500

import (
	"os"
	"unsafe"
	"syscall"
)

var iS500 *S500 = nil

func IS500() (*S500) {
	if (iS500 == nil) {
		iS500 = &S500{nil}
		iS500.hFile, _ = os.OpenFile("/dev/mem", os.O_RDWR, 0)
	}

	return iS500
}

func FreeS500() {
	if iS500 != nil { iS500.hFile.Close()	}
}

/*****************************************************************************/

type S500 struct {
	hFile		*os.File
	
	//epollFD		int
	//epollEvent	map[uint32]func()
}

func (this *S500) GetMMap(BaseAddr uint32) ([]uint8, bool) {
	if this.hFile == nil { return nil, false }

	mem, err := syscall.Mmap(int(this.hFile.Fd()), int64(BaseAddr & 0xFFFFF000), 
		os.Getpagesize(), syscall.PROT_READ | syscall.PROT_WRITE, syscall.MAP_SHARED)
		
	return mem, (err == nil)
}

func (this *S500) FreeMMap(hMem []uint8) {
	syscall.Munmap(hMem)
}

func (this *S500) Register(hMem []uint8, offset uint32) (*uint32, bool) {
	if hMem == nil { return nil, false }
	return  (*uint32)(unsafe.Pointer(&hMem[offset])), true
}

/*****************************************************************************/
/*func (this *S500) setupEpoll() {
	var err error
	this.epollEvent = make(map[uint32] chan bool)
	
	this.epollFD, err = syscall.EpollCreate1(0)
	if err != nil {
		fmt.Println("CreateEpoll Error")
		return
	}
	
	go func () {
		var epollEvents [this.epollEvent.Count]syscall.EpollEvent
		for {
			Len, err := syscall.EpollWait(this.epollFD, epollEvents[ : ], -1)
			if err != nil {
				if err == syscall.EAGAIN { continue }
				fmt.Println("EpollWait Error")
				return
			}
			
			for I := 0; I < Len; I++ {
				if X, ok := this.epollEvent[int(epollEvents[I].Fd)]; ok {
					X <- true
				}
			}
		}
	}()
}

func (this *S500) AddISR(RegAddr uint32, Interrupt chan bool) (error) {
	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN | (syscall.EPOLLET & 0xFFFFFFFF) | syscall.EPOLLPRI
	
	this.epollEvent[RegAddr] = Interrupt
	err := syscall.SetNonblock(Interrupt, true)
	if err != nil { return err }
	
	event.Fd = int32(RegAddr)
	err = syscall.EpollCtl(this.epollFD, syscall.EPOLL_CTL_ADD, RegAddr, &event)
	if err != nil { return err }
	
	return nil
}

func (this *S500) DelISR(RegAddr uint32) (error) {
	err := syscall.EpollCtl(this.epollFD, syscall.EPOLL_CTL_DEL, RegAddr, nil)
	if err != nil { return err }
	
	err = syscall.SetNonblock(RegAddr, false)
	if err != nil { return err }
	
	delete(this.epollEvent, RegAddr)
	return nil
}*/
