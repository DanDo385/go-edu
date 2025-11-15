# Project 36: Caching Reverse Proxy

## 1. What Is This About?

### Real-World Scenario

You're running a web application with a database-heavy backend. Each request hits the database, which is slow and expensive:

**Problem**:
- Your API endpoint `/api/users/123` takes 200ms (database query)
- You serve 1000 requests/sec
- Total database load: 1000 queries/sec
- Database struggles under load, response times increase to 2 seconds

**❌ Naive approach:** Every request hits the backend
```
User 1 → Backend → Database (200ms)
User 2 → Backend → Database (200ms)
User 3 → Backend → Database (200ms)
...1000 requests → 1000 database queries
```

**Result**: Database is overwhelmed, requests are slow, users are frustrated.

**✅ Solution:** Caching reverse proxy
```
User 1 → Proxy → Backend → Database (200ms) [CACHE MISS, cache response]
User 2 → Proxy → [CACHE HIT! Return cached response] (1ms)
User 3 → Proxy → [CACHE HIT! Return cached response] (1ms)
...997 more requests → All cache hits!

Total database queries: 1 instead of 1000
Response time: 1ms instead of 200ms
```

**Result**: 99.9% reduction in database load, 200x faster responses!

This project teaches you how to build a **production-grade caching reverse proxy** that:
- **Intercepts requests**: Acts as a middleman between clients and backend
- **Caches responses**: Stores frequently-accessed data in memory
- **Reduces load**: Drastically reduces backend and database load
- **Speeds up responses**: Serves cached data in microseconds
- **Handles invalidation**: Smart cache eviction strategies

### What You'll Learn

1. **Reverse proxies**: What they are and how they work
2. **HTTP caching**: Cache-Control headers, ETags, expiration
3. **httputil.ReverseProxy**: Go's built-in reverse proxy
4. **Cache invalidation**: LRU, TTL, manual invalidation
5. **HTTP middleware**: Layering caching on top of proxying
6. **Concurrency-safe caching**: Thread-safe cache implementations

### The Challenge

Build a caching reverse proxy that:
- Forwards requests to a backend server
- Caches GET responses in memory
- Respects Cache-Control headers (max-age, no-cache, no-store)
- Supports TTL (time-to-live) expiration
- Implements LRU eviction when cache is full
- Provides cache statistics (hits, misses, hit rate)
- Thread-safe for concurrent requests

---

## 2. First Principles: Reverse Proxies and Caching

### What is a Reverse Proxy?

A **reverse proxy** sits between clients and backend servers, forwarding requests and responses.

**Client → Reverse Proxy → Backend Server**

**Forward proxy** (like VPN): Client uses proxy to access internet
```
Client → Forward Proxy → Internet
(Hides client's IP)
```

**Reverse proxy**: Server uses proxy to handle client requests
```
Client → Reverse Proxy → Backend
(Hides backend servers)
```

### Why Use Reverse Proxies?

**1. Load balancing**
```
                    → Backend 1
Client → Proxy  → → Backend 2
                    → Backend 3
```
Distribute load across multiple servers.

**2. Caching**
```
Client → Proxy [Cache] → Backend
```
Cache responses to reduce backend load.

**3. SSL/TLS termination**
```
Client --HTTPS-→ Proxy --HTTP-→ Backend
```
Proxy handles encryption, backend handles logic.

**4. Security**
```
Client → Proxy [Firewall/WAF] → Backend
```
Hide backend servers, filter malicious requests.

**5. Compression**
```
Client ←-gzip-← Proxy ←-raw-← Backend
```
Compress responses to save bandwidth.

### How HTTP Caching Works

**HTTP Cache-Control headers** tell proxies/browsers how to cache responses.

**Example response**:
```http
HTTP/1.1 200 OK
Cache-Control: max-age=3600, public
Content-Type: application/json

{"id": 123, "name": "Alice"}
```

**Cache-Control directives**:

