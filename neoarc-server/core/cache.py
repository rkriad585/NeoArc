import hashlib
import json
import time


class AliasCache:
    def __init__(self, ttl_seconds=30):
        self._cache: dict[str, dict] = {}
        self._ttl = ttl_seconds

    def get(self, key: str):
        entry = self._cache.get(key)
        if entry and time.time() - entry['time'] < self._ttl:
            return entry['data'], entry['etag']
        return None, None

    def set(self, key: str, data: dict):
        etag = hashlib.md5(json.dumps(data, sort_keys=True).encode()).hexdigest()
        self._cache[key] = {'data': data, 'etag': etag, 'time': time.time()}
        return etag

    def invalidate(self, key: str = None):
        if key:
            self._cache.pop(key, None)
        else:
            self._cache.clear()


alias_cache = AliasCache(ttl_seconds=30)
