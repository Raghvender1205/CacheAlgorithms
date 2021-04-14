import java.util.LinkedHashMap;
import java.util.Map;

public class FIFOCache<S, T> {
    LinkedHashMap<S, T> cache;
    int capacity;
    FIFOCache(int capacity) {
        cache = new LinkedHashMap<>(capacity);
        this.capacity = capacity;
    } 

    T get(S key) {
        if (!cache.containsKey(key)) {
            return null;
        }
        return cache.get(key);
    }

    void put(S key, T value) {
        if (cache.containsKey(key)) {
            cache.replace(key, value);
            return;
        } else if (cache.size() == capacity) {
            Map.Entry<S, T> entry = cache.entrySet().iterator().next();
            cache.remove(entry.getKey());
        }    
        cache.put(key, value);
    }

}