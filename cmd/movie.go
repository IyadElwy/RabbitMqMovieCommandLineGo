package main

type MovieDataList struct {
	Results []MovieData `json:"results"`
}

type MovieData struct {
	Title    string `json:"title"`
	Overview string `json:"overview"`
}
