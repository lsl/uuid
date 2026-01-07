# uuid

Fast UUID v7 generation.

Based on the conclusions from [benchmarking UUID v7](https://lgsn.dev/uuidbench).

Combines buffered `rand.Read` calls from [coolaj86/uuidv7](https://github.com/coolaj86/uuidv7) and sequence counting from [gofrs/uuid](https://github.com/gofrs/uuid) to produce a faster UUID v7 generation.

## Usage

```go
import "github.com/lsl/uuid"
u := uuid.NewV7() // uuid.UUIDv7
```

