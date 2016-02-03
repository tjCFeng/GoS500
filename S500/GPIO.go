/*
	Author: tjCFeng(LiuYang)
	EMail: tjCFeng@163.com
	[2016.02.02]
*/

package S500

import (
	
)

	const BaseGPIO = 0xB01B0000
	
	const (PA, PB, PC, PD, PE = 0, 1, 2, 3, 4)
	const (FunIN, FunOUT, Fun3, Fun4, Fun5, Fun6 = 0, 1, 2, 3, 4, 5)
	
	type PORT struct {
		Port	uint8
		hMem	[]uint8
		
		OUTEN	*uint32
		INEN	*uint32
		DAT		*uint32
		
		MFP		[4]*uint32
	}
	
func CreatePort(PORTx uint8) (*PORT, bool) {
	var Result bool = false
	
	if (PORTx > PE) { return nil, Result }
	
	port := &PORT{Port: PORTx}
	port.hMem, Result = IS500().GetMMap(BaseGPIO)
	if !Result { return nil, Result }
	
	Reg := (BaseGPIO & 0x00000FFF) + uint32(PORTx) * 0x0C
	port.OUTEN, Result = IS500().Register(port.hMem, Reg + 0x00)
	port.INEN, Result = IS500().Register(port.hMem, Reg + 0x04)
	port.DAT, Result = IS500().Register(port.hMem, Reg + 0x08)
	
	port.MFP[0], Result = IS500().Register(port.hMem, Reg + 0x40)
	port.MFP[1], Result = IS500().Register(port.hMem, Reg + 0x44)
	port.MFP[2], Result = IS500().Register(port.hMem, Reg + 0x48)
	port.MFP[3], Result = IS500().Register(port.hMem, Reg + 0x4C)
	
	return port, Result
}

func FreePort(port *PORT) {
	if (port.hMem != nil) { IS500().FreeMMap(port.hMem) }
}

/******************************************************************************/
	type GPIO struct {
		Port	*PORT
		Pin		uint8
		Bit		uint32
	}
	
func CreateGPIO(PORTx uint8, PINx uint8) (*GPIO, bool) {
	var Result bool = false
	
	switch (PORTx) {
		case PA:
		case PB:
		case PC:
		case PD:
		case PE:
	}
	
	gpio := &GPIO{}
	gpio.Port, Result = CreatePort(PORTx)
	if !Result { return nil, Result }
	
	gpio.Pin = PINx
	gpio.Bit = (0x1 << PINx)
	
	return gpio, Result
}

func FreeGPIO(gpio *GPIO) {
	FreePort(gpio.Port)
}

func (this *GPIO) SetFun(Fun uint8) {
	*this.Port.OUTEN &^= this.Bit
	*this.Port.INEN &^= this.Bit
	
	switch (Fun) {
		case FunIN: *this.Port.INEN |= this.Bit
		case FunOUT: *this.Port.OUTEN |= this.Bit
		case Fun3: 
		case Fun4:
		case Fun5:
		case Fun6:
	}
}

func (this *GPIO) SetData(Data bool) {
	switch (Data) {
		case true: *this.Port.DAT |= this.Bit
		case false: *this.Port.DAT &^= this.Bit
	}
}

func (this *GPIO) GetData() bool {
	return (*this.Port.DAT & this.Bit) == this.Bit
}

func (this *GPIO) Flip() {
	*this.Port.DAT ^= this.Bit
}

