# Steampipe Plugin for IPGeolocation.io

Use SQL to query IP geolocation, threat intelligence, ASN details, and abuse contacts from [IPGeolocation.io](https://ipgeolocation.io).

```sql
select
  ip,
  country_name,
  city,
  isp,
  timezone_name
from
  ipgeolocation_ip
where
  ip = '8.8.8.8';
```

## Quick Start

Install the plugin:

```sh
steampipe plugin install ipgeolocation/ipgeolocation
```

Configure your API key in `~/.steampipe/config/ipgeolocation.spc`:

```hcl
connection "ipgeolocation" {
  plugin = "ipgeolocation/ipgeolocation"

  # Get your free API key at https://app.ipgeolocation.io/dashboard
  # Can also be set with the IPGEOLOCATION_API_KEY environment variable.
  # api_key = "a1b2c3d4e5f6..."
}
```

Run a query:

```sh
steampipe query "select ip, country_name, city from ipgeolocation_ip where ip = '8.8.8.8'"
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads) ≥ 0.20
- [Go](https://golang.org/dl/) 1.26+

Clone the repo, build, and install:

```sh
git clone https://github.com/IPGeolocation/steampipe-plugin-ipgeolocation
cd steampipe-plugin-ipgeolocation
make setup
```

`make setup` runs `go mod tidy`, builds the plugin binary, and installs both the binary and the sample config file.

To rebuild after making changes:

```sh
make install
```

Run tests:

```sh
make test
```

## Documentation

Full table documentation is available on [Steampipe Hub](https://hub.steampipe.io/plugins/ipgeolocation/ipgeolocation/tables) and in the [`docs/`](docs/tables) directory.

| Table | Endpoint | Description | Credits |
|---|---|---|---|
| [ipgeolocation_ip](docs/tables/ipgeolocation_ip.md) | `/v3/ipgeo` | Geolocation, network, timezone, security, company, abuse, hostname | 1–4 |
| [ipgeolocation_security](docs/tables/ipgeolocation_security.md) | `/v3/security` | VPN, proxy, Tor, relay, bot, spam, and threat score | 2 |
| [ipgeolocation_abuse](docs/tables/ipgeolocation_abuse.md) | `/v3/abuse` | Abuse contact emails, phone numbers, org, and CIDR route | 1 |
| [ipgeolocation_asn](docs/tables/ipgeolocation_asn.md) | `/v3/asn` | ASN details, routes, peers, upstreams, downstreams, WHOIS | 1 |

## License

This plugin is licensed under the [Apache 2.0 License](LICENSE).