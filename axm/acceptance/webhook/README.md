# Webhook acceptance tests

This package contains the temporary Go receiver, Cloudflare tunnel integration,
and Apple Business API webhook lifecycle tests. The receiver binds to loopback
on an ephemeral port; its public HTTPS URL exists only while the tunnel runs.

## Local receiver tests

```sh
go test -race ./axm/acceptance/webhook
```

Without `AXM_WEBHOOK_E2E=true`, the live lifecycle is skipped. Local tests do not
download binaries, start tunnels, or contact Apple. They exercise HMAC signatures,
invalid requests, the 1 MiB body limit, concurrent receipt, and cancellation.

## Live lab lifecycle

Supply an Apple Business API account with webhook permissions using the SDK's
normal environment variables. Keep credential values outside the repository:

```sh
export APPLE_KEY_ID='<key-id>'
export APPLE_ISSUER_ID='<client-id>'
export APPLE_PRIVATE_KEY_PATH='/absolute/path/to/key.p8'
# Alternatively, set APPLE_PRIVATE_KEY_PEM to the PEM contents.
AXM_WEBHOOK_E2E=true ./axm/acceptance/webhook/run.sh
```

The runner requires Bash, Go, curl, and tar on macOS. It downloads the official
`cloudflared` release into temporary storage for macOS/Linux arm64/amd64. To use
an existing binary, set `CLOUDFLARED_BIN` to its absolute path. The public tunnel
must permit unauthenticated Apple HTTP requests; the receiver authenticates them
using the random per-run HMAC secret and `x-apple-signature` header.

For repeated lab runs, an existing unexpired OAuth bearer token can be supplied
through `AXM_WEBHOOK_ACCESS_TOKEN` to avoid issuing a fresh token each time.
The standard `APPLE_*` credentials still configure the SDK client. Keep the
token private and outside the repository; it is never printed by these tests.

The ten-minute run first checks webhook API access, failing before it starts
infrastructure when the account lacks permission. It then:

1. Starts the receiver, waits for tunnel registration, allows five seconds for
   DNS propagation, and checks public DNS publication,
   then checks the temporary HTTPS endpoint. Checking DNS publication first avoids
   caching an early “hostname not found” response in the system resolver.
2. Creates a uniquely named webhook with a random signing secret.
3. Lists and reads it, then disables and re-enables it with PATCH.
4. Sends a ping and verifies signed `webhookPingCreated` receipt.
5. Polls delivery records, includes the event, and correlates its payload to the
   received notification until delivery succeeds.
6. Replays that delivery, checks another signed receipt, and confirms the replay
   succeeded for the original event.
7. Deletes the test webhook and confirms that reading it returns 404.

The signed ping must carry the ID returned by `CreatePingV1`. Apple's stored
ping event payload can contain an empty ID; in that case the test matches its
nonzero creation timestamp and notification type to the signed receipt. It then
checks the delivery's event relationship and `SUCCEEDED` state. A nonempty,
different stored ping ID never matches.
Both delivery polling phases request `include=event`; the live API omits the
event relationship without it, even when `event` is in the delivery field list.

Each delivery phase polls every two seconds for up to 90 seconds. The default
Apple request timeout is 30 seconds (`AXM_REQUEST_TIMEOUT` overrides it).
Delivery polling honors `Retry-After` when Apple responds with HTTP 429, while
retaining the phase deadline.
Delivery queries include both date filters required by Apple's live service,
with a start before the test and an end refreshed on every poll.
Deletion cleanup uses a separate 30-second context, including on test failure
or a handled SIGINT/SIGTERM. The runner removes its temporary binary and stops
the receiver and tunnel. Cleanup never targets pre-existing webhooks or changes
device assignments; it runs even when `AXM_SKIP_CLEANUP=true`.

Logs show operation/status evidence without credentials, signing secrets, or raw
notification bodies. A forced process kill or machine failure cannot execute
cleanup: if interrupted that way, inspect the lab for this run's uniquely named
`go-sdk-webhook-...` resource before deleting it manually.

Cloudflare Quick Tunnels have no uptime guarantee. A tunnel or Apple permission
failure is a failed live test, not evidence that delivery succeeded. See the
[Cloudflare instructions](https://developers.cloudflare.com/tunnel/get-started/quick-tunnels/)
and [Apple's event and signature documentation](https://developer.apple.com/documentation/apple-school-and-business-manager-api/understanding-apple-business-api-webhook-events).

Tunnel readiness also requires HTTPS access to `cloudflare-dns.com/dns-query`.
This checks only the temporary public hostname; Apple credentials and signing
secrets are never sent to that endpoint. Normal system DNS and TLS verification
remain in use for the subsequent HTTPS health check.

If the initial webhook list request returns `403 FORBIDDEN_ERROR`, review the
API account's role in Apple Business under Settings > API. Successful device or
MDM-server requests do not demonstrate access to the webhook endpoints. See
[Apple's API account and role instructions](https://support.apple.com/en-ca/guide/business/axm33189f66a/web).
The permission preflight fails before starting the receiver or creating a webhook.

To check HTTPS connectivity and signed synthetic notification receipt with an
existing `cloudflared` binary, without using Apple credentials or creating
webhook resources:

```sh
AXM_WEBHOOK_TUNNEL_CHECK=true CLOUDFLARED_BIN='/absolute/path/to/cloudflared' \
  go test -v ./axm/acceptance/webhook -run '^TestTunnelHTTPS$' -timeout 2m
```
