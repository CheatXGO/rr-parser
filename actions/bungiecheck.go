package actions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"main/structures"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var mu sync.Mutex

func CheckStats(usercommand, u string, t bool, cfg structures.Config) string {
	var wg sync.WaitGroup
	var bPlayer structures.BungoSearch
	var bAnswer structures.BungoAnswer
	var cnormal, chard, cgg int
	cnch := make(chan int)
	chch := make(chan int)
	cggch := make(chan int)
	cquit := make(chan struct{})
	defer close(cnch)
	defer close(chch)
	defer close(cggch)
	defer close(cquit)
	nck := u
	if !t {
		return u //user not registered in DB
	}
	raids := make(map[string]string)
	cnormal, chard, cgg = 0, 0, 0 //reset counters
	for key, val := range Raids {
		raids[strings.ToLower(key)] = strings.Trim(strings.ToLower(val), "\n")
	}
	if _, ok := raids[strings.ToLower(usercommand)]; ok { //if /my raid in raid list
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://www.bungie.net/Platform/User/Search/Prefix/"+u+"/0/", nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Set("X-API-Key", cfg.Bapi)
		res, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(body, &bPlayer)
		if err != nil {
			e := structures.MyError{Fun: "bungiecheck json.Unmarshal &bPlayer", Err: "can't decode bPlayer answer"}
			log.Fatalln(e.Error())
		}
		if len(bPlayer.Response.SearchResults) == 0 {
			return "Nickname: " + u + " not found in Bungie database"
		} else {
			for i := range bPlayer.Response.SearchResults {
				for j := range bPlayer.Response.SearchResults[i].DestinyMemberships {
					if bPlayer.Response.SearchResults[i].DestinyMemberships[j].CrossSaveOverride == bPlayer.Response.SearchResults[i].DestinyMemberships[j].MembershipType {
						u = bPlayer.Response.SearchResults[i].DestinyMemberships[j].MembershipId
					}
				}
			}
			if bPlayer.Response.SearchResults[0].DestinyMemberships[0].CrossSaveOverride == 0 { //if not crossplayed
				u = bPlayer.Response.SearchResults[0].DestinyMemberships[0].MembershipId
			}
			req, err = http.NewRequest("GET", "https://kuqw4o1zj5.execute-api.us-west-2.amazonaws.com/api/player/"+u, nil)
			if err != nil {
				log.Fatalln(err)
			}
			res, err = client.Do(req)
			if err != nil {
				log.Fatalln(err)
			}
			body, err = ioutil.ReadAll(res.Body)
			err = json.Unmarshal(body, &bAnswer)
			if err != nil {
				e := structures.MyError{Fun: "bungiecheck json.Unmarshal &bPlayer", Err: "can't decode bPlayer answer"}
				log.Fatalln(e.Error())
			}
			for i := range bAnswer.Response.Activities {
				wg.Add(1)
				go fastquery(&wg, client, cfg, raids, bAnswer.Response.Activities[i].ActivityHash, bAnswer.Response.Activities[i].Values.Clears, usercommand, cnch, chch, cggch)
			}
			go func(wg *sync.WaitGroup, cquit chan struct{}) {
				wg.Wait()
				cquit <- struct{}{}
			}(&wg, cquit)
			for {
				select {
				case c1 := <-cnch:
					cnormal = cnormal + c1
				case c2 := <-chch:
					chard = chard + c2
				case c3 := <-cggch:
					cgg = cgg + c3
				case <-cquit:
					//break
					return nck + " " + usercommand + " N: " + strconv.Itoa(cnormal) + " H: " + strconv.Itoa(chard) + " G: " + strconv.Itoa(cgg)
				}
			}
		}
	} else {
		return "Please fill right raid abbreviation after /my"
	}
}

func fastquery(wg *sync.WaitGroup, client *http.Client, cfg structures.Config, raids map[string]string, i int64, j int, usercommand string, cnch, chch, cggch chan int) {
	var bRaid structures.RaidAnswer
	defer wg.Done()
	reqs, err := http.NewRequest("GET", "https://www.bungie.net/Platform/Destiny2/Manifest/DestinyActivityDefinition/"+strconv.FormatInt(i, 10), nil)
	if err != nil {
		log.Fatalln(err)
	}
	reqs.Header.Set("X-API-Key", cfg.Bapi)
	resq, err := client.Do(reqs)
	if err != nil {
		log.Fatalln(err)
	}
	bodys, err := ioutil.ReadAll(resq.Body)
	err = json.Unmarshal(bodys, &bRaid)
	if err != nil {
		e := structures.MyError{Fun: "bungiecheck json.Unmarshal &bRaid", Err: "can't decode bRaid answer"}
		log.Fatalln(e.Error())
	}
	for key, val := range raids {
		if val == strings.ToLower(bRaid.Response.OriginalDisplayProperties.Name) {
			if key == strings.ToLower(usercommand) {
				if bRaid.Response.Matchmaking.RequiresGuardianOath == false {
					if bRaid.Response.Tier == 0 {
						cnch <- j
					} else {
						chch <- j
					}
				} else {
					cggch <- j
				}
			}
		}
	}
}
