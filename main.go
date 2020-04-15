package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/caojiehz/enumIfAddrs/ifAddr"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	PublicIPURL = "http://149.28.31.219:6001/smartip"
)

func main() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("net.Interfaces error: %s\n", err.Error())
		return
	}

	var ifAddrs []ifAddr.LocalIfAddr
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) == 0 || strings.Contains(netInterfaces[i].Name, "docker"){
			continue
		}

		addrs, _ := netInterfaces[i].Addrs()
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					var addr ifAddr.LocalIfAddr
					addr.IfName = netInterfaces[i].Name
					addr.MTU = netInterfaces[i].MTU
					addr.IPv4 = ipnet.IP.String()
					ifAddrs = append(ifAddrs, addr)
				}
			}
		}

	}

	get_public_ip := func(local ifAddr.LocalIfAddr) (publicIP ifAddr.PublicIP, err error ){
		transport := &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				//本地地址  ipaddr是本地外网IP
				lAddr, err := net.ResolveTCPAddr(netw, local.IPv4+":0")
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
		HostName, _ := os.Hostname()
		req.Header("HostName", HostName)
		req.Header("IfName", local.IfName)
		req.Header("IfAddr", local.IPv4)

		body, err := req.String()
		if err != nil {
			//fmt.Printf("get %s error: %s\n", PublicIPURL, err.Error())
			return
		}

		err = json.Unmarshal([]byte(body), &publicIP)
		if err != nil{
			fmt.Printf("unmarshal %s error: %s\n", body, err.Error())
			return
		}
		return
	}

	var addrs []ifAddr.IfAddrInfo
	for _, item := range ifAddrs {
		var addr ifAddr.IfAddrInfo
		addr.LocalIfAddr = item
		addr.PublicIP, err = get_public_ip(item)
		if err == nil{
			addrs = append(addrs, addr)
			fmt.Printf("%v\n", addr)
		}
	}
}
