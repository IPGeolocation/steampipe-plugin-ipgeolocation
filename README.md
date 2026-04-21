# Steampipe Plugin — IPGeolocation.io

Query IP geolocation, network, security, timezone, company, abuse, and hostname data using SQL via the [IPGeolocation.io](https://ipgeolocation.io) v3 API.

## Quick Start

### 1. Install the plugin

```sh
steampipe plugin install ipgeolocation
```

### 2. Configure credentials

```sh
cp config/ipgeolocation.spc ~/.steampipe/config/ipgeolocation.spc
# edit the file and set your api_key
```

### 3. Query away

```sql
-- Where is 8.8.8.8?
select ip, country_name, city, isp
from ipgeolocation_ip
where ip = '8.8.8.8';

-- Is it a VPN/proxy/Tor?
select ip, is_vpn, is_proxy, is_tor, threat_score
from ipgeolocation_ip
where ip = '185.220.101.45';
```

## Requirements

| Requirement | Version |
|---|---|
| Steampipe | ≥ 0.20 |
| Steampipe Plugin SDK | v5 |
| Go | ≥ 1.21 |

## Building from Source

```sh
git clone https://github.com/turbot/steampipe-plugin-ipgeolocation
cd steampipe-plugin-ipgeolocation
make
```

## Configuration Reference

```hcl
connection "ipgeolocation" {
  plugin  = "ipgeolocation"

  # Required for most modules. Free tier covers basic location data.
  # Sign up at https://app.ipgeolocation.io/
  api_key = "YOUR_API_KEY_HERE"
}
```

## Tables

| Table | Description |
|---|---|
| [ipgeolocation_ip](docs/tables/ipgeolocation_ip.md) | Geolocation data for a given IP address or domain |

## API Plan Notes

| Module | Free | Paid |
|---|---|---|
| Location (country, city, lat/lon) | ✅ | ✅ |
| Timezone | ✅ | ✅ |
| Network / ASN | ✅ | ✅ |
| Hostname | ❌ | ✅ |
| Security (VPN, proxy, threat score) | ❌ | ✅ |
| Company | ❌ | ✅ |
| Abuse contact | ❌ | ✅ |

Free-plan columns will return `null` when the field is not available on your subscription tier.
