package arex

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Auth Represents authentication helpers provided by Postman
type Auth struct {

	// The attributes for API Key Authentication.
	Apikey []*AuthAttribute `json:"apikey,omitempty"`

	// The attributes for [AWS Auth](http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html).
	Awsv4 []*AuthAttribute `json:"awsv4,omitempty"`

	// The attributes for [Basic Authentication](https://en.wikipedia.org/wiki/Basic_access_authentication).
	Basic []*AuthAttribute `json:"basic,omitempty"`

	// The helper attributes for [Bearer Token Authentication](https://tools.ietf.org/html/rfc6750)
	Bearer []*AuthAttribute `json:"bearer,omitempty"`

	// The attributes for [Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication).
	Digest []*AuthAttribute `json:"digest,omitempty"`

	// The attributes for [Akamai EdgeGrid Authentication](https://developer.akamai.com/legacy/introduction/Client_Auth.html).
	Edgegrid []*AuthAttribute `json:"edgegrid,omitempty"`

	// The attributes for [Hawk Authentication](https://github.com/hueniverse/hawk)
	Hawk   []*AuthAttribute `json:"hawk,omitempty"`
	Noauth interface{}      `json:"noauth,omitempty"`

	// The attributes for [NTLM Authentication](https://msdn.microsoft.com/en-us/library/cc237488.aspx)
	Ntlm []*AuthAttribute `json:"ntlm,omitempty"`

	// The attributes for [OAuth2](https://oauth.net/1/)
	Oauth1 []*AuthAttribute `json:"oauth1,omitempty"`

	// Helper attributes for [OAuth2](https://oauth.net/2/)
	Oauth2 []*AuthAttribute `json:"oauth2,omitempty"`
	Type   string           `json:"type"`
}

