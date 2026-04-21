# Table: ipgeolocation_asn

ASN lookup using the [IPGeolocation.io](https://ipgeolocation.io) `/v3/asn` endpoint.

Look up by IP address **or** ASN number. Returns the owning organisation, registration metadata (RIR, allocation date, status), announced route counts, and optionally the full list of routes, peers, upstreams, downstreams, and raw WHOIS text.

> **Note:** Costs **1 credit** per request.

---

## Examples

### Look up the ASN for an IP address

```sql
select
  ip,
  asn_number,
  asn_name,
  organization,
  country,
  type,
  rir
from
  ipgeolocation_asn
where
  ip = '8.8.8.8';
```

### Look up an ASN directly by number

```sql
select
  asn_number,
  asn_name,
  organization,
  country,
  type,
  domain,
  date_allocated,
  rir
from
  ipgeolocation_asn
where
  asn_number = 'AS15169';
```

### Check how many routes an ASN announces

```sql
select
  asn_number,
  organization,
  num_of_ipv4_routes,
  num_of_ipv6_routes
from
  ipgeolocation_asn
where
  asn_number = 'AS13335';
```

### List all announced IP prefixes (routes)

```sql
select
  asn_number,
  organization,
  jsonb_array_elements_text(routes) as prefix
from
  ipgeolocation_asn
where
  asn_number = 'AS15169';
```

### Explore peering relationships

```sql
select
  asn_number,
  organization,
  peers
from
  ipgeolocation_asn
where
  asn_number = 'AS12';
```

### Expand peers into individual rows

```sql
select
  asn_number,
  organization,
  peer ->> 'as_number'   as peer_asn,
  peer ->> 'description' as peer_name,
  peer ->> 'country'     as peer_country
from
  ipgeolocation_asn,
  jsonb_array_elements(peers) as peer
where
  asn_number = 'AS12';
```

### Show upstream transit providers

```sql
select
  asn_number,
  organization,
  upstream ->> 'as_number'   as upstream_asn,
  upstream ->> 'description' as upstream_name,
  upstream ->> 'country'     as upstream_country
from
  ipgeolocation_asn,
  jsonb_array_elements(upstreams) as upstream
where
  asn_number = 'AS12';
```

### Get raw WHOIS text

```sql
select
  asn_number,
  whois_response
from
  ipgeolocation_asn
where
  asn_number = 'AS15169';
```

---

## Column Reference

| Column | Type | Description |
|---|---|---|
| ip | text | IP used for lookup (when queried by IP) |
| asn_number | text | AS number e.g. "AS15169" |
| asn_name | text | Short registered ASN name |
| organization | text | Organisation owning the ASN |
| country | text | ISO alpha-2 country of registration |
| type | text | ASN type: ISP, HOSTING, EDUCATION, GOVERNMENT, BUSINESS |
| domain | text | Organisation's primary domain |
| date_allocated | text | Allocation date (YYYY-MM-DD) |
| allocation_status | text | e.g. "assigned" |
| rir | text | Regional Internet Registry (ARIN, RIPE, APNIC, LACNIC, AFRINIC) |
| num_of_ipv4_routes | bigint | Number of IPv4 prefixes announced |
| num_of_ipv6_routes | bigint | Number of IPv6 prefixes announced |
| routes | jsonb | Array of CIDR prefixes announced |
| peers | jsonb | Array of peer ASNs |
| upstreams | jsonb | Array of upstream/transit ASNs |
| downstreams | jsonb | Array of downstream customer ASNs |
| whois_response | text | Raw WHOIS text |
| raw | jsonb | Full raw API response |