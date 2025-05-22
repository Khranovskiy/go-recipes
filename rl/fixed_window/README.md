# Rate Limiting Algorithms

This package implements rate limiting algorithms for controlling the rate of operations in a system.

## Fixed Window Rate Limiting

The current implementation uses a Fixed Window rate limiting algorithm, which:
- Enforces a maximum number of operations per second
- Uses a fixed time window (1 second)
- Spaces out operations evenly within the window
- Provides simple throttling mechanism

## Rate Limiting Algorithms Comparison

### Leaky Bucket Algorithm

**How it works:**
- Think of it as a bucket with a hole at the bottom
- Requests (water) pour into the bucket at any rate
- The bucket leaks at a constant rate (processing rate)
- If the bucket is full, new requests are rejected/dropped

**Pros:**
1. Provides a smooth, constant output rate
2. Good for protecting downstream systems from bursty traffic
3. Predictable resource usage
4. Simple to implement
5. Memory efficient (only needs to track bucket capacity)

**Cons:**
1. Can be too restrictive for bursty but legitimate traffic
2. No concept of "saving up" capacity for future bursts
3. May drop requests even when system could handle them
4. Less flexible than token bucket

**Best used when:**
- You need to ensure a constant processing rate
- Protecting downstream systems from traffic spikes
- Working with systems that can't handle bursty traffic
- Need predictable resource usage
- Examples: API gateways, database access control, network traffic shaping

### Token Bucket Algorithm

**How it works:**
- Tokens are added to a bucket at a constant rate
- Each request consumes a token
- If there are no tokens, requests are rejected
- Tokens can accumulate up to a maximum capacity
- Allows for bursty traffic up to the bucket capacity

**Pros:**
1. Allows for bursty traffic while maintaining average rate
2. More flexible than leaky bucket
3. Better for handling legitimate traffic spikes
4. Can "save up" capacity for future bursts
5. More suitable for modern web applications

**Cons:**
1. Less predictable output rate
2. Can still allow bursts that might overwhelm systems
3. Slightly more complex to implement
4. Needs to track both token generation rate and bucket capacity

**Best used when:**
- You need to allow for bursty traffic
- Working with modern web applications
- Need to handle legitimate traffic spikes
- Want to provide better user experience
- Examples: API rate limiting, web application throttling, CDN rate limiting

## Key Differences

1. **Traffic Handling:**
   - Leaky Bucket: Forces constant output rate, no bursts allowed
   - Token Bucket: Allows bursts up to bucket capacity

2. **Request Rejection:**
   - Leaky Bucket: Rejects when bucket is full
   - Token Bucket: Rejects when no tokens are available

3. **Resource Management:**
   - Leaky Bucket: More predictable resource usage
   - Token Bucket: More flexible but less predictable

4. **Implementation Complexity:**
   - Leaky Bucket: Simpler to implement
   - Token Bucket: More complex but more flexible

## When to Choose Which

Choose **Leaky Bucket** when:
- You need strict rate limiting
- System resources are limited
- Predictable output rate is crucial
- Protecting downstream systems is the priority

Choose **Token Bucket** when:
- You want to allow for legitimate traffic bursts
- User experience is a priority
- System can handle occasional spikes
- You need more flexible rate limiting

In modern web applications, Token Bucket is often preferred because it provides a better balance between protection and flexibility. However, Leaky Bucket might be more appropriate in systems where strict rate control is necessary, such as in network traffic shaping or when protecting sensitive backend systems. 