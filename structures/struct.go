package structures

import "fmt"

type MyError struct {
	Err string
	Fun string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("Func: %s Error: %s", e.Fun, e.Err)
}

type Config struct { // for data from config.json
	TelegramBotToken string `json:"TelegramBotToken"` // tg bot api token
	TgBotIP          string `json:"TgBotIP"`          // bot IP for webhook
	TgBotPort        string `json:"TgBotPort"`        // bot port for webhook
	Pip              string `json:"pip"`              // postgres IP
	Pport            string `json:"pport"`            // postgres port
	Pdb              string `json:"pdb"`              // postgres database
	Plg              string `json:"plg"`              // postgres login
	Ppw              string `json:"ppw"`              // postgres password
	Bapi             string `json:"bapi"`             //bungie API key
}

type BotMessage struct {
	Message struct {
		MessageId int `json:"message_id"`
		From      struct {
			Username string `json:"username"`
			Id       int    `json:"id"`
			IsBot    bool   `json:"is_bot"`
		}
		Chat struct {
			Id int `json:"id"`
		}
		Date int    `json:"date"`
		Text string `json:"text"`
	}
}

type SendMessage struct {
	Id   int    `json:"chat_id"`
	Text string `json:"text"`
}

type TelegUser struct { // user operations reg upd with DB
	Id       int
	Idtelegr int
	Tgnick   string
	D2nick   string
}

type BungoSearch struct {
	Response struct {
		SearchResults []struct {
			BungieGlobalDisplayName     string `json:"bungieGlobalDisplayName"`
			BungieGlobalDisplayNameCode int    `json:"bungieGlobalDisplayNameCode"`
			BungieNetMembershipId       string `json:"bungieNetMembershipId"`
			DestinyMemberships          []struct {
				CrossSaveOverride           int    `json:"crossSaveOverride"`
				ApplicableMembershipTypes   []int  `json:"applicableMembershipTypes"`
				IsPublic                    bool   `json:"isPublic"`
				MembershipType              int    `json:"membershipType"`
				MembershipId                string `json:"membershipId"`
				DisplayName                 string `json:"displayName"`
				BungieGlobalDisplayName     string `json:"bungieGlobalDisplayName"`
				BungieGlobalDisplayNameCode int    `json:"bungieGlobalDisplayNameCode"`
			} `json:"destinyMemberships"`
		} `json:"searchResults"`
		Page    int  `json:"page"`
		HasMore bool `json:"hasMore"`
	} `json:"Response"`
}

type BungoAnswer struct {
	Response struct {
		MembershipId string `json:"membershipId"`
		Activities   []struct {
			ActivityHash int64 `json:"activityHash"`
			Values       struct {
				Clears           int `json:"clears"`
				FullClears       int `json:"fullClears"`
				SherpaCount      int `json:"sherpaCount"`
				FastestFullClear struct {
					Value      int    `json:"value"`
					InstanceId string `json:"instanceId"`
				} `json:"fastestFullClear"`
				FlawlessDetails struct {
					InstanceId   string `json:"instanceId"`
					AccountCount int    `json:"accountCount"`
				} `json:"flawlessDetails,omitempty"`
				FlawlessActivities []struct {
					InstanceId   string `json:"instanceId"`
					AccountCount int    `json:"accountCount"`
				} `json:"flawlessActivities,omitempty"`
			} `json:"values"`
		} `json:"activities"`
		SpeedRank struct {
			Value   int    `json:"value"`
			Tier    string `json:"tier"`
			Subtier string `json:"subtier"`
		} `json:"speedRank"`
		ClearsRank struct {
			Value   int    `json:"value"`
			Tier    string `json:"tier"`
			Subtier string `json:"subtier"`
		} `json:"clearsRank"`
	} `json:"response"`
}

type JsonData struct {
	ActivityHash struct {
		OriginalDisplayProperties struct {
			Name string `json:"name"`
		} `json:"originalDisplayProperties"`
		Tier        int `json:"tier"`
		Matchmaking struct {
			RequiresGuardianOath bool `json:"requiresGuardianOath"`
		} `json:"matchmaking"`
	}
}

type RaidAnswer struct {
	Response struct {
		OriginalDisplayProperties struct {
			Name string `json:"name"`
		} `json:"originalDisplayProperties"`
		Tier        int `json:"tier"`
		Matchmaking struct {
			RequiresGuardianOath bool `json:"requiresGuardianOath"`
		} `json:"matchmaking"`
	} `json:"Response"`
}
