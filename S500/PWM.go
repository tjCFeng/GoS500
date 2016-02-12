/*
	Author: tjCFeng(LiuYang)
	EMail: tjCFeng@163.com
	[2016.02.05]
*/

package S500

import (

)

	const (PWM_0, PWM_1, PWM_2, PWM_3, PWM_4, PWM_5 = 0, 1, 2, 3, 4, 5)
	const (PWM_IC_32K, PWM_HOSC_24M = 0, 1)
	
	type PWM struct {
		gpio	*GPIO;
		pwmx	uint8

		PWM_CTL	*uint32
		PWM_CLK	*uint32
	}
	
	/*
	GPIOA14	PWM4
	GPIOA15	PWM5
	GPIOA16	PWM0
	GPIOA17	PWM1
	GPIOA18	PWM4
	GPIOA19	PWM2
	GPIOA20	PWM3
	GPIOB03	PWM0/PWM4
	GPIOB04	PWM1/PWM5
	GPIOB05	PWM0
	GPIOB06	PWM1
	GPIOB07	PWM2
	GPIOB08	PWM3
	GPIOB09	PWM2
	GPIOB30	PWM2/PWM4
	GPIOB31	PWM3
	GPIOC31	PWM0
	GPIOD10	PWM1
	GPIOD16	PWM5
	GPIOD17	PWM4
	GPIOD28	PWM4
	GPIOD29	PWM5
	*/
	
func CreatePWM(gpio *GPIO, PWMx uint8) (*PWM, bool) {
	var Result bool = false
	
	if (PWMx > PWM_5) { return nil, Result }

	pwm := &PWM{pwmx: PWMx}
	switch (PWMx) {
		case PWM_0: fallthrough
		case PWM_1: fallthrough
		case PWM_2: fallthrough
		case PWM_3: pwm.PWM_CTL, Result = IS500().Register(gpio.Port.hMem, 0x50 + uint32(PWMx * 4))
		case PWM_4: fallthrough
		case PWM_5: pwm.PWM_CTL, Result = IS500().Register(gpio.Port.hMem, 0x68 + uint32(PWMx * 4))
	}
	if !Result {
		FreeGPIO(gpio)
		return nil, Result
	}
	pwm.gpio = gpio
	
	IMFP().SetPWM(gpio, PWMx)
	ICMU().SetPWMSRC(PWMx, PWM_HOSC_24M)
	return pwm, true
}

func FreePWM(pwm *PWM) {
	FreeGPIO(pwm.gpio)
}


func (this *PWM) SetPolarity(Polarity bool) {
	switch (Polarity) {
		case true: *this.PWM_CTL |= (0x1 << 20)
		case false: *this.PWM_CTL &^= (0x1 << 20)
	}
}

func (this *PWM) SetSRC(SRC uint8) {
	ICMU().SetPWMSRC(this.pwmx, SRC)
}

func (this *PWM) SetDIV(DIV uint16) {
	ICMU().SetPWMDIV(this.pwmx, DIV)
}

func (this *PWM) SetPeriod(Period uint16) {
	Val := *this.PWM_CTL
	Val &^= 0x3FF
	*this.PWM_CTL = Val + uint32(Period & 0x3FF)
}

func (this *PWM) SetDuty(Duty uint16) {
	Val := *this.PWM_CTL
	Val &^= (0x3FF << 10)
	*this.PWM_CTL = Val + uint32((Duty & 0x3FF) << 10)
}

func (this *PWM) Start() {
	ICMU().SetPWMCLK(this.pwmx, true)
}

func (this *PWM) Stop() {
	ICMU().SetPWMCLK(this.pwmx, false)
}
