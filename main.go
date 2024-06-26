package main

import (
	"database/sql"
	"log"
	"math"
	"net"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func ping(ipv4_addr string) {
	c, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatalf("ListenPacket: %v", err)
	}
	defer c.Close()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  0,
			Data: []byte("Gopher"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	if _, err := c.WriteTo(wb, &net.UDPAddr{IP: net.ParseIP(ipv4_addr)}); err != nil {
		log.Fatal(err)
	}

	rb := make([]byte, 1500)
	n, _, err := c.ReadFrom(rb)
	if err != nil {
		log.Fatalf("ReadFrom: %v", err)
	}
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rb[:n])
	if err != nil {
		log.Fatalf("ICMP Parse: %v", err)
	}
	if rm.Type == ipv4.ICMPTypeEchoReply {
		log.Println("Received ICMP Echo Reply")
	}
}

func main() {
	log.Println("Started")
	var client = redistimeseries.NewClient("localhost:6379", "nohelp", nil)
	var keyname = "mytest"
        _, haveit := client.Info(keyname)
	if haveit != nil {
			client.CreateKeyWithOptions(keyname, redistimeseries.DefaultCreateOptions)
			client.CreateKeyWithOptions(keyname+"_avg", redistimeseries.DefaultCreateOptions)
			client.CreateRule(keyname, redistimeseries.AvgAggregation, 60, keyname+"_avg")
        }

	// db, err := sql.Open("sqlite3", "/grafana/metrics.db")
	// if err != nil {
	// 	log.Fatalf("SQL Open: %v", err)
	// }
	// defer db.Close()

	// sqlStmt := `
	// create table metrics (ts INTEGER NOT NULL PRIMARY KEY, value REAL);
	// delete from metrics;
	// `
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Printf("DB Exec: %v %s\n", err, sqlStmt)
	// }

	for {
		dt_start := time.Now()
		ping("10.204.227.154")
		dt_end := time.Now()
		rtt := float64(dt_end.Sub(dt_start))
		rtt2 := math.Round(rtt/1000) / 1000
		log.Println("Elapsed:", rtt2)

		// stmt, err := db.Prepare("insert into metrics(ts, value) values(?, ?)")
		// if err != nil {
		// 	log.Fatalf("DB Prepare: %v", err)
		// }
		// defer stmt.Close()

		// dt_unix := dt_start.Unix()
		// stmt.Exec(dt_unix, rtt2)

		// Add sample with timestamp from server time and value 100
		// TS.ADD mytest * 100
		_, err := client.AddAutoTs(keyname, rtt2)
		if err != nil {
		        fmt.Println("Error:", err)
		}

		// time.Sleep(time.Second)
	}
	// 	dt := time.Now()
	// 	unix := dt.Unix()
	// 	stmt.Exec(unix, i%10)
	// for i := 0; i < 100; i++ {
	// 	log.Println(i)
	// 	stmt, err := db.Prepare("insert into metrics(ts, value) values(?, ?)")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer stmt.Close()

	// 	dt := time.Now()
	// 	unix := dt.Unix()
	// 	stmt.Exec(unix, i%10)
	// 	time.Sleep(time.Second)
	// }

}
