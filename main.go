package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"goDeviceScanner/speaking"
	"net"
	"os"
	"os/exec"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db            *sql.DB
	speakingVoice string = "daniel"
)

func sayHome(spokenName string) {
	cmd := exec.Command("say", "-v", speakingVoice, spokenName, "is home")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func sayLeft(spokenName string) {
	cmd := exec.Command("say", "-v", speakingVoice, spokenName, "has left")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func hasArrived(lan string) string {
	query := `
		SELECT status, spoken_name
		FROM people
		WHERE lan_id = ?
	`
	var status string
	var spokenName string
	err := db.QueryRow(query, lan).Scan(&status, &spokenName)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			// fmt.Println("User not found in db")
			return ""
		} else {
			panic(err)
		}
	}
	if status == "away" { // Person has arrived
		update := `
			UPDATE people
			SET status = 'here'
			WHERE lan_id = ?
		`
		res, err := db.Exec(update, lan)
		if err != nil {
			panic(err)
		}
		numRows, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Deemed arrived %v\n", numRows)
		return spokenName
	}
	return ""
}

func hasLeft(names []string) [][]string {
	query := `
		SELECT lan_id, spoken_name
		FROM people
		WHERE status = 'here'
	`
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	var peopleThatLeftNames [][]string
	for rows.Next() {
		var lanID string
		var spokenName string
		err := rows.Scan(&lanID, &spokenName)
		if err != nil {
			panic(err)
		}
		isHere := false
		for _, personHere := range names { // Check if they are in the list of people we see
			if personHere == lanID {
				isHere = true
			}
		}
		if !isHere {
			peopleThatLeftNames = append(peopleThatLeftNames, []string{lanID, spokenName})
			fmt.Println("Deemed away", lanID)
		}
	}
	for _, personAway := range peopleThatLeftNames {
		updateQ := `
			UPDATE people
			SET status = 'away'
			WHERE lan_id = ?
		`
		_, err := db.Exec(updateQ, personAway[0])
		if err != nil {
			panic(err)
		}
	}
	return peopleThatLeftNames
}

func UpdateByLan(lan string, state string) {
	query := `
		UPDATE people
		SET status = ?
		WHERE lan_id = ?
	`
	res, err := db.Exec(query, state, lan)
	if err != nil {
		panic(err)
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated %v rows\n", numRows)
}

func SetupDB() {
	query := `
		CREATE TABLE IF NOT EXISTS people (
			fname TEXT NOT NULL,
			lname TEXT NOT NULL,
			spoken_name TEXT NOT NULL,
			status TEXT NOT NULL,
			lan_id TEXT NOT NULL PRIMARY KEY
		);
		-- INSERT INTO people (fname, lname, spoken_name, lan_id, status) VALUES ('Dan', 'Goodman', 'Dan', 'dans-iphone-x.lan.', 'away');
	`
	_, err := db.Exec(query)
	if err != nil {

		panic(err)
	}
}

// TODO: Update to read file so we can update state... or just switch to sqlite or something better than a json file

func main() {
	// Figure out base address
	// ifaces, err := net.Interfaces()

	// if err != nil {
	// 	panic(err)
	// }
	// for in := range ifaces {
	// 	fmt.Println(ifaces[in])
	// }

	// Check GCP Connection
	speaking.TestAuth()

	// Open DB
	dbConn, err := sql.Open("sqlite3", "./people.db")
	if err != nil {
		panic(err)
	}
	db = dbConn

	// SetupDB()

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "192.168.86.1:53")
		},
	}
	for true {
		if os.Getenv("DEBUG") == "true" {
			fmt.Println("Running")
		}
		var foundAddrs []string
		for i := 0; i < 255; i++ {
			names, err := r.LookupAddr(context.TODO(), (fmt.Sprintf("192.168.86.%v", i)))
			if err != nil {
				continue
			}
			for _, name := range names { // Check if people have arrived
				if os.Getenv("DEBUG") == "true" {
					fmt.Println(name)
				}
				foundAddrs = append(foundAddrs, name)
				if arrivedName := hasArrived(name); arrivedName != "" {
					speaking.Say(arrivedName, "arrived")
				}
			}
		}
		// Check who is away
		peeps := hasLeft(foundAddrs)
		for _, person := range peeps {
			speaking.Say(person[1], "left")
		}
		time.Sleep(time.Millisecond * 2500)
	}

	db.Close()
}
