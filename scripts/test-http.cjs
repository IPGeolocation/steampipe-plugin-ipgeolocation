'use strict';

const { spawn } = require('node:child_process');
const { once } = require('node:events');
const path = require('node:path');

async function main() {
  const proxy = spawn('postman', [
    'simulation', 'run', 'postman/simulations/steampipe/simulation.yaml',
    '--proxy', '--proxy-port', '7012', '--no-history', '--log-format', 'json',
  ], {
    cwd: path.resolve(__dirname, '..'), stdio: ['ignore', 'pipe', 'pipe'],
    detached: process.platform !== 'win32',
  });
  const closed = new Promise((resolve) => proxy.once('close', resolve));
  let output = '';
  let errors = '';
  let timer;
  try {
    const certificate = await new Promise((resolve, reject) => {
      timer = setTimeout(() => reject(new Error(`Postman proxy did not start: ${errors}`)), 20000);
      proxy.once('error', reject);
      proxy.once('exit', (code) => reject(new Error(`Postman proxy exited (${code}): ${errors}`)));
      proxy.stderr.on('data', (chunk) => { errors += chunk; });
      proxy.stdout.on('data', (chunk) => {
        output += chunk;
        const match = output.match(/NODE_EXTRA_CA_CERTS=([^\r\n]+)/);
        if (match) resolve(match[1].trim());
      });
    });
    clearTimeout(timer);
    const checks = spawn('go', ['test', './ipgeolocation', '-run', '^TestHTTPSIntegration$', '-count=1', '-v'], {
      cwd: path.resolve(__dirname, '..'), stdio: 'inherit',
      env: {
        ...process.env,
        HTTPS_PROXY: 'http://127.0.0.1:7012', HTTP_PROXY: 'http://127.0.0.1:7012',
        NO_PROXY: '', SSL_CERT_FILE: certificate, STEAMPIPE_HTTP_TEST: '1',
        // shortcut: allow Postman 1.73's negative test-certificate serials until the CLI fixes them.
        GODEBUG: [process.env.GODEBUG, 'x509negativeserial=1'].filter(Boolean).join(','),
      },
    });
    const [code] = await once(checks, 'exit');
    process.exitCode = code === 0 ? 0 : 1;
    console.log(output);
    if (process.exitCode) console.error(errors);
  } finally {
    clearTimeout(timer);
    if (proxy.pid) {
      // The npm-installed CLI launches a child; stop the entire proxy process group.
      const stop = (signal) => {
        try {
          if (process.platform === 'win32') proxy.kill(signal);
          else process.kill(-proxy.pid, signal);
        } catch (err) {
          if (err.code !== 'ESRCH') throw err;
        }
      };
      const force = setTimeout(() => stop('SIGKILL'), 2000);
      stop('SIGINT');
      await closed;
      clearTimeout(force);
    }
  }
}

main().catch((err) => { console.error(err.message); process.exitCode = 1; });
