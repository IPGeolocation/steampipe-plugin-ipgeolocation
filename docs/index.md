---
organization: ipgeolocation
category: ["saas", "internet", "security"]
icon_url: "/images/plugins/ipgeolocation/ipgeolocation.svg"
brand_color: "#FF6B35"
display_name: "IPGeolocation.io"
short_name: "ipgeolocation"
description: "Use Steampipe to query IP geolocation, threat intelligence, ASN details, and abuse contacts from IPGeolocation.io."
og_description: "Query IP geolocation, security signals, ASN relationships, and abuse contacts with SQL using Steampipe."
og_image: "/images/plugins/ipgeolocation/ipgeolocation-social.png"
engines: ["steampipe", "sqlite", "postgres", "export"]
---

# IPGeolocation.io + Steampipe

[IPGeolocation.io](https://ipgeolocation.io) provides real-time IP intelligence — geolocation, VPN/proxy/Tor detection, ASN data, and abuse contacts — for any IPv4 or IPv6 address via a fast REST API.

[Steampipe](https://steampipe.io) is an open-source tool for querying cloud APIs using SQL.

For example:

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

```
+----------+--------------+---------------+------------------+-------------------+
| ip       | country_name | city          | isp              | timezone_name     |
+----------+--------------+---------------+------------------+-------------------+
| 8.8.8.8  | United States| Mountain View | Google LLC       | America/Chicago   |
+----------+--------------+---------------+------------------+-------------------+
```

## Documentation

- **[Table: ipgeolocation_ip](https://hub.steampipe.io/plugins/ipgeolocation/ipgeolocation/tables/ipgeolocation_ip)** — full geolocation, network, timezone, and hostname data for any IP
- **[Table: ipgeolocation_security](https://hub.steampipe.io/plugins/ipgeolocation/ipgeolocation/tables/ipgeolocation_security)** — VPN, proxy, Tor, relay, bot, spam, and threat score signals
- **[Table: ipgeolocation_abuse](https://hub.steampipe.io/plugins/ipgeolocation/ipgeolocation/tables/ipgeolocation_abuse)** — abuse contact emails, phone numbers, and responsible organisation
- **[Table: ipgeolocation_asn](https://hub.steampipe.io/plugins/ipgeolocation/ipgeolocation/tables/ipgeolocation_asn)** — ASN details, announced routes, peers, upstreams, downstreams, and WHOIS

## Get Started

### Install the plugin

```sh
steampipe plugin install ipgeolocation/ipgeolocation
```

### Configure credentials

```sh
steampipe plugin configure ipgeolocation
```

Or create `~/.steampipe/config/ipgeolocation.spc` manually:

```hcl
connection "ipgeolocation" {
  plugin = "ipgeolocation/ipgeolocation"

  # Get your API key at https://app.ipgeolocation.io/dashboard
  # Free tier works without a key for basic geolocation.
  # api_key = "YOUR_API_KEY_HERE"
}
```

You can also set the key via environment variable:

```sh
export IPGEOLOCATION_API_KEY="YOUR_API_KEY_HERE"
```

### Run your first query

```sh
steampipe query "select ip, country_name, city from ipgeolocation_ip where ip = '1.1.1.1'"
```

## Example Queries

### Geolocate an IP

```sql
select
  ip,
  country_name,
  state_prov,
  city,
  latitude,
  longitude,
  isp
from
  ipgeolocation_ip
where
  ip = '8.8.8.8';
```

### Detect VPN, proxy, and Tor usage

```sql
select
  ip,
  threat_score,
  is_vpn,
  vpn_provider_names,
  is_proxy,
  is_tor,
  is_known_attacker
from
  ipgeolocation_security
where
  ip = '185.220.101.45';
```

### Find the abuse contact for a suspicious IP

```sql
select
  ip,
  name,
  organization,
  emails,
  phone_numbers,
  route,
  country
from
  ipgeolocation_abuse
where
  ip = '91.128.103.196';
```

### Look up ASN details by IP

```sql
select
  ip,
  as_number,
  organization,
  country,
  type,
  num_of_ipv4_routes,
  num_of_ipv6_routes
from
  ipgeolocation_asn
where
  ip = '8.8.8.8';
```

### Look up ASN details by ASN number

```sql
select
  as_number,
  organization,
  country,
  type,
  domain,
  rir
from
  ipgeolocation_asn
where
  asn = '15169';
```

### Expand ASN peers into individual rows

```sql
select
  as_number,
  organization,
  peer ->> 'as_number'   as peer_asn,
  peer ->> 'description' as peer_name,
  peer ->> 'country'     as peer_country
from
  ipgeolocation_asn,
  jsonb_array_elements(peers) as peer
where
  asn = '15169';
```

### Check if an IP is in a cloud provider

```sql
select
  ip,
  is_cloud_provider,
  cloud_provider_name,
  is_vpn,
  threat_score
from
  ipgeolocation_security
where
  ip = '13.32.0.1';
```

Sign up at [ipgeolocation.io](https://app.ipgeolocation.io) to get your API key.