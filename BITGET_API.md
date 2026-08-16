# Bitget UTA v3 API Reference — Exchange-Specific Quirks

This document records live-verified behavior of the Bitget UTA v3 API that is missing from or
contradicts the official docs (https://www.bitget.com/api-doc/uta/). Verified 2026-08-13 against
`https://api.bitget.com` on SPOT, USDT-FUTURES and COIN-FUTURES.

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

### WebSocket kline channel

The v3 WS Candlesticks Channel (`topic: kline`) accepts **only** the 10 documented base intervals —
the extended/utc values above are REST-only. For other timeframes over WS use the legacy v2 WS
(see "WebSocket kline intervals: v3 vs legacy v2" below).

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

Verified live 2026-08-16 on SPOT and USDT-FUTURES (BTCUSDT), every result reproduced on a second run.

### v3 WS (`wss://ws.bitget.com/v3/ws/public`, topic `kline`)

Accepts **only** the 10 documented intervals: `1m, 3m, 5m, 15m, 30m, 1H, 4H, 6H, 12H, 1D`.
Every other value — `2H`, the utc-suffixed variants (`6Hutc`/`12Hutc`/`1Dutc`/...), `3D`, `1W`,
`1M`, `3H`, `8H` and any case variant (`2h`, `1DUTC`, `1dutc`, `12hutc`) — is rejected with
error `30001` `"{...} doesn't exist"`. Unlike REST, the v3 WS does **not** inherit the v2 values.
Native `6H`/`12H`/`1D` open on the UTC+8 grid, same as REST (confirmed on WS).

Consequence: over the v3 WS, UTC-aligned candles exist only for intervals up to `4H`.

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

## WebSocket Candlesticks Channel quirks

- The **initial push** after subscribing is a snapshot carrying recent candle history
  (500 items, oldest-first): the current candle is the **last** item of `data`.
  The docs show a single-item example only.
- While trades occur, updates are pushed once per second; otherwise once per interval.
- Error events arrive with a **numeric** `code` field (the docs show a string),
  e.g. `30001` = topic does not exist.
