## v0.1.1 [2026-04-24]
_Minor bug fixes + suggestions by Steampipe reviewer_

- Fixed bug in `ipgeolocation_asn` table where `as_number` is being written as `asn_number`
- Fixed return values of `ipgeolocation_ip` where some of fields are pointing to non-existent fields of API
- Removed `raw` field from all tables
- Fixed typos in table documentation
- Made the `ip` field required in all tables 

## v0.1.0 [2026-04-21]

_Initial release_

- New tables added
  - `ipgeolocation_ip`
  - `ipgeolocation_security`
  - `ipgeolocation_abuse`
  - `ipgeolocation_asn`