package configurt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/go-github/v27/github"
	"golang.org/x/oauth2"
)

// Client manages fetching and syncing config values
type Client struct {
	githubClient    *github.Client
	Owner           string
	Repo            string
	Filename        string
	configMap       map[string]interface{}
	RefreshInterval time.Duration // mutex to protect access to configMap
	clientMu        sync.Mutex
}

// NewClient returns a new Configurt client.
func NewClient(username string, githubAccessToken string, configRepo string, configFileName string, refreshInterval time.Duration) *Client {

	client := Client{
		githubClient:    getGitHubClient(githubAccessToken),
		Owner:           username,
		Repo:            configRepo,
		Filename:        configFileName,
		configMap:       make(map[string]interface{}),
		RefreshInterval: refreshInterval,
	}

	client.fetch()
	go client.refresh()
	return &client
}

func getGitHubClient(githubAccessToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

// Uses github client to fetch content from the config file
func (c *Client) fetch() {
	getOpts := github.RepositoryContentGetOptions{Ref: "master"}
	filecontent, _, _, err := c.githubClient.Repositories.GetContents(context.Background(), c.Owner, c.Repo, c.Filename, &getOpts)

	if err != nil {
		log.Panic(err)
	}
	data, decodeErr := base64.StdEncoding.DecodeString(*filecontent.Content)
	if decodeErr != nil {
		log.Panic(decodeErr)
	}
	c.clientMu.Lock()
	defer c.clientMu.Unlock()
	// log.Printf("filecontent: %v", string(data))
	jsonErr := json.Unmarshal(data, &c.configMap)
	if jsonErr != nil {
		log.Panic(jsonErr)
	}
	// log.Printf("map contents: %v", &c.configMap)
}

// Get the config value given a config key, returns a generic interface{} type
func (c *Client) Get(configKey string) interface{} {
	c.clientMu.Lock()
	defer c.clientMu.Unlock()
	return c.configMap[configKey]
}

// A background thread that refreshes the config values at regular intervals
func (c *Client) refresh() {
	if c.RefreshInterval > 0 {
		for {
			time.Sleep(c.RefreshInterval)
			log.Printf("refreshing")
			c.fetch()
		}
	}
}

// GetAsString returns config value as string
func (c *Client) GetAsString(configKey string) string {
	val := c.Get(configKey)
	return val.(string)
}

// GetAsFloat returns config value as float64
func (c *Client) GetAsFloat(configKey string) float64 {
	val := c.Get(configKey)
	return val.(float64)
}

// GetAsFloat returns config value as int
func (c *Client) GetAsInt(configKey string) int {
	val := c.Get(configKey)
	return int(val.(float64))
}

// GetAsStringArray returns config value as string array
func (c *Client) GetAsStringArray(configKey string) []string {
	val := c.Get(configKey)
	array := val.([]interface{})
	stringArray := []string{}
	for _, v := range array {
		stringval := v.(string)
		stringArray = append(stringArray, stringval)
	}
	return stringArray

}

// GetAsFloatArray returns config value as float64 array
func (c *Client) GetAsFloatArray(configKey string) []float64 {
	val := c.Get(configKey)
	array := val.([]interface{})
	floatArray := []float64{}
	for _, v := range array {
		floatv := v.(float64)
		floatArray = append(floatArray, floatv)
	}
	return floatArray
}

func (c *Client) GetAsIntArray(configKey string) []int {
	val := c.Get(configKey)
	array := val.([]interface{})
	intArray := []int{}
	for _, v := range array {
		intv := int(v.(float64))
		intArray = append(intArray, intv)
	}
	return intArray
}
