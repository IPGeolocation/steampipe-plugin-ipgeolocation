# Table: ipgeolocation_security

IP threat intelligence using the [IPGeolocation.io](https://ipgeolocation.io) `/v3/security` endpoint.

Returns VPN, proxy, Tor, relay, residential proxy, bot, spam, known-attacker, and cloud-provider signals — each with confidence scores, provider names, and last-seen dates. Also includes basic location and network context.

> **Note:** This endpoint requires a **paid plan** API key and costs **2 credits per request**.

---

## Examples

### Check a single suspicious IP

```sql
select
  ip,
  threat_score,
  is_vpn,
  is_proxy,
  is_tor,
  is_known_attacker,
  is_bot,
  is_spam
from
  ipgeolocation_security
where
  ip = '185.220.101.45';
```

### Identify VPN provider details

```sql
select
  ip,
  is_vpn,
  vpn_provider_names,
  vpn_confidence_score,
  vpn_last_seen
from
  ipgeolocation_security
where
  ip = '2.56.188.34';
```

### Check proxy details with confidence

```sql
select
  ip,
  is_proxy,
  proxy_provider_names,
  proxy_confidence_score,
  proxy_last_seen,
  is_residential_proxy
from
  ipgeolocation_security
where
  ip = '2.56.188.34';
```

### Cloud / hosting provider check

```sql
select
  ip,
  is_cloud_provider,
  cloud_provider_name,
  asn,
  isp
from
  ipgeolocation_security
where
  ip = '13.32.0.1';
```

### Relay detection (Apple Private Relay, iCloud)

```sql
select
  ip,
  is_relay,
  relay_provider_name,
  is_anonymous,
  threat_score
from
  ipgeolocation_security
where
  ip = '17.58.97.1';
```

### Quick risk triage — high threat score

```sql
select
  ip,
  threat_score,
  is_vpn,
  is_proxy,
  is_tor,
  is_residential_proxy,
  is_known_attacker,
  country_name,
  isp
from
  ipgeolocation_security
where
  ip = '91.108.56.1'
  and threat_score > 50;
```

### Inspect full raw response

```sql
select
  ip,
  raw
from
  ipgeolocation_security
where
  ip = '8.8.8.8';
```

---

## Column Reference

| Column | Type | Description |
|---|---|---|
| ip | text | Queried IP address |
| threat_score | bigint | 0–100 overall risk score |
| is_anonymous | bool | Any anonymization detected |
| is_vpn | bool | Known VPN exit node |
| vpn_provider_names | jsonb | Array of VPN provider names |
| vpn_confidence_score | bigint | VPN detection confidence 0–100 |
| vpn_last_seen | text | Last seen as VPN (YYYY-MM-DD) |
| is_proxy | bool | Known proxy |
| proxy_provider_names | jsonb | Array of proxy provider names |
| proxy_confidence_score | bigint | Proxy detection confidence 0–100 |
| proxy_last_seen | text | Last seen as proxy (YYYY-MM-DD) |
| is_residential_proxy | bool | Residential proxy |
| is_tor | bool | Tor exit node |
| is_relay | bool | Relay network (e.g. iCloud Private Relay) |
| relay_provider_name | text | Relay provider name |
| is_cloud_provider | bool | Cloud or hosting IP |
| cloud_provider_name | text | Cloud provider name |
| is_known_attacker | bool | Flagged for attack activity |
| is_bot | bool | Bot activity detected |
| is_spam | bool | On spam block lists |
| raw | jsonb | Full raw API response |