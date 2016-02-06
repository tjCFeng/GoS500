/*
	Author: tjCFeng(LiuYang)
	EMail: tjCFeng@163.com
	[2016.02.06]
*/

package S500

import (

)

	const BaseTWI0 = 0xB0170000
	const BaseTWI1 = 0xB0174000
	const BaseTWI2 = 0xB0178000
	const BaseTWI3 = 0xB017C000
	
	const (TWI_0, TWI_1, TWI_2, TWI_3 = 0, 1, 2, 3)
	const (TWI_Slave, TWI_Master = 0, 1)
	const (TWI_SpeedLow, TWI_SpeedHigh = 0, 1)
	
	/*
	SPI0_SS/GPIOC23		TWI3_SCL
	SPI0_MOSI/GPIOC25		TWI3_SDA
	
	UART0_RX/GPIOC26		TWI0_SDA
	UART0_TX/GPIOC27		TWI0_SCL
	
	TWI0_SCLK/GPIOC28		TWI0_SCL/TWI1_SCL
	TWI0_SDATA/GPIOC29	TWI0_SDA/TWI1_SDA
	
	PCM1_SYNC/GPIOD30		TWI3_SCL
	PCM1_OUT/GPIOD31		TWI3_SDA
	
	TWI1_SCLK/GPIOE0		TWI1_SCL
	TWI1_SDATA/GPIOE1		TWI1_SDA
	
	TWI2_SCLK/GPIOE2		TWI2_SCL
	TWI2_SDATA/GPIOE3		TWI2_SDA
	*/
	
	type TWI struct {
		sda		*GPIO;
		scl		*GPIO
		twix	uint8
		hMem	[]uint8

		TWI_CTL	*uint32
		TWI_DIV	*uint32
		TWI_STAT	*uint32
		TWI_ADDR	*uint32
		TWI_TXDAT	*uint32
		TWI_RXDAT	*uint32
		TWI_CMD	*uint32
		TWI_FIFOCTL	*uint32
		TWI_FIFOSTAT	*uint32
		TWI_DATCNT	*uint32
		TWI_RCNT	*uint32
	}
	
	
func CreateTWI(SDA *GPIO, SCL *GPIO, TWIx uint8) (*TWI, bool) {
	var Result bool = false
	
	if (TWIx > TWI_3) { return nil, Result }

	twi := &TWI{twix: TWIx}
	switch (TWIx) {
		case TWI_0: twi.hMem, Result = IS500().GetMMap(BaseTWI0)
		case TWI_1: twi.hMem, Result = IS500().GetMMap(BaseTWI1)
		case TWI_2: twi.hMem, Result = IS500().GetMMap(BaseTWI2)
		case TWI_3: twi.hMem, Result = IS500().GetMMap(BaseTWI3)
	}
	if !Result { return nil, Result }
	
	Reg := uint32(BaseGPIO & 0x00000FFF)
	twi.TWI_CTL, Result = IS500().Register(twi.hMem, Reg + 0x00)
	twi.TWI_DIV, Result = IS500().Register(twi.hMem, Reg + 0x04)
	twi.TWI_STAT, Result = IS500().Register(twi.hMem, Reg + 0x08)
	twi.TWI_ADDR, Result = IS500().Register(twi.hMem, Reg + 0x0C)
	twi.TWI_TXDAT, Result = IS500().Register(twi.hMem, Reg + 0x10)
	twi.TWI_RXDAT, Result = IS500().Register(twi.hMem, Reg + 0x14)
	twi.TWI_CMD, Result = IS500().Register(twi.hMem, Reg + 0x18)
	twi.TWI_FIFOCTL, Result = IS500().Register(twi.hMem, Reg + 0x1C)
	twi.TWI_FIFOSTAT, Result = IS500().Register(twi.hMem, Reg + 0x20)
	twi.TWI_DATCNT, Result = IS500().Register(twi.hMem, Reg + 0x24)
	twi.TWI_RCNT, Result = IS500().Register(twi.hMem, Reg + 0x28)
	
	twi.sda = SDA
	twi.scl = SCL
	IMFP().SetTWI(SDA, SCL, TWIx)
	ICMU().SetTWICLK(TWIx, true)
	
	*twi.TWI_CMD |= (0x1 << 8)
	*this.TWI_CTL |= (0x1 << 7)
	
	return twi, Result
}

func FreeTWI(twi *TWI) {
	if (twi.hMem != nil) { IS500().FreeMMap(twi.hMem) }
	if (twi.sda != nil) { FreeGPIO(twi.sda) }
	if (twi.scl != nil) { FreeGPIO(twi.scl) }
}

func (this *TWI) Start() {
	*this.TWI_CTL |= (0x1 << 7)
}

func (this *TWI) Stop() {
	*this.TWI_CTL &^= (0x1 << 7)
}

func (this *TWI) SetMode(Mode uint8) {
	
	switch (Mode) {
		case TWI_Slave: *this.TWI_CMD &^= (0x1 << 11)
		case TWI_Master: *this.TWI_CMD |= (0x1 << 11)
	}
}

func (this *TWI) SetSpeed(Speed uint8) {
	switch (Speed) {
		case TWI_SpeedLow: *this.TWI_CTL &^= (0x1 << 10)
		case TWI_SpeedHigh: *this.TWI_CTL |= (0x1 << 10)
		
	}
}

func (this *TWI) SetSlaveAddr(SlaveAddr uint8) {
	*this.TWI_ADDR = uint32(SlaveAddr)
}
