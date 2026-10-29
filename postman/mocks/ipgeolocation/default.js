const http = require('node:http');

http.createServer((req, res) => {
  const url = new URL(req.url, 'http://localhost');
  const reply = (status, body) => { res.writeHead(status, { 'Content-Type': 'application/json' }); res.end(JSON.stringify(body)); };
  if (!['/v3/ipgeo', '/v3/security', '/v3/abuse', '/v3/asn'].includes(url.pathname)) return reply(501, { message: 'Endpoint not modeled' });
  if (req.method !== 'GET') return reply(400, { message: 'Expected GET' });
  const key = url.searchParams.get('apiKey');
  const statuses = { BAD_KEY: 401, FORBIDDEN: 403, RATE_LIMIT: 429, SERVER_ERROR: 503 };
  if (statuses[key]) return reply(statuses[key], { message: 'Test failure response' });
  if (key === 'BAD_JSON') { res.writeHead(200); return res.end('{'); }
  if (key === 'SLOW') return setTimeout(() => { if (!res.destroyed) reply(200, {}); }, 250);
  if (key !== 'TEST_KEY' && !(key === null && url.pathname === '/v3/ipgeo')) return reply(401, { message: 'Invalid test key' });
  const ip = url.searchParams.get('ip') || '198.51.100.42';
  if (url.pathname === '/v3/ipgeo') {
    if (url.searchParams.get('include') !== '*') return reply(400, { message: 'Expected include=*' });
    return reply(200, { ip, location: { city: 'New Delhi' }, time_zone: { offset: 5.5, offset_with_dst: 5.5 } });
  }
  if (url.pathname === '/v3/security') return reply(200, { ip, security: { is_proxy: false } });
  if (url.pathname === '/v3/abuse') return reply(200, { ip, abuse: { emails: ['abuse@example.test'] } });
  if (url.searchParams.get('include') !== 'peers,upstreams,downstreams,routes,whois_response') return reply(400, { message: 'Missing ASN relationship fields' });
  if (url.searchParams.get('asn') !== '15169' && !url.searchParams.get('ip')) return reply(400, { message: 'Expected an IP or normalized ASN' });
  reply(200, { ip, asn: { as_number: 'AS15169' } });
}).listen(process.env.PORT || 4500);
