🧊 Snowflake UUID Generator (Go)

A high-performance, thread-safe **Snowflake ID generator** written in Go.  
It generates unique 64-bit IDs inspired by Twitter’s Snowflake algorithm — ideal for distributed systems, databases, and high-concurrency environments.

---

⚙️ Features

- 64-bit **unique IDs** combining timestamp, datacenter ID, worker ID, and sequence.
- **Configurable epoch** for compact IDs and version control.
- **Shard-based concurrency** (lock-minimized with round-robin scheduling).
- **Thread-safe** — safe for concurrent use.
- Production-grade error handling for clock rollback and sequence overflow.
- Modular Go package ready for reuse and scaling.

---

🧮 ID Bit Layout

| Bits | Field          | Description                                |
|------|----------------|--------------------------------------------|
| 41   | Timestamp      | Milliseconds since custom epoch            |
| 5    | Datacenter ID  | Supports up to 32 datacenters              |
| 5    | Worker ID      | Supports up to 32 workers per datacenter   |
| 12   | Sequence       | 4096 unique IDs per millisecond per worker |

Total: 63 bits(1 bit reserved for sign)

---

🧠 Formula

```
ID = (timestamp << timeShift) | (datacenterId << dcShift) | (workerId << workerShift) | sequence
```

Each part ensures uniqueness:
- **Timestamp** prevents collisions across time.
- **Datacenter & Worker** isolate IDs per node.
- **Sequence** differentiates multiple requests in the same millisecond.

---

🗓️ Default Epoch

```
January 1, 2025 (UTC)
1735689600000 ms since Unix epoch
```

You can customize this in the `GeneratorConfig`.

---

📦 Project Structure

```
uuid-generator/
├── cmd/
│   └── server/
│       └── main.go          # Example usage / entry point
├── pkg/
│   └── snowflake/
│       ├── generator.go     # Core Snowflake generator logic
│       ├── encode.go        # (optional) for future encoding utilities
│       └── decode.go        # (optional) for ID breakdown
├── generator-test.go         # Test file with concurrent ID generation
└── go.mod                    # Module definition
```

---

🚀 Usage

Import the package

```go
import "github.com/gokul1063/uuid-generator/pkg/snowflake"
```

Initialize the generator

```go
gen, err := snowflake.NewGenerator(snowflake.GeneratorConfig{
    DatacenterId: 1,
    WorkerId:     1,
    Shards:       1, // Set to 1 for single-node setups
})
if err != nil {
    log.Fatalf("failed to initialize generator: %v", err)
}
```

Generate IDs

```go
id, err := gen.Next()
if err != nil {
    log.Println("error generating ID:", err)
}
fmt.Println("Generated ID:", id)
```

---

🧪 Run Test

To test ID generation and log output:

```bash
go run generator-test.go
```

Logs are automatically stored in:

```
./logs/
```

---

🧰 Configuration Options

| Field          | Type   | Description                                  |
|----------------|--------|----------------------------------------------|
| `Epoch`        | int64  | Custom epoch in milliseconds (optional)      |
| `DatacenterId` | uint64 | Datacenter identifier (0–31)                 |
| `WorkerId`     | uint64 | Worker identifier (0–31)                     |
| `Shards`       | int    | Internal concurrency shards (recommended: 1) |

---


⚠️ Notes (For future use)

- **Single shard** (Shards = 1) ensures guaranteed uniqueness on a single machine.
- If you use **multiple shards or nodes**, assign **unique WorkerIds** or DatacenterIds per node.
- The generator waits up to 5ms if the system clock moves backwards.

---

## 🏗️ Future Enhancements

- Add REST API endpoint for UUID generation.
- Implement persistence-based ID tracking for debugging.
- Optional Base58 or ULID encoding for shorter string forms.

---

## 🧑‍💻 Author

**Gokul R**  
Backend team of Aspirenet 
Building scalable distributed tools in Go 🧠⚡



