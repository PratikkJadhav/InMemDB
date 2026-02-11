

# InMemDB

**InMemDB** is a lightweight, high-performance in-memory key-value store built in Go. Designed to closely mirror the architecture of Redis, it implements a **single-threaded event loop** model to ensure command atomicity without using complex locking mechanisms (mutexes).

It features a custom implementation of the **RESP (Redis Serialization Protocol)**, making it compatible with standard Redis clients, and it supports three eviction policies.


<img width="2076" height="1288" alt="InMemDB" src="https://github.com/user-attachments/assets/5c66e031-1461-42fd-aa95-8da4a3384605" />


## Key Features

### 1. Architecture

* **Single-Threaded Design:** InMemDB processes commands sequentially on a main thread. This mimics Redis's core architecture, eliminating race conditions and context-switching for data operations.
* **RESP Protocol:** Built a custom parser for the Redis Serialization Protocol, allowing seamless communication with standard CLI tools.

### 2. Supported Commands

InMemDB supports a subset of core Redis commands across different data types and system operations:

* **String Operations:** `SET`, `GET`, `INCR`, `DEL`
* **Key Management:** `TTL`, `EXPIRE` (Time-To-Live support)
* **Hashes:** `HSET`, `HGET`, `HGETALL`, `HDEL`
* **System & Persistence:** `INFO`, `BGREWRITEAOF` (Background AOF Rewrite), `LRU`

### 3. Eviction

To handle memory pressure, InMemDB implements several eviction strategies that trigger when the key limit is reached:

* **`simple-first`:** FIFO (First-In-First-Out) eviction strategy.
* **`allkeys-random`:** Randomly selects and removes keys to free space.
* **`allkeys-lru`:** Approximated Least Recently Used eviction, sampling keys to discard those that haven't been accessed for the longest time. ( I implemented 40% bulk eviction ) 

---

## Some Challenges I faced while implementation of InMemDB:

**Memory Tracking vs. Key Counting**
One of the significant challenges I faced during development was implementing a memory-based eviction limit (e.g., `maxmemory 100mb`). Redis achieves this using a custom memory allocator wrapper (`zmalloc`) to track the exact byte size of every object.

* *Attempt:* I attempted to implement a similar wrapper around Go's memory allocation to track the exact size of keys and values.
* *Current Solution:* Due to the complexity of tracking low-level memory usage in a garbage-collected language like Go without using unsafe pointers, the current version defaults to a **Key-Count Limit** (currently set to `max_keys = 100` ) rather than a raw byte limit.

---

## Future Work:

* **Memory-Aware Eviction:** Revisit the memory tracking implementation to replace the "Max Keys" limit with a true "Max Memory" limit, potentially using a custom allocator.
* **LFU Policy:** Implement the **Least Frequently Used** eviction algorithm to handle access patterns where keys are accessed often but not necessarily recently.
* **Pub/Sub:** Add support for Publish/Subscribe messaging patterns (`PUBLISH`, `SUBSCRIBE`).

---

##  How to Run

1. **Clone the repository:**
```bash
git clone https://github.com/PratikkJadhav/InMemDB.git
cd InMemDB

```


2. **Start the Server:**
```bash
go run main.go --port 7379

```


3. **Connect with Redis CLI:**
```bash
redis-cli -p 7379


