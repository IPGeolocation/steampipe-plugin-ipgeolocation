# Table: ipgeolocation_ip

Look up geolocation, network, timezone, security, company, abuse, and hostname data for any IPv4, IPv6, or domain (paid only) using the [IPGeolocation.io](https://ipgeolocation.io) v3 API.

## Configuration

Configure your API key in `~/.steampipe/config/ipgeolocation.spc`:

```hcl
connection "ipgeolocation" {
  plugin  = "ipgeolocation/ipgeolocation"

  # API key from https://app.ipgeolocation.io/dashboard (Required)
  api_key = "YOUR_API_KEY_HERE"
}
```

A free-tier key from [ipgeolocation.io](https://app.ipgeolocation.io/dashboard) covers basic location data. Paid plans unlock security, hostname, company, abuse data.

---

## Examples

### Look up a specific IP address

```sql
select
  ip,
  country_name,
  city,
  state_prov,
  latitude,
  longitude,
  isp
from
  ipgeolocation_ip
where
  ip = '8.8.8.8';
```

### Look up your own public IP

```sql
select
  ip,
  country_name,
  city,
  timezone_name
from
  ipgeolocation_ip;
```

### Check security signals for a suspicious IP (paid plan)

```sql
select
  ip,
  threat_score,
  is_vpn,
  is_proxy,
  is_tor,
  is_bot,
  is_known_attacker
from
  ipgeolocation_ip
where
  ip = '185.220.101.45';
```

### Get full ASN and network details

```sql
select
  ip,
  asn,
  organization,
  connection_type
from
  ipgeolocation_ip
where
  ip = '1.1.1.1';
```

### Get timezone information

```sql
select
  ip,
  timezone_name,
  timezone_current_time,
  timezone_offset,
  timezone_is_dst
from
  ipgeolocation_ip
where
  ip = '203.0.113.1';
```

### Get company and abuse contact data (paid plan)

```sql
select
  ip,
  company_name,
  company_domain,
  company_type,
  abuse_email,
  abuse_phone
from
  ipgeolocation_ip
where
  ip = '104.21.0.1';
```

---

## Column Reference

| Column | Type | Description |
|---|---|---|
| ip | text | Queried IP address |
| hostname | text | Resolved hostname (paid) |
| continent_code | text | Two-letter continent code |
| continent_name | text | Continent name |
| country_code2 | text | ISO alpha-2 country code |
| country_code3 | text | ISO alpha-3 country code |
| country_name | text | Country name |
| country_name_official | text | Official country name |
| country_capital | text | Capital city |
| state_prov | text | State / province |
| state_code | text | ISO 3166-2 state code |
| district | text | District / county |
| city | text | City |
| locality | text | Locality (Paid) |
| zipcode | text | ZIP / postal code |
| latitude | text | Latitude |
| longitude | text | Longitude |
| accuracy_radius | text | Accuracy radius (km) (Paid) |
| confidence | text | Confidence level (low, medium, high) (Paid) |
| is_eu | bool | True if EU member state |
| country_flag | text | Country flag image URL |
| country_emoji | text | Country flag emoji |
| geoname_id | text | GeoNames ID |
| dma_code | text | DMA code (US only) (Paid) |
| asn | text | Autonomous System Number |
| organization | text | ASN organization |
| connection_type | text | Connection type |
| timezone_name | text | IANA timezone name |
| timezone_offset | bigint | UTC offset (seconds) |
| timezone_offset_with_dst | bigint | UTC offset incl. DST |
| timezone_current_time | text | Current local time |
| timezone_is_dst | bool | DST active? |
| threat_score | bigint | Threat score 0-100 (paid) |
| is_vpn | bool | VPN exit node? (paid) |
| is_proxy | bool | Proxy? (paid) |
| is_tor | bool | Tor exit node? (paid) |
| is_bot | bool | Bot IP? (paid) |
| is_spam | bool | On spam lists? (paid) |
| is_cloud_provider | bool | Cloud/hosting IP? (paid) |
| is_known_attacker | bool | Known attacker? (paid) |
| is_anonymous | bool | Anonymous IP? (paid) |
| is_residential_proxy | bool | Residential proxy? (paid) |
| company_name | text | Company name (paid) |
| company_domain | text | Company domain (paid) |
| company_type | text | Company type (paid) |
| abuse_name | text | Abuse contact name (paid) |
| abuse_email | text | Abuse email (paid) |
| abuse_address | text | Abuse postal address (paid) |
| abuse_phone | text | Abuse phone (paid) |
| abuse_network | text | Abuse CIDR block (paid) |
