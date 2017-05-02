export VCAP_SERVICES='{
    "speech_to_text": [
	{
	    "name": "speech-to-text-service",
	    "label": "speech_to_text",
	    "plan": "standard",
	    "credentials": {
		"url": "https://stream.watsonplatform.net/speech-to-text/api",
		"username": "",
		"password": ""
	    }
	}
    ]
}'
