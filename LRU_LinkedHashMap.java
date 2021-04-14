import java.util.LinkedHashMap;
import java.util.Map;

// LRU Caching using LinkedHashMap

public class LRU_LinkedHashMap<S, T> {
    LinkedHashMap<S, T> cache;
    int capacity;

    LRU_LinkedHashMap(int capacity) {
        cache = new LinkedHashMap<S, T>(capacity);
        this.capacity = capacity;
    }

    T get(S key) {
        if (!cache.containsKey(key)) {
            return null;
        }
        T val = cache.remove(key);
        cache.put(key, val);
        return val;
    }

    void put(S key, T value) {
        if (cache.containsKey(key)) {
            cache.remove(key);
        } else if (cache.size() == capacity) {
            Map.Entry<S, T> entry = cache.entrySet().iterator().next();
            cache.remove(entry.getKey());
        }

        cache.put(key, value);
    }
} 