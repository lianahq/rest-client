REST Client
===========

REST client in each programming language to use for LianaTech RESTful services.

Usage
=====

PHP

    <?php
    $user_id = <API_USER>;
    $secret_key = <API_SECRET>;
    $url = <API_URL>;
    $realm = <API_REALM>;
    $api_version = 1;

    require_once 'rest-client/php/RestClient.php';

    // Call V1 getRecipient

    // Create client
    $clientV1 = new \LianaTech\RestClient($user_id, $secret_key, $url, $api_version, $realm);

    // Define parameters
    $params = array(
        12345 // recipient id
    );

    // Call endpoint
    try {
        $res = $clientV1->call('getRecipient', $params);
    } catch (LianaTech\RestClientAuthorizationException $e) {
        echo "\n\tERROR: Authorization failed\n\n";
    } catch (LianaTech\APIException $e) {
        echo "\n\tERROR: API Exception: " . $e->getMessage() . "\n\n";
    } catch (Exception $e) {
        echo "\n\tERROR: " . $e->getMessage() . "\n\n";
    }

    $res = $res['result'];
    echo sprintf("V1 - getRecipient: Got recipient id: %s email: %s\n", $res['recipient']['id'], $res['recipient']['email']);


    // Call V2 import mailinglist

    // Create client
    $clientV2 = new \LianaTech\RestClient($user_id, $secret_key, $url, 2, $realm);

    // Define parameters
    $csv_data = "email;some-property\nrecipient-1@example.com;some-value-1\nrecipient-2@example.com;some-value-2";

    $params = array(
        'name' => 'new list',
        'data' => base64_encode($csv_data),
        'type' => 'csv',
        'truncate' => false
    );

    // Call endpoint
    try {
        $res = $clientV2->call('import/mailinglist', $params);
    } catch (LianaTech\RestClientAuthorizationException $e) {
        echo "\n\tERROR: Authorization failed\n\n";
    } catch (LianaTech\APIException $e) {
        echo "\n\tERROR: API Exception: " . $e->getMessage() . "\n\n";
    } catch (Exception $e) {
        echo "\n\tERROR: " . $e->getMessage() . "\n\n";
    }

    $res = $res['result'];
    echo sprintf("V2 - import/mailinglist: List created/updated, list_id %s\n", $res['list_id']);


    // Call V3 events

    // Create client
    $clientV3 = new \LianaTech\Restclient($user_id, $secret_key, $url, 3, $realm);

    // Call endpoint
    try {
        $res = $clientV3->call('events?at_start=2023-06-20T12:00:00&at_end=2023-06-22T12:00:00', array(), 'GET');
    } catch (LianaTech\RestClientAuthorizationException $e) {
        echo "\n\tERROR: Authorization failed\n\n";
    } catch (LianaTech\APIException $e) {
        echo "\n\tERROR: API Exception: " . $e->getMessage() . "\n\n";
    } catch (Exception $e) {
        echo "\n\tERROR: " . $e->getMessage() . "\n\n";
    }

    echo "V3 - events: Got following results:\n";
    foreach ($res['items'] as $item) {
        echo sprintf("\tat: %s, recipient: %s, event: %s\n", $item['at'], $item['recipient']['email'], $item['type']);
    }

Python

    import sys
    sys.path.append('./rest-client/python')

    if sys.version_info[0] == 3:
        import base64

    from RestClient import RestClient, APIException

    api_user = <API_USER_ID>
    api_secret = <API_SECRET>
    api_url = <API_URL>
    api_realm = <API_REALM>
    api_version = 1

    # Call V1 getRecipient

    # Create client
    client_v1 = RestClient(api_user, api_secret, api_url, api_version, api_realm)

    # Define parameters
    params = [
        12345 # recipient id
    ]

    # Call endpoint
    try:
        data = client_v1.call('getRecipient', params)
    except APIException as e:
        response = client_v1.get_http_response()
        print('API call failed: ' + str(e))
        print(response.status_code)
        print(response.headers)
        print(response.text)
        exit(1)

    # Print response data
    print("V1 - getRecipient: Got recipient id: {}, email: {}\n".format(data['recipient']['id'], data['recipient']['email']))


    # Call V2 import mailinglist

    # Create client
    client_v2 = RestClient(api_user, api_secret, api_url, 2, api_realm)

    # Define parameters
    csv_data = "email;some-property\nrecipient-1@example.com;some-value-1\nrecipient-2@example.com;some-value-2"

    if sys.version_info[0] == 3:
        csv_data = base64.b64encode(bytes(csv_data, 'utf-8')).decode('ascii')
    else:
        csv_data = csv_data.encode('base64')

    params = {
        'name' : 'new list',
        'data' : csv_data,
        'type' : 'csv',
        'truncate' : False
    }

    # Call endpoint
    try:
        data = client_v2.call('import/mailinglist', params)
    except APIException as e:
        response = client_v2.get_http_response()
        print('API call failed: ' + str(e))
        print(response.status_code)
        print(response.headers)
        print(response.text)
        exit(1)

    print("V2 - import/mailinglist: List created/updated, list_id {}\n".format(data['list_id']))


    # Call V3 events

    # Create client
    client_v3 = RestClient(api_user, api_secret, api_url, 3, api_realm)

    # Call endpoint
    try:
        data = client_v3.call('events?at_start=2023-06-20T12:00:00&at_end=2023-06-22T12:00:00', [], 'GET')
    except APIException as e:
        response = client_v3.get_http_response()
        print('API call failed: ' + str(e))
        print(response.status_code)
        print(response.headers)
        print(response.text)
        exit(1)

    print("V3 - events: Got following results:\n")
    for item in data['items']:
        print("\tat: {}, recipient: {}, event: {}\n".format(item['at'], item['recipient']['email'], item['type']))


