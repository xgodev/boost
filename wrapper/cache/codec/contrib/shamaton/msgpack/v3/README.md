# shamaton/msgpack v3 codec

`wrapper/cache` codec that encodes values as MessagePack with
[`github.com/shamaton/msgpack/v3`](https://github.com/shamaton/msgpack).
The v2 codec (`../v2`) stays; both implement `cache.Codec[T]`.

```go
import (
	"github.com/xgodev/boost/wrapper/cache"
	msgpack "github.com/xgodev/boost/wrapper/cache/codec/contrib/shamaton/msgpack/v3"
)

manager := cache.NewManager[Order]("orders", msgpack.New[Order](), drv)
```

## Moving from v2

- **Compatible on the wire.** v3 reads the entries v2 wrote and v2 reads the
  entries v3 writes (strings, bytes, structs, slices, times), so a service can
  switch codecs without flushing its cache, and instances on v2 and v3 can
  share one cache during a rollout (`codec_test.go`: `TestCodec_ReadsV2`,
  `TestCodec_V2Reads`, `TestCodec_Time_DecodesAsUTC`).
- **`time.Time` decodes in UTC.** v2 decoded a timestamp into `time.Local`;
  v3 decodes it into `time.UTC`. The instant is the same (`Equal` holds), only
  `Location()` differs, so code that formats the time or compares with `==`
  sees a change. `msgpack.SetDecodedTimeAsLocal()` from
  `github.com/shamaton/msgpack/v3` restores the v2 behaviour process-wide.
