## v0.1.4 [2026-10-01]

- Updated Go dependencies to fix security vulnerabilities (gRPC, OpenTelemetry, golang.org/x/crypto, x/net, x/text, go-getter and others). No changes to tables, columns or query results.
- Builds now require Go 1.26, and release binaries are built with the latest Go 1.26 patch release.

## v0.1.2 [2026-05-12]

- fix abuse and asn docs description
- fix Makefile issues that was targeting local plugins and added netgo tag for building binaries.

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