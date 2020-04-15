package ifAddr

import (
	"strings"
)

//ISP类型
const (
	TEL  = 1 << iota
	UNI  = 1 << iota
	CMCC = 1 << iota
	EDU  = 1 << iota
	NA   = 1 << iota
	BGP  = TEL | UNI | CMCC | EDU | NA
)

type NetIsp struct {
	ISP int //ISP库中查找出来的原始ISP字符串
}

func (s *NetIsp) SetLocISP(isp string) {
	s.ISP = 0

	//国内IP库能识别电信/联通/移动/教育网
	isps := strings.Split(isp, "/")
	for _, item := range isps {
		switch item {
		case "电信":
			s.ISP |= TEL
		case "联通":
			s.ISP |= UNI
		case "移动":
			s.ISP |= CMCC
		case "教育网":
			s.ISP |= EDU
		case "BGP":
			s.ISP |= BGP
		case "TEL":
			s.ISP |= TEL
		case "UNI":
			s.ISP |= UNI
		case "CMCC":
			s.ISP |= CMCC
		case "EDU":
			s.ISP |= EDU
		default:
			s.ISP |= NA
		}
	}
}
func (s *NetIsp) SetISP(isps []string) {
	s.ISP = 0

	//国内IP库能识别电信/联通/移动/教育网
	for _, item := range isps {
		switch item {
		case "电信":
			s.ISP |= TEL
		case "联通":
			s.ISP |= UNI
		case "移动":
			s.ISP |= CMCC
		case "教育网":
			s.ISP |= EDU
		case "BGP":
			s.ISP |= BGP
		case "TEL":
			s.ISP |= TEL
		case "UNI":
			s.ISP |= UNI
		case "CMCC":
			s.ISP |= CMCC
		case "EDU":
			s.ISP |= EDU
		default:
			s.ISP |= NA
		}
	}
}

func (s *NetIsp) GetISP() int {
	return s.ISP
}

func (s NetIsp) HasTel() bool {
	return s.ISP & TEL != 0
}

func (s NetIsp) HasCmcc() bool {
	return s.ISP & CMCC != 0
}

func (s NetIsp) HasUni() bool {
	return s.ISP & UNI != 0
}

func (s NetIsp) HasEdu() bool {
	return s.ISP & EDU != 0
}

func (s NetIsp) HasNa() bool {
	return s.ISP & NA != 0
}

func (s *NetIsp) IsNa() bool {
	t := s.ISP
	n := 0
	for t > 0 {
		t = t & (t - 1)
		n++
	}

	if n == 1 && (s.ISP&NA) == NA {
		return true
	}

	return false
}

func (s NetIsp) IsBGP() bool {
	return s.ISP == BGP
}

func (s NetIsp) IsTripleLine() bool {
	return s.HasTel() && s.HasCmcc() && s.HasUni()
}

func (s *NetIsp) NeedProxy() bool {
	if s.ISP&TEL == TEL || s.ISP&UNI == UNI || s.ISP&CMCC == CMCC {
		return false
	}
	return true
}

func (s NetIsp) String() string {
	vec := []string{}
	if (s.ISP & BGP) == BGP {
		vec = append(vec, "BGP")
	} else {
		if (s.ISP & TEL) == TEL {
			vec = append(vec, "TEL")
		}
		if (s.ISP & UNI) == UNI {
			vec = append(vec, "UNI")
		}
		if (s.ISP & CMCC) == CMCC {
			vec = append(vec, "CMCC")
		}
		if (s.ISP & EDU) == EDU {
			vec = append(vec, "EDU")
		}
		if (s.ISP & NA) == NA {
			vec = append(vec, "NA")
		}
	}
	return strings.Join(vec, "|")
}