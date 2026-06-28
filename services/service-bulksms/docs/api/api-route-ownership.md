# API route ownership

| Route                     | Method | Owner team        | Notes                          |
| ------------------------- | ------ | ----------------- | ------------------------------ |
| `/healthz`                | GET    | platform          | Liveness.                      |
| `/readyz`                 | GET    | platform          | Readiness.                     |
| TODO                      |        |                   |                                |

Every new route must add a row here in the same PR that adds the route.
