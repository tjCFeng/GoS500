/*
	Author: tjCFeng(LiuYang)
	EMail: tjCFeng@163.com
	[2016.02.09]
*/

package S500

import (
	
)

	const BaseSPI0 = 0xB0200000
	const BaseSPI1 = 0xB0204000
	const BaseSPI2 = 0xB0208000
	const BaseSPI3 = 0xB020C000
	
	const (SPI_0, SPI_1, SPI_2, SPI_3 = 0, 1, 2, 3)
	
	type SPI struct {
		mosi	*GPIO
		miso	*GPIO
		sclk	*GPIO
		cs		*GPIO
		spix	uint8
		hMem	[]uint8

		SPI_CTL	*uint32
		SPI_DIV	*uint32
		SPI_STAT	*uint32
		SPI_RXDAT	*uint32
		SPI_TXDAT	*uint32
		SPI_TCNT	*uint32
		SPI_SEED	*uint32
		SPI_TXCR	*uint32
		SPI_RXCR	*uint32
	}
	
	/*
	 * ETH_TXD0/GPIOA14		SPI2_SCLK
	 * ETH_TXD1/GPIOA15		SPI2_SS
	 * ETH_CRS_DIV/GPIOA18	SPI2_MISO
	 * ETH_REF_CLK/GPIOA21	SPI2_MOSI
	 * 
	 * ETH_TXEN/GPIOA16		SPI3_SCLK
	 * ETH_RXER/GPIOA17		SPI3_MOSI
	 * ETH_RXD1/GPIOA19		SPI3_SS
	 * ETH_RXD0/GPIOA20		SPI3_MISO
	 * 
	 * DSI_DP0/GPIOC6			SPI0_MISO
	 * DSI_DN0/GPIOC7			SPI0_MOSI
	 * DSI_DP2/GPIOC8			SPI0_SCLK
	 * DSI_DN2/GPIOC9			SPI0_SS
	 * 
	 * TWI3_SCLK/GPIOC22		SPI0_SCLK
	 * I2S_LRCLK1/GPIOC23	SPI0_SS
	 * I2S_MCLK1/GPIOC24		SPI0_MISO
	 * TWI3_SDATA/GPIOC25	SPI0_MOSI
	 * 
	 * UART0_RX/GPIOC26		SPI1_MISO
	 * UART0_TX/GPIOC27		SPI1_SS
	 * TWI0_SCLK/GPIOC28		SPI1_SCLK
	 * TWI0_SDATA/GPIOC29	SPI1_MOSI
	 * 
	 * NAND_ALE/GPIOD12		SPI2_MISO
	 * NAND_CLE/GPIOD13		SPI2_MOSI
	 * NAND_CE0B/GPIOD14		SPI2_SCLK
	 * NAND_CE1B/GPIOD15		SPI2_SS
	 */
	 
	 const STAT_TCOM = 4 //(0x1 << 2)
	 const SS_HL = 0x10 //(0x1 << 4)

func CreateSPI(MOSI *GPIO, MISO *GPIO, SCLK *GPIO, CS *GPIO, SPIx uint8) (*SPI, bool) {
	var Result bool = false
	
	if (SPIx > SPI_3) { return nil, Result }

	var BaseAddr uint32 = 0
	spi := &SPI{spix: SPIx}
	switch (SPIx) {
		case SPI_0: BaseAddr = BaseSPI0
		case SPI_1: BaseAddr = BaseSPI1
		case SPI_2: BaseAddr = BaseSPI2
		case SPI_3: BaseAddr = BaseSPI3
	}

	spi.hMem, Result = IS500().GetMMap(BaseAddr)
	if !Result { return nil, Result }
	Reg := uint32(BaseAddr & 0x00000FFF)
	spi.SPI_CTL, Result = IS500().Register(spi.hMem, Reg + 0x00)
	spi.SPI_DIV, Result = IS500().Register(spi.hMem, Reg + 0x04)
	spi.SPI_STAT, Result = IS500().Register(spi.hMem, Reg + 0x08)
	spi.SPI_RXDAT, Result = IS500().Register(spi.hMem, Reg + 0x0C)
	spi.SPI_TXDAT, Result = IS500().Register(spi.hMem, Reg + 0x10)
	spi.SPI_TCNT, Result = IS500().Register(spi.hMem, Reg + 0x14)
	spi.SPI_SEED, Result = IS500().Register(spi.hMem, Reg + 0x18)
	spi.SPI_TXCR, Result = IS500().Register(spi.hMem, Reg + 0x1C)
	spi.SPI_RXCR, Result = IS500().Register(spi.hMem, Reg + 0x20)
	
	if (MOSI != nil) && (MISO != nil) && (SCLK != nil) && (CS != nil) {
		spi.mosi = MOSI
		spi.miso = MISO
		spi.sclk = SCLK
		spi.cs = CS
		IMFP().SetSPI(MOSI, MISO, SCLK, CS, SPIx)
	}
	ICMU().SetSPICLK(SPIx, true)
	
	
	*spi.SPI_DIV = 0x3FF
	*spi.SPI_CTL &^= (0x3 << 8)
	*spi.SPI_CTL |= SS_HL
	*spi.SPI_STAT = 0xFFFFFFFF
	
	
	return spi, Result
}

func FreeSPI(spi *SPI) {
	if (spi.hMem != nil) { IS500().FreeMMap(spi.hMem) }
	if (spi.mosi != nil) { FreeGPIO(spi.mosi) }
	if (spi.miso != nil) { FreeGPIO(spi.miso) }
	if (spi.sclk != nil) { FreeGPIO(spi.sclk) }
	if (spi.cs != nil) { FreeGPIO(spi.cs) }
}

func (this *SPI) Open() {
	*this.SPI_CTL |= (0x1 << 18)
}

func (this *SPI) Close() {
	*this.SPI_CTL &^= (0x1 << 18)
}

func (this *SPI) ss(HL bool) {
	switch (HL) {
		case false: *this.SPI_CTL &^= SS_HL
		case true: *this.SPI_CTL |= SS_HL
	}
}

func (this *SPI) SetDIV(DIV uint16) {
	*this.SPI_DIV = uint32(DIV & 0x3FF)
}

func (this *SPI) waitReady() bool {
	Timeout := 0xFFFF
	
	for {
		if ((*this.SPI_STAT & STAT_TCOM) == STAT_TCOM) {
			*this.SPI_STAT |= STAT_TCOM
			return true
		} else {
			Timeout--
			if (Timeout == 0) { return false }
		}
	}
}

func (this *SPI) RW8Bit(Data []uint8) ([]uint8, uint32) {
	*this.SPI_CTL &^= (0x3 << 8) //8bit
	*this.SPI_STAT |= (0x1 << 5) + (0x1 << 4) //Reset FIFO

	Len := uint32(len(Data))
	Buf := make([]uint8, Len)
	Count := uint32(0)
	
	this.ss(false)
	for Count = 0; Count < Len; Count++ {
		*this.SPI_TXDAT = uint32(Data[Count])
		if !this.waitReady() { break }
		Buf[Count] = uint8(*this.SPI_RXDAT)
	}
	defer this.ss(true)
	
	return  Buf, Count
}
