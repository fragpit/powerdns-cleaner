package main

type PowerDNSZone struct {
	Account          string   `json:"account"`
	APIRectify       bool     `json:"api_rectify"`
	Catalog          string   `json:"catalog"`
	DNSSEC           bool     `json:"dnssec"`
	EditedSerial     int64    `json:"edited_serial"`
	ID               string   `json:"id"`
	Kind             string   `json:"kind"`
	LastCheck        int      `json:"last_check"`
	MasterTSIGKeyIDs []string `json:"master_tsig_key_ids"`
	Masters          []string `json:"masters"`
	Name             string   `json:"name"`
	NotifiedSerial   int      `json:"notified_serial"`
	NSEC3Narrow      bool     `json:"nsec3narrow"`
	NSEC3Param       string   `json:"nsec3param"`
	RRSets           []RRSet  `json:"rrsets"`
	Serial           int64    `json:"serial"`
	SlaveTSIGKeyIDs  []string `json:"slave_tsig_key_ids"`
	SOAEdit          string   `json:"soa_edit"`
	SOAEditAPI       string   `json:"soa_edit_api"`
	URL              string   `json:"url"`
}

type RRSet struct {
	Comments []Comment `json:"comments"`
	Name     string    `json:"name"`
	Records  []Record  `json:"records"`
	TTL      int       `json:"ttl"`
	Type     string    `json:"type"`
}

type Comment struct {
	Content    string `json:"content"`
	Account    string `json:"account"`
	ModifiedAt int    `json:"modified_at"`
}

type Record struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}
