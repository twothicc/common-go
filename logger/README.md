# Logger

This logger module handles the initialization and configuration of a zap logger, with utility functions that handles adding and propagating values to be logged via the Go context.

# Usage

At the start of your program, call `InitLogger(<isTest>)` where `<isTest>` refers to whether the development environment is test. This will initialize and configure the global zap logger.

---

Ideally, at the start of every request handler, call `logger.NewLogContext(ctx, zap.Stringer("reqId", reqId))` where `reqId` can be a UUID representing a specific request.

**Note**: Only need to be called once.
**Note**: Can be called again to add more log values.

---

Call e.g. `logger.WithContext(ctx).Info("hello world", zap.String("msg", "hello world"))` to log with context.
