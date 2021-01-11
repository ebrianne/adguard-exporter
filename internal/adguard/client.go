package adguard

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ebrianne/adguard-exporter/internal/metrics"
	"github.com/mitchellh/mapstructure"
)

var (
	port               uint16
	statusURLPattern   = "%s://%s:%d/control/status"
	statsURLPattern    = "%s://%s:%d/control/stats"
	logstatsURLPattern = "%s://%s:%d/control/querylog?limit=%s&response_status=\"all\""
	m                  map[string]int
)

// Client struct is a AdGuard  client to request an instance of a AdGuard  ad blocker.
type Client struct {
	httpClient http.Client
	interval   time.Duration
	logLimit   string
	protocol   string
	hostname   string
	port       uint16
	username   string
	password   string
}

// NewClient method initializes a new AdGuard  client.
func NewClient(protocol, hostname, username, password, adport string, interval time.Duration, logLimit string) *Client {
	if protocol != "http" && protocol != "https" {
		log.Printf("protocol %s is invalid. Must be http or https.", protocol)
		os.Exit(1)
	}

	if httpport == "" {
		port = 80
		if protocol == "https" {
			port = 443
		}
	} else {
		port = uint16(adport)
	}

	return &Client{
		protocol: protocol,
		hostname: hostname,
		port:     port,
		username: username,
		password: password,
		interval: interval,
		logLimit: logLimit,
		httpClient: http.Client{
			Transport: &http.Transport{TLSClientConfig: GetTlsConfig()},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// Scrape method authenticates and retrieves statistics from AdGuard  JSON API
// and then pass them as Prometheus metrics.
func (c *Client) Scrape() {
	for range time.Tick(c.interval) {

		allstats := c.getStatistics()
		//Set the metrics
		c.setMetrics(allstats.status, allstats.stats, allstats.logStats)

		log.Printf("New tick of statistics: %s", allstats.stats.ToString())
	}
}

// Function to set the prometheus metrics
func (c *Client) setMetrics(status *Status, stats *Stats, logstats *LogStats) {
	//Status
	var isRunning int = 0
	if status.Running == true {
		isRunning = 1
	}
	metrics.Running.WithLabelValues(c.hostname).Set(float64(isRunning))

	var isProtected int = 0
	if status.ProtectionEnabled == true {
		isProtected = 1
	}
	metrics.ProtectionEnabled.WithLabelValues(c.hostname).Set(float64(isProtected))

	//Stats
	metrics.AvgProcessingTime.WithLabelValues(c.hostname).Set(float64(stats.AvgProcessingTime))
	metrics.DnsQueries.WithLabelValues(c.hostname).Set(float64(stats.DnsQueries))
	metrics.BlockedFiltering.WithLabelValues(c.hostname).Set(float64(stats.BlockedFiltering))
	metrics.ParentalFiltering.WithLabelValues(c.hostname).Set(float64(stats.ParentalFiltering))
	metrics.SafeBrowsingFiltering.WithLabelValues(c.hostname).Set(float64(stats.SafeBrowsingFiltering))
	metrics.SafeSearchFiltering.WithLabelValues(c.hostname).Set(float64(stats.SafeSearchFiltering))

	for l := range stats.TopQueries {
		for domain, value := range stats.TopQueries[l] {
			metrics.TopQueries.WithLabelValues(c.hostname, domain).Set(float64(value))
		}
	}

	for l := range stats.TopBlocked {
		for domain, value := range stats.TopBlocked[l] {
			metrics.TopBlocked.WithLabelValues(c.hostname, domain).Set(float64(value))
		}
	}

	for l := range stats.TopClients {
		for source, value := range stats.TopClients[l] {
			metrics.TopClients.WithLabelValues(c.hostname, source).Set(float64(value))
		}
	}

	//LogQuery
	m = make(map[string]int)
	logdata := logstats.Data
	for i := range logdata {
		dnsanswer := logdata[i].Answer
		if dnsanswer != nil && len(dnsanswer) > 0 {
			for j := range dnsanswer {
				var dnsType string
				//Check the type of dnsanswer[j].Value, if string leave it be, otherwise get back the object to get the correct DNS type
				switch v := dnsanswer[j].Value.(type) {
				case string:
					dnsType = dnsanswer[j].Type
					m[dnsType] += 1
				case map[string]interface{}:
					var dns65 Type65
					mapstructure.Decode(v, &dns65)
					dnsType = "TYPE" + strconv.Itoa(dns65.Hdr.Rrtype)
					m[dnsType] += 1
				default:
					continue
				}
			}
		}
	}

	for key, value := range m {
		metrics.QueryTypes.WithLabelValues(c.hostname, key).Set(float64(value))
	}

	//clear the map
	for k := range m {
		delete(m, k)
	}
}

// Function to get the general stats
func (c *Client) getStatistics() *AllStats {

	var status Status
	statusURL := fmt.Sprintf(statusURLPattern, c.protocol, c.hostname, c.port)
	body := c.MakeRequest(statusURL)
	err := json.Unmarshal(body, &status)
	if err != nil {
		log.Println("Unable to unmarshal Adguard log statistics to log statistics struct model", err)
	}

	var stats Stats
	statsURL := fmt.Sprintf(statsURLPattern, c.protocol, c.hostname, c.port)
	body = c.MakeRequest(statsURL)
	err = json.Unmarshal(body, &stats)
	if err != nil {
		log.Println("Unable to unmarshal Adguard statistics to statistics struct model", err)
	}

	var logstats LogStats
	logstatsURL := fmt.Sprintf(logstatsURLPattern, c.protocol, c.hostname, c.port, c.logLimit)
	body = c.MakeRequest(logstatsURL)
	err = json.Unmarshal(body, &logstats)
	if err != nil {
		log.Println("Unable to unmarshal Adguard log statistics to log statistics struct model", err)
	}

	var allstats AllStats
	allstats.status = &status
	allstats.stats = &stats
	allstats.logStats = &logstats

	return &allstats
}

func (c *Client) MakeRequest(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("An error has occurred when creating HTTP statistics request", err)
	}

	req.Host = c.hostname
	req.Header.Add("User-Agent", "Mozilla/5.0")

	if c.isUsingPassword() {
		c.authenticateRequest(req)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal("An error has occurred during login to Adguard", err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("An error occured in the request, Status Code ", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Unable to read Adguard statistics HTTP response", err)
	}

	return body
}

func (c *Client) isUsingPassword() bool {
	return len(c.password) > 0
}

func (c *Client) authenticateRequest(req *http.Request) {
	req.SetBasicAuth(c.username, c.password)
}

func GetTlsConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}