| Directive | Meaning |
|-----------|---------|
| `max-age=3600` | Cache for 3600 seconds (1 hour) |
| `public` | Can be cached by anyone (proxies, browsers) |
| `private` | Only browser can cache (not proxies) |
| `no-cache` | Must revalidate with server before using cache |
| `no-store` | Never cache (sensitive data) |
| `must-revalidate` | Must check if stale before using |

**Caching flow**:
```
1. Request arrives: GET /api/users/123
2. Check cache: Is this URL cached and fresh?
3a. CACHE HIT: Return cached response immediately
3b. CACHE MISS: Forward to backend
4. Backend responds
5. Store in cache (if cacheable)
6. Return response to client
```

### What is Cache Invalidation?

**Cache invalidation**: Removing stale/outdated data from cache.

> "There are only two hard things in Computer Science: cache invalidation and naming things." — Phil Karlton

**Strategies**:

**1. Time-based (TTL - Time To Live)**
```
Cache entry: {key: "/api/users/123", value: {...}, expires: 1700000000}
Current time: 1700000001
→ Entry expired, remove from cache
```

**2. Size-based (LRU - Least Recently Used)**
```
Cache is full (max 100 entries)
New entry arrives
→ Evict least recently used entry
```

**3. Manual invalidation**
```
POST /api/users/123 (update user)
→ Invalidate cache for /api/users/123
```

**4. Event-based**
```
Database update event
→ Invalidate all related cache entries
```

### What is LRU Cache?

**LRU (Least Recently Used)** evicts items that haven't been accessed recently.

**Data structure**: Hash map + doubly linked list

**Operations**:
- `Get(key)`: O(1) — Move item to front
- `Put(key, value)`: O(1) — Add to front, evict from back if full

**Example**:
```
Cache capacity: 3

PUT(A, 1)  →  [A]
PUT(B, 2)  →  [B, A]
PUT(C, 3)  →  [C, B, A]
GET(A)     →  [A, C, B]  (A moved to front)
PUT(D, 4)  →  [D, A, C]  (B evicted, was least recently used)
```

**Why LRU?**
- **Hot data stays cached**: Frequently accessed items remain
- **Cold data evicted**: Rarely used items are removed
- **O(1) operations**: Efficient for high-traffic systems

---

## 3. Breaking Down the Solution

### Step 1: Understanding httputil.ReverseProxy

Go's `net/http/httputil` package provides a **ReverseProxy** type.

**Basic usage**:
```go
proxy := httputil.NewSingleHostReverseProxy(backendURL)
http.Handle("/", proxy)
http.ListenAndServe(":8080", nil)
```

**That's it!** The proxy forwards all requests to `backendURL`.

**How it works**:
1. Client sends request to proxy (`:8080`)
2. Proxy modifies request:
   - Changes `Host` header to backend host
   - Forwards headers, body, method
3. Proxy sends request to backend
4. Backend responds
5. Proxy forwards response to client

**Customizing the proxy**:
```go
proxy := &httputil.ReverseProxy{
    Director: func(req *http.Request) {
        req.URL.Scheme = "http"
        req.URL.Host = "backend:8000"
        req.Header.Set("X-Forwarded-For", req.RemoteAddr)
    },
    ModifyResponse: func(resp *http.Response) error {
        // Modify response before sending to client
        resp.Header.Set("X-Proxy", "my-proxy")
        return nil
    },
}
```

**Director**: Modifies the outgoing request
**ModifyResponse**: Modifies the response before returning

### Step 2: Cache Entry Structure

```go
type CacheEntry struct {
    Response   *http.Response  // Cached response
    Body       []byte          // Response body (must cache separately)
    Expiry     time.Time       // When this entry expires
    AccessTime time.Time       // Last access time (for LRU)
}
```

**Why cache body separately?**

`http.Response.Body` is an `io.ReadCloser` — once read, it's consumed.
```go
// This DOESN'T work:
resp, _ := http.Get(url)
cache[url] = resp

client1.Write(resp.Body)  // Works
client2.Write(resp.Body)  // FAILS! Body already consumed
```

