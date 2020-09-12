package twitch

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetTopStreamNames(amount int) ([]string, error) {
	url := "https://gql.twitch.tv/gql"
	method := "POST"

	query := fmt.Sprintf("{\"query\":\"query {streams(first: %d) {edges {node {broadcaster {channel {name}}}}}}\",\"variables\":{}}", amount)
	payload := strings.NewReader(query)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Client-ID", "kimne78kx3ncx6brgo4mv6wki5h1ko")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	// parse json response
	json, err := gabs.ParseJSON(body)
	if err != nil {
		return nil, err
	}

	var channelIds []string

	for _, child := range json.Path("data.streams").Search("edges").Children() {
		channelIds = append(channelIds, child.Path("node.broadcaster.channel.name").Data().(string))
	}

	return channelIds, nil
}
