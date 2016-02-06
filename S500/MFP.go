/*
	Author: tjCFeng(LiuYang)
	EMail: tjCFeng@163.com
	[2016.02.05]
*/

package S500

import (
	
)

	const BaseMFP = BaseGPIO
	
	var iMFP *MFP = nil

	type MFP struct {
		hMem		[]uint8
		MFP_CTL	[4]*uint32
	}

func IMFP() (*MFP) {
	
	if (iMFP == nil) {
		iMFP = &MFP{}
		
		mfp := &MFP{}
		mfp.hMem, _ = IS500().GetMMap(BaseMFP)
	
		Reg := uint32(BaseMFP & 0x00000FFF)
		iMFP.MFP_CTL[0], _ = IS500().Register(mfp.hMem, Reg + 0x40)
		iMFP.MFP_CTL[1], _ = IS500().Register(mfp.hMem, Reg + 0x44)
		iMFP.MFP_CTL[2], _ = IS500().Register(mfp.hMem, Reg + 0x48)
		iMFP.MFP_CTL[3], _ = IS500().Register(mfp.hMem, Reg + 0x4C)
	}

	return iMFP
}

func FreeMFP() {
	if (iMFP != nil) { IS500().FreeMMap(iMFP.hMem) }
}

func (this *MFP) CloseGPIO(gpio *GPIO) {
	*gpio.Port.OUTEN &^= gpio.Bit
	*gpio.Port.INEN &^= gpio.Bit
}

func (this *MFP) SetGPIO(gpio *GPIO, Fun uint8) {
	if (Fun > FunOUT) { return }
	this.CloseGPIO(gpio)
	
	switch (Fun) {
		case FunIN: *gpio.Port.INEN |= gpio.Bit
		case FunOUT: *gpio.Port.OUTEN |= gpio.Bit
	}
}

func (this *MFP) SetPWM(gpio *GPIO, PWMx uint8) {
	if PWMx > PWM_5 { return }
	this.CloseGPIO(gpio)
	
	switch (PWMx) {
		case PWM_3:
			if (gpio.Port.Port == PB) && (gpio.Pin == 31) {
				*this.MFP_CTL[1] &^= (0x7 << 14)
				*this.MFP_CTL[1] |= (0x3 << 14)
			}
	}
}

func (this *MFP) SetTWI(sda *GPIO, scl *GPIO, TWIx uint8) {
	if TWIx > TWI_3 { return }
	this.CloseGPIO(sda)
	this.CloseGPIO(scl)
	
	switch (TWIx) {
		case TWI_2: 
	}
}
