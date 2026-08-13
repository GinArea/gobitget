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

The WS Candlesticks Channel (`topic: kline`) documents only the 10 base intervals.
The extended/utc values above are verified for **REST only**.

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

## WebSocket Candlesticks Channel quirks

- The **initial push** after subscribing is a snapshot carrying recent candle history
  (500 items, oldest-first): the current candle is the **last** item of `data`.
  The docs show a single-item example only.
- While trades occur, updates are pushed once per second; otherwise once per interval.
- Error events arrive with a **numeric** `code` field (the docs show a string),
  e.g. `30001` = topic does not exist.
