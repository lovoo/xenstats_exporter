package main

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/nilshell/xmlrpc"
	xsclient "github.com/xenserver/go-xenserver-client"
)

// ApiCaller hold information about the Xenserver and the credentials
type ApiCaller struct {
	Server       string
	Username     string
	Password     string
	xenAPIClient *xsclient.XenAPIClient
}

// NewApiCaller Creates a new ApiCaller
func NewApiCaller(host, username, password string) *ApiCaller {
	return &ApiCaller{
		Server:   host,
		Username: username,
		Password: password,
	}
}

// ApiObject of type ..
type ApiObject xsclient.XenAPIObject

// NewXenAPIClient -
func (d *ApiCaller) newXenAPIClient() (c xsclient.XenAPIClient, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	url, err := url.Parse("https://" + d.Server)
	if err != nil {
		return c, err
	}

	rpcClient, err := xmlrpc.NewClient(url.String(), tr)

	return xsclient.XenAPIClient{
		Host:     d.Server,
		Url:      url.String(),
		RPC:      rpcClient,
		Username: d.Username,
		Password: d.Password,
	}, err
}

// GetXenAPIClient returns
func (d *ApiCaller) GetXenAPIClient() (*xsclient.XenAPIClient, error) {
	var err error
	if d.xenAPIClient == nil {
		c, err := d.newXenAPIClient()
		if err != nil {
			return nil, err
		}
		if err := c.Login(); err != nil {
			return nil, err
		}
		d.xenAPIClient = &c
	}
	return d.xenAPIClient, err
}

// GetSpecificValue -
func (d *ApiCaller) GetSpecificValue(apikey string, params string) (interface{}, error) {
	result := xsclient.APIResult{}
	err := d.xenAPIClient.APICall(&result, apikey, params)
	return result.Value, err
}

// GetMultiValues -
func (d *ApiCaller) GetMultiValues(apikey string, params ...string) (apiObjects []*ApiObject, err error) {
	result := xsclient.APIResult{}

	if len(params) > 0 {
		err = d.xenAPIClient.APICall(&result, apikey, params[0])
	} else {
		err = d.xenAPIClient.APICall(&result, apikey)
	}

	if err != nil {
		return apiObjects, err
	}

	for _, elem := range result.Value.([]interface{}) {
		apiObject := new(ApiObject)
		apiObject.Ref = elem.(string)
		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects, err
}
