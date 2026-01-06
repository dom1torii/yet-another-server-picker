package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Response struct {
	Success bool           `json:"success"`
	Pops    map[string]Pop `json:"pops"`
}

type Pop struct {
	Desc   string  `json:"desc"`
	Relays []Relay `json:"relays"`
}

type Relay struct {
	Ipv4      string `json:"ipv4"`
	PortRange [2]int `json:"port_range"`
}

func FetchRelays() (Response, error) {
	url := "https://api.steampowered.com/ISteamApps/GetSDRConfig/v1?appid=730"
	log.Println("Fetching API...")

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("Error fetching:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error fetching:", err)
	}

	log.Println("Success fetching API. Status:", resp.Status)

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return Response{}, err
	}

	response.Pops = filterPops(response.Pops)

	return response, nil
}

// skip useless pops that have no relays
func filterPops(pops map[string]Pop) map[string]Pop {
	filteredPops := make(map[string]Pop)
	for key, pop := range pops {
		if len(pop.Relays) > 0 {
			filteredPops[key] = pop
		}
	}
	return filteredPops
}
