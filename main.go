package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"time"

	_ "embed"
)

func sayHome() {
	cmd := exec.Command("say", "-v", "daniel", "Dan is home")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

type aBoy struct {
	People []struct {
		Lan   string `json:"lan"`
		Name  string `json:"name"`
		State string `json:"state"`
	} `json:"people"`
}

func (b *aBoy) inDerDo(lanName string) bool {
	for _, i := range b.People {
		if lanName == i.Lan {
			return true
		}
	}
	return false
}

func (b *aBoy) UpdateByLan(lan string, state string) {
	for i, _ := range b.People {
		if lan == b.People[i].Lan {
			b.People[i].State = "hey"
		}
	}
}

// TODO: Update to read file so we can update state... or just switch to sqlite or something better than a json file
//go:embed people.json
var theBoy []byte

func main() {
	// Figure out base address
	// ifaces, err := net.Interfaces()

	// if err != nil {
	// 	panic(err)
	// }
	// for in := range ifaces {
	// 	fmt.Println(ifaces[in])
	// }

	var myStuff aBoy
	json.Unmarshal(theBoy, &myStuff)
	fmt.Println(myStuff)
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
			if myStuff.inDerDo(name) {
				go sayHome()
				myStuff.UpdateByLan(name, "here")
				dat, _ := json.Marshal(myStuff)
				fmt.Println(string(dat))
			}
		}
	}
}