**Solution**: Read body into `[]byte`, create new readers for each request.
```go
body, _ := io.ReadAll(resp.Body)
entry := CacheEntry{Response: resp, Body: body}

// Later, create new reader for each client:
resp.Body = io.NopCloser(bytes.NewReader(entry.Body))
```

### Step 3: Cache Structure

```go
type Cache struct {
    mu       sync.RWMutex          // Protect concurrent access
    entries  map[string]*CacheEntry // Cache storage
    lru      *list.List             // LRU linked list
    maxSize  int                    // Max entries
    ttl      time.Duration          // Default TTL
    hits     int64                  // Cache hit counter
    misses   int64                  // Cache miss counter
}
```

**Fields explained**:
- `mu`: Read/write mutex for thread safety
- `entries`: Hash map for O(1) lookups
- `lru`: Doubly linked list for LRU ordering
- `maxSize`: Maximum cache entries (e.g., 1000)
- `ttl`: Default time-to-live (e.g., 5 minutes)
- `hits/misses`: Statistics

**Thread safety**:
```go
func (c *Cache) Get(key string) (*CacheEntry, bool) {
    c.mu.RLock()         // Multiple readers allowed
    defer c.mu.RUnlock()

    entry, exists := c.entries[key]
    return entry, exists
}

func (c *Cache) Set(key string, entry *CacheEntry) {
    c.mu.Lock()          // Exclusive write access
    defer c.mu.Unlock()

    c.entries[key] = entry
}
```

### Step 4: Caching Middleware

```go
func (c *Cache) CachingProxy(backend *httputil.ReverseProxy) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Only cache GET requests
        if r.Method != http.MethodGet {
            backend.ServeHTTP(w, r)
            return
        }

        // Generate cache key
        key := r.URL.String()

        // Try cache
        if entry, exists := c.Get(key); exists && !c.isExpired(entry) {
            c.serveFromCache(w, entry)
            atomic.AddInt64(&c.hits, 1)
            return
        }

        // Cache miss, forward to backend
        atomic.AddInt64(&c.misses, 1)
        c.proxyAndCache(w, r, backend, key)
    })
}
```

**Flow**:
1. Check if method is GET (only cache idempotent requests)
2. Generate cache key (URL)
3. Check cache
4. If hit and fresh → serve from cache
5. If miss → forward to backend, cache response

### Step 5: Response Capture

To cache a response, we need to **capture** it before sending to client.

**Problem**: `backend.ServeHTTP(w, r)` writes directly to `w` — we can't intercept.

**Solution**: `ResponseWriter` wrapper
```go
type ResponseRecorder struct {
    http.ResponseWriter
    Status int
    Body   *bytes.Buffer
    Header http.Header
}

func (rr *ResponseRecorder) WriteHeader(status int) {
    rr.Status = status
}

func (rr *ResponseRecorder) Write(b []byte) (int, error) {
    return rr.Body.Write(b)  // Capture body
}
```

**Usage**:
```go
recorder := &ResponseRecorder{
    ResponseWriter: w,
    Body:          new(bytes.Buffer),
    Header:        make(http.Header),
}

backend.ServeHTTP(recorder, r)  // Backend writes to recorder

// Now we have the response!
entry := &CacheEntry{
    Body:   recorder.Body.Bytes(),
    Expiry: time.Now().Add(c.ttl),
}
c.Set(key, entry)

// Forward to client
w.WriteHeader(recorder.Status)
w.Write(recorder.Body.Bytes())
```

### Step 6: LRU Eviction

When cache is full, evict least recently used entry.

