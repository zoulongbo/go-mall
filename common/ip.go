package common

import (
	"errors"
	"net"
)

func GetEntranceIp()(string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addr {
		//检查ip地址是否回环地址
		if ipnet, ok :=address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}
	return "", errors.New("获取ip异常")
}