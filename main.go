package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"time"
)

func sayHome() {
	cmd := exec.Command("say", "-v", "daniel", "Dan is home")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	// Figure out base address
	// ifaces, err := net.Interfaces()

	// if err != nil {
	// 	panic(err)
	// }
	// for in := range ifaces {
	// 	fmt.Println(ifaces[in])
	// }
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "192.168.86.1:53")
		},
	}
	for i := 0; i < 255; i++ {
		names, err := r.LookupAddr(context.TODO(), (fmt.Sprintf("192.168.86.%v", i)))
		if err != nil {
			continue
		}
		// fmt.Println(names)
		for _, name := range names {
			fmt.Println(name)
			if name == "dans-iphone-x.lan." {
				go sayHome()
			}
		}
	}
}
