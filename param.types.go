package main

import "encoding/json"

type Param struct {
	Command  string `json:"command"`
	Verbose  bool   `json:"verbose"`
	DataPath string `json:"data_path"`

	FetchParam *FetchParam `json:"fetch_param,omitempty"`
	ListParam  *ListParam  `json:"list_param,omitempty"`
	CopyParam  *CopyParam  `json:"copy_param,omitempty"`
	ItermParam *ItermParam `json:"iterm_param,omitempty"`
}

func (p *Param) toJSON() string {
	b, _ := json.Marshal(p)
	return string(b)
}

type FetchParam struct {
	UseSample bool `json:"use_sample"`
}

type ListParam struct {
	Hostgroup string `json:"hostgroup"`
	Type      string `json:"type"`
}

type CopyParam struct {
	DestinationPath string `json:"dest_path"`
}

type ItermParam struct {
	Hostgroup  string `json:"hostgroup"`
	NumOfHost  int    `json:"num_of_host"`
	Username   string `json:"username"`
	IsPickHost bool   `json:"is_pick_host"`
}
