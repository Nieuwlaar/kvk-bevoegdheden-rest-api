Machtiging REST API
--------------------

Simple REST API build on top of kvk extract lib. 

## Setup
Make sure you have cloned the dependency project https://github.com/privacybydesign/kvk-extract in the same folder as this project.

## To run locally
```
sh restart.sh
```
This will run the API. You can remove the cert and key from restart script, the API will fallback on cached data. 

You can request the only endpoint of the API with a POST request on http://localhost:3333/api/bevoegdheid with the following body:
```
{
	"kvkNummer": "90000021",
	"functionaris": {
		"dateOfBirth": "14-12-1979",
		"firstNames": "Kerry",
		"lastName": "Rone"
	}
}
```
