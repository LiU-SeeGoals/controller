# Actions Package Documentation 🐦

This package defines several actions that can be translated into byte sequences for communication purposes. Each action type implements the `TranslateReal()` method, which returns the byte representation of the action.

## GoTo Action

The `GoTo` action moves a robot to a specified position and goal.

### Byte Breakdown for `GoTo`

- **Byte 0**: Message length (19 bytes)
- **Byte 1**: Action type (2)
- **Byte 2**: Robot ID (1-6)
- **Bytes 3-12**: Position vector
  - **Bytes 3-4**: First element of `Pos` as `int16`
  - **Bytes 5-6**: Second element of `Pos` as `int16`
  - **Bytes 7-10**: Third element of `Pos` as `float32` (big-endian)
- **Bytes 11-19**: Goal vector
  - **Bytes 11-12**: First element of `Goal` as `int16`
  - **Bytes 13-14**: Second element of `Goal` as `int16`
  - **Bytes 15-18**: Third element of `Goal` as `float32` (big-endian)

---

## Stop Action

The `Stop` action commands a robot to stop.

### Byte Breakdown for `Stop`

- **Byte 0**: Message length (3 bytes)
- **Byte 1**: Action type (This action has the id 0)
- **Byte 2**: Robot ID

## Kick Action

The `Kick` action commands a robot to kick.

### Byte Breakdown for `Kick`

- **Byte 0**: Message length (4 bytes)
- **Byte 1**: Action type (1)
- **Byte 2**: Robot ID
- **Byte 3**: Kicking speed (0-255)

---

## Init Action

The `Init` action initializes a robot.

### Byte Breakdown for `Init`

- **Byte 0**: Message length (3 bytes)
- **Byte 1**: Action type (3)
- **Byte 2**: Robot ID (1-6)

