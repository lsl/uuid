# uuid

Fast uuidv7 generation.

Based on the conclusions of [benchmarking uuidv7 packages](https://lgsn.dev/uuidbench). Combines UUID generation approaches of [coolaj86/uuidv7](https://github.com/coolaj86/uuidv7) and [gofrs/uuid](https://github.com/gofrs/uuid) to buffer `rand.Read` calls and sequence generations within the same millisecond.

## Usage

```
  import github.com/lsl/uuid
	u := uuid.NewV7() // uuid.UUIDv7
```

