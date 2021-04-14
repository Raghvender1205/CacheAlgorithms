/**                 Least Frequently Used Cache 
 * We count how many times an item was
 * accessed and evict the item that has least access count.
*/
import java.util.Map;
import java.util.LinkedHashMap;
import java.util.Objects;

public class LFUCache<S, T> {
    public class Node<S, T> {
        S key;
        T value;
        Integer count;

        Node (S a, T b, Integer c) {
            key = a;
            value = b;
            count = c;
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) {
                return true;
            }

            if (o == null || getClass() != o.getClass()) {
                return false;
            }

            Node<S, T> node = (Node<S, T>) o;
            return key.equals(node.key) && value.equals(node.value) && count.equals(node.count);
        }

        @Override
        public int hashCode() {
            return Objects.hash(key, value, count);
        }
    }
    LinkedHashMap<S, Node<S, T>> cache;
    int capacity;

    public LFUCache(int capacity) {
        this.capacity = capacity;
        this.cache = new LinkedHashMap<>(capacity);
    }

    public void put(S key, T value) {
        if (cache.containsValue(key)) {
            Node<S, T> node = cache.get(key);
            cache.remove(key);
            node.value = value;
            node.count ++;
            cache.put(key, node);
            return;
        }

        if (isFull()) {
            if (this.capacity == 0) {
                return;
            }
            S minKey = null;
            int minFreq = Integer.MAX_VALUE;
            for (Map.Entry<S, Node<S, T>> entry : cache.entrySet()) {
                if (minFreq > entry.getValue().count) {
                    minFreq = entry.getValue().count;
                    minKey = entry.getKey();
                }
            }
            cache.remove(minKey);
        }

        Node<S, T> newNode = new Node<S, T>(key, value, 0);
        cache.put(key, newNode);
    }

    public T get(S key) {
        if (cache.containsKey(key)) {
            // Cache hit
            Node<S, T> node = cache.get(key);
            node.count ++;
            cache.remove(key, node);
            return node.value;
        }
        return null; // Cache miss
    }

    public boolean isFull() {
        if (this.cache.size() == this.capacity) {
            return true;
        }
        return false;
    }
}