**Using container/list**:
```go
import "container/list"

type Cache struct {
    lru     *list.List
    entries map[string]*list.Element
}

type lruEntry struct {
    key   string
    value *CacheEntry
}

func (c *Cache) Get(key string) (*CacheEntry, bool) {
    if elem, exists := c.entries[key]; exists {
        c.lru.MoveToFront(elem)  // Mark as recently used
        return elem.Value.(*lruEntry).value, true
    }
    return nil, false
}

func (c *Cache) Set(key string, entry *CacheEntry) {
    if len(c.entries) >= c.maxSize {
        // Evict LRU
        oldest := c.lru.Back()
        c.lru.Remove(oldest)
        delete(c.entries, oldest.Value.(*lruEntry).key)
    }

    elem := c.lru.PushFront(&lruEntry{key, entry})
    c.entries[key] = elem
}
```

---

## 4. Complete Solution Walkthrough

### Cache Implementation

```go
type Cache struct {
    mu       sync.RWMutex
    entries  map[string]*CacheEntry
    lru      *list.List
    lruMap   map[string]*list.Element
    maxSize  int
    ttl      time.Duration
    hits     int64
    misses   int64
}

func NewCache(maxSize int, ttl time.Duration) *Cache {
    return &Cache{
        entries: make(map[string]*CacheEntry),
        lru:     list.New(),
        lruMap:  make(map[string]*list.Element),
        maxSize: maxSize,
        ttl:     ttl,
    }
}
```

### Get Operation

```go
func (c *Cache) Get(key string) (*CacheEntry, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // Check if exists
    elem, exists := c.lruMap[key]
    if !exists {
        return nil, false
    }

    // Check if expired
    entry := elem.Value.(*lruEntry).value
    if time.Now().After(entry.Expiry) {
        // Expired, remove from cache
        c.lru.Remove(elem)
        delete(c.lruMap, key)
        delete(c.entries, key)
        return nil, false
    }

    // Mark as recently used
    c.lru.MoveToFront(elem)
    entry.AccessTime = time.Now()

    return entry, true
}
```

**Key operations**:
1. Check existence in LRU map
2. Check expiration
3. Move to front (mark as recently used)
4. Update access time

### Set Operation

```go
func (c *Cache) Set(key string, entry *CacheEntry) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // Check if already exists
    if elem, exists := c.lruMap[key]; exists {
        // Update existing entry
        elem.Value.(*lruEntry).value = entry
        c.lru.MoveToFront(elem)
        return
    }

    // Check if cache is full
    if c.lru.Len() >= c.maxSize {
        // Evict LRU
        oldest := c.lru.Back()
        if oldest != nil {
            lruEntry := oldest.Value.(*lruEntry)
            c.lru.Remove(oldest)
            delete(c.lruMap, lruEntry.key)
            delete(c.entries, lruEntry.key)
        }
    }

    // Add new entry
    elem := c.lru.PushFront(&lruEntry{key: key, value: entry})
    c.lruMap[key] = elem
    c.entries[key] = entry
}
```

**Key operations**:
1. Update if exists
2. Evict LRU if full
3. Add to front (most recently used)

### Caching Proxy Handler

```go
func (c *Cache) Handler(backend *httputil.ReverseProxy) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Only cache GET requests
        if r.Method != http.MethodGet {
            backend.ServeHTTP(w, r)
            return
        }

        key := r.URL.String()

        // Try cache
        if entry, exists := c.Get(key); exists {
            c.serveFromCache(w, entry)
            atomic.AddInt64(&c.hits, 1)
            return
        }

        // Cache miss
        atomic.AddInt64(&c.misses, 1)

        // Capture response
        recorder := &ResponseRecorder{
            ResponseWriter: w,
            body:          new(bytes.Buffer),
            header:        make(http.Header),
        }

        backend.ServeHTTP(recorder, r)

        // Cache if cacheable
        if c.isCacheable(recorder) {
            entry := &CacheEntry{
                Body:       recorder.body.Bytes(),
                StatusCode: recorder.status,
                Header:     recorder.header,
                Expiry:     c.calculateExpiry(recorder.header),
                AccessTime: time.Now(),
            }
            c.Set(key, entry)
        }

        // Send to client
        c.serveFromRecorder(w, recorder)
    })
}
```

