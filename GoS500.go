package main

import (
	"github.com/tjCFeng/GoS500/S500"
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
}
