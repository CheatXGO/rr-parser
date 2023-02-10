package actions

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/CheatXGO/rr-parser/structures"
)

func CheckStats(usercommand, u string, t bool, cfg structures.Config, page int) string {
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
	if !t {
		return u //user not registered in DB
	}
	raids := make(map[string]string)
	cnormal, chard, cgg = 0, 0, 0 //reset counters
	for key, val := range Raids {
		raids[strings.ToLower(key)] = strings.Trim(strings.ToLower(val), "\n")
	}
	if _, ok := raids[strings.ToLower(usercommand)]; ok { //if /my raid in raid list
		var (
			bungoDisplayName  = u
			bungoNameCode     int
			bungoMembershipId string
			err               error
		)

		if i := strings.Index(bungoDisplayName, "#"); i > 0 {
			bungoNameCode, err = strconv.Atoi(strings.ReplaceAll(bungoDisplayName[i+1:], " ", ""))
			if err != nil {
				return "Part after #: " + bungoDisplayName[i+1:] + " can contains only digits"
			}
			bungoDisplayName = bungoDisplayName[:i]
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://www.bungie.net/Platform/User/Search/Prefix/"+bungoDisplayName+"/"+strconv.Itoa(page)+"/", nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Set("X-API-Key", cfg.Bapi)
		res, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		body, err := io.ReadAll(res.Body)
		err = json.Unmarshal(body, &bPlayer)
		if err != nil {
			e := structures.MyError{Func: "bungiecheck json.Unmarshal &bPlayer", Err: "can't decode bPlayer answer"}
			log.Fatalln(e.Error())
		}

		if len(bPlayer.Response.SearchResults) == 0 {
			return "Nickname: " + u + " not found in Bungie database"
		}

		if len(bPlayer.Response.SearchResults) == 1 && bungoNameCode == 0 {
			bungoNameCode = bPlayer.Response.SearchResults[0].BungieGlobalDisplayNameCode
		}

		if len(bPlayer.Response.SearchResults) > 1 && bungoNameCode == 0 {
			return "Found more than one bungie IDs associated with nickname: " + u + " Please fill nickname in bungie ID format: " + u + "#3011 for example."
		}

		for i := range bPlayer.Response.SearchResults {
			if bPlayer.Response.SearchResults[i].BungieGlobalDisplayNameCode == bungoNameCode {
				for j := range bPlayer.Response.SearchResults[i].DestinyMemberships {
					if bPlayer.Response.SearchResults[i].DestinyMemberships[j].CrossSaveOverride == bPlayer.Response.SearchResults[i].DestinyMemberships[j].MembershipType ||
						bPlayer.Response.SearchResults[i].DestinyMemberships[j].CrossSaveOverride == 0 { //if not crosssaved
						bungoDisplayName = bPlayer.Response.SearchResults[i].BungieGlobalDisplayName
						bungoMembershipId = bPlayer.Response.SearchResults[i].DestinyMemberships[j].MembershipId
					}
				}
			}
		}

		if bungoMembershipId == "" && bPlayer.Response.HasMore {
			return CheckStats(usercommand, u, t, cfg, page+1)
		}

		if bungoMembershipId == "" {
			return "Nickname: " + u + " not found in Bungie database"
		}

		req, err = http.NewRequest("GET", "https://kuqw4o1zj5.execute-api.us-west-2.amazonaws.com/api/player/"+bungoMembershipId, nil)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		body, err = io.ReadAll(res.Body)
		err = json.Unmarshal(body, &bAnswer)
		if err != nil {
			e := structures.MyError{Func: "bungiecheck json.Unmarshal &bPlayer", Err: "can't decode bPlayer answer"}
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
				return bungoDisplayName + "#" + strconv.Itoa(bungoNameCode) + " " + usercommand + " N: " + strconv.Itoa(cnormal) + " H: " + strconv.Itoa(chard) + " G: " + strconv.Itoa(cgg)
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
	bodys, err := io.ReadAll(resq.Body)
	err = json.Unmarshal(bodys, &bRaid)
	if err != nil {
		e := structures.MyError{Func: "bungiecheck json.Unmarshal &bRaid", Err: "can't decode bRaid answer"}
		log.Fatalln(e.Error())
	}
	for key, val := range raids {
		if val == strings.ToLower(bRaid.Response.OriginalDisplayProperties.Name) {
			if key == strings.ToLower(usercommand) {
				if bRaid.Response.Matchmaking.RequiresGuardianOath == false {
					if bRaid.Response.Tier == 0 && bRaid.Response.DisplayProperties.Name != "King's Fall: Master" {
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