### Serving from Cache

```go
func (c *Cache) serveFromCache(w http.ResponseWriter, entry *CacheEntry) {
    // Copy headers
    for key, values := range entry.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }

    // Add cache status header
    w.Header().Set("X-Cache", "HIT")

    // Write status and body
    w.WriteHeader(entry.StatusCode)
    w.Write(entry.Body)
}
```

### Determining Cacheability

```go
func (c *Cache) isCacheable(recorder *ResponseRecorder) bool {
    // Only cache successful responses
    if recorder.status != http.StatusOK {
        return false
    }

    // Check Cache-Control
    cacheControl := recorder.header.Get("Cache-Control")

    // no-store means never cache
    if strings.Contains(cacheControl, "no-store") {
        return false
    }

    // private means only browser can cache, not proxy
    if strings.Contains(cacheControl, "private") {
        return false
    }

    return true
}
```

### Calculating Expiry

```go
func (c *Cache) calculateExpiry(header http.Header) time.Time {
    cacheControl := header.Get("Cache-Control")

    // Parse max-age
    if strings.Contains(cacheControl, "max-age=") {
        parts := strings.Split(cacheControl, "max-age=")
        if len(parts) > 1 {
            ageStr := strings.Split(parts[1], ",")[0]
            if age, err := strconv.Atoi(strings.TrimSpace(ageStr)); err == nil {
                return time.Now().Add(time.Duration(age) * time.Second)
            }
        }
    }

    // Use default TTL
    return time.Now().Add(c.ttl)
}
```

---

## 5. Key Concepts Explained

### Concept 1: Cache Stampede

**Problem**: Cache expires, multiple requests arrive simultaneously, all hit backend.

```
Cache expires at T=0
T=0: Request 1 → MISS → Backend
T=0: Request 2 → MISS → Backend
T=0: Request 3 → MISS → Backend
...100 requests → 100 backend hits!
```

**Solution**: Cache locking (request coalescing)
```go
type Cache struct {
    locks map[string]*sync.Mutex
}

func (c *Cache) Get(key string) (*CacheEntry, bool) {
    // Check cache
    if entry, exists := c.entries[key]; exists {
        return entry, true
    }

    // Acquire lock for this key
    lock := c.getLock(key)
    lock.Lock()
    defer lock.Unlock()

    // Check again (another goroutine might have cached it)
    if entry, exists := c.entries[key]; exists {
        return entry, true
    }

    // Only one goroutine reaches here
    // Fetch from backend and cache
}
```

### Concept 2: Cache Warming

**Cache warming**: Pre-populate cache with expected data.

```go
func (c *Cache) WarmUp(urls []string) {
    for _, url := range urls {
        resp, err := http.Get(url)
        if err != nil {
            continue
        }
        body, _ := io.ReadAll(resp.Body)
        c.Set(url, &CacheEntry{
            Body:   body,
            Expiry: time.Now().Add(c.ttl),
        })
        resp.Body.Close()
    }
}
```

**Use case**: Warm cache on server startup to avoid cold start.

### Concept 3: Multi-Tier Caching

**L1**: In-memory (fast, small)
**L2**: Redis (medium speed, large)
**L3**: Database (slow, unlimited)

```go
func (c *Cache) Get(key string) (*CacheEntry, bool) {
    // L1: In-memory
    if entry, exists := c.memCache.Get(key); exists {
        return entry, true
    }

    // L2: Redis
    if entry, exists := c.redisCache.Get(key); exists {
        c.memCache.Set(key, entry)  // Promote to L1
        return entry, true
    }

    // L3: Database
    entry := c.db.Query(key)
    c.redisCache.Set(key, entry)  // Store in L2
    c.memCache.Set(key, entry)    // Store in L1
    return entry, true
}
```

### Concept 4: Cache Tags

**Problem**: Invalidate related entries when data changes.

