# Table: ipgeolocation_abuse

Abuse contact lookup using the [IPGeolocation.io](https://ipgeolocation.io) `/v3/abuse` endpoint.

Returns the responsible organisation, role, emails, phone numbers, postal address, CIDR route, and registered country for any IPv4 or IPv6 address — everything you need to file an abuse report.

> **Note:** Requires a **paid plan** API key. Costs **1 credit** per request.

---

## Examples

### Look up the abuse contact for an IP

```sql
select
  ip,
  name,
  organization,
  emails,
  phone_numbers,
  country,
  route
from
  ipgeolocation_abuse
where
  ip = '1.0.0.0';
```

### Get only the abuse email addresses

```sql
select
  ip,
  emails
from
  ipgeolocation_abuse
where
  ip = '8.8.8.8';
```

---

## Column Reference

| Column | Type | Description |
|---|---|---|
| ip | text | Queried IP address |
| route | text | CIDR block covering the IP |
| country | text | ISO alpha-2 country of registrant |
| name | text | Abuse contact or IRT name |
| organization | text | Responsible organisation |
| kind | text | Contact kind (e.g. "group") |
| address | text | Postal address |
| emails | jsonb | Array of abuse email addresses |
| phone_numbers | jsonb | Array of abuse phone numbers |