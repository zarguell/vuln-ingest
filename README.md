# Vulnerability Ingestion Microservice

This is a microservice written in Go that serves an API endpoint to ingest vulnerabilities into an SQLlite database. The API accepts `POST` requests with a JSON body containing the details about the vulnerability, such as the title, CVE or CWE, evidence, and where it was found. This microservice (and documentation) was written entirely by chatGPT.

## Features

- API endpoint to ingest vulnerabilities
- Accepts `POST` requests with a JSON body containing vulnerability details
- Option to require an API token for authentication
- Reads a JSON config file to set API token requirement and service port
- Hashes API tokens using the bcrypt algorithm and stores them in an SQLlite database

## Usage

To use this microservice, run the following command:

```bash
go run vulnerability_ingestion.go -config config.json
```

Replace `config.json` with the path to your config file. The config file must be a valid JSON file with the following structure:

```json
{
"require_token": true,
"port": 8080,
"db_file": "vulnerabilities.db",
"tokens_db_file": "tokens.db"
}
```

- `require_token`: set to `true` if an API token is required for authentication, or `false` if not.
- `port`: the port number the microservice will listen on.
- `db_file`: the name of the SQLlite database file to store vulnerabilities in.
- `tokens_db_file`: the name of the SQLlite database file to store API tokens in.

## API Token Database

This microservice uses an SQLlite database to store the hashed API tokens. The database is specified in the config file as `tokens_db_file` and contains a single table named `tokens`. The `tokens` table has two columns:

- `id`: an auto-incrementing integer that serves as the primary key for each record
- `hash`: a binary large object (BLOB) that stores the hashed API token

## Vulnerabilities Database

This microservice uses an SQLlite database to store the ingested vulnerabilities. The database is specified in the config file as `db_file` and contains a single table named `vulnerabilities`. The `vulnerabilities` table has the following columns:

- `id`: an auto-incrementing integer that serves as the primary key for each record
- `title`: the title of the vulnerability
- `cve`: the Common Vulnerabilities and Exposures (CVE) identifier for the vulnerability
- `cwe`: the Common Weakness Enumeration (CWE) identifier for the vulnerability
- `evidence`: a string containing any evidence or details about the vulnerability
- `found_where`: a string indicating where the vulnerability was found

## API Endpoint

The API endpoint for this microservice is `/ingest`. It accepts `POST` requests with a JSON body containing the following keys:

- `title`: the title of the vulnerability
- `cve`: the Common Vulnerabilities and Exposures (CVE) identifier for the vulnerability
- `cwe`: the Common Weakness Enumeration (CWE) identifier for the vulnerability
- `evidence`: a string containing any evidence or details about the vulnerability
- `found_where`: a string indicating where the vulnerability was found
- `token (optional)`: the API token to authenticate the request, if required

The API returns the following HTTP status codes:

- `201 Created`: If the vulnerability was successfully ingested into the database
- `400 Bad Request`: If the request is missing required keys or the token key is incorrect
- `401 Unauthorized`: If the require_token option is set to true in the config file and the token key is not provided or is incorrect
- `500 Internal Server Error`: If an error occurs while trying to ingest the vulnerability into the database

## Sample JSON Body for POST Request

The following is an example of a JSON body for a POST request to the microservice's API endpoint:

```json
{
    "title": "SQL Injection Vulnerability",
    "cve": "CVE-2021-23456",
    "cwe": "CWE-89",
    "evidence": "This vulnerability was found during a penetration test",
    "found_where": "Penetration Test",
    "token": "optional_token_value"
}
```

Note that the token key is optional, and its value should only be provided if the require_token option is set to true in the config file.

## API Token Generation Helper Script

A helper script is provided to generate API tokens for this microservice. The script is written in Go and generates a random API token, hashes it using the bcrypt algorithm, and inserts it into the tokens table in the API tokens database. The raw API token is then output to stdout.

To use this script, run the following command:

```bash
go run generate_token.go -tokens_db_file tokens.db
```

Replace `tokens.db` with the path to your API tokens database file. The script will generate a random API token, hash it, insert it into the tokens table, and output the raw API token to `stdout`.

## SQLlite Database

This script uses an SQLlite database to store the hashed API tokens. The database is named tokens.db and contains a single table named tokens. The tokens table has two columns:

- `id:` an auto-incrementing integer that serves as the primary key for each record 
- `hash`: a binary large object (BLOB) that stores the hashed API token

## Output

The raw API token will be outputted to `stdout` in plain text format. This raw token can be used as the API token in the configuration file for the micro service described in the previous sections.


## Conclusion

This microservice provides an API endpoint to ingest vulnerabilities into an SQLlite database, with the option to require an API token for authentication. The microservice can be configured using a JSON config file, and the API tokens are stored in a separate SQLlite database, hashed using the bcrypt algorithm. A helper script is provided to generate API tokens for this microservice.