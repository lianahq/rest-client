package restclient

import (
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"crypto/tls"
	"crypto/sha256"
	"crypto/hmac"
	"context"
	"io/ioutil"
	"testing"
)

// Check that the sign function works as expected
func TestSignature(test *testing.T) {

	apiUser := 123
	apiSecret := "apisecret"
	apiUrl := "https://api.local"
	apiVersion := 1
	apiRealm := "REALM"

	apiclient, err := NewRestClient(apiUser, apiSecret, apiUrl, apiVersion, apiRealm)

	if err != nil {
		test.Errorf("failed to create RestClient:")
	}

	signature, err := apiclient.sign(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		"POST",
		"cdefd9b4ca40e984f3482ed3c7ae077a",
		"application/json",
		"Wed, 21 Oct 2015 07:28:00 GMT",
		"[\"var1\",\"var2\"]",
		"/api/v1/unit/test",
	));

	if err != nil {
		test.Errorf("failed to generate signature")
	}

	if signature != "9f29f8706f8e5720eb4950650e11f70c7d4680d59876e1250968212f742da062" {
		test.Errorf("invalid signature: %s", signature)
	}
}


// API returns succeed false and an error message
func TestCallFail(test *testing.T) {
	handler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		reqtype := req.Header.Get("Content-Type")
		if reqtype != "application/json" {
			test.Errorf("invalid Content-Type header: %s", reqtype)
		}

		reqmd5 := req.Header.Get("Content-MD5")
		if reqmd5 != "cdefd9b4ca40e984f3482ed3c7ae077a" {
			test.Errorf("invalid MD5 header: %s", reqmd5)
		}

		reqdate := req.Header.Get("Date")
		reqauth := req.Header.Get("Authorization")

		sign := hmac.New(sha256.New, []byte("apisecret"))
		sign.Write([]byte(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
			"POST",
			"cdefd9b4ca40e984f3482ed3c7ae077a",
			"application/json",
			reqdate,
			"[\"var1\",\"var2\"]",
			"/api/v1/unit/test",
		)))

		if reqauth != fmt.Sprintf("%s%s", "REALM 123:", hex.EncodeToString(sign.Sum(nil))) {
			test.Errorf("invalid Authorization header: %s", reqauth)
		}

		payload, _ := ioutil.ReadAll(req.Body)
		if string(payload) != "[\"var1\",\"var2\"]" {
			test.Errorf("invalid request payload: %s", string(payload))
		}

		resp.Write([]byte("{\"succeed\":false,\"message\":\"unit test fail\"}"))
	})

	httpclient, srvteardown := createDummyServer(handler);

	defer srvteardown()

	apiUser := 123
	apiSecret := "apisecret"
	apiUrl := "https://api.local"
	apiVersion := 1
	apiRealm := "REALM"

	apiclient, err := NewRestClient(apiUser, apiSecret, apiUrl, apiVersion, apiRealm)

	if err != nil {
		test.Errorf("failed to contruct a restclient: %s", err.Error())
	}

	apiclient.httpclient = httpclient

	data, err := apiclient.Call("unit/test", []string{"var1", "var2"})

	if data != nil {
		test.Error("request was supposed to fail")
	}

	if err == nil {
		test.Error("request did not produce an error")
	}

	if err.Error() != "/api/v1/unit/test: unit test fail" {
		test.Errorf("request did not yield proper error message: %s", err.Error())
	}
}

