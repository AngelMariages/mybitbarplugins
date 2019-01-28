package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Response struct {
	Time        int64 `json:"now"`
	Origin      string
	Destination string
}

func main() {
	//Do the delta time
	data := Response{}

	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			json.Unmarshal(body, &data)
			t := time.Unix(int64(data.Time/1000), 0)
			fmt.Printf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
			fmt.Printf("\n---\nOrigin: %s, Destination: %s", data.Origin, data.Destination)
			return
		}
	}
	fmt.Println("Error", err)
}
