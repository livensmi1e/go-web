## Go web service template

> [!NOTE]
>
> **Purpose**: Use it whenever I want to start backend for something quickly — an idea, or a side project — without spending [much] time for internal structure.

### Development

```
make help
```

### Deployment

Make `.env.prod` with variables similar to `.env.dev`.

### Performance (need to improve alot)

```bash
.\hey.exe -n 10000 -c 50 -m POST -H "Content-Type: application/json" -d '{\"email\":\"user@email.com\",\"password\":\"password\"}' http://localhost:8000/api/auth/login

Summary:
  Total:        117.3703 secs
  Slowest:      1.5058 secs
  Fastest:      0.0920 secs
  Average:      0.5803 secs
  Requests/sec: 85.2005

  Total data:   2990000 bytes
  Size/request: 299 bytes

Response time histogram:
  0.092 [1]     |
  0.233 [51]    |■
  0.375 [414]   |■■■■
  0.516 [3013]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.658 [3946]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.799 [1850]  |■■■■■■■■■■■■■■■■■■■
  0.940 [549]   |■■■■■■
  1.082 [131]   |■
  1.223 [34]    |
  1.364 [8]     |
  1.506 [3]     |


Latency distribution:
  10% in 0.4180 secs
  25% in 0.4836 secs
  50% in 0.5645 secs
  75% in 0.6631 secs
  90% in 0.7661 secs
  95% in 0.8373 secs
  99% in 0.9908 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0016 secs, 0.0920 secs, 1.5058 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0073 secs
  req write:    0.0001 secs, 0.0000 secs, 0.0548 secs
  resp wait:    0.5782 secs, 0.0918 secs, 1.5056 secs
  resp read:    0.0003 secs, 0.0000 secs, 0.1491 secs

Status code distribution:
  [200] 10000 responses
```

### Improvement

-   Utilize concurrency
-   Apply cache mechanism for req/resp

### References

-   https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/ [Highly recommended]
-   https://github.com/manakuro/golang-clean-architecture/tree/main [Highly recommended]
-   https://github.com/evrone/go-clean-template/tree/master
-   https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
-   https://golang.cafe/blog/golang-functional-options-pattern.html
