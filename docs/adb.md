# ADB Swipe Examples from Presets (screen 1080x2400, 300px)

| Content Direction | Command                                                                 |
|-------------------|-------------------------------------------------------------------------|
| **Right**         | `adb -s RF8RC00M8MF shell input touchscreen swipe 540 1200 240 1200 400` |
| **Left**          | `adb -s RF8RC00M8MF shell input touchscreen swipe 540 1200 840 1200 400` |
| **Up**            | `adb -s RF8RC00M8MF shell input touchscreen swipe 540 1200 540 1500 400` |
| **Down**          | `adb -s RF8RC00M8MF shell input touchscreen swipe 540 1200 540 900 400`  |

**Format:**  
`adb -s <device_id> shell input touchscreen swipe <x1> <y1> <x2> <y2> <duration_ms>`

- `<x1> <y1>` — starting point
- `<x2> <y2>` — ending point
- `<duration_ms>` — swipe duration

> For other directions/offsets — change coordinates accordingly.