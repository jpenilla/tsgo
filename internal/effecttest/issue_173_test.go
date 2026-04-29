package effecttest_test

import (
	"testing"

	"github.com/microsoft/typescript-go/shim/fourslash"

	_ "github.com/effect-ts/tsgo/etscheckerhooks"
	_ "github.com/effect-ts/tsgo/etslshooks"
	_ "github.com/effect-ts/tsgo/etstesthooks"
)

// Original issue #173: this in directly yielded argument.
func TestIssue173_DirectYieldedThisDoesNotReportFalsePositiveTS2683(t *testing.T) {
	t.Parallel()

	const content = `// @Filename: /tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "target": "ESNext",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "plugins": [
      {
        "name": "@effect/language-service"
      }
    ]
  }
}
// @Filename: /repro.ts
import { Effect, Exit, Scope } from "effect"

class Repro {
  readonly #scope: Scope.Closeable = Scope.makeUnsafe()

  run(): Effect.Effect<void> {
    return Effect.gen({ self: this }, function* () {
      yield* Scope.close([|this|].#scope, Exit.void)
    })
  }
}`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyNonSuggestionDiagnostics(t, nil)
}

// Original issue #173: this hoisted out of direct yield.
func TestIssue173_HoistedThisScopeDoesNotReportTS2683(t *testing.T) {
	t.Parallel()

	const content = `// @Filename: /tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "target": "ESNext",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "plugins": [
      {
        "name": "@effect/language-service"
      }
    ]
  }
}
// @Filename: /repro.ts
import { Effect, Exit, Scope } from "effect"

class Repro {
  readonly #scope: Scope.Closeable = Scope.makeUnsafe()

  run(): Effect.Effect<void> {
    return Effect.gen({ self: this }, function* () {
      const scope = this.#scope
      yield* Scope.close(scope, Exit.void)
    })
  }
}`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyNonSuggestionDiagnostics(t, nil)
}

// this in callee of a multi-arg call inside Effect.gen.
func TestIssue173_ThisInCalleePosition(t *testing.T) {
	t.Parallel()

	const content = `// @Filename: /tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "target": "ESNext",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "plugins": [
      {
        "name": "@effect/language-service"
      }
    ]
  }
}
// @Filename: /repro.ts
import { Effect } from "effect"

class Repro {
  #inner(x: number, y: number): Effect.Effect<number> {
    return Effect.succeed(x + y)
  }

  run(): Effect.Effect<number> {
    return Effect.gen({ self: this }, function* () {
      return yield* this.#inner(1, 2)
    })
  }
}`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyNonSuggestionDiagnostics(t, nil)
}

// this inside a callback passed to Effect.tryPromise inside Effect.gen.
func TestIssue173_ThisInsideEffectTryPromiseCallback(t *testing.T) {
	t.Parallel()

	const content = `// @Filename: /tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "target": "ESNext",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "plugins": [
      {
        "name": "@effect/language-service"
      }
    ]
  }
}
// @Filename: /repro.ts
import { Effect } from "effect"

class Repro {
  readonly #value = 42

  run(): Effect.Effect<number> {
    return Effect.gen({ self: this }, function* () {
      return yield* Effect.tryPromise(() => Promise.resolve(this.#value))
    })
  }
}`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyNonSuggestionDiagnostics(t, nil)
}

// this inside Effect.tryPromise callback calling a two-arg method.
func TestIssue173_ThisInMultiArgCallInsideCallback(t *testing.T) {
	t.Parallel()

	const content = `// @Filename: /tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "target": "ESNext",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "plugins": [
      {
        "name": "@effect/language-service"
      }
    ]
  }
}
// @Filename: /repro.ts
import { Effect } from "effect"

class Repro {
  #inner(x: number, y: number): number {
    return x + y
  }

  run(): Effect.Effect<number> {
    return Effect.gen({ self: this }, function* () {
      return yield* Effect.sync(() => this.#inner(1, 2))
    })
  }
}`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyNonSuggestionDiagnostics(t, nil)
}

// this in callee with .pipe() chain inside Effect.gen.
func TestIssue173_ThisInCalleeWithPipeCall(t *testing.T) {
	t.Parallel()

	const content = `// @Filename: /tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "target": "ESNext",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "plugins": [
      {
        "name": "@effect/language-service"
      }
    ]
  }
}
// @Filename: /repro.ts
import { Effect } from "effect"

class Repro {
  #inner(x: number): Effect.Effect<number> {
    return Effect.succeed(x)
  }

  run(): Effect.Effect<number> {
    return Effect.gen({ self: this }, function* () {
      return yield* this.#inner(1).pipe(Effect.map((n) => n + 1))
    })
  }
}`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyNonSuggestionDiagnostics(t, nil)
}

// this used in a curried-call argument outside Effect.gen (window typing scenario).
// Reproduces the pattern from sessions.ts where this in `this.#lockFor(id)`
// was typed as Window & typeof globalThis in a regular class method.
func TestIssue173_ThisInPrivateMethodCallOutsideGen(t *testing.T) {
	t.Parallel()

	const content = `// @Filename: /tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "target": "ESNext",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "plugins": [
      {
        "name": "@effect/language-service"
      }
    ]
  }
}
// @Filename: /repro.ts
import { Effect } from "effect"

class Repro {
  #getWrapper(key: string): (e: Effect.Effect<void>) => Effect.Effect<void> {
    return (e) => e
  }

  discard(key: string): Effect.Effect<void> {
    return this.#getWrapper(key)(
      Effect.gen({ self: this }, function* () {
        return
      }),
    )
  }
}`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyNonSuggestionDiagnostics(t, nil)
}

// Generic 3-arg function with arrow callback inside nested Effect.gen({ self: this }).
func TestIssue173_Generic3ArgCallInsideNestedEffectGen(t *testing.T) {
	t.Parallel()

	const content = `// @Filename: /tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "target": "ESNext",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "plugins": [
      {
        "name": "@effect/language-service"
      }
    ]
  }
}
// @Filename: /repro.ts
import { Effect } from "effect"

function runMcpRequest<T>(
  name: string,
  label: string,
  request: () => PromiseLike<T>,
): Effect.Effect<T> {
  return Effect.succeed() as any
}

class Repro {
  run(server: { readonly name: string }): Effect.Effect<void> {
    return Effect.gen({ self: this }, function* () {
      return yield* Effect.gen({ self: this }, function* () {
        yield* runMcpRequest(server.name, "connect", () =>
          Promise.resolve(42),
        )
      })
    })
  }
}`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyNonSuggestionDiagnostics(t, nil)
}
