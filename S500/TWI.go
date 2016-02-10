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
	
	const I2C_FIFOSTAT_CECB = (0x1 << 0)
	const I2C_FIFOSTAT_RNB = (0x1 << 1)
	
	const I2C_CMD_SBE			= (0x1 << 0)
	const I2C_CMD_AS_MASK	= (0x7 << 1)
	const I2C_CMD_RBE			= (0x1 << 4)
	const I2C_CMD_SAS_MASK	= (0x7 << 5)
	const I2C_CMD_DE			= (0x1 << 8)
	const I2C_CMD_NS			= (0x1 << 9)
	const I2C_CMD_SE			= (0x1 << 10)
	const I2C_CMD_MSS			= (0x1 << 11)
	const I2C_CMD_WRS			= (0x1 << 12)
	const I2C_CMD_EXEC		= (0x1 << 15)
	const I2C_CMD_X = I2C_CMD_EXEC | I2C_CMD_MSS | I2C_CMD_SE | I2C_CMD_DE | I2C_CMD_SBE
	
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

	var BaseAddr uint32 = 0
	twi := &TWI{twix: TWIx}
	switch (TWIx) {
		case TWI_0: BaseAddr = BaseTWI0
		case TWI_1: BaseAddr = BaseTWI1
		case TWI_2: BaseAddr = BaseTWI2
		case TWI_3: BaseAddr = BaseTWI3
	}
	
	twi.hMem, Result = IS500().GetMMap(BaseAddr)
	if !Result { return nil, Result }
	
	Reg := uint32(BaseAddr & 0x00000FFF)
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

	*twi.TWI_DIV = 0x3F
	*twi.TWI_CTL |= 0x80
	*twi.TWI_FIFOCTL = 0x06
	*twi.TWI_DATCNT = 0
	
	return twi, Result
}

func FreeTWI(twi *TWI) {
	if (twi.hMem != nil) { IS500().FreeMMap(twi.hMem) }
	if (twi.sda != nil) { FreeGPIO(twi.sda) }
	if (twi.scl != nil) { FreeGPIO(twi.scl) }
}

func (this *TWI) SetMode(Mode uint8) {
	switch (Mode) {
		case TWI_Slave: *this.TWI_CMD &^= (0x1 << 11)
		case TWI_Master: *this.TWI_CMD |= (0x1 << 11)
	}
}

func (this *TWI) SetSpeed(Speed uint8) {
	switch (Speed) {
		case TWI_SpeedLow: *this.TWI_DIV = 0x3F //100K
		case TWI_SpeedHigh: *this.TWI_DIV = 0x10 //400K
	}
}

func (this *TWI) Write(Addr uint8, Reg uint8, Data []uint8) bool {
	var I uint32 = 0
	var Len uint32 = uint32(len(Data))
	
	*this.TWI_CTL |= 0x80
	*this.TWI_DATCNT = Len
	*this.TWI_TXDAT = uint32(Addr << 1)
	*this.TWI_TXDAT = uint32(Reg)
	for I = 0; I < Len; I++ { *this.TWI_TXDAT = uint32(Data[I]) }
	*this.TWI_CMD = I2C_CMD_X + (0x2 << 1)
	
	for I = 0; I < 0xFFFF; I++ {
		if (*this.TWI_FIFOSTAT & I2C_FIFOSTAT_RNB) == I2C_FIFOSTAT_RNB {
			*this.TWI_FIFOSTAT |= I2C_FIFOSTAT_RNB
			*this.TWI_FIFOCTL = 0x06
			for {
				if (*this.TWI_FIFOCTL & 0x06) != 0x06 { break }
			}
			return false
		}
		if (*this.TWI_FIFOSTAT & I2C_FIFOSTAT_CECB) == I2C_FIFOSTAT_CECB { return true }
	}
	
	return false
}

func (this *TWI) Read(Addr uint8, Reg uint8, DataLen uint8) ([]uint8, bool) {
	var I uint32 = 0
	var Len uint32 = uint32(DataLen)
	var Data []uint8 = make([]uint8, DataLen)
	
	*this.TWI_CTL |= 0x80
	*this.TWI_DATCNT = Len
	*this.TWI_TXDAT = uint32(Addr << 1)
	*this.TWI_TXDAT = uint32(Reg)
	*this.TWI_TXDAT = uint32((Addr << 1) | 1)
	*this.TWI_CMD = I2C_CMD_X + I2C_CMD_RBE + I2C_CMD_NS + (0x1 << 5) + (0x2 << 1)
		
	for I = 0; I < 0xFFFF; I++ {
		if (*this.TWI_FIFOSTAT & I2C_FIFOSTAT_RNB) == I2C_FIFOSTAT_RNB {
			*this.TWI_FIFOSTAT |= I2C_FIFOSTAT_RNB
			*this.TWI_FIFOCTL = 0x06
			for {
				if (*this.TWI_FIFOCTL & 0x06) != 0x06 { break }
			}
			return nil, false
		}
		if (*this.TWI_FIFOSTAT & I2C_FIFOSTAT_CECB) == I2C_FIFOSTAT_CECB {
			for I = 0; I < Len; I++ {
				Data[I] = uint8(*this.TWI_RXDAT)
			}
			return Data, true
		}
	}
	
	return nil, false
}
