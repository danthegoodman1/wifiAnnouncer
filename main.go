package main

import (
	"context"
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
	"wifiAnnouncer/configParser"
	"wifiAnnouncer/speaking"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db          *sql.DB
	scannerMode bool
)

func hasArrived(lan string) string {
	DebugLog("--- Checking if", lan, "has just arrived")
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
			return ""
		} else {
			panic(err)
		}
	}
	if status == "away" { // Person has arrived
		DebugLog("User arrived:", lan)
		update := `
			UPDATE people
			SET status = 'here'
			WHERE lan_id = ?
		`
		_, err := db.Exec(update, lan)
		if err != nil {
			panic(err)
		}
		DebugLog("Set user arrived in db", lan)
		// numRows, err := res.RowsAffected()
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Printf("Deemed arrived %v\n", numRows)
		return spokenName
	}
	return ""
}

func hasLeft(names []string) [][]string {
	DebugLog("### Checking if anyone has left")
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
		DebugLog(lanID, "is considered here right now")
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
			DebugLog("Deemed away", lanID)
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
			spoken_name TEXT NOT NULL,
			status TEXT NOT NULL,
			lan_id TEXT NOT NULL PRIMARY KEY
		);
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
	configParser.ParseConfig()
	parseFlags()

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "192.168.86.1:53")
		},
	}

	if scannerMode {
		fmt.Println("### Scanning for devices!")
		for i := 0; i < 255; i++ {
			names, err := r.LookupAddr(context.TODO(), (fmt.Sprintf("%s.%d", configParser.InterfaceToPrefix(), i)))
			if err != nil {
				continue
			}
			for _, name := range names { // Check if people have arrived
				fmt.Println(name)
			}
		}
		os.Exit(0)
	}

	// Check GCP Connection
	speaking.TestAuth()
	DebugLog("Connecting to GCP valid!")

	// Check if people.db exists

	createdDB := false

	if _, err := os.Stat("./people.db"); os.IsNotExist(err) {
		DebugLog("Creating sqlite file")
		createdDB = true
		file, err := os.Create("./people.db") // Create SQLite file
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	// Open DB
	dbConn, err := sql.Open("sqlite3", "./people.db")
	if err != nil {
		panic(err)
	}
	db = dbConn

	if createdDB {
		SetupDB()
	}

	verifyRegisteredDevices()

	for {
		DebugLog("Starting to scan")
		var foundAddrs []string
		for i := 0; i < 255; i++ {
			names, err := r.LookupAddr(context.TODO(), (fmt.Sprintf("192.168.86.%v", i)))
			if err != nil {
				continue
			}
			for _, name := range names { // Check if people have arrived
				DebugLog("Scan Found:", name)
				foundAddrs = append(foundAddrs, name)
				if arrivedName := hasArrived(name); arrivedName != "" {
					if speaking.IsInConfig(name) {
						DebugLog(name, "is in config, speaking")
						speaking.Say(arrivedName, configParser.Config.ArrivedSuffix)
					}
				}
			}
		}
		// Check who is away
		peeps := hasLeft(foundAddrs)
		for _, person := range peeps {
			if speaking.IsInConfig(person[0]) {
				DebugLog(person, "is in config, speaking")
				speaking.Say(person[1], configParser.Config.LeftSuffix)
			}
		}
		time.Sleep(time.Millisecond * 2500)
	}
}

func parseFlags() {
	flag.BoolVar(&scannerMode, "scannerMode", false, "Performs a network scan and exits")
	flag.Parse()
}

func verifyRegisteredDevices() {
	// Check all included ones exist
	for _, watchedDevice := range configParser.Config.RegisteredDevices {
		DebugLog("Trying to insert user", watchedDevice.Hostname)
		query := `
			INSERT INTO people (spoken_name, lan_id, status) VALUES (?, ?, ?) ON CONFLICT DO NOTHING;
		`
		_, err := db.Exec(query, watchedDevice.Name, watchedDevice.Hostname, watchedDevice.DefaultState)
		if err != nil {
			panic(err)
		}
	}

	// Check there are no extra ones in DB (delete ones that aren't in config)
	query := `
		SELECT lan_id FROM people;
	`
	row, err := db.Query(query)

	if err != nil {
		panic(err)
	}

	defer row.Close()
	for row.Next() {
		var lanID string
		row.Scan(&lanID)
		DebugLog("Checking if removed", lanID)
		inConfig := false
		for _, watchedDevices := range configParser.Config.RegisteredDevices {
			// We found it in the config
			if lanID == watchedDevices.Hostname {
				inConfig = true
			}
		}

		// Must have been removed from config, delete
		if !inConfig {
			DebugLog("Not found in config, removing", lanID)
			_, err = db.Query(`
				DELETE FROM people
				WHERE lan_id = ?
				;
			`, lanID)

			if err != nil {
				panic(err)
			}
		}
	}
}

func DebugLog(message ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		fmt.Println(message...)
	}
}
