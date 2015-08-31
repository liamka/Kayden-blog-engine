package models

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Title string `json:"title"`
  	Subtitle string `json:"subtitle"`
  	URI string `json:"uri"`
  	Description string `json:"description"`
  	Keywords string `json:"keywords"`
  	Mysql string `json:"mysql"`
  	CookieName string `json:"cookieName"`
  	Pass string `json:"pass"`
  	Social   []struct {
	    Url   string `json:"url"`
	    Title string `json:"title"`
	} `json:"social"`
}

type Vars struct {
	Title   	string
	Subtitle 	string
	URI 		string
	Description string
	Keywords 	string
	Social   	[]*Soc
}

type Soc struct {
	Url string
	Title     string
}

var SocData struct{
		Social []*Soc `json:"Social"`
	}

func Conf() Config {
	var с Config
	configFile, _ := ioutil.ReadFile("config.json")
	json.Unmarshal([]byte(configFile), &с) 
	return с
}

func Values(vars Config) Vars {
	varz, _ := json.Marshal(vars)
	json.Unmarshal(varz, &SocData)
	p := Vars{
		Title:   	vars.Title,
		Subtitle:   vars.Subtitle,
		URI:   vars.URI,
		Description:   vars.Description,
		Keywords:   vars.Keywords,
		Social:   	[]*Soc{},
	}
	p.Social = append(p.Social, SocData.Social...)
	return p
}