**Example**: User updates profile
```
Cache entries:
- /api/users/123
- /api/users/123/posts
- /api/feed (contains user 123's data)
```

All need invalidation when user 123 changes!

**Solution**: Tag-based invalidation
```go
type CacheEntry struct {
    Body   []byte
    Tags   []string  // e.g., ["user:123"]
}

func (c *Cache) InvalidateTag(tag string) {
    for key, entry := range c.entries {
        if contains(entry.Tags, tag) {
            delete(c.entries, key)
        }
    }
}

// Usage:
c.InvalidateTag("user:123")  // Invalidates all entries tagged with user:123
```

### Concept 5: Stale-While-Revalidate

**Pattern**: Serve stale cache while refreshing in background.

```go
func (c *Cache) Get(key string) (*CacheEntry, bool) {
    entry, exists := c.entries[key]
    if !exists {
        return nil, false
    }

    if time.Now().After(entry.Expiry) {
        // Expired, but serve stale
        go c.refresh(key)  // Refresh in background
        entry.Header.Set("X-Cache", "STALE")
        return entry, true
    }

    return entry, true
}
```

**Benefits**: Always fast response, eventual consistency.

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Cache Aside (Lazy Loading)

```go
func GetUser(id int) (*User, error) {
    key := fmt.Sprintf("user:%d", id)

    // Check cache
    if user, exists := cache.Get(key); exists {
        return user.(*User), nil
    }

    // Not in cache, fetch from DB
    user, err := db.QueryUser(id)
    if err != nil {
        return nil, err
    }

    // Store in cache
    cache.Set(key, user)
    return user, nil
}
```

### Pattern 2: Write-Through Cache

```go
func UpdateUser(user *User) error {
    // Write to DB
    if err := db.UpdateUser(user); err != nil {
        return err
    }

    // Write to cache
    key := fmt.Sprintf("user:%d", user.ID)
    cache.Set(key, user)

    return nil
}
```

### Pattern 3: Cache Invalidation on Write

```go
func UpdateUser(user *User) error {
    // Write to DB
    if err := db.UpdateUser(user); err != nil {
        return err
    }

    // Invalidate cache
    key := fmt.Sprintf("user:%d", user.ID)
    cache.Delete(key)

    return nil
}
```

### Pattern 4: Time-Based Invalidation

```go
func (c *Cache) StartCleanup(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            c.mu.Lock()
            now := time.Now()
            for key, entry := range c.entries {
                if now.After(entry.Expiry) {
                    delete(c.entries, key)
                    delete(c.lruMap, key)
                }
            }
            c.mu.Unlock()
        }
    }()
}
```

### Pattern 5: Circuit Breaker + Cache Fallback

```go
func GetData(key string) ([]byte, error) {
    // Try backend
    data, err := backend.Get(key)

    if err == nil {
        cache.Set(key, data)
        return data, nil
    }

    // Backend failed, try cache
    if cached, exists := cache.Get(key); exists {
        log.Println("Backend down, serving stale cache")
        return cached, nil
    }

    return nil, err
}
```

---

## 7. Real-World Applications

### Content Delivery Networks (CDNs)

**Use case**: Cache static assets (images, CSS, JS) at edge locations

**Companies**: Cloudflare, Fastly, Akamai

```go
// CDN edge server
cache := NewCache(10000, 24*time.Hour)  // Cache for 24 hours
proxy := httputil.NewSingleHostReverseProxy(originServer)
http.ListenAndServe(":80", cache.Handler(proxy))
```

### API Gateways

**Use case**: Cache API responses to reduce backend load

**Companies**: Kong, AWS API Gateway, Apigee

```go
// API Gateway with caching
cache := NewCache(1000, 5*time.Minute)
backend := httputil.NewSingleHostReverseProxy(apiServer)
http.ListenAndServe(":8080", cache.Handler(backend))
```

### Load Balancers

**Use case**: Distribute load + cache responses

