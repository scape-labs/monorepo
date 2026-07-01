# Bounded contexts

> TODO: flesh out each context with its aggregate roots, key events, and which service is the publisher vs subscriber.

| Service     | Aggregate roots (examples)                          | Publishes                                | Subscribes to                                |
| ----------- | --------------------------------------------------- | ---------------------------------------- | -------------------------------------------- |
| `bulksms`   | Campaign, Message, Recipient                        | `message.sent`, `message.failed`         | `billing.credited` (from wirepay)            |
| `wirepay`   | Account, Transfer, Balance                          | `transfer.completed`, `billing.credited` | `flow.task.completed`                        |
| `emicro`    | Loan, Repayment, Disbursement                       | `loan.disbursed`, `loan.repaid`          | `wirepay.transfer.completed`                 |
| `flow`      | Workflow, Task                                      | `flow.task.completed`                    | `*` (it's the orchestrator)                  |
| `flow-ussd` | Session, Menu                                       | (none — it's a UI)                       | n/a                                          |

## Rules of engagement

1. **No service imports another service's `internal/` package.** Cross-service communication is over AMQP or typed RPC (see `kit/rpc`).
2. **Shared types live in `libraries/`.** Money, tenant, ids, audit events.
3. **Tier 0 services are the source of truth.** If you need their data, subscribe to their events.
