/*
	Author: tjCFeng(LiuYang)
	EMail: tjCFeng@163.com
	[2016.02.05]
*/

package S500

import (

)

	const BaseCMU = 0xB0160000

	var iCMU *CMU = nil

	type CMU struct {
		hMem	[]uint8
		
		PWM4CLK	*uint32
		PWM5CLK	*uint32
		PWM0CLK	*uint32
		PWM1CLK	*uint32
		PWM2CLK	*uint32
		PWM3CLK	*uint32
		DEVCLKEN	[2]*uint32
		DEVRST		[2]*uint32
	}

func ICMU() (*CMU) {
	if (iCMU == nil) {
		iCMU = &CMU{}
		iCMU.hMem, _ = IS500().GetMMap(BaseCMU)
	
		Reg := uint32(BaseCMU & 0x00000FFF)
		iCMU.PWM4CLK, _ = IS500().Register(iCMU.hMem, Reg + 0x68)
		iCMU.PWM5CLK, _ = IS500().Register(iCMU.hMem, Reg + 0x6C)
		iCMU.PWM0CLK, _ = IS500().Register(iCMU.hMem, Reg + 0x70)
		iCMU.PWM1CLK, _ = IS500().Register(iCMU.hMem, Reg + 0x74)
		iCMU.PWM2CLK, _ = IS500().Register(iCMU.hMem, Reg + 0x78)
		iCMU.PWM3CLK, _ = IS500().Register(iCMU.hMem, Reg + 0x7C)
		iCMU.DEVCLKEN[0], _ = IS500().Register(iCMU.hMem, Reg + 0xA0)
		iCMU.DEVCLKEN[1], _ = IS500().Register(iCMU.hMem, Reg + 0xA4)
		iCMU.DEVRST[0], _ = IS500().Register(iCMU.hMem, Reg + 0xA0)
		iCMU.DEVRST[1], _ = IS500().Register(iCMU.hMem, Reg + 0xAC)
	}

	return iCMU
}

func FreeCMU() {
	if iCMU != nil { IS500().FreeMMap(iCMU.hMem) }
}

func (this *CMU) SetPWMCLK(PWMx uint8, Enable bool) {
	switch (PWMx) {
		case PWM_3:
			switch (Enable) {
				case false: *this.DEVCLKEN[1] &^= (0x1 << 26)
				case true: *this.DEVCLKEN[1] |= (0x1 << 26)
			}
	}
}

func (this *CMU) SetPWMSRC(PWMx uint8, SRC uint8) {
	if SRC > PWM_HOSC_24M { return }
	switch (PWMx) {
		case PWM_3:
			switch (SRC) {
				case 0: *this.PWM3CLK &^= (0x1 << 12)
				case 1: *this.PWM3CLK |= uint32(SRC << 12)
			}
	}
}

func (this *CMU) SetPWMDIV(PWMx uint8, DIV uint16) {
	
	switch (PWMx) {
		case PWM_3: 
			SRC := *this.PWM3CLK & (0x1 << 12)
			*this.PWM3CLK = SRC + uint32(DIV)
	}
}

func (this *CMU) SetTWICLK(TWIx uint8, Enable bool) {
	switch (TWIx) {
		case TWI_0:
			switch (Enable) {
				case false: *this.DEVCLKEN[1] &^= (0x1 << 14)
				case true: *this.DEVCLKEN[1] |= (0x1 << 14) 
			}
		case TWI_2:
			switch (Enable) {
				case false: *this.DEVCLKEN[1] &^= (0x1 << 30)
				case true: *this.DEVCLKEN[1] |= (0x1 << 30)
			}
		case TWI_3:
			switch (Enable) {
				case false: *this.DEVCLKEN[1] &^= (0x1 << 31)
				case true: *this.DEVCLKEN[1] |= (0x1 << 31)
			}
	}
}

func (this *CMU) SetSPICLK(SPIx uint8, Enable bool) {
	*this.DEVRST[1] &^= (0x1 << (8 + SPIx))
	switch (Enable) {
		case false: *this.DEVCLKEN[1] &^= (0x1 << (10 + SPIx))
		case true:
			*this.DEVCLKEN[1] |= (0x1 << (10 + SPIx))
			*this.DEVRST[1] |= (0x1 << (8 + SPIx))
	}
}
