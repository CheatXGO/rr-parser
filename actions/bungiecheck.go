package actions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"main/structures"
	"net/http"
	"strconv"
	"strings"
)

var bPlayer structures.BungoSearch
var bAnswer structures.BungoAnswer
var bRaid structures.RaidAnswer
var cnormal, chard, cgg int

func CheckStats(usercommand, u string, t bool, cfg structures.Config) string {
	nck := u
	if !t {
		return u //user not registered in DB
	} else {
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
					req, err := http.NewRequest("GET", "https://www.bungie.net/Platform/Destiny2/Manifest/DestinyActivityDefinition/"+strconv.FormatInt(bAnswer.Response.Activities[i].ActivityHash, 10), nil)
					if err != nil {
						log.Fatalln(err)
					}
					req.Header.Set("X-API-Key", cfg.Bapi)
					res, err := client.Do(req)
					if err != nil {
						log.Fatalln(err)
					}
					body, err := ioutil.ReadAll(res.Body)
					err = json.Unmarshal(body, &bRaid)
					if err != nil {
						e := structures.MyError{Fun: "bungiecheck json.Unmarshal &bRaid", Err: "can't decode bRaid answer"}
						log.Fatalln(e.Error())
					}
					for key, val := range raids {
						if val == strings.ToLower(bRaid.Response.OriginalDisplayProperties.Name) {
							if key == strings.ToLower(usercommand) {
								if bRaid.Response.Matchmaking.RequiresGuardianOath == false {
									if bRaid.Response.Tier == 0 {
										cnormal = cnormal + bAnswer.Response.Activities[i].Values.Clears
									} else {
										chard = chard + bAnswer.Response.Activities[i].Values.Clears
									}
								} else {
									cgg = cgg + bAnswer.Response.Activities[i].Values.Clears
								}
							}
						}
					}
				}
			}
		} else {
			return "Please fill right raid abbreviation after /my"
		}
	}
	return nck + " " + usercommand + " N: " + strconv.Itoa(cnormal) + " H: " + strconv.Itoa(chard) + " G: " + strconv.Itoa(cgg)
}
