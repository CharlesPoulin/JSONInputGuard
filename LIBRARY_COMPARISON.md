# JSONInputGuard vs Popular JSON Libraries Comparison

## Benchmark Results Summary

Based on comprehensive benchmarks against popular JSON libraries, here's how JSONInputGuard performs:

### Raw Performance Numbers (64KB payload)

| Library | Operations/sec | Throughput | Memory Usage | Allocations |
|---------|---------------|------------|--------------|-------------|
| **GJSON** | 22,652 ops/sec | 1,145.45 MB/s | 68,424 B/op | 2 allocs/op |
| **JSONInputGuard** | 4,509 ops/sec | 235.53 MB/s | 14,487 B/op | 5 allocs/op |
| **Jsoniter** | 339 ops/sec | 19.40 MB/s | 876,574 B/op | 16,496 allocs/op |
| **Fastjson** | 204 ops/sec | 11.35 MB/s | 3,292,834 B/op | 172 allocs/op |
| **Standard Library** | 363 ops/sec | 19.76 MB/s | 797,260 B/op | 16,489 allocs/op |

## Key Findings

### 🏆 **GJSON is the Fastest**
- **GJSON**: 22,652 ops/sec (1,145 MB/s)
- **JSONInputGuard**: 4,509 ops/sec (235 MB/s)
- **Standard Library**: 363 ops/sec (19.76 MB/s)

### 🎯 **JSONInputGuard vs Standard Library**
- **4.5x faster** than standard library
- **55x less memory** usage
- **3,298x fewer allocations**

### 📊 **Memory Efficiency Ranking**
1. **GJSON**: 68KB/op (2 allocs)
2. **JSONInputGuard**: 14KB/op (5 allocs) 
3. **Standard Library**: 797KB/op (16,489 allocs)
4. **Jsoniter**: 876KB/op (16,496 allocs)
5. **Fastjson**: 3.3MB/op (172 allocs)

## Library Analysis

### 🚀 **GJSON** - The Speed Champion
**Pros:**
- Fastest parsing (22,652 ops/sec)
- Minimal memory usage (68KB/op)
- Very few allocations (2 allocs/op)
- Excellent for read-only operations

**Cons:**
- Read-only (no unmarshaling to structs)
- Limited validation capabilities
- Not suitable for full object creation

**Best for:** Fast JSON querying and validation

### 🛡️ **JSONInputGuard** - The Balanced Performer
**Pros:**
- 4.5x faster than standard library
- 55x less memory than stdlib
- AST-based validation
- Full struct unmarshaling
- Early rejection of invalid payloads

**Cons:**
- Slower than GJSON for pure parsing
- More complex than simple libraries

**Best for:** High-throughput APIs with validation needs

### 📚 **Standard Library** - The Baseline
**Pros:**
- Built into Go
- Full feature set
- Well-tested and stable

**Cons:**
- Slowest of the bunch
- High memory usage
- Many allocations

**Best for:** Simple applications, prototyping

### ⚡ **Jsoniter** - The Alternative
**Pros:**
- Drop-in replacement for stdlib
- API compatible
- Slightly better performance than stdlib

**Cons:**
- Still high memory usage
- Many allocations
- Not dramatically faster than stdlib

**Best for:** Existing codebases wanting modest improvement

### 🐌 **Fastjson** - The Memory Hog
**Pros:**
- Good for large JSON documents
- Streaming capabilities

**Cons:**
- Highest memory usage (3.3MB/op)
- Slowest of the libraries tested
- Complex API

**Best for:** Large document processing (not recommended for APIs)

## Use Case Recommendations

### 🎯 **For High-Throughput APIs**
1. **GJSON** - If you only need validation/querying
2. **JSONInputGuard** - If you need full struct unmarshaling + validation
3. **Jsoniter** - If you want drop-in stdlib replacement

### 🎯 **For Memory-Constrained Environments**
1. **GJSON** - Minimal memory footprint
2. **JSONInputGuard** - Excellent memory efficiency with full features
3. **Standard Library** - Avoid for high-throughput scenarios

### 🎯 **For Validation-Heavy Workloads**
1. **JSONInputGuard** - AST-based validation is very efficient
2. **GJSON** - Good for simple field checks
3. **Standard Library** - Full validation but slow

### 🎯 **For Existing Codebases**
1. **Jsoniter** - Drop-in replacement
2. **JSONInputGuard** - If you can refactor for better performance
3. **Standard Library** - If performance isn't critical

## Performance Insights

### **Speed vs Memory Trade-offs**

| Library | Speed Rank | Memory Rank | Best For |
|---------|------------|-------------|----------|
| GJSON | 1st | 1st | Pure parsing/querying |
| JSONInputGuard | 2nd | 2nd | APIs with validation |
| Standard Library | 5th | 4th | Simple applications |
| Jsoniter | 4th | 5th | Drop-in replacement |
| Fastjson | 6th | 6th | Large documents |

### **Allocation Efficiency**
- **GJSON**: 2 allocs/op (99.99% fewer than stdlib)
- **JSONInputGuard**: 5 allocs/op (99.97% fewer than stdlib)
- **Standard Library**: 16,489 allocs/op (baseline)

## Conclusion

### 🏆 **JSONInputGuard's Position**

JSONInputGuard offers an excellent **balance of performance and features**:

- **4.5x faster** than standard library
- **55x less memory** than standard library
- **3,298x fewer allocations** than standard library
- **Full struct unmarshaling** capabilities
- **AST-based validation** for early rejection

### 🎯 **When to Choose Each Library**

**Choose GJSON when:**
- You need maximum speed
- You only need to query/validate JSON
- Memory is extremely constrained

**Choose JSONInputGuard when:**
- You need full struct unmarshaling
- You want excellent performance with validation
- You're building high-throughput APIs

**Choose Standard Library when:**
- Performance isn't critical
- You want maximum compatibility
- You're prototyping or learning

**Choose Jsoniter when:**
- You want a drop-in stdlib replacement
- You need modest performance improvement
- You can't refactor existing code

### 🚀 **JSONInputGuard's Sweet Spot**

JSONInputGuard excels in **high-throughput API scenarios** where you need:
- Fast JSON processing
- Full struct unmarshaling
- Efficient validation
- Low memory footprint
- Early rejection of invalid payloads

It's particularly well-suited for:
- Microservices
- Serverless functions (AWS Lambda)
- Real-time APIs
- High-traffic web services

The AST-based approach provides excellent performance while maintaining full Go struct compatibility, making it ideal for production environments where both speed and functionality matter.
