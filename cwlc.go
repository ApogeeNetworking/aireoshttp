package cwlc

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

// Client properties of an ISE Instance
type Client struct {
	BaseURL  string
	username string
	password string
	IP       string
	http     *http.Client
	cookie   *http.Cookie
}

// New creates an Instance of an ISE client
func New(host, user, pass string, ignoreSSL bool) *Client {
	return &Client{
		BaseURL:  "http://" + host,
		username: user,
		password: pass,
		IP:       host,
		http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: 8 * time.Second,
		},
	}
}

// MakeReq performs an HTTP Request for our Client
func (c *Client) MakeReq(path string) (*http.Response, error) {
	if c.cookie == nil {
		return nil, fmt.Errorf("you must be authenticated first")
	}
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	// Add sessionId Cookie
	req.AddCookie(c.cookie)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	return res, nil
}

// Login establishes a session with a Cisco WLC
func (c *Client) Login() error {
	ep := "/screens/preframeset.html"
	ep += "?username=" + c.username + "&password=" + c.password
	req, err := http.NewRequest("GET", c.BaseURL+ep, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.SetBasicAuth(c.username, c.password)
	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	scookie := res.Cookies()[0]
	re := regexp.MustCompile(`sessionId=(\w+)`)
	f := re.FindAllStringSubmatch(scookie.String(), -1)
	c.cookie = &http.Cookie{
		Name:  "sessionId",
		Value: f[0][1],
	}
	req.AddCookie(c.cookie)
	// "Login One More Time" ... This seems to Activate the Cookie
	res, err = c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	return nil
}

// AP an AP's properties on WLC
type AP struct {
	Name    string `json:"Nm"`
	MacAddr string `json:"Mc"`
}

// GetAps retrieves APs from a Cisco WLC
func (c *Client) GetAps() ([]AP, error) {
	res, err := c.MakeReq("/data/ap-attributes-slot1.html")
	if err != nil {
		return nil, err
	}
	type wlcApRes struct {
		NumOfAps int  `json:"T"`
		Aps      []AP `json:"Data"`
	}
	defer res.Body.Close()
	var data wlcApRes
	json.NewDecoder(res.Body).Decode(&data)
	return data.Aps, nil
}

// ApDetail Cisco Access Point Details
type ApDetail struct {
	Name        string `json:"Name"`
	MacAddr     string `json:"MacAddress"`
	IPAddr      string `json:"IPAddress"`
	RemoteSW    string `json:"CDP"`
	RemoteIntf  string `json:"LLDP"`
	Model       string `json:"Model"`
	Serial      string `json:"SerialNumber"`
	Group       string `json:"Groups"`
	SpeedDuplex string `json:"EthSpeed"`
}

// ApDetailRes ...
type ApDetailRes struct {
	Data ApDetail `json:"Data"`
}

// GetAp retrieves Detailed information of AP from Cisco WLC
func (c *Client) GetApDetails(mac string) (ApDetail, error) {
	ep := "/data/rfdashboard/apview_general.html"
	ep += "?deviceMacAddress=" + mac
	res, err := c.MakeReq(ep)
	if err != nil {
		return ApDetail{}, err
	}
	defer res.Body.Close()
	var ap ApDetailRes
	json.NewDecoder(res.Body).Decode(&ap)
	return ap.Data, nil
}

// WlcInfo properties of a Cisco WLC
type WlcInfo struct {
	Name      string `json:"sysname"`
	Model     string `json:"prodid"`
	Platform  string `json:"platform"`
	Serial    string `json:"serial"`
	Version   string `json:"version"`
	ApInUse   int    `json:"apinuse"`
	ApMax     int    `json:"apmax"`
	WlanCount int    `json:"wlanCount"`
	UpTime    [4]int // Days, Hours, Minutes, Seconds
}

// GetSysInfo retrieves Detailed information of AP from Cisco WLC
func (c *Client) GetSysInfo() (WlcInfo, error) {
	ep := "/data/system_information.html"
	res, err := c.MakeReq(ep)
	if err != nil {
		return WlcInfo{}, err
	}
	defer res.Body.Close()
	var sys WlcInfo
	json.NewDecoder(res.Body).Decode(&sys)
	return sys, nil
}