**Companies**: nginx, HAProxy, Envoy

```go
backends := []*url.URL{backend1, backend2, backend3}
lb := NewLoadBalancer(backends)
cache := NewCache(5000, 10*time.Minute)
http.ListenAndServe(":80", cache.Handler(lb))
```

### Microservices Caching Layer

**Use case**: Cache inter-service communication

**Companies**: Netflix, Uber, Airbnb

```go
// Service A calls Service B through caching proxy
serviceB := "http://service-b:8000"
proxy := httputil.NewSingleHostReverseProxy(parseURL(serviceB))
cache := NewCache(500, 2*time.Minute)
client := &http.Client{Transport: cache.Handler(proxy)}
```

### WordPress/Blog Caching

**Use case**: Cache dynamically-generated pages

**Companies**: WordPress VIP, Medium

```go
// Cache blog posts
cache := NewCache(1000, 1*time.Hour)
wordpress := httputil.NewSingleHostReverseProxy(wordpressServer)
http.ListenAndServe(":80", cache.Handler(wordpress))
```

---

## 8. Common Mistakes to Avoid

### Mistake 1: Caching POST/PUT/DELETE

**❌ Wrong**:
```go
// Caching ALL requests
cache.Set(r.URL.String(), response)
```

**Problem**: POST/PUT/DELETE are not idempotent — caching them is dangerous.

**✅ Correct**:
```go
// Only cache GET requests
if r.Method != http.MethodGet {
    backend.ServeHTTP(w, r)
    return
}
cache.Set(r.URL.String(), response)
```

### Mistake 2: Not Considering Cache-Control

**❌ Wrong**:
```go
// Always cache
cache.Set(key, response)
```

**Problem**: Response might have `Cache-Control: no-store` (sensitive data).

**✅ Correct**:
```go
cacheControl := response.Header.Get("Cache-Control")
if strings.Contains(cacheControl, "no-store") {
    return  // Don't cache
}
cache.Set(key, response)
```

### Mistake 3: Ignoring Response Body Consumption

**❌ Wrong**:
```go
cache[url] = response  // Response.Body is io.Reader, consumed once
client1.Write(response.Body)  // Works
client2.Write(response.Body)  // FAILS! Body already read
```

**✅ Correct**:
```go
body, _ := io.ReadAll(response.Body)
cache[url] = body

// For each client:
response.Body = io.NopCloser(bytes.NewReader(body))
```

### Mistake 4: No Cache Size Limit

**❌ Wrong**:
```go
// Unlimited cache
cache[key] = value
```

**Problem**: Memory exhaustion if cache grows unbounded.

**✅ Correct**:
```go
if len(cache) >= maxSize {
    evictLRU()
}
cache[key] = value
```

### Mistake 5: No Thread Safety

**❌ Wrong**:
```go
cache[key] = value  // Race condition!
```

**Problem**: Concurrent reads/writes cause data races.

**✅ Correct**:
```go
mu.Lock()
cache[key] = value
mu.Unlock()
```

### Mistake 6: Caching Errors

**❌ Wrong**:
```go
response, _ := backend.Get(url)
cache.Set(url, response)  // Caches 500 errors!
```

**✅ Correct**:
```go
response, _ := backend.Get(url)
if response.StatusCode == http.StatusOK {
    cache.Set(url, response)
}
```

### Mistake 7: Forgetting to Clone Response

**❌ Wrong**:
```go
func serveFromCache(w http.ResponseWriter, entry *CacheEntry) {
    w.Write(entry.Body)
    // No headers!
}
```

**✅ Correct**:
```go
func serveFromCache(w http.ResponseWriter, entry *CacheEntry) {
    for k, v := range entry.Header {
        w.Header()[k] = v
    }
    w.WriteHeader(entry.StatusCode)
    w.Write(entry.Body)
}
```

---

## 9. Stretch Goals

### Goal 1: Implement Cache Statistics Dashboard ⭐⭐

