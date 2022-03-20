# go-sparseset

This library provides an implementation of [sparse
sets](https://dl.acm.org/doi/10.1145/176454.176484) and utilities to operate on
sparse sets, including, efficient traversal, joining, sorting, etc.

Values are stored contiguously in memory therefore traversing a set is
particularly fast because it takes advantage of the underlying CPU caching.

```go
// Construction
set := sparseset.New[string](4096, 1<<20)
*set.Add(id) = value
...

// Properties.
length := set.Length()

// Removal.
set.Remove(id)

// Lookup.
if value, ok := set.Get(id); ok {
  // Do something with value...
} else {
  // Set does not contain id...
}

// Traversal.
for iterator := sparseset.Iterate(set); ; {
  key, value, ok := iterator.Next()
  if !ok {
    break
  }

  // Do something with key and value...
}

// Joining 2 sets.
for iterator := sparseset.Join(set1, set2); ; {
  key, value1, value2, ok := iterator.Next()
  if !ok {
    break
  }

  // Do something with key, value1, and value2...
}

// Joining 3 sets.
for iterator := sparseset.Join3(set1, set2, set3); ; {
  key, value1, value2, value3, ok := iterator.Next()
  if !ok {
    break
  }

  // Do something with key, value1, value2, and value3...
}
```