// AuthAttribute Represents an attribute for any authorization method provided by Postman. For example `username` and `password` are set as auth attributes for Basic Authentication method.
type AuthAttribute struct {
	Key   string      `json:"key"`
	Type  string      `json:"type,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// Cert An object containing path to file certificate, on the file system
type Cert struct {

	// The path to file containing key for certificate, on the file system
	Src interface{} `json:"src,omitempty"`
}

// Certificate A representation of an ssl certificate
type Certificate struct {

	// An object containing path to file certificate, on the file system
	Cert *Cert `json:"cert,omitempty"`

	// An object containing path to file containing private key, on the file system
	Key *Key `json:"key,omitempty"`

	// A list of Url match pattern strings, to identify Urls this certificate can be used for.
	Matches []string `json:"matches,omitempty"`

	// A name for the certificate for user reference
	Name string `json:"name,omitempty"`

	// The passphrase for the certificate
	Passphrase string `json:"passphrase,omitempty"`
}

// Cookie A Cookie, that follows the [Google Chrome format](https://developer.chrome.com/extensions/cookies)
type Cookie struct {

	// The domain for which this cookie is valid.
	Domain string `json:"domain"`

	// When the cookie expires.
	Expires interface{} `json:"expires,omitempty"`

	// Custom attributes for a cookie go here, such as the [Priority Field](https://code.google.com/p/chromium/issues/detail?id=232693)
	Extensions []interface{} `json:"extensions,omitempty"`

	// True if the cookie is a host-only cookie. (i.e. a request's URL domain must exactly match the domain of the cookie).
	HostOnly bool `json:"hostOnly,omitempty"`

	// Indicates if this cookie is HTTP Only. (if True, the cookie is inaccessible to client-side scripts)
	HTTPOnly bool   `json:"httpOnly,omitempty"`
	MaxAge   string `json:"maxAge,omitempty"`

	// This is the name of the Cookie.
	Name string `json:"name,omitempty"`

	// The path associated with the Cookie.
	Path string `json:"path"`

	// Indicates if the 'secure' flag is set on the Cookie, meaning that it is transmitted over secure connections only. (typically HTTPS)
	Secure bool `json:"secure,omitempty"`

	// True if the cookie is a session cookie.
	Session bool `json:"session,omitempty"`

	// The value of the Cookie.
	Value string `json:"value,omitempty"`
}

// Event Defines a script associated with an associated event name
type Event struct {

	// Indicates whether the event is disabled. If absent, the event is assumed to be enabled.
	Disabled bool `json:"disabled,omitempty"`

	// A unique identifier for the enclosing event.
	ID string `json:"id,omitempty"`

	// Can be set to `test` or `prerequest` for test scripts or pre-request scripts respectively.
	Listen string  `json:"listen"`
	Script *Script `json:"script,omitempty"`
}

// Header Represents a single HTTP Header
type Header struct {
	Description interface{} `json:"description,omitempty"`

	// If set to true, the current header will not be sent with requests.
	Disabled bool `json:"disabled,omitempty"`

	// This holds the LHS of the HTTP Header, e.g ``Content-Type`` or ``X-Custom-Header``
	Key string `json:"key"`

	// The value (or the RHS) of the Header is stored in this field.
	Value string `json:"value"`
}

// Info Detailed description of the info block
type Info struct {
	Description interface{} `json:"description,omitempty"`

	// A collection's friendly name is defined by this field. You would want to set this field to a value that would allow you to easily identify this collection among a bunch of other collections, as such outlining its usage or content.
	Name string `json:"name"`

	// Every collection is identified by the unique value of this field. The value of this field is usually easiest to generate using a UID generator function. If you already have a collection, it is recommended that you maintain the same id since changing the id usually implies that is a different collection than it was originally.
	//  *Note: This field exists for compatibility reasons with Collection Format V1.*
	PostmanID string `json:"_postman_id,omitempty"`

	// This should ideally hold a link to the Postman schema that is used to validate this collection. E.g: https://schema.getpostman.com/collection/v1
	Schema  string      `json:"schema"`
	Version interface{} `json:"version,omitempty"`
}

// Item Items are entities which contain an actual HTTP request, and sample responses attached to it.
type Item struct {
	Description interface{} `json:"description,omitempty"`
	Event       []*Event    `json:"event,omitempty"`

	// A unique ID that is used to identify collections internally
	ID string `json:"id,omitempty"`

	// A human readable identifier for the current item.
	Name                    string                   `json:"name,omitempty"`
	ProtocolProfileBehavior *ProtocolProfileBehavior `json:"protocolProfileBehavior,omitempty"`
	Request                 Request                  `json:"request"`
	Response                []*Response              `json:"response,omitempty"`
	Variable                []*Variable              `json:"variable,omitempty"`
}

// ItemGroup One of the primary goals of Postman is to organize the development of APIs. To this end, it is necessary to be able to group requests together. This can be achived using 'Folders'. A folder just is an ordered set of requests.
type ItemGroup struct {
	Auth        interface{} `json:"auth,omitempty"`
	Description interface{} `json:"description,omitempty"`
	Event       []*Event    `json:"event,omitempty"`

	// Items are entities which contain an actual HTTP request, and sample responses attached to it. Folders may contain many items.
	// Item []interface{} `json:"item"`
	Item []Item `json:"item"`

	// A folder's friendly name is defined by this field. You would want to set this field to a value that would allow you to easily identify this folder.
	Name                    string                   `json:"name,omitempty"`
	ProtocolProfileBehavior *ProtocolProfileBehavior `json:"protocolProfileBehavior,omitempty"`
	Variable                []*Variable              `json:"variable,omitempty"`
}

// Key An object containing path to file containing private key, on the file system
type Key struct {

	// The path to file containing key for certificate, on the file system
	Src interface{} `json:"src,omitempty"`
}

// ProtocolProfileBehavior Set of configurations used to alter the usual behavior of sending the request
type ProtocolProfileBehavior struct {
}

// ProxyConfig Using the Proxy, you can configure your custom proxy into the postman for particular url match
type ProxyConfig struct {

	// When set to true, ignores this proxy configuration entity
	Disabled bool `json:"disabled,omitempty"`

	// The proxy server host
	Host string `json:"host,omitempty"`

	// The Url match for which the proxy config is defined
	Match string `json:"match,omitempty"`

	// The proxy server port
	Port int `json:"port,omitempty"`

	// The tunneling details for the proxy config
	Tunnel bool `json:"tunnel,omitempty"`
}

// Request A request represents an HTTP request. If a string, the string is assumed to be the request URL and the method is assumed to be 'GET'.
type Request struct {
	Auth        Auth         `json:"auth,omitempty"`
	Body        Body         `json:"body,omitempty"`
	Certificate *Certificate `json:"certificate,omitempty"`
	Description interface{}  `json:"description,omitempty"`
	Header      []Header     `json:"header,omitempty"`
	Method      interface{}  `json:"method,omitempty"`
	Proxy       *ProxyConfig `json:"proxy,omitempty"`
	URL         interface{}  `json:"url,omitempty"`
}

// Response A response represents an HTTP response.
type Response struct {

	// The raw text of the response.
	Body interface{} `json:"body,omitempty"`

	// The numerical response code, example: 200, 201, 404, etc.
	Code   int         `json:"code,omitempty"`
	Cookie []*Cookie   `json:"cookie,omitempty"`
	Header interface{} `json:"header,omitempty"`

	// A unique, user defined identifier that can  be used to refer to this response from requests.
	ID              string      `json:"id,omitempty"`
	OriginalRequest interface{} `json:"originalRequest,omitempty"`

	// The time taken by the request to complete. If a number, the unit is milliseconds. If the response is manually created, this can be set to `null`.
	ResponseTime interface{} `json:"responseTime,omitempty"`

	// The response status, e.g: '200 OK'
	Status string `json:"status,omitempty"`

	// Set of timing information related to request and response in milliseconds
	Timings interface{} `json:"timings,omitempty"`
}

// Body This field contains the data usually contained in the request body.
type Body struct {

	// When set to true, prevents request body from being sent.
	Disabled bool          `json:"disabled,omitempty"`
	File     *File         `json:"file,omitempty"`
	Formdata []interface{} `json:"formdata,omitempty"`
	Graphql  interface{}   `json:"graphql,omitempty"`

	// Postman stores the type of data associated with this request in this field.
	Mode interface{} `json:"mode,omitempty"`

	// Additional configurations and options set for various body modes.
	Options    interface{}            `json:"options,omitempty"`
	Raw        string                 `json:"raw,omitempty"`
	Urlencoded []*URLEncodedParameter `json:"urlencoded,omitempty"`
}

// File some kind
type File struct {
	Content string      `json:"content,omitempty"`
	Src     interface{} `json:"src,omitempty"`
}

// URLEncodedParameter todo
type URLEncodedParameter struct {
	Description interface{} `json:"description,omitempty"`
	Disabled    bool        `json:"disabled,omitempty"`
	Key         string      `json:"key"`
	Value       string      `json:"value,omitempty"`
}

// Root ???????????????
type Root struct {
	Auth  interface{} `json:"auth,omitempty"`
	Event []*Event    `json:"event,omitempty"`
	Info  *Info       `json:"info"`

	// Items are the basic unit for a Postman collection. You can think of them as corresponding to a single API endpoint. Each Item has one request and may have multiple API responses associated with it.
	Item                    []interface{}            `json:"item"`
	ProtocolProfileBehavior *ProtocolProfileBehavior `json:"protocolProfileBehavior,omitempty"`
	Variable                []*Variable              `json:"variable,omitempty"`
}

// Script A script is a snippet of Javascript code that can be used to to perform setup or teardown operations on a particular response.
type Script struct {
	Exec interface{} `json:"exec,omitempty"`

	// A unique, user defined identifier that can  be used to refer to this script from requests.
	ID string `json:"id,omitempty"`

	// Script name
	Name string      `json:"name,omitempty"`
	Src  interface{} `json:"src,omitempty"`

	// Type of the script. E.g: 'text/javascript'
	Type string `json:"type,omitempty"`
}

// Variable Using variables in your Postman requests eliminates the need to duplicate requests, which can save a lot of time. Variables can be defined, and referenced to from any part of a request.
type Variable struct {
	Description interface{} `json:"description,omitempty"`
	Disabled    bool        `json:"disabled,omitempty"`

	// A variable ID is a unique user-defined value that identifies the variable within a collection. In traditional terms, this would be a variable name.
	ID string `json:"id,omitempty"`

	// A variable key is a human friendly value that identifies the variable within a collection. In traditional terms, this would be a variable name.
	Key string `json:"key,omitempty"`

	// Variable name
	Name string `json:"name,omitempty"`

	// When set to true, indicates that this variable has been set by Postman
	System bool `json:"system,omitempty"`

	// A variable may have multiple types. This field specifies the type of the variable.
	Type string `json:"type,omitempty"`

	// The value that a variable holds in this collection. Ultimately, the variables will be replaced by this value, when say running a set of requests from a collection
	Value interface{} `json:"value,omitempty"`
}

//MarshalJSON marshal to json
func (strct *Auth) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "apikey" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"apikey\": ")
	tmp, err := json.Marshal(strct.Apikey)
	if err != nil {
		return nil, err
	}
	buf.Write(tmp)

	comma = true
	// Marshal the "awsv4" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"awsv4\": ")
	if tmp, err := json.Marshal(strct.Awsv4); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "basic" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"basic\": ")
	if tmp, err := json.Marshal(strct.Basic); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "bearer" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"bearer\": ")
	if tmp, err := json.Marshal(strct.Bearer); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "digest" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"digest\": ")
	if tmp, err := json.Marshal(strct.Digest); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "edgegrid" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"edgegrid\": ")
	if tmp, err := json.Marshal(strct.Edgegrid); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "hawk" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"hawk\": ")
	if tmp, err := json.Marshal(strct.Hawk); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "noauth" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"noauth\": ")
	if tmp, err := json.Marshal(strct.Noauth); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "ntlm" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"ntlm\": ")
	if tmp, err := json.Marshal(strct.Ntlm); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "oauth1" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"oauth1\": ")
	if tmp, err := json.Marshal(strct.Oauth1); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "oauth2" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"oauth2\": ")
	if tmp, err := json.Marshal(strct.Oauth2); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Type" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "type" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"type\": ")
	if tmp, err := json.Marshal(strct.Type); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON sss
func (strct *Auth) UnmarshalJSON(b []byte) error {
	typeReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "apikey":
			if err := json.Unmarshal([]byte(v), &strct.Apikey); err != nil {
				return err
			}
		case "awsv4":
			if err := json.Unmarshal([]byte(v), &strct.Awsv4); err != nil {
				return err
			}
		case "basic":
			if err := json.Unmarshal([]byte(v), &strct.Basic); err != nil {
				return err
			}
		case "bearer":
			if err := json.Unmarshal([]byte(v), &strct.Bearer); err != nil {
				return err
			}
		case "digest":
			if err := json.Unmarshal([]byte(v), &strct.Digest); err != nil {
				return err
			}
		case "edgegrid":
			if err := json.Unmarshal([]byte(v), &strct.Edgegrid); err != nil {
				return err
			}
		case "hawk":
			if err := json.Unmarshal([]byte(v), &strct.Hawk); err != nil {
				return err
			}
		case "noauth":
			if err := json.Unmarshal([]byte(v), &strct.Noauth); err != nil {
				return err
			}
		case "ntlm":
			if err := json.Unmarshal([]byte(v), &strct.Ntlm); err != nil {
				return err
			}
		case "oauth1":
			if err := json.Unmarshal([]byte(v), &strct.Oauth1); err != nil {
				return err
			}
		case "oauth2":
			if err := json.Unmarshal([]byte(v), &strct.Oauth2); err != nil {
				return err
			}
		case "type":
			if err := json.Unmarshal([]byte(v), &strct.Type); err != nil {
				return err
			}
			typeReceived = true
		}
	}
	// check if type (a required property) was received
	if !typeReceived {
		return errors.New("\"type\" is required but was not present")
	}
	return nil
}

// MarshalJSON xx
func (strct *AuthAttribute) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "Key" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "key" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"key\": ")
	if tmp, err := json.Marshal(strct.Key); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "type" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"type\": ")
	if tmp, err := json.Marshal(strct.Type); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "value" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"value\": ")
	if tmp, err := json.Marshal(strct.Value); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON todo
func (strct *AuthAttribute) UnmarshalJSON(b []byte) error {
	keyReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "key":
			if err := json.Unmarshal([]byte(v), &strct.Key); err != nil {
				return err
			}
			keyReceived = true
		case "type":
			if err := json.Unmarshal([]byte(v), &strct.Type); err != nil {
				return err
			}
		case "value":
			if err := json.Unmarshal([]byte(v), &strct.Value); err != nil {
				return err
			}
		}
	}
	// check if key (a required property) was received
	if !keyReceived {
		return errors.New("\"key\" is required but was not present")
	}
	return nil
}

// MarshalJSON todo
func (strct *Cookie) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "Domain" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "domain" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"domain\": ")
	if tmp, err := json.Marshal(strct.Domain); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "expires" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"expires\": ")
	if tmp, err := json.Marshal(strct.Expires); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "extensions" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"extensions\": ")
	if tmp, err := json.Marshal(strct.Extensions); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "hostOnly" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"hostOnly\": ")
	if tmp, err := json.Marshal(strct.HostOnly); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "httpOnly" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"httpOnly\": ")
	if tmp, err := json.Marshal(strct.HTTPOnly); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "maxAge" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"maxAge\": ")
	if tmp, err := json.Marshal(strct.MaxAge); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err := json.Marshal(strct.Name); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Path" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "path" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"path\": ")
	if tmp, err := json.Marshal(strct.Path); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "secure" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"secure\": ")
	if tmp, err := json.Marshal(strct.Secure); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "session" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"session\": ")
	if tmp, err := json.Marshal(strct.Session); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "value" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"value\": ")
	if tmp, err := json.Marshal(strct.Value); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Cookie) UnmarshalJSON(b []byte) error {
	domainReceived := false
	pathReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "domain":
			if err := json.Unmarshal([]byte(v), &strct.Domain); err != nil {
				return err
			}
			domainReceived = true
		case "expires":
			if err := json.Unmarshal([]byte(v), &strct.Expires); err != nil {
				return err
			}
		case "extensions":
			if err := json.Unmarshal([]byte(v), &strct.Extensions); err != nil {
				return err
			}
		case "hostOnly":
			if err := json.Unmarshal([]byte(v), &strct.HostOnly); err != nil {
				return err
			}
		case "httpOnly":
			if err := json.Unmarshal([]byte(v), &strct.HTTPOnly); err != nil {
				return err
			}
		case "maxAge":
			if err := json.Unmarshal([]byte(v), &strct.MaxAge); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		case "path":
			if err := json.Unmarshal([]byte(v), &strct.Path); err != nil {
				return err
			}
			pathReceived = true
		case "secure":
			if err := json.Unmarshal([]byte(v), &strct.Secure); err != nil {
				return err
			}
		case "session":
			if err := json.Unmarshal([]byte(v), &strct.Session); err != nil {
				return err
			}
		case "value":
			if err := json.Unmarshal([]byte(v), &strct.Value); err != nil {
				return err
			}
		}
	}
	// check if domain (a required property) was received
	if !domainReceived {
		return errors.New("\"domain\" is required but was not present")
	}
	// check if path (a required property) was received
	if !pathReceived {
		return errors.New("\"path\" is required but was not present")
	}
	return nil
}

func (strct *Event) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "disabled" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"disabled\": ")
	if tmp, err := json.Marshal(strct.Disabled); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "id" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"id\": ")
	if tmp, err := json.Marshal(strct.ID); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Listen" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "listen" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"listen\": ")
	if tmp, err := json.Marshal(strct.Listen); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "script" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"script\": ")
	if tmp, err := json.Marshal(strct.Script); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Event) UnmarshalJSON(b []byte) error {
	listenReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "disabled":
			if err := json.Unmarshal([]byte(v), &strct.Disabled); err != nil {
				return err
			}
		case "id":
			if err := json.Unmarshal([]byte(v), &strct.ID); err != nil {
				return err
			}
		case "listen":
			if err := json.Unmarshal([]byte(v), &strct.Listen); err != nil {
				return err
			}
			listenReceived = true
		case "script":
			if err := json.Unmarshal([]byte(v), &strct.Script); err != nil {
				return err
			}
		}
	}
	// check if listen (a required property) was received
	if !listenReceived {
		return errors.New("\"listen\" is required but was not present")
	}
	return nil
}

func (strct *Header) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "description" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"description\": ")
	if tmp, err := json.Marshal(strct.Description); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "disabled" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"disabled\": ")
	if tmp, err := json.Marshal(strct.Disabled); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Key" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "key" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"key\": ")
	if tmp, err := json.Marshal(strct.Key); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Value" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "value" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"value\": ")
	if tmp, err := json.Marshal(strct.Value); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Header) UnmarshalJSON(b []byte) error {
	keyReceived := false
	valueReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "description":
			if err := json.Unmarshal([]byte(v), &strct.Description); err != nil {
				return err
			}
		case "disabled":
			if err := json.Unmarshal([]byte(v), &strct.Disabled); err != nil {
				return err
			}
		case "key":
			if err := json.Unmarshal([]byte(v), &strct.Key); err != nil {
				return err
			}
			keyReceived = true
		case "value":
			if err := json.Unmarshal([]byte(v), &strct.Value); err != nil {
				return err
			}
			valueReceived = true
		}
	}
	// check if key (a required property) was received
	if !keyReceived {
		return errors.New("\"key\" is required but was not present")
	}
	// check if value (a required property) was received
	if !valueReceived {
		return errors.New("\"value\" is required but was not present")
	}
	return nil
}

func (strct *Info) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "description" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"description\": ")
	if tmp, err := json.Marshal(strct.Description); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Name" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err := json.Marshal(strct.Name); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "_postman_id" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"_postman_id\": ")
	if tmp, err := json.Marshal(strct.PostmanID); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Schema" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "schema" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"schema\": ")
	if tmp, err := json.Marshal(strct.Schema); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "version" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"version\": ")
	if tmp, err := json.Marshal(strct.Version); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshal
func (strct *Info) UnmarshalJSON(b []byte) error {
	nameReceived := false
	schemaReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "description":
			if err := json.Unmarshal([]byte(v), &strct.Description); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
			nameReceived = true
		case "_postman_id":
			if err := json.Unmarshal([]byte(v), &strct.PostmanID); err != nil {
				return err
			}
		case "schema":
			if err := json.Unmarshal([]byte(v), &strct.Schema); err != nil {
				return err
			}
			schemaReceived = true
		case "version":
			if err := json.Unmarshal([]byte(v), &strct.Version); err != nil {
				return err
			}
		}
	}
	// check if name (a required property) was received
	if !nameReceived {
		return errors.New("\"name\" is required but was not present")
	}
	// check if schema (a required property) was received
	if !schemaReceived {
		return errors.New("\"schema\" is required but was not present")
	}
	return nil
}

// MarshalJSON just marshal on items?
func (strct *Item) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "description" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"description\": ")
	if tmp, err := json.Marshal(strct.Description); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "event" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"event\": ")
	if tmp, err := json.Marshal(strct.Event); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "id" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"id\": ")
	if tmp, err := json.Marshal(strct.ID); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err := json.Marshal(strct.Name); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "protocolProfileBehavior" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"protocolProfileBehavior\": ")
	if tmp, err := json.Marshal(strct.ProtocolProfileBehavior); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Request" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "request" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"request\": ")
	if tmp, err := json.Marshal(strct.Request); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "response" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"response\": ")
	if tmp, err := json.Marshal(strct.Response); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "variable" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"variable\": ")
	if tmp, err := json.Marshal(strct.Variable); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Item) UnmarshalJSON(b []byte) error {
	requestReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "description":
			if err := json.Unmarshal([]byte(v), &strct.Description); err != nil {
				return err
			}
		case "event":
			if err := json.Unmarshal([]byte(v), &strct.Event); err != nil {
				return err
			}
		case "id":
			if err := json.Unmarshal([]byte(v), &strct.ID); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		case "protocolProfileBehavior":
			if err := json.Unmarshal([]byte(v), &strct.ProtocolProfileBehavior); err != nil {
				return err
			}
		case "request":
			if err := json.Unmarshal([]byte(v), &strct.Request); err != nil {
				return err
			}
			requestReceived = true
		case "response":
			if err := json.Unmarshal([]byte(v), &strct.Response); err != nil {
				return err
			}
		case "variable":
			if err := json.Unmarshal([]byte(v), &strct.Variable); err != nil {
				return err
			}
		}
	}
	// check if request (a required property) was received
	if !requestReceived {
		return errors.New("\"request\" is required but was not present")
	}
	return nil
}

func (strct *ItemGroup) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "auth" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"auth\": ")
	if tmp, err := json.Marshal(strct.Auth); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "description" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"description\": ")
	if tmp, err := json.Marshal(strct.Description); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "event" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"event\": ")
	if tmp, err := json.Marshal(strct.Event); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Item" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "item" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"item\": ")
	if tmp, err := json.Marshal(strct.Item); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err := json.Marshal(strct.Name); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "protocolProfileBehavior" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"protocolProfileBehavior\": ")
	if tmp, err := json.Marshal(strct.ProtocolProfileBehavior); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "variable" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"variable\": ")
	if tmp, err := json.Marshal(strct.Variable); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *ItemGroup) UnmarshalJSON(b []byte) error {
	itemReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "auth":
			if err := json.Unmarshal([]byte(v), &strct.Auth); err != nil {
				return err
			}
		case "description":
			if err := json.Unmarshal([]byte(v), &strct.Description); err != nil {
				return err
			}
		case "event":
			if err := json.Unmarshal([]byte(v), &strct.Event); err != nil {
				return err
			}
		case "item":
			if err := json.Unmarshal([]byte(v), &strct.Item); err != nil {
				return err
			}
			itemReceived = true
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		case "protocolProfileBehavior":
			if err := json.Unmarshal([]byte(v), &strct.ProtocolProfileBehavior); err != nil {
				return err
			}
		case "variable":
			if err := json.Unmarshal([]byte(v), &strct.Variable); err != nil {
				return err
			}
		}
	}
	// check if item (a required property) was received
	if !itemReceived {
		return errors.New("\"item\" is required but was not present")
	}
	return nil
}

func (strct *Root) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "auth" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"auth\": ")
	if tmp, err := json.Marshal(strct.Auth); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "event" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"event\": ")
	if tmp, err := json.Marshal(strct.Event); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Info" field is required
	if strct.Info == nil {
		return nil, errors.New("info is a required field")
	}
	// Marshal the "info" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"info\": ")
	if tmp, err := json.Marshal(strct.Info); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Item" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "item" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"item\": ")
	if tmp, err := json.Marshal(strct.Item); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "protocolProfileBehavior" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"protocolProfileBehavior\": ")
	if tmp, err := json.Marshal(strct.ProtocolProfileBehavior); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "variable" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"variable\": ")
	if tmp, err := json.Marshal(strct.Variable); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Root) UnmarshalJSON(b []byte) error {
	infoReceived := false
	itemReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "auth":
			if err := json.Unmarshal([]byte(v), &strct.Auth); err != nil {
				return err
			}
		case "event":
			if err := json.Unmarshal([]byte(v), &strct.Event); err != nil {
				return err
			}
		case "info":
			if err := json.Unmarshal([]byte(v), &strct.Info); err != nil {
				return err
			}
			infoReceived = true
		case "item":
			if err := json.Unmarshal([]byte(v), &strct.Item); err != nil {
				return err
			}
			itemReceived = true
		case "protocolProfileBehavior":
			if err := json.Unmarshal([]byte(v), &strct.ProtocolProfileBehavior); err != nil {
				return err
			}
		case "variable":
			if err := json.Unmarshal([]byte(v), &strct.Variable); err != nil {
				return err
			}
		}
	}
	// check if info (a required property) was received
	if !infoReceived {
		return errors.New("\"info\" is required but was not present")
	}
	// check if item (a required property) was received
	if !itemReceived {
		return errors.New("\"item\" is required but was not present")
	}
	return nil
}

func exportAREXToPostman(appid, start string) interface{} {
	var startTime time.Time
	startTime, err := time.Parse("2022-02-22", start)
	if err != nil {
		startTime = time.Time{}
	}

	rl := queryServletmocker(context.TODO(), appid, startTime)
	root := getPostmanRoot(rl)
	jsonData, err := root.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	structData := make(map[string]interface{})
	json.Unmarshal(jsonData, &structData)
	// res, err := json.MarshalIndent(structData, " ", " ")
	return structData
}

func getPostmanRoot(mockers []*servletmocker) *Root {
	hset := make(map[string]struct{})
	filterItem := func(item *Item) bool {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%v", item.Request.Method))
		sb.WriteString(item.Name)
		sb.WriteString(fmt.Sprintf("%v", item.Request.URL))
		sb.WriteString(item.Request.Body.Raw)
		if _, _exist := hset[sb.String()]; _exist {
			return true
		}
		hset[sb.String()] = struct{}{}
		return false
	}

	FirstUpper := func(s string) string {
		if s == "" {
			return ""
		}
		return strings.ToUpper(s[:1]) + s[1:]
	}

	convertMockerToItem := func(mocker *servletmocker) *Item {
		var item Item
		item.Name = mocker.AppID + mocker.Path
		var request Request
		request.Method = mocker.Method
		headerlist := make([]Header, 0)
		for mkey, mvalue := range mocker.RequestHeaders {
			headerlist = append(headerlist, Header{Key: FirstUpper(mkey), Value: mvalue, Disabled: false})
		}
		request.Header = append(headerlist, Header{Key: "arex-record-id", Value: mocker.ID})
		request.URL = "http://{{app_arexed_url}}" + mocker.Path
		request.Body.Mode = "raw"
		request.Body.Options = "{\"raw\":{\"language\":\"json\"}}"
		if mocker.Request != "" {
			bytes, err := unBase64andZstdString(mocker.Request)
			if err != nil {
				request.Body.Raw = err.Error()
			} else {
				text := string(bytes)
				request.Body.Raw = string(text)
				if strings.Contains(text, "ew0") {
					unzipData, err := base64.StdEncoding.DecodeString(text)
					if err == nil {
						request.Body.Raw = string(unzipData)
					}
				}
			}
		}
		item.Request = request

		if filterItem(&item) {
			return nil
		}

		return &item
	}

	if len(mockers) == 0 {
		return nil
	}

	var info Info
	info.Description = "arex auto generation test case"
	info.Name = "ArexImported"
	info.Schema = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"

	var root Root
	root.Info = &info

	items := make([]interface{}, 0)
	for _, mocker := range mockers {
		item := convertMockerToItem(mocker)
		if item == nil {
			continue
		}
		items = append(items, item)
	}

	root.Item = items
	return &root
}