Add HTTP endpoint to view cache stats.

**Hint**:
```go
func (c *Cache) StatsHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        stats := map[string]interface{}{
            "hits":     atomic.LoadInt64(&c.hits),
            "misses":   atomic.LoadInt64(&c.misses),
            "size":     c.lru.Len(),
            "max_size": c.maxSize,
            "hit_rate": float64(c.hits) / float64(c.hits+c.misses),
        }
        json.NewEncoder(w).Encode(stats)
    }
}
```

### Goal 2: Add Manual Cache Invalidation API ⭐⭐

Allow cache invalidation via HTTP endpoint.

**Hint**:
```go
// DELETE /cache/api/users/123 → Invalidate specific key
func (c *Cache) InvalidateHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        key := r.URL.Path[len("/cache"):]
        c.Delete(key)
        w.WriteHeader(http.StatusNoContent)
    }
}
```

### Goal 3: Implement Conditional Requests (ETags) ⭐⭐⭐

Support `If-None-Match` for bandwidth savings.

**Hint**:
```go
func (c *Cache) serveFromCache(w http.ResponseWriter, r *http.Request, entry *CacheEntry) {
    etag := entry.Header.Get("ETag")
    if r.Header.Get("If-None-Match") == etag {
        w.WriteHeader(http.StatusNotModified)
        return
    }

    w.Header().Set("ETag", etag)
    w.Write(entry.Body)
}
```

### Goal 4: Add Metrics with Prometheus ⭐⭐⭐

Export cache metrics to Prometheus.

**Hint**:
```go
var (
    cacheHits = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "cache_hits_total",
        Help: "Total cache hits",
    })
    cacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "cache_misses_total",
        Help: "Total cache misses",
    })
)

func init() {
    prometheus.MustRegister(cacheHits, cacheMisses)
}
```

### Goal 5: Implement Distributed Caching with Redis ⭐⭐⭐⭐

Use Redis as shared cache across multiple proxy instances.

**Hint**:
```go
import "github.com/go-redis/redis/v8"

type RedisCache struct {
    client *redis.Client
}

func (rc *RedisCache) Get(key string) (*CacheEntry, bool) {
    val, err := rc.client.Get(ctx, key).Bytes()
    if err != nil {
        return nil, false
    }

    var entry CacheEntry
    json.Unmarshal(val, &entry)
    return &entry, true
}
```

---

## How to Run

```bash
# Run the caching proxy
go run ./minis/36-caching-reverse-proxy/cmd/proxy/main.go

# Run tests
go test ./minis/36-caching-reverse-proxy/...

# Run with verbose output
go test -v ./minis/36-caching-reverse-proxy/...

# Run with race detector
go test -race ./minis/36-caching-reverse-proxy/...

# Test against real backend
# Terminal 1: Start backend
go run backend_server.go

# Terminal 2: Start proxy
go run cmd/proxy/main.go

# Terminal 3: Make requests
curl http://localhost:8080/api/data
curl http://localhost:8080/api/data  # Should be cached
```

---

## Summary

**What you learned**:
- ✅ Reverse proxies intercept and forward HTTP requests
- ✅ HTTP caching dramatically reduces backend load
- ✅ Cache-Control headers control caching behavior
- ✅ LRU eviction keeps hot data cached
- ✅ httputil.ReverseProxy makes proxy building easy
- ✅ Thread safety is critical for concurrent caching

**Why this matters**:
Caching is one of the most effective performance optimizations. A well-designed caching layer can reduce backend load by 90%+ and response times by 100x. Every major website uses caching proxies.

**Key formulas**:
- **Hit rate**: `hits / (hits + misses)`
- **Cache effectiveness**: `(original_latency - cached_latency) / original_latency`
- **Load reduction**: `1 - (1 / hit_rate)`

**Next steps**:
- Explore Redis for distributed caching
- Learn about CDN architectures
- Study cache invalidation patterns in microservices

Cache wisely!
