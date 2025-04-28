package types

import "fmt"

type Scheme string

const (
	Vmess  Scheme = "vmess"
	Vless  Scheme = "vless"
	SS     Scheme = "ss"
	Trojan Scheme = "trojan"
	Socks  Scheme = "socks"
)

var SupportedSchemes = []Scheme{
	Vmess,
	Vless,
	SS,
	Trojan,
	Socks,
}

var schemeToStr = map[Scheme]string{
	Vmess:  "vmess",
	Vless:  "vless",
	SS:     "ss",
	Trojan: "trojan",
	Socks:  "socks",
}

var strToScheme = map[string]Scheme{
	"vmess":  Vmess,
	"vless":  Vless,
	"ss":     SS,
	"trojan": Trojan,
	"socks":  Socks,
}

// 枚举 → String
func (s Scheme) String() string {
	return schemeToStr[s]
}

// String → 枚举
func ParseScheme(str string) (Scheme, error) {
	if s, ok := strToScheme[str]; ok {
		return s, nil
	}
	return "", fmt.Errorf("无效状态: %s", str)
}
