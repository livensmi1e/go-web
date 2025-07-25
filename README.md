## Go web service template

> [!NOTE]
>
> **Purpose**: I use it whenever I want to start backend for something quickly — an idea, a tool, or a side project — without worrying [much] about architecture structure.

### Development

Reference `Makefile` for more details.

```
make install
make infra-dev-up
make server-dev
```

### Deployment

Make `.env.prod` with variables similar to `.env.dev`.

### References

-   https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/ [Highly recommended]
-   https://github.com/manakuro/golang-clean-architecture/tree/main [Highly recommended]
-   https://github.com/evrone/go-clean-template/tree/master
-   https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
-   https://golang.cafe/blog/golang-functional-options-pattern.html
