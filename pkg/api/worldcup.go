package api

import (
	"encoding/json"
	"net/http"
	"time"
)

const espnBase = "https://site.api.espn.com/apis"

var espnClient = &http.Client{Timeout: 10 * time.Second}

type wcTeam struct {
	Name string `json:"name"`
	Abbr string `json:"abbr"`
	Logo string `json:"logo"`
}

type wcFixture struct {
	ID        string `json:"id"`
	Date      string `json:"date"`
	Status    string `json:"status"`
	Detail    string `json:"detail"`
	Venue     string `json:"venue"`
	HomeTeam  wcTeam `json:"home_team"`
	AwayTeam  wcTeam `json:"away_team"`
	HomeScore string `json:"home_score"`
	AwayScore string `json:"away_score"`
}

type wcStandingEntry struct {
	Team wcTeam `json:"team"`
	GP   string `json:"gp"`
	W    string `json:"w"`
	D    string `json:"d"`
	L    string `json:"l"`
	GF   string `json:"gf"`
	GA   string `json:"ga"`
	GD   string `json:"gd"`
	Pts  string `json:"pts"`
}

type wcGroup struct {
	Name    string            `json:"name"`
	Entries []wcStandingEntry `json:"entries"`
}

func WorldCupFixtures() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := espnClient.Get(espnBase + "/site/v2/sports/soccer/fifa.world/scoreboard?dates=20260611-20261231&limit=200")
		if err != nil {
			writeError(w, http.StatusBadGateway, "failed to fetch fixtures")
			return
		}
		defer resp.Body.Close()

		var raw struct {
			Events []struct {
				ID           string `json:"id"`
				Date         string `json:"date"`
				Competitions []struct {
					Status struct {
						Type struct {
							State  string `json:"state"`
							Detail string `json:"detail"`
						} `json:"type"`
					} `json:"status"`
					Venue struct {
						FullName string `json:"fullName"`
					} `json:"venue"`
					Competitors []struct {
						HomeAway string `json:"homeAway"`
						Score    string `json:"score"`
						Team     struct {
							DisplayName  string `json:"displayName"`
							Abbreviation string `json:"abbreviation"`
							Logo         string `json:"logo"`
						} `json:"team"`
					} `json:"competitors"`
				} `json:"competitions"`
			} `json:"events"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			writeError(w, http.StatusBadGateway, "failed to parse fixtures")
			return
		}

		fixtures := make([]wcFixture, 0, len(raw.Events))
		for _, e := range raw.Events {
			if len(e.Competitions) == 0 {
				continue
			}
			comp := e.Competitions[0]
			f := wcFixture{
				ID:     e.ID,
				Date:   e.Date,
				Status: comp.Status.Type.State,
				Detail: comp.Status.Type.Detail,
				Venue:  comp.Venue.FullName,
			}
			for _, c := range comp.Competitors {
				t := wcTeam{
					Name: c.Team.DisplayName,
					Abbr: c.Team.Abbreviation,
					Logo: c.Team.Logo,
				}
				if c.HomeAway == "home" {
					f.HomeTeam = t
					f.HomeScore = c.Score
				} else {
					f.AwayTeam = t
					f.AwayScore = c.Score
				}
			}
			fixtures = append(fixtures, f)
		}

		writeJSON(w, http.StatusOK, fixtures)
	}
}

func WorldCupStandings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := espnClient.Get(espnBase + "/v2/sports/soccer/fifa.world/standings")
		if err != nil {
			writeError(w, http.StatusBadGateway, "failed to fetch standings")
			return
		}
		defer resp.Body.Close()

		var raw struct {
			Children []struct {
				Name      string `json:"name"`
				Standings struct {
					Entries []struct {
						Team struct {
							DisplayName  string `json:"displayName"`
							Abbreviation string `json:"abbreviation"`
							Logos        []struct {
								Href string `json:"href"`
							} `json:"logos"`
						} `json:"team"`
						Stats []struct {
							Abbreviation string `json:"abbreviation"`
							DisplayValue string `json:"displayValue"`
						} `json:"stats"`
					} `json:"entries"`
				} `json:"standings"`
			} `json:"children"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			writeError(w, http.StatusBadGateway, "failed to parse standings")
			return
		}

		groups := make([]wcGroup, 0, len(raw.Children))
		for _, child := range raw.Children {
			g := wcGroup{Name: child.Name}
			for _, e := range child.Standings.Entries {
				logo := ""
				if len(e.Team.Logos) > 0 {
					logo = e.Team.Logos[0].Href
				}
				entry := wcStandingEntry{
					Team: wcTeam{
						Name: e.Team.DisplayName,
						Abbr: e.Team.Abbreviation,
						Logo: logo,
					},
				}
				for _, s := range e.Stats {
					switch s.Abbreviation {
					case "GP":
						entry.GP = s.DisplayValue
					case "W":
						entry.W = s.DisplayValue
					case "D":
						entry.D = s.DisplayValue
					case "L":
						entry.L = s.DisplayValue
					case "F":
						entry.GF = s.DisplayValue
					case "A":
						entry.GA = s.DisplayValue
					case "GD":
						entry.GD = s.DisplayValue
					case "P":
						entry.Pts = s.DisplayValue
					}
				}
				g.Entries = append(g.Entries, entry)
			}
			groups = append(groups, g)
		}

		writeJSON(w, http.StatusOK, groups)
	}
}
