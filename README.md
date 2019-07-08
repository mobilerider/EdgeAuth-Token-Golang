# EdgeAuth-Token-Golang: Akamai Edge Authorization Token for Golang

EdgeAuth-Token-Golang is Akamai Edge Authorization Token in the HTTP Cookie, Query String and Header for a client.
You can configure it in the Property Manager at https://control.akamai.com.
It's the behaviors which is Auth Token 2.0 Verification and Segmented Media Protection.

EdgeAuth-Token-Golang supports for Golang 1.x

<div style="text-align:center"><img src=https://github.com/AstinCHOI/akamai-asset/blob/master/edgeauth/edgeauth.png?raw=true/></div>


## Installation

To install Akamai Edge Authorization Token with dep:  

```Shell
$ dep ensure -add github.com/mobilerider/EdgeAuth-Token-Golang
```
  

## Examples
#### ACL (Access Control List) parameter option

```Golang
package main

import "github.com/mobilerider/EdgeAuth-Token-Golang/edgeauth

func main {
    acl := "/s/c/m/f/5/9/2/f5929e909d4/*"
    key := "52a152a152a152a152a152a152a1"

    config := &edgeauth.Config{
		Algo:           crypto.SHA256,
		Key:            sampleKey,
		DurationWindow: 300 * time.Second,
	}

	client := edgeauth.NewClient(config)

    token, err := client.GenerateToken(acl, false)
    
    if err != nil {
		// Handle error
    }
    
    // Generated token value would in the form of:
    // exp=1562609231~acl=/s/c/m/f/5/9/2/f5929e909d4/*~hmac=7a6bd5d5abdad74bda765b4da67b7ad54b6a4d6ba54d67b4ad76b4

    // You will probably use that token as a value of a cookie 
    // or query string parameter, the name of the parameter is 
    // configured via PM
}

```
* ACL can use the wildcard(\*, ?) in the path.
* Don't use '!' in your path because it's ACL Delimiter
* Use 'escapeEarly=false' as default setting but it doesn't matter turning on/off 'Escape token input' option in the Property Manager


## Usage

#### EdgeAuth Config Class

```Golang
type Config struct {
	Algo           crypto.Hash
	Key            string
	Salt           string
	FieldDelimiter string
	ACLDelimiter   string
	StartTime      time.Time
	EndTime        time.Time
	DurationWindow time.Duration
	IP             string
	SessionID      string
	Payload        string
	Verbose        bool
}
```

| Parameter | Description |
|-----------|-------------|
| options.key | Secret required to generate the token. It must be hexadecimal digit string with even-length. |
| Algo  | Algorithm to use to generate the token. ('sha1', 'sha256', or 'md5') [ Default: 'sha256' ] |
| Salt | Additional data validated by the token but NOT included in the token body. (It will be deprecated) |
| StartTime | What is the start time? (Use string 'now' for the current time) |
| EndTime | When does this token expire? endTime overrides windowSeconds |
| DurationWindow | How long is this token valid for? |
| FieldDelimiter | Character used to delimit token body fields. [ Default: ~ ] |
| AclDelimiter | Character used to delimit acl. [ Default: ! ] |
| EscapeEarly | Not implemented yet. |
| Verbose | Not implemented yet. |


#### EdgeAuth's Method

```Golang
client.GenerateToken(value, isUrl) {}

// both return the authorization token string.
```

| Parameter | Description |
|-----------|-------------|
| url | Single URL path (String) |
| acl | Access Control List can use the wildcard(\*, ?). It can be String (single path) or Array (multi paths) |


## Others

If you use the **Segmented Media Protection** behavior in AMD(Adaptive Media Delivery) Product, **tokenName(options.tokenName)** should be '**hdnts**'.

<div style="text-align:center"><img src=https://github.com/AstinCHOI/akamai-asset/blob/master/edgeauth/segmented_media_protection.png?raw=true/></div>

### TODOs

1- Implement EscapeEarly option
2- Implement Verbose option

### Important!!: This is NOT an official library from Akamai Technologies.