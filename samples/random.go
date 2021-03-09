package samples

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/tkircsi/pcbook/pb"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomFloat64(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomFloat32(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomStringFromSet(e ...string) string {
	l := len(e)
	if l == 0 {
		return ""
	}
	ix := rand.Intn(l)
	return e[ix]
}

func randomID() string {
	return uuid.New().String()
}

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QUERTY
	case 2:
		return pb.Keyboard_QUERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet(
			"Core i9-10980HK",
			"Core i7-10700TE",
			"Core i7-10875H",
			"Core i9-10885H",
			"Core i7-10870H",
			"Xeon E-2286M",
			"Core i7-11375H",
			"Core i5-10500H",
			"Core i7-10850H",
			"Core i7-10750H",
			"Core i7-8700B",
			"Xeon E-2276M",
			"Core i5-8500B",
			"Core i7-8569U",
			"Core i5-9500TE",
			"Core i7-9850HL",
			"Core i7-10810U",
			"Core i5-10200H",
			"Core i7-8559U",
			"Core i7-1065G7",
			"Core i5-10400H",
			"Core i5-10300H",
			"Core i5-1035G7",
		)
	}
	return randomStringFromSet(
		"AMD Ryzen 9 5980HS",
		"AMD Ryzen 9 5900HX",
		"AMD Ryzen 9 5900HS",
		"AMD Ryzen 7 5800HS",
		"AMD Ryzen 7 5800H",
		"AMD Ryzen 9 4900HS",
		"AMD Ryzen 7 5800U",
		"AMD Ryzen 7 4800H",
		"AMD Ryzen 9 4900H",
		"AMD Ryzen 7 4800HS",
		"AMD Ryzen 7 Extreme Edition",
		"AMD Ryzen 7 4800U",
		"AMD Ryzen 5 4600H",
		"AMD Ryzen 5 PRO 4400GE",
		"AMD Ryzen 5 4600HS",
		"AMD Ryzen 5 4600U",
		"AMD Ryzen 7 4700U",
		"AMD Ryzen 5 5500U",
	)
}

func randomGPUBrand() string {
	return randomStringFromSet("NVIDIA", "AMD")
}

func randomGPUName(brand string) string {
	if brand == "NVIDIA" {
		return randomStringFromSet(
			"GeForce GTX 1650",
			"GeForce GTX 1660 SUPER",
			"GeForce GTX 1660 Ti",
			"GeForce GTX 1050 Ti",
			"GeForce RTX 2060 SUPER",
			"GeForce RTX 2070",
			"GeForce RTX 3060 Laptop GPU",
			"GeForce RTX 2070 SUPER",
			"GeForce RTX 3060 Ti",
			"GeForce GT 1030",
			"GeForce RTX 3070",
		)
	}
	return randomStringFromSet(
		"Radeon RX 590",
		"Radeon RX 580",
		"Radeon RX 5600 XT",
		"Radeon RX 5500 XT",
		"Radeon RX 5700",
		"Radeon RX 5700 XT",
		"Radeon RX 6800 XT",
		"Radeon RX 6800",
		"Radeon Pro WX 4100",
		"Radeon Pro WX 8200",
	)
}

func randomScreenResolution() *pb.Screen_Resolution {
	height := randomInt(1080, 4320)
	width := height * 16 / 9
	return &pb.Screen_Resolution{
		Height: uint32(height),
		Width:  uint32(width),
	}
}

func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

func randomLaptopBrand() string {
	return randomStringFromSet(
		"Lenovo",
		"HP",
		"Dell",
		"Apple",
		"Asus",
	)
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Lenovo":
		return randomStringFromSet(
			"IdeaPad S145 81N300MNRM Notebook",
			"Legion 5 82AU005CHV Notebook",
			"Ideapad 3 81W00095HV Notebook",
			"V15 82C500JFRM Notebook",
			"IdeaPad 3 81W100MBRM Notebook",
			"ThinkBook 15 G2 20VG006CRM Notebook",
			"V15 82C7001HHV Notebook",
		)
	case "HP":
		return randomStringFromSet(
			"HP 15-dw1018nq 2G2C0EA Notebook",
			"HP 250 G7 14Z83EA Notebook",
			"HP Pavilion 15-ec0014nh 9HF69EA Notebook",
			"HP 15-dw1017nq 2G2B8EA Notebook",
			"HP 250 G8 2M2L6ES Notebook",
			"HP 14s-dq1009nh 8BW28EA Notebook",
			"HP 15-da2048nq 2L9N8EA Notebook",
		)
	case "Dell":
		return randomStringFromSet(
			"Dell Vostro 3500 N3004VN3500EMEA01_2105_UBU Notebook",
			"Dell Latitude 3410 DL341015898541UBU Notebook",
			"Dell Vostro 3501 N6503VN3501EMEA01 Notebook",
			"Dell Latitude 7480 L7480-64 Notebook",
			"Dell Inspiron 3501 3501FI3UA1 Notebook",
			"Dell Inspiron 3793 DI3793I341UHDUBU Notebook",
		)
	case "Apple":
		return randomStringFromSet(
			"Apple MacBook Air 13.3 M1 Chip 8GB 256GB MGN63 Notebook",
			"Apple MacBook Air 13 Mid 2017 MQD32 Notebook",
			"Apple MacBook Air 13.3 MGND3ZE/A Notebook",
			"Apple MacBook Air 13 2020 MWTJ2 Notebook",
			"Apple MacBook Pro 16 MVVJ2 Notebook",
			"Apple MacBook Pro 16 MVVM2 Notebook",
		)
	case "Asus":
		return randomStringFromSet(
			"ASUS M509DA-BR988 Notebook",
			"ASUS ZenBook 14 UX431FL-AN014T Notebook",
			"ASUS VivoBook 15 M509DA-BR1421 Notebook",
			"ASUS VivoBook X X543MA-DM1216 Notebook",
			"ASUS TUF Gaming FA506II-HN306 Notebook",
			"ASUS VivoBook S14 M433IA-EB400 Notebook",
		)
	default:
		return "Unknown"
	}
}
