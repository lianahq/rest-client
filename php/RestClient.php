<?php
namespace LianaTech;

/**
 * PHP RestClient for LianaTech RESTful services
 *
 * Example:
 *
 * 	$client = new RestCLient(...);
 *
 */
class RestClient {

	protected $api_key = null;
	protected $api_user = null;
	protected $api_url = null;
	protected $api_version = null;
	protected $api_realm = null;

	public function __construct($api_user, $api_key, $api_url, $api_version, $api_realm) {
		$this->api_user = $api_user;
		$this->api_key = $api_key;
		$this->api_url = $api_url;
		$this->api_version = intval($api_version);
		$this->api_realm = $api_realm;
	}

	public function call($path, $args = array(), $method = 'POST') {
		$args = $method !== 'GET' ? json_encode($args) : '';
		return $this->request($path, $args, $method);
	}

	protected function sign(array $message) {
		return hash_hmac('sha256', implode("\n", $message), $this->api_key);
	}

	protected function request($path, $contents, $method) {
		$md5 = md5($contents);
		$datetime = new \DateTime('now', new \DateTimeZone('Europe/Helsinki'));
		$timestamp = $datetime->format('c');
		$url = $this->api_url . '/api/v'. $this->api_version .'/' . $path;

		$message = array(
			$method,
			$md5,
			'application/json',
			$timestamp,
			$contents,
			'/api/v'. $this->api_version .'/' . $path
		);

		$signature = $this->sign($message);
		$user_id = $this->api_user;

		$ch = curl_init();
		curl_setopt($ch, CURLOPT_URL, $url);
		curl_setopt($ch, CURLOPT_HTTPHEADER, array(
			'Content-Type: application/json',
			"Content-MD5: {$md5}",
			"Date: {$timestamp}",
			"Authorization: {$this->api_realm} {$user_id}:{$signature}"
		));
		curl_setopt($ch, CURLOPT_USERAGENT, 'PHP-LMAPI');
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($ch, CURLOPT_TIMEOUT, 10);
		curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $method);
		curl_setopt($ch, CURLOPT_POSTFIELDS, $contents);
		$result = curl_exec($ch);
		$http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
		curl_close($ch);

		$result = json_decode($result, true);

		switch ($http_code) {
			case 401:
				throw new RestClientAuthorizationException;
				break;
			default:
				if ($http_code >= 400) {
					throw new APIException(sprintf(
						'API response with status code %s %s',
						$http_code,
						$result['message'] ? $result['message'] : ''
					));
				}
		}

		if ($result === null) {
			throw new APIException('API did not return a valid json string');
		}

		return $result;
	}

}

class RestClientAuthorizationException extends \Exception {}
class APIException extends \Exception {}
