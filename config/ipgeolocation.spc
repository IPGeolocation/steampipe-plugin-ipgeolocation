connection "ipgeolocation" {
  plugin = "local/ipgeolocation"

  # API key from https://app.ipgeolocation.io/dashboard
  # Free tier works without a key for basic geolocation.
  # Paid plans unlock hostname, security, company, abuse, and user-agent modules.
  # Can also be set with the IPGEOLOCATION_API_KEY environment variable.
  api_key = "a3fbecf6536a470b8d3426419f21ae99"
}
