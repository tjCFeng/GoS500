package Device

import (
	"../S500"
)

	const DS3231_ADDR = 0x68

	/* RTC register addresses */
	const RTC_SEC_REG_ADDR	= 0x0
	const RTC_MIN_REG_ADDR	= 0x1
	const RTC_HR_REG_ADDR	= 0x2
	const RTC_DAY_REG_ADDR	= 0x3
	const RTC_DATE_REG_ADDR	= 0x4
	const RTC_MON_REG_ADDR	= 0x5
	const RTC_YR_REG_ADDR	= 0x6
	const RTC_CTL_REG_ADDR	= 0x0E
	const RTC_STAT_REG_ADDR	= 0x0F


	/* RTC control register bits */
	const RTC_CTL_BIT_A1IE	= 0x1 // Alarm 1 interrupt enable
	const RTC_CTL_BIT_A2IE	= 0x2 // Alarm 2 interrupt enable
	const RTC_CTL_BIT_INTCN	= 0x4 // Interrupt control
	const RTC_CTL_BIT_RS1	= 0x8 // Rate select 1
	const RTC_CTL_BIT_RS2	= 0x10 // Rate select 2
	const RTC_CTL_BIT_DOSC	= 0x80 // Disable Oscillator

	/* RTC status register bits */
	const RTC_STAT_BIT_A1F	= 0x1 // Alarm 1 flag
	const RTC_STAT_BIT_A2F	= 0x2 // Alarm 2 flag
	const RTC_STAT_BIT_OSF	= 0x80 // Oscillator stop flag

	var iDS3231 *DS3231 = nil

func IDS3231() (*DS3231) {
	if (iDS3231 == nil) {
		iDS3231 = &DS3231{}
		iDS3231.sda, _ = S500.CreateGPIO(S500.PE, 3)
		iDS3231.scl, _ = S500.CreateGPIO(S500.PE, 2)
		iDS3231.twi, _ = S500.CreateTWI(iDS3231.sda, iDS3231.scl, S500.TWI_2)
		iDS3231.twi.SetSpeed(S500.TWI_SpeedLow)
	}

	return iDS3231
}

func FreeDS3231() {
	S500.FreeTWI(iDS3231.twi)
}

	type DS3231 struct {
		twi		*S500.TWI
		sda		*S500.GPIO
		scl		*S500.GPIO
	}
	
func (this *DS3231) bcd2bin(bcd uint8) uint8 {
	return ((bcd & 0x0F) + (bcd >> 4) * 10)
}

func (this *DS3231) bin2bcd(bin uint8) uint8 {
	return ((bin / 10) << 4 + (bin % 10))
}
	
func (this *DS3231) Reset() {
	this.twi.Write(DS3231_ADDR, RTC_CTL_REG_ADDR, []uint8{RTC_CTL_BIT_RS1 | RTC_CTL_BIT_RS2})
}

func (this *DS3231) Write(Year uint8, Month uint8, Day uint8, Hour uint8, Minute uint8, Second uint8, Week uint8) bool {
	DT:= []uint8{this.bin2bcd(Second), this.bin2bcd(Minute), this.bin2bcd(Hour), 
				this.bin2bcd(Week), this.bin2bcd(Day), this.bin2bcd(Month), this.bin2bcd(Year)}
	return this.twi.Write(DS3231_ADDR, RTC_SEC_REG_ADDR, DT)
}

func (this *DS3231) Read() ([]uint8, bool) { //Year, Month, Day, Hour, Minute, Second, Week
	Buf, Result := this.twi.Read(DS3231_ADDR, RTC_SEC_REG_ADDR, 7)
	if !Result { return []uint8{0, 0, 0, 0, 0, 0, 0}, Result }
	
	DT := []uint8{0, 0, 0, 0, 0, 0, 0}
	DT[0] = this.bcd2bin(Buf[6])
	DT[1] = this.bcd2bin(Buf[5])
	DT[2] = this.bcd2bin(Buf[4])
	DT[3] = this.bcd2bin(Buf[2])
	DT[4] = this.bcd2bin(Buf[1])
	DT[5] = this.bcd2bin(Buf[0])
	DT[6] = this.bcd2bin(Buf[3])
	return DT, Result
}

func (this *DS3231) Temperature() float32 {
	Result := this.twi.Write(DS3231_ADDR, RTC_CTL_REG_ADDR, []uint8{0x20})
	if !Result { return 255.0}
	Buf, _ := this.twi.Read(DS3231_ADDR, 0x11, 2)
	//if Buf[0] > 0x85 { Buf[0] = 0xFF - Buf[0]}
	return float32(this.bcd2bin(Buf[0])) + float32(this.bcd2bin(Buf[1])) / 100
}
