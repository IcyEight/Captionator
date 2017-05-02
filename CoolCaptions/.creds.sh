export VCAP_SERVICES='{
    "speech_to_text": [
	{
	    "name": "speech-to-text-service",
	    "label": "speech_to_text",
	    "plan": "standard",
	    "credentials": {
		"url": "https://stream.watsonplatform.net/speech-to-text/api",
		"username": "d6b9e278-7192-410f-8a04-4c3121a5db23",
		"password": "mkGChO1ClOvu"
	    }
	}
    ],
    "tone_analyzer": [
       {
          "name": "tone-analyzer",
          "label": "tone_analyzer",
          "credentials": {
             "url": "https://gateway.watsonplatform.net/tone-analyzer/api",
             "password": "bafdec9d-2dd5-4c4e-adc9-6e30ab881b35",
             "username": "mM35y5lZlaVO"
          }
       }
    ]
}'
