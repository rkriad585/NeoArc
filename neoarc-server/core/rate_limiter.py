import time
from collections import defaultdict


class RateLimiter:
    def __init__(self):
        self._store: dict[str, list[float]] = defaultdict(list)

    def hit(self, key: str, max_attempts: int = 5, window_seconds: int = 60) -> bool:
        now = time.time()
        cutoff = now - window_seconds
        self._store[key] = [t for t in self._store[key] if t > cutoff]
        if len(self._store[key]) >= max_attempts:
            return True
        self._store[key].append(now)
        return False

    def clear(self):
        self._store.clear()


rate_limiter = RateLimiter()
