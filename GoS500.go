package main

import (
	"bufio"
	"os"
	"github.com/tjCFeng/GoS500/S500"
	"github.com/tjCFeng/GoS500/Device"
)

func main() {
	defer S500.FreeS500()
	
	//echo "none" > /sys/class/leds/green:GPIOB12/trigger
	PB12, _ := S500.CreateGPIO(S500.PB, 12)
	defer S500.FreeGPIO(PB12)
	PB12.SetFun(S500.FunOUT)
	PB12.Flip()
	
	//echo "none" > /sys/class/leds/blue:GPIOB31/trigger
	PB31, _ := S500.CreateGPIO(S500.PB, 31)
	defer S500.FreeGPIO(PB31)
	PB31.SetFun(S500.FunOUT)
	PB31.Flip()
	
	PB31, _ := S500.CreateGPIO(S500.PB, 31)
	PWM3, _ := S500.CreatePWM(PB31, S500.PWM_3)
	defer S500.FreePWM(PWM3)
	PWM3.SetSRC(S500.PWM_HOSC_24M)
	PWM3.SetDIV(0xFF)
	PWM3.SetPolarity(true)
	PWM3.SetPeriod(1000)
	PWM3.SetDuty(500)
	
	Device.ISSD1306().Open()
	defer Device.FreeSSD1306()
	Device.ISSD1306().Writes(0, 0, []uint8("LeMaker Guitar"))
	Device.ISSD1306().Writes(0, 1, []uint8(" -- www.ICKey.cn"))
	
	reader := bufio.NewReader(os.Stdin)
	for {
		key, _, _:= reader.ReadLine()
		switch string(key)  {
			case "exit": return
			default: continue
		}
	}
}
