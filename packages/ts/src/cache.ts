/**
 * Default in-memory LRU cache backend for the Transport's TTL response cache.
 * Plug-in users implement CacheBackend to swap in Redis/memcached.
 */

import type { CacheBackend, CacheEntry } from './transport_types.js';

export class InMemoryLRUCache implements CacheBackend {
  private readonly max: number;
  private readonly entries = new Map<string, CacheEntry>();

  constructor(max = 256) {
    this.max = max;
  }

  get(key: string): CacheEntry | null {
    const entry = this.entries.get(key);
    if (entry === undefined) return null;
    if (Date.now() - entry.storedAt > entry.ttlMs) {
      this.entries.delete(key);
      return null;
    }
    // why: re-insert to mark as recently used.
    this.entries.delete(key);
    this.entries.set(key, entry);
    return entry;
  }

  set(key: string, entry: CacheEntry): void {
    if (this.entries.has(key)) this.entries.delete(key);
    this.entries.set(key, entry);
    while (this.entries.size > this.max) {
      const firstKey = this.entries.keys().next().value;
      if (firstKey === undefined) break;
      this.entries.delete(firstKey);
    }
  }

  delete(key: string): void {
    this.entries.delete(key);
  }

  clear(): void {
    this.entries.clear();
  }
}