Golang

	package main

	import (
		mailerapi "example/golang"
		"fmt"
		b64 "encoding/base64"
	)
	func main() {

		// API Credentials
		apiUser := <API_USER>
		apiSecret := "<API_SECRET>"
		apiUrl := "<API_URL>"
		apiVersion := <API_VERSION>
		apiRealm := "<API_REALM>"

		// Call V1 getRecipient

		// Create client
		mailerV1, err := mailerapi.NewRestClient(apiUser, apiSecret, apiUrl, apiVersion, apiRealm)

		// Define parameters
		gerRecipientParams := []interface{}{
			12345, // recipient id
		}

		// Call endpoint
		resp, err := mailerV1.Call("getRecipient", gerRecipientParams)

		if err != nil {
			panic("API returned an error: " + err.Error())
		}

		// Go thru response
		v1RespMap, ok := resp.(map[string]interface{})
		if ok {
			recipient := v1RespMap["recipient"].(map[string]interface{})
			fmt.Printf("V1 - getRecipient: Got recipient id: %+v email: %+v\n", recipient["id"], recipient["email"])
		}


		// Call V2 import mailinglist

		// Create client
		apiVersion = 2
		mailerV2, err := mailerapi.NewRestClient(apiUser, apiSecret, apiUrl, apiVersion, apiRealm)

		// Define parameters

		csvData := "email;some-property\nrecipient-1@example.com;some-value-1\nrecipient-2@example.com;some-value-2"

		// v2 functions take params with keys, tags used here since lowercase type is restricted
		importMailingListParams := struct{
			ListName string `json:"name"`
			Data string `json:"data"`
			Type string `json:"type"`
			Truncate bool `json:"truncate"`
		}{
			ListName: "new list from go client",
			Data: b64.StdEncoding.EncodeToString([]byte(csvData)),
			Type: "csv",
			Truncate: false,
		}

		// Call endpoint
		resp, err = mailerV2.Call("import/mailinglist", importMailingListParams)

		if err != nil {
			panic("API returned an error: " + err.Error())
		}

		// Go thru response
		v2RespMap, ok := resp.(map[string]interface{})
		if ok {
			fmt.Printf("V2 - import/mailinglist: List import started, list_id %+v\n", v2RespMap["list_id"])
		}


		// Call V3 events

		// Create client
		apiVersion = 3
		mailerV3, err := mailerapi.NewRestClient(apiUser, apiSecret, apiUrl, apiVersion, apiRealm)

		// Call endpoint
		resp, err = mailerV3.Call("events?at_start=2025-08-21T12:00:00&at_end=2025-08-22T12:00:00", nil, "GET")

		if err != nil {
			panic("API returned an error: " + err.Error())
		}

		// Go thru response
		v3RespMap, ok := resp.([]interface{})
		if ok {
			fmt.Printf("V3 - events: Got following results:\n")
			for _, v := range v3RespMap {
				event := v.(map[string]interface{})
				event_recipient := event["recipient"].(map[string]interface{})
				fmt.Printf("\tat: %+v, type: %+v recipient %+v: %+v\n", event["at"], event["type"], event_recipient["id"], event_recipient["email"])
			}
		}
	}



Development
===========

1. Clone this repository (and go to folder)

2. [Install composer](https://github.com/composer/composer)

3. Install required PHP dependencies (it will read composer.json file)

   `php composer.phar install`

4. Running tasks (currently only unit tests)

   `make test`
