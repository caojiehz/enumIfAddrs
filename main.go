package main

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/wangtuanjie/ip17mon"
	"net"
	"net/http"
	"time"
)

const (
	PublicIPURL = "http://127.0.0.1:6001/smartip"
)

type IfAddr struct {
	IfName string
	MTU    int
	IPv4   string
}

type SmartIP struct {
	IP string
	ip17mon.LocationInfo
}

type IfAddrPublic struct {
	IfAddr
	SmartIP
}

func (addr IfAddrPublic)InnerIP()string{
	return addr.IPv4
}

func (addr IfAddrPublic)OuterIP()string{
	return addr.IP
}

func (addr IfAddrPublic)IsOuterIP()bool{
	if addr.IP == ""{
		return false
	}
	return addr.IPv4 == addr.IP
}

func main() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("net.Interfaces failed: %s", err.Error())
		return
	}

	var ifAddrs []IfAddr
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) == 0 {
			continue
		}

		addrs, _ := netInterfaces[i].Addrs()
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					var addr IfAddr
					addr.IfName = netInterfaces[i].Name
					addr.MTU = netInterfaces[i].MTU
					addr.IPv4 = ipnet.IP.String()
					ifAddrs = append(ifAddrs, addr)
				}
			}
		}

	}

	get_public_ip := func(local string) string {
		transport := &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				//本地地址  ipaddr是本地外网IP
				lAddr, err := net.ResolveTCPAddr(netw, local+":0")
				if err != nil {
					return nil, err
				}
				//被请求的地址
				rAddr, err := net.ResolveTCPAddr(netw, addr)
				if err != nil {
					return nil, err
				}
				conn, err := net.DialTCP(netw, lAddr, rAddr)
				if err != nil {
					return nil, err
				}
				deadline := time.Now().Add(5 * time.Second)
				conn.SetDeadline(deadline)
				return conn, nil
			}}

		req := httplib.Get(PublicIPURL).SetTransport(transport)
		str, err := req.String()
		if err != nil {
			fmt.Println(err)
		}
		return str
	}
	for _, item := range ifAddrs {
		fmt.Printf("%+v, public: %s\n", item, get_public_ip(item.IPv4))
	}
}
