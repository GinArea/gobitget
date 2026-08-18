# Bitget UTA v3 API Reference — Exchange-Specific Quirks

This document records live-verified behavior of the Bitget UTA v3 API (and of the legacy v2 WS
used for candle streaming) that is missing from or contradicts the official docs
(https://www.bitget.com/api-doc/uta/). Everything was verified against `https://api.bitget.com`;
the scope of each verification (categories, symbols) is noted per section.

---

## Candlestick Intervals and UTC Alignment

Applies to `GET /api/v3/market/candles` (REST). Grids are **identical across all categories**
(SPOT, USDT-FUTURES, COIN-FUTURES).

The official v3 docs list only 10 intervals (`1m, 3m, 5m, 15m, 30m, 1H, 4H, 6H, 12H, 1D`).
In reality the endpoint also accepts undocumented values inherited from the v2 API, including
`utc`-suffixed variants. Native intervals of 6 hours and longer open on the **UTC+8 grid**
(00:00 Asia/Shanghai = 16:00 UTC); the `utc`-suffixed variants open on the **UTC grid**.
Intervals up to 4H need no utc variant: their period divides the 8-hour offset, so both grids
coincide.

### Reference table: which interval to use for UTC-aligned candle boundaries

| Period    | Native interval | Native grid (candle open)        | UTC-aligned? | Use for UTC boundaries |
|-----------|-----------------|----------------------------------|--------------|------------------------|
| 1 minute  | `1m`            | grids coincide                   | yes          | `1m`                   |
| 3 minutes | `3m`            | grids coincide                   | yes          | `3m`                   |
| 5 minutes | `5m`            | grids coincide                   | yes          | `5m`                   |
| 15 minutes| `15m`           | grids coincide                   | yes          | `15m`                  |
| 30 minutes| `30m`           | grids coincide                   | yes          | `30m`                  |
| 1 hour    | `1H`            | grids coincide                   | yes          | `1H`                   |
| 2 hours   | `2H` †          | grids coincide                   | yes          | `2H` †                 |
| 4 hours   | `4H`            | grids coincide                   | yes          | `4H`                   |
| 6 hours   | `6H`            | UTC+8 grid (04/10/16/22 UTC)     | **no**       | **`6Hutc`** †          |
| 12 hours  | `12H`           | UTC+8 grid (04:00/16:00 UTC)     | **no**       | **`12Hutc`** †         |
| 1 day     | `1D`            | 00:00 UTC+8 (16:00 UTC)          | **no**       | **`1Dutc`** †          |
| 3 days    | `3D` †          | 00:00 UTC+8 grid                 | **no**       | **`3Dutc`** †          |
| 1 week    | `1W` †          | Monday 00:00 UTC+8 (Sun 16:00 UTC)| **no**      | **`1Wutc`** † (Monday 00:00 UTC) |
| 1 month   | `1M` †          | 1st day 00:00 UTC+8              | **no**       | **`1Mutc`** † (1st day 00:00 UTC) |

† — undocumented in the v3 docs (inherited from the v2 API), verified live.

### Intervals that do NOT exist

`3H`, `8H`, `1Hutc`, `4Hutc` — the request fails with HTTP 400, body
`{"code":"40020","msg":"Parameter <value> error"}`.

### Interval case matters

Minutes are lowercase (`1m`), hours/days/weeks/months are uppercase (`1H`, `1D`, `1W`, `1M`),
the `utc` suffix is lowercase (`1Dutc`).

---

## Other verified `GET /api/v3/market/candles` quirks

- **`limit`**: default 100, maximum 1000 (the official docs table swaps these values);
  `limit > 1000` fails with parameter error 40020.
- **Time window bounds**: `startTime` is **exclusive**, `endTime` is **inclusive** —
  the response contains candles with `startTime < ts <= endTime`.
- **Ordering**: rows are sorted oldest-first.
- **Row shape**: `[ts, open, high, low, close, volume, turnover]`, every element is a JSON string;
  identical for all categories.
- **`type=mark`** candles have `volume` and `turnover` equal to 0.
- **History depth** is limited and varies by interval: e.g. 1m candles are unavailable a few
  months back, and the monthly (`1M`) history of BTCUSDT USDT-FUTURES held only 4 rows —
  a request may return fewer rows than `limit` without an error.

## WebSocket kline intervals: v3 vs legacy v2

Verified live on SPOT and USDT-FUTURES (BTCUSDT).

### v3 WS (`wss://ws.bitget.com/v3/ws/public`, topic `kline`)

Accepts **only** the 10 documented intervals: `1m, 3m, 5m, 15m, 30m, 1H, 4H, 6H, 12H, 1D`.
Every other value — `2H`, the utc-suffixed variants (`6Hutc`/`12Hutc`/`1Dutc`/...), `3D`, `1W`,
`1M`, `3H`, `8H` and any case variant (`2h`, `1DUTC`, `1dutc`, `12hutc`) — is rejected with
error `30001` `"{...} doesn't exist"`. Unlike REST, the v3 WS does **not** inherit the v2 values.
Native `6H`/`12H`/`1D` open on the UTC+8 grid, same as REST (confirmed on WS).

Consequence: over the v3 WS, UTC-aligned candles exist only for intervals up to `4H`.
For this reason `bitgetv3` deliberately does not wrap the v3 kline channel (the v3 ticker
channel IS wrapped): WS candle streaming goes exclusively through the legacy v2 WS
(`WsPublicV2.Candle`).

### Legacy v2 WS (`wss://ws.bitget.com/v2/ws/public`, channel `candle<tf>`)

Subscription format (differs from v3):

```json
{"op":"subscribe","args":[{"instType":"SPOT","channel":"candle1Dutc","instId":"BTCUSDT"}]}
```

Accepted channels (both SPOT and USDT-FUTURES):

| Channels | Grid |
|----------|------|
| `candle1m, 3m, 5m, 15m, 30m, 1H, 2H, 4H, 8H` | UTC-aligned (period divides 8h; note `8H`: 00/08/16 UTC) |
| `candle6H, 12H, 1D, 3D, 1W, 1M` | UTC+8 grid (weeks Monday, months 1st, 00:00 UTC+8) |
| `candle6Hutc, 12Hutc, 1Dutc, 3Dutc, 1Wutc, 1Mutc` | UTC grid (weeks Monday, months 1st, 00:00 UTC) |

Rejected: `candle3H` — error `30016` `"Param error"` (note: a different code than the v3 WS 30001).
`8H` exists on the v2 WS but NOT in the v3 REST/WS at all.

v2 push format: same envelope shape (`arg`/`action`/`data`, `snapshot` then `update`), keepalive is
the same text `ping`/`pong`. The initial snapshot carries history like v3 (500 rows or the full
available depth, oldest-first, current candle last). Rows are arrays of **8 strings**:
`[ts, open, high, low, close, baseVolume, quoteVolume, usdtVolume]` — the field order was
disambiguated on a BTC-quoted pair (ETHBTC): index 6 is the quote-coin volume, index 7 the USDT
volume (they coincide for USDT-quoted symbols). Ack and error events echo the numeric `code`
(e.g. 30016) like v3. In `bitgetv3` this WS is wrapped by `WsPublicV2` (`ws_public_v2.go`).

### UTC-aligned WS candles: which source to use per timeframe

| Timeframe | v3 WS (`kline` interval) | Legacy v2 WS channel |
|-----------|--------------------------|----------------------|
| 1m–30m, 1H, 4H | native value (UTC-aligned) | `candle<tf>` (UTC-aligned) |
| 2h | — (rejected) | **`candle2H`** |
| 6h UTC | — (`6H` is UTC+8) | **`candle6Hutc`** |
| 8h | — | **`candle8H`** |
| 12h UTC | — (`12H` is UTC+8) | **`candle12Hutc`** |
| 1D UTC | — (`1D` opens 16:00 UTC) | **`candle1Dutc`** |
| 3D / 1W / 1M UTC | — | **`candle3Dutc` / `candle1Wutc` / `candle1Mutc`** |

## Private WebSocket (v3)

Endpoint: `wss://ws.bitget.com/v3/ws/private`. Verified live.

### Login

```json
{"op":"login","args":[{"apiKey":"...","passphrase":"...","timestamp":"1786962582","sign":"..."}]}
```

- `sign` = `base64(HMAC_SHA256(secret, timestamp + "GET" + "/user/verify"))`.
- `timestamp` is documented in unix **seconds**; the server actually accepts a millisecond
  timestamp too (both verified live). `bitgetv3` sends seconds per the docs.
- Success ack: `{"event":"login","code":0,"connId":"..."}` — `code` is a JSON **number**
  (the docs show strings), same as all other v3 WS events.
- A **rejected login arrives as an `event:"error"`, not as a login ack with a non-zero code**:
  `{"event":"error","code":30015,"msg":"Invalid sign"}`. The connection is NOT closed by the
  server (no close frame within 8s), unlike OKX which closes with 4001 — so a client must stop
  its own reconnect loop or it will retry the same bad login forever (`bitgetv3` cancels the
  reconnect job and fires `WithOnLoginFailed`).
- A subscription sent before login fails with
  `{"event":"error","code":30004,"msg":"User not logged in/User must be logged in"}`.

### Position channel (`topic: position`)

Subscription args: `{"instType":"UTA","topic":"position"}`.

- `instType` must be the literal uppercase `UTA`; lowercase `uta` is rejected with 30001
  "doesn't exist" (note: public v3 channels use lowercase instTypes — the opposite).
- The `symbol` field is optional: omitting it, sending `"symbol":""` and sending a real symbol
  are all accepted. However the **push envelope `arg` never echoes the symbol** — even when the
  subscription was made with one (the subscribe ack does echo it). A client that routes pushes
  by the full arg must therefore register the subscription without a symbol. `bitgetv3`
  subscribes account-wide (symbolless) only.
- A **snapshot is pushed immediately after subscribing even when the account is flat**:
  `{"action":"snapshot","arg":{"instType":"UTA","topic":"position"},"data":[],"ts":...}`.
  Each (re)subscription — including the automatic resubscribe after a reconnect — yields a
  fresh snapshot.
- There are no periodic pushes: after the snapshot, data arrives only on position events
  (see the docs event list). Keepalive is the same text `ping`/`pong` as the public WS
  (connection verified to survive a 2.5-minute idle window with pings).
- **Position events arrive with `action: "update"`** — only the initial subscription push
  is a `snapshot`. The docs example shows `snapshot` only, and the order channel (below)
  uses `snapshot` even for events — the two channels are inconsistent (both verified live).
- **Field set verified live against a real position push** (open + close of a minimal
  0.0001 BTCUSDT long): the push carries **no undocumented extras** (25 keys),
  and of the 26 documented fields **`marginRate` never arrives** — `WsPosition.MarginRate`
  always parses as 0.
- Close sequence: a reduce-only close produces two updates — first still `opening` with the
  size moved to `frozen`, then `ended` with zeroed sizes and most numerics collapsed to
  empty strings (`liqPrice`, `mmr`, `profitRate`, `breakEvenPrice`, `openFeeTotal`,
  `closeFeeTotal`, `totalFundingFee`, `cashDividend`).
- Numeric fields arrive as JSON strings and may be **empty strings** — parse them leniently.
  `openFeeTotal` arrives **negative** on an open position.

### Order channel (`topic: order`)

Subscription args: `{"instType":"UTA","topic":"order"}` (same uppercase-UTA, symbolless form as
the position channel; the push envelope `arg` carries no symbol either). Verified live.

- **No push on first-time subscription** (per the docs, confirmed live) — unlike the position
  channel, which sends a snapshot (possibly empty) after every (re)subscription.
- The push `action` is **always `snapshot`**, even for order events (place/fill/cancel).
- `category` in push data is **lowercase** (`usdt-futures`), unlike the uppercase REST
  `category` values — same asymmetry as public-WS instTypes.
- Push sequence is race-dependent: a fast-filling market order may arrive as a single
  `filled` push, or as a `new` push (cumExecQty 0) followed by a `filled` one — both
  observed live for identical orders. A zero-fill ioc limit order arrives as a single
  `cancelled` push with `cancelReason: "ioc_not_full_cancel"` and an empty `feeDetail`.
- `tradeSide` is documented as `open`/`close` but the live push sends detailed variants
  (one-way mode: `buy_single`/`sell_single`), like the REST order endpoints.
- `timeInForce` is engine-normalized: market orders push `"ioc"`.
- The live push carries fields **absent from the docs table**: `tpTriggerBy`, `slTriggerBy`,
  `takeprofit`, `stoploss` (note: these two keys are all-lowercase, unlike the REST
  `takeProfit`/`stopLoss`), `tpOrderType`, `slOrderType`, `tpLimitPrice`, `slLimitPrice`
  and `matchType` (numeric-as-string, meaning unknown, observed `"0"`).
- Numeric fields may be empty strings (`price` of a market order, the tp/sl block) —
  parse them leniently.

### Account channel (`topic: account`)

Subscription args: `{"instType":"UTA","topic":"account"}` (same uppercase-UTA, symbolless form
as the other private channels). Verified live with a real 0.0001 BTCUSDT open/close cycle.

- A **snapshot is pushed on first-time subscription** (per the docs, confirmed live), and after
  every (re)subscription including the automatic resubscribe after a reconnect — same behavior
  as the position channel, unlike the order channel.
- **Balance events arrive with `action: "update"`** (order fills, settlement, transfers) — like
  the position channel, unlike the order channel which uses `snapshot` for events.
- **Field set matches the docs exactly**: 8 top-level keys, 9 coin-item keys,
  **no undocumented extras and nothing missing** — the only private channel with zero
  discrepancies. Note the docs key `unrealisedPnL` (capital L).
- The push `data` list carries a **single item** (the whole account).
- Values are **absolute balances, not deltas**: an `update` repeats the full current state.
  The coin list carries only non-zero balances (same as REST `GET /api/v3/account/assets`).
  Known limitation of the verification: whether an `update` restricts the coin list to changed
  coins is not established (indistinguishable on a single-coin account) — do not assume the
  list is complete beyond the non-zero rule.
- One update per fill event was observed: the open fill pushes non-zero `imr`/`mmr`/`mgnRatio`
  and the fee-reduced balance; the close fill pushes them back to zero.
- The WS item differs from REST `AccountAssets`: `totalEquity` here vs `accountEquity` there;
  no `usdtEquity`/`btcEquity`/`positionValue`/`leverage` here; the coin item has `borrow`
  (absent from REST) and plural `debts` (REST has `debt`).

---

## FD Broker attribution (`X-CHANNEL-API-CODE`)

Source: "Bitget FD Broker Dashboard Guidance" (not part of the public API docs). Orders placed
via API are attributed to our broker account when the request carries the Channel API Code
(found in the Broker Dashboard, Personal Center → Basic Info).

**REST**: add the code as an HTTP header on order-placing requests:

```http
X-CHANNEL-API-CODE: <channel-api-code>
```

Supported endpoints (per the guidance; v3 subset relevant to this library):

- `/api/v3/trade/place-order`
- `/api/v3/trade/place-batch`
- `/api/v3/trade/modify-order`
- (plus the classic v2 mix/spot place/modify endpoints)

The header is not part of the signature (the pre-sign string covers only
`timestamp + method + path + body`), so it can be added independently of signing.

Implementation: `Client.channelApiCode` (default `ChannelApiCode` from `url.go`, override via
`WithChannelApiCode`) is set on every signed POST in `request.go`. An empty code disables the
header.

**WebSocket**: the guidance defines an `apiCode` field at the same level as `op`/`args` in the
place-order channel. Not implemented — this library has no WS order placement.
