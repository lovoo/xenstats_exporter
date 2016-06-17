package main

import (
	"crypto/tls"
	"net/http"

	"github.com/nilshell/xmlrpc"
	xsclient "github.com/xenserver/go-xenserver-client"
)

// XenAPIClient -
type XenAPIClient struct {
	xsclient.XenAPIClient
}

// Driver hold information about the Xenserver and the credentials
type Driver struct {
	Server       string
	Username     string
	Password     string
	xenAPIClient *XenAPIClient
}

// NewDriver Creates a new Driver
func NewDriver() *Driver {
	return &Driver{}
}

// ApiObject of type ..
type ApiObject xsclient.XenAPIObject

// NewXenAPIClient -
func NewXenAPIClient(host, username, password string) (c XenAPIClient) {
	c.Host = host
	c.Url = "https://" + host
	c.Username = username
	c.Password = password
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	c.RPC, _ = xmlrpc.NewClient(c.Url, tr)
	return
}

// GetXenAPIClient returns
func (d *Driver) GetXenAPIClient() (*XenAPIClient, error) {
	if d.xenAPIClient == nil {
		c := NewXenAPIClient(d.Server, d.Username, d.Password)
		if err := c.Login(); err != nil {
			return nil, err
		}
		d.xenAPIClient = &c
	}
	return d.xenAPIClient, nil
}

// GetSpecificValue -
func (d *Driver) GetSpecificValue(apikey string, params string) (interface{}, error) {
	result := xsclient.APIResult{}
	err := d.xenAPIClient.APICall(&result, apikey, params)
	if err != nil {
		return result.Value.(interface{}), err
	}
	return result.Value.(interface{}), err
}

// GetMultiValues -
func (d *Driver) GetMultiValues(apikey string, params ...string) (apiObjects []*ApiObject, err error) {
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
