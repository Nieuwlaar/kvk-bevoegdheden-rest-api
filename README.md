KVK Bevoegdheden REST API
--------------------

Simple REST API build on top of KVK Bevoegdheden lib. 

## Setup
Make sure you have cloned the dependency project https://github.com/kvk-innovatie/kvk-bevoegdheden in the same folder as this project.

## To run locally
```
sh restart.sh
```
or Visual Studio Code launch.json
```
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceRoot}"
        }
    ]
}
```
This will run the API. You can remove the cert and key from restart script, the API will fallback on cached data. 


### Base URL
```
http://localhost:3333/api
```
## Endpoints

### 1. LPID

**Endpoint:**
```
POST /api/lpid/{kvkNummer}
```
**Description:**
Fetches the LPID (Legal Person Identification Data) details for the specified `kvkNummer`.

**Request Parameters:**
- `kvkNummer` (path): The KVK number of the entity.

**Response:**
```
{
    "data": {
        "id": "NLNHR.{kvkNummer}",
        "legal_person_name": "{Legal Person Name}",
        "legal_form": "{Legal Form}"
    },
    "metadata": {
        "issuing_authority_name": "Kamer van Koophandel",
        "issuer_id": "NLNHR.59581883",
        "issuing_country": "NL",
        "issuance_date": "2022-06-15T15:35:52.687Z",
        "expiry_date": "2025-06-15T15:35:52.687Z",
        "schema": "http://schema.example.com",
        "revocation_information": "http://revoke.example.com"
    }
}
```
### 2. Company Certificate
**Endpoint:**
```
POST /api/company-certificate/{kvkNummer}
```
**Description:**
Fetches the company certificate details for the specified kvkNummer.

**Request Parameters:**
- `kvkNummer` (path): The KVK number of the entity.


**Response:**
```
{
    "data": {
        "id": "NLNHR.{kvkNummer}",
        "legal_person_name": "{Legal Person Name}",
        "legal_form": "{Legal Form}",
        "registration_number": "{KVK Number}",
        "registered_country": "NL",
        "registered_office": "{Registered Office Address}",
        "postal_address": "{Postal Address}",
        "electronic_address": "{Email Address}",
        "date_of_registration": "{Date of Registration}",
        "capital_subscribed": "{Capital Subscribed}",
        "status": "{Status}",
        "authorized_persons": [
            {
                "full_name": "{Full Name}",
                "date_of_birth": "{Date of Birth}",
                "interpretatie": {
                    "isAuthorized": "{Yes/No/Not determined}"
                }
            }
        ],
        "object": "{Object}"
    },
    "metadata": {
        "issuing_authority_name": "Kamer van Koophandel",
        "issuer_id": "NLNHR.59581883",
        "issuing_country": "NL",
        "issuance_date": "2022-06-15T15:35:52.687Z",
        "expiry_date": "2025-06-15T15:35:52.687Z",
        "schema": "http://schema.example.com",
        "revocation_information": "http://revoke.example.com"
    }
}
```
### 3. Natural Person Signatory Right
**Endpoint:**
```
POST /api/signatory-rights/{kvkNummer}
```
**Description:**
Checks if a natural person has signatory rights for the specified kvkNummer.

**Request Parameters:**
- `kvkNummer` (path): The KVK number of the entity.
- Request Body:
```
{
    "geslachtsnaam": "Klaassen",     // Surname
    "voornamen": "Jan",              // First names
    "geboortedatum": "01-01-2000",   // Date of Birth
    "voorvoegselGeslachtsnaam": ""   // Prefix of Surname
}
```
**Response:**
```
{
    "data": {
        "full_name": "{Matched Full Name}",
        "date_of_birth": "{Date of Birth}",
        "is_authorized": true/false,
        "id": "NLNHR.{kvkNummer}",
        "legal_person_name": "{Legal Person Name}",
        "legal_form": "{Legal Form}"
    },
    "metadata": {
        "issuing_authority_name": "Kamer van Koophandel",
        "issuer_id": "NLNHR.59581883",
        "issuing_country": "NL",
        "issuance_date": "2022-06-15T15:35:52.687Z",
        "expiry_date": "2025-06-15T15:35:52.687Z",
        "schema": "http://schema.example.com",
        "revocation_information": "http://revoke.example.com"
    }
}
```

## Examples
6 example responses from the HR dataservice XML files are provided. 
- 90000001
- 90000002
- 90000003
- 90000004
- 90000005
- 90000006
You can use/edit these to test this rest-api.