// API returns succeed true and array of values
func TestCallSuccess1(test *testing.T) {
	handler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		reqtype := req.Header.Get("Content-Type")
		if reqtype != "application/json" {
			test.Errorf("invalid Content-Type header: %s", reqtype)
		}

		reqmd5 := req.Header.Get("Content-MD5")
		if reqmd5 != "0a9b61dec51f0560d8bd2a4740dbfe4e" {
			test.Errorf("invalid MD5 header: %s", reqmd5)
		}

		reqdate := req.Header.Get("Date")

		reqauth := req.Header.Get("Authorization")

		sign := hmac.New(sha256.New, []byte("apisecret2"))
		sign.Write([]byte(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
			"POST",
			"0a9b61dec51f0560d8bd2a4740dbfe4e",
			"application/json",
			reqdate,
			"[\"var3\"]",
			"/api/v1/unit/test2",
		)))

		if reqauth != fmt.Sprintf("%s%s", "REALM 1234:", hex.EncodeToString(sign.Sum(nil))) {
			test.Errorf("invalid Authorization header: %s", reqauth)
		}

		payload, _ := ioutil.ReadAll(req.Body)
		if string(payload) != "[\"var3\"]" {
			test.Errorf("invalid request payload: %s", string(payload))
		}

		resp.Write([]byte("{\"succeed\":true,\"result\":[\"val1\",\"val2\"]}"))
	})

	httpclient, srvteardown := createDummyServer(handler);

	defer srvteardown()

	apiUser := 1234
	apiSecret := "apisecret2"
	apiUrl := "https://api2.local"
	apiVersion := 1
	apiRealm := "REALM"

	apiclient, err := NewRestClient(apiUser, apiSecret, apiUrl, apiVersion, apiRealm)

	if err != nil {
		test.Errorf("failed to contruct a restclient: %s", err.Error())
	}

	apiclient.httpclient = httpclient

	data, err := apiclient.Call("unit/test2", []string{"var3"})

	values, _ := data.([]interface{})
	value1, _ := values[0].(string)
	value2, _ := values[1].(string)

	if value1 != "val1" {
		test.Error("request yielded an unexpected response #1")
	}

	if value2 != "val2" {
		test.Error("request yielded an unexpected response #2")
	}

	if err != nil {
		test.Errorf("request produced an error: %s", err.Error())
	}
}

// API request with multiple params and API returns an array of objects
func TestCallSuccess2(test *testing.T) {
	handler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		reqtype := req.Header.Get("Content-Type")
		if reqtype != "application/json" {
			test.Errorf("invalid Content-Type header: %s", reqtype)
		}

		reqmd5 := req.Header.Get("Content-MD5")
		if reqmd5 != "2f7da26fc0796322186a72244f8b8eb4" {
			test.Errorf("invalid MD5 header: %s", reqmd5)
		}

		reqdate := req.Header.Get("Date")

		reqauth := req.Header.Get("Authorization")

		sign := hmac.New(sha256.New, []byte("apisecret3"))
		sign.Write([]byte(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
			"POST",
			"2f7da26fc0796322186a72244f8b8eb4",
			"application/json",
			reqdate,
			"[\"var4\",\"var5\"]",
			"/api/v1/unit/test3",
		)))

		if reqauth != fmt.Sprintf("%s%s", "REALM 1234:", hex.EncodeToString(sign.Sum(nil))) {
			test.Errorf("invalid Authorization header: %s", reqauth)
		}

		payload, _ := ioutil.ReadAll(req.Body)
		if string(payload) != "[\"var4\",\"var5\"]" {
			test.Errorf("invalid request payload: %s", string(payload))
		}

		resp.Write([]byte("{\"succeed\":true,\"result\":[{\"key1\":\"val1\"},{\"key2\":\"val2\"}]}"))
	})

	httpclient, srvteardown := createDummyServer(handler);

	defer srvteardown()

	apiUser := 1234
	apiSecret := "apisecret3"
	apiUrl := "https://api3.local"
	apiVersion := 1
	apiRealm := "REALM"

	apiclient, err := NewRestClient(apiUser, apiSecret, apiUrl, apiVersion, apiRealm)

	if err != nil {
		test.Errorf("failed to contruct a restclient: %s", err.Error())
	}

	apiclient.httpclient = httpclient

	data, err := apiclient.Call("unit/test3", []string{"var4", "var5"})

	object_arr, _ := data.([]interface{})
	object1, _ := object_arr[0].(map[string]interface{})
	object2, _ := object_arr[1].(map[string]interface{})

	if object1["key1"].(string) != "val1" {
		test.Error("request yielded an unexpected response #1")
	}

	if object2["key2"].(string) != "val2" {
		test.Error("request yielded an unexpected response #2")
	}

	if err != nil {
		test.Errorf("request produced an error: %s", err.Error())
	}
}

// Creates a dummy server instance which serves mockup HTTP responses
func createDummyServer(handler http.Handler) (*http.Client, func()) {
	server := httptest.NewTLSServer(handler)

	dummyclient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return dummyclient, server.Close
}
