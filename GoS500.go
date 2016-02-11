package main

import (
	"fmt"
	"bufio"
	"os"
	"time"
	"github.com/tjCFeng/GoS500/S500"
	"github.com/tjCFeng/GoS500/Device"
)

func main() {
	defer S500.FreeS500()
	
	/*GPIO*/
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
	
	/*PWM*/
	PB31, _ := S500.CreateGPIO(S500.PB, 31)
	PWM3, _ := S500.CreatePWM(PB31, S500.PWM_3)
	defer S500.FreePWM(PWM3)
	PWM3.SetSRC(S500.PWM_HOSC_24M)
	PWM3.SetDIV(0xFF)
	PWM3.SetPolarity(true)
	PWM3.SetPeriod(1000)
	PWM3.SetDuty(500)
	
	/*SPI*/
	Device.ISSD1306().Open()
	defer Device.FreeSSD1306()
	Device.ISSD1306().Writes(0, 0, []uint8("LeMaker Guitar"))
	Device.ISSD1306().Writes(0, 1, []uint8(" -- www.ICKey.cn"))
	
	/*TWI*/
	SDA, _ :=S500.CreateGPIO(S500.PE, 3)
	SCL, _ :=S500.CreateGPIO(S500.PE, 2)
	TWI2, _ := S500.CreateTWI(SDA, SCL, S500.TWI_2)
	defer S500.FreeTWI(TWI2)
	TWI2.SetSpeed(S500.TWI_SpeedLow)

	HMC5883ADDR := uint8(0x1E)
	TWI2.Write(HMC5883ADDR, 0x00, []uint8{0x14})
	TWI2.Write(HMC5883ADDR, 0x01, []uint8{0x20})
	TWI2.Write(HMC5883ADDR, 0x02, []uint8{0x00})

	go func () {
		for {
			Data, ok := TWI2.Read(HMC5883ADDR, 0x03, 6)
			fmt.Println(Data, ok)
			time.Sleep(1 * time.Second)
		}
	}()
	
	
	//fmt.Println(Device.IDS3231().Write(16, 2, 11, 20, 24, 30, 4))
	fmt.Println(Device.IDS3231().Read())
	fmt.Println(Device.IDS3231().Temperature())
	
	
	reader := bufio.NewReader(os.Stdin)
	for {
		key, _, _:= reader.ReadLine()
		switch string(key)  {
			case "exit": return
			default: continue
		}
	}
}
