package model

import (
	"encoding/json"
	"io"
)

type PostSearchMatches map[string][]string

type PostSearchResults struct {
	*PostList
	Matches PostSearchMatches `json:"matches"`
}

func MakePostSearchResults(posts *PostList, matches PostSearchMatches) *PostSearchResults {
	return &PostSearchResults{
		posts,
		matches,
	}
}

func (o *PostSearchResults) ToJson() string {
	copy := *o
	copy.PostList.StripActionIntegrations()
	b, err := json.Marshal(&copy)
	if err != nil {
		return ""
	}
	return string(b)
}

func PostSearchResultsFromJson(data io.Reader) *PostSearchResults {
	var o *PostSearchResults
	json.NewDecoder(data).Decode(&o)
	return o
}
