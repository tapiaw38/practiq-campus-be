# Practiq Campus BE architecture

## Request flow

Use this flow for every operation:

```text
handler -> usecase -> repository -> domain
```

- A handler (Gin) parses the HTTP request, calls a usecase and writes the response. It does not contain business logic.
- A usecase executes one operation per file (`create.go`, `get.go`, `list.go`, `update.go`, `delete.go`) and depends on `appcontext.Factory` to reach repositories.
- A repository implements persistence for one domain aggregate. Usecases depend on repository interfaces, not concrete implementations.
- Domain entities in `internal/domain` are structs. They may carry pure predicate methods over their own fields (e.g. a time-based state check) and the small methods needed to satisfy a standard interface such as `error`; they must not query the database, call a usecase or depend on transport types.

## File layout

```text
internal/domain/<entity>.go
internal/usecases/<domain>/<operation>.go
internal/usecases/<domain>/output_types.go
internal/adapters/datasources/repositories/<domain>/
internal/adapters/web/handlers/<domain>/
```

- One usecase interface, struct and constructor per operation: `<Operation>Usecase`, `<operation>Usecase`, `New<Operation>Usecase`.
- Keep usecase input types in the same file as the operation (`<Operation>Input`) and shared response DTOs in `output_types.go`.
- Register new usecases and repositories in `internal/usecases/usecases.go` and `internal/adapters/datasources/repositories/repositories.go`.
- Wrap persistence and validation errors with `apperrors.ApplicationError` and the matching entry in `internal/platform/errors/mappings`.

## Route registration

`RegisterRoutes` builds three groups off `/api`, and the order of the statements is load-bearing because Gin copies the parent's handler chain when a group is created, not when a route is added:

- `superAdminOnly` is created before `api.Use(middlewares.RequireTenant(...))` on purpose. Platform administration — listing institutions, enabling one, changing its status — is how a tenant comes to exist, so those routes have to stay reachable with no tenant selected.
- `teacherOnly` must be created *after* that `api.Use(...)` line. Moving it up would drop `RequireTenant` from every teacher route, so a global auth role would decide access before Campus knew which institution was selected.
- Routes whose access depends on enrollment rather than role (submissions, quiz attempts, forum posts, messaging) hang off the plain `api` group and check access inside their usecase.
- Self-enrolment has no route by design: an institution decides who attends what, so enrolment comes from the course's teacher or its admin.

## Platform boundaries

- Put cross-cutting infrastructure in `internal/platform/<area>` (`appcontext`, `auth`, `config`, `database`, `errors`, `identity`, `revocation`, `storage`, `tenantcontext`, `unlock`).
- A usecase or handler must not open its own database connection or JWT parser; it goes through `appcontext.Factory` and `internal/platform/auth`.
- Keep tenant resolution in `tenantcontext`; do not read tenant IDs from ad hoc request fields.

## Code comments

- Do not add explanatory comments to application code.
- Prefer explicit names, small functions and typed boundaries over comments.
- Keep rationale in documentation or commit messages.
- Exceptions: required legal notices, generated-file directives and tool directives.
