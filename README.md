# Project Description

`pdnsutil` fetches DNS records from the corporate PowerDNS and allows them to be deleted based on a specific pattern.

The records are received in JSON format with the following structure:

```json
{
    "account": "account",
    "api_rectify": false,
    "catalog": "",
    "dnssec": false,
    "edited_serial": 2024050701,
    "id": "example.com.",
    "kind": "Master",
    "last_check": 0,
    "master_tsig_key_ids": [],
    "masters": [],
    "name": "example.com.",
    "notified_serial": 0,
    "nsec3narrow": false,
    "nsec3param": "",
    "rrsets": [
        {
            "comments": [],
            "name": "host1.example.com.",
            "records": [
                {
                    "content": "192.168.0.2",
                    "disabled": false
                }
            ],
            "ttl": 0,
            "type": "A"
        },
       {
            "comments": [],
            "name": "host2.example.com.",
            "records": [
                {
                    "content": "192.168.0.3",
                    "disabled": false
                }
            ],
            "ttl": 0,
            "type": "A"
        }
    ],
    "serial": 2024050701,
    "slave_tsig_key_ids": [],
    "soa_edit": "",
    "soa_edit_api": "",
    "url": "/api/v1/servers/localhost/zones/example.com."
}
```

so all the records are in the `rrsets` field.

the `pdnsutil` must have following subcommands:

- `pdnsutil list-records`
- `pdnsutil delete-records`

The root command must be `pdnsutil` and it must accept following flags:

- `-h/--help` - help
- `-d/--debug` - debug
- `-a/--api-url` - PowerDNS API URL (required)
- `-k/--api-key` - PowerDNS API key (required)
- `-z/--zone` - zone name (required)

all required flags must also be read from the environment variables:

- `PDNS_API_URL` - PowerDNS API URL
- `PDNS_API_KEY` - PowerDNS API key
- `PDNS_ZONE` - zone name

for example:

```bash
pdnsutil \
  --zone "example.com."
  --api-url "https://pdns.example.com/api/v1/"
  --api-key "1234567890"
  list-records
```

Subcommands must accept following flags:

- `-f/--filter` - filter for records
- `-e/--exlcude` - exclude from filtered
- `--force` - don't ask confirmation

Filter is a regular expression for the record name.

For example:

```bash
pdnsutil --zone example.com. --api-url https://pdns.example.com/api/v1/ --api-key 1234567890 list-records --filter 'host-.*'
```

# Curl requests
## Fetch records

```
curl -X GET \
  -H 'Content-Type: application/json' \
  -H 'X-API-Key: 1111111' \
  https://pdns.example.com/api/v1/servers/localhost/zones/example.com

```

## Delete record

```
curl -X PATCH \
  -H 'Content-Type: application/json' \
  -H 'X-API-Key: 1111111' \
  -d '{"rrsets": [{"name": "host1.example.com.", "type": "A", "ttl": 3600, "changetype": "DELETE"}]}' \
  https://pdns.example.com/api/v1/servers/localhost/zones/example.com

```

# Subcommands
## pdnsutil list-records

subcommand must fetch all records from PowerDNS into the PowerDNSZone struct.

Then it must filter records by the filter flag and print them to the stdout.

## pdnsutil delete-records

subcommand must fetch all records from PowerDNS into the PowerDNSZone struct.
Then it must filter records by the filter flag and save them to the slice of records.
Then it must print the slice of records to the stdout.
It must ask for confirmation before deleting the records.
Then it must delete the records from filtered slice from PowerDNS.

# Build

`go build -o pdnsutil` 
