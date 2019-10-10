package main

type AcceptMessage struct {
	Type       string        `json:"type"`
	IceServers []interface{} `json:"ice_servers,omitempty"`
}

type RejectMessage struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type AcceptMetadataMessage struct {
	Type       string        `json:"type"`
	Metadata   interface{}   `json:"authz_metadata,omitempty"`
	IceServers []interface{} `json:"ice_servers,omitempty"`
}
