:i count 60
:b testcase 23
integers/arithmetics.yo
:i returncode 0
:b stdout 74
69
420
69
420
69
1
0
0
1
1
0
1
0
0
1
1
0
1
0
1
0
0
1
69
420
69
420
69
420

:b stderr 0

:b testcase 26
integers/typed-literals.yo
:i returncode 0
:b stdout 18
69
420
1337
80085

:b stderr 0

:b testcase 37
integers/untyped-literal-auto-cast.yo
:i returncode 0
:b stdout 18
69
420
1337
80085

:b stderr 0

:b testcase 31
integers/error-type-mismatch.yo
:i returncode 1
:b stdout 0

:b stderr 72
integers/error-type-mismatch.yo:1:17: ERROR: Expected type i32, got i64

:b testcase 32
integers/error-invalid-suffix.yo
:i returncode 1
:b stdout 0

:b stderr 85
integers/error-invalid-suffix.yo:2:7: ERROR: Invalid suffix 'i69' to integer literal

:b testcase 53
integers/error-untyped-literal-auto-cast-too-large.yo
:i returncode 1
:b stdout 0

:b stderr 114
integers/error-untyped-literal-auto-cast-too-large.yo:2:16: ERROR: Integer literal '420' is too large for type i8

:b testcase 11
booleans.yo
:i returncode 0
:b stdout 64
1
0
0
1
69
420
1
69
420
0
69
0
69
0
69
1
69
1
69
420
1
69
420
0

:b stderr 0

:b testcase 8
block.yo
:i returncode 0
:b stdout 8
3
2
4
3

:b stderr 0

:b testcase 12
condition.yo
:i returncode 0
:b stdout 8
1
2
3
4

:b stderr 0

:b testcase 7
loop.yo
:i returncode 0
:b stdout 20
0
1
2
3
4
5
6
7
8
9

:b stderr 0

:b testcase 30
global-variables/definition.yo
:i returncode 0
:b stdout 12
69
420
1337

:b stderr 0

:b testcase 36
global-variables/definition-forms.yo
:i returncode 0
:b stdout 12
69
420
1337

:b stderr 0

:b testcase 30
global-variables/assignment.yo
:i returncode 0
:b stdout 7
69
420

:b stderr 0

:b testcase 35
global-variables/error-undefined.yo
:i returncode 1
:b stdout 0

:b stderr 76
global-variables/error-undefined.yo:2:11: ERROR: Undefined identifier 'foo'

:b testcase 38
global-variables/error-redefinition.yo
:i returncode 1
:b stdout 0

:b stderr 152
global-variables/error-redefinition.yo:2:5: ERROR: Redefinition of global identifier 'x'
global-variables/error-redefinition.yo:1:5: NOTE: Defined here

:b testcase 46
global-variables/error-assignment-undefined.yo
:i returncode 1
:b stdout 0

:b stderr 86
global-variables/error-assignment-undefined.yo:2:5: ERROR: Undefined identifier 'foo'

:b testcase 47
global-variables/error-assignment-not-memory.yo
:i returncode 1
:b stdout 0

:b stderr 97
global-variables/error-assignment-not-memory.yo:2:5: ERROR: Cannot assign to value not in memory

:b testcase 70
global-variables/error-assignment-does-not-match-type-in-definition.yo
:i returncode 1
:b stdout 0

:b stderr 112
global-variables/error-assignment-does-not-match-type-in-definition.yo:1:14: ERROR: Expected type bool, got i64

:b testcase 54
global-variables/error-cannot-define-with-unit-type.yo
:i returncode 1
:b stdout 0

:b stderr 103
global-variables/error-cannot-define-with-unit-type.yo:3:5: ERROR: Cannot define variable with type ()

:b testcase 29
local-variables/assignment.yo
:i returncode 0
:b stdout 7
69
420

:b stderr 0

:b testcase 29
local-variables/definition.yo
:i returncode 0
:b stdout 12
69
420
1337

:b stderr 0

:b testcase 35
local-variables/definition-forms.yo
:i returncode 0
:b stdout 12
69
420
1337

:b stderr 0

:b testcase 28
local-variables/shadowing.yo
:i returncode 0
:b stdout 7
69
420

:b stderr 0

:b testcase 28
local-variables/reference.yo
:i returncode 0
:b stdout 7
69
420

:b stderr 0

:b testcase 42
local-variables/error-use-outside-scope.yo
:i returncode 1
:b stdout 0

:b stderr 81
local-variables/error-use-outside-scope.yo:7:11: ERROR: Undefined identifier 'x'

:b testcase 61
local-variables/error-use-outside-scope-despite-same-depth.yo
:i returncode 1
:b stdout 0

:b stderr 100
local-variables/error-use-outside-scope-despite-same-depth.yo:6:15: ERROR: Undefined identifier 'x'

:b testcase 69
local-variables/error-assignment-does-not-match-type-in-definition.yo
:i returncode 1
:b stdout 0

:b stderr 111
local-variables/error-assignment-does-not-match-type-in-definition.yo:2:18: ERROR: Expected type bool, got i64

:b testcase 53
local-variables/error-cannot-define-with-unit-type.yo
:i returncode 1
:b stdout 0

:b stderr 102
local-variables/error-cannot-define-with-unit-type.yo:4:9: ERROR: Cannot define variable with type ()

:b testcase 33
pointers/reference-dereference.yo
:i returncode 0
:b stdout 7
69
420

:b stderr 0

:b testcase 42
pointers/reference-dereference-multiple.yo
:i returncode 0
:b stdout 31
69
420
420
420
420
69
420
1337

:b stderr 0

:b testcase 37
pointers/multiple-level-type-parse.yo
:i returncode 0
:b stdout 7
69
420

:b stderr 0

:b testcase 22
pointers/arithmetic.yo
:i returncode 0
:b stdout 2
1

:b stderr 0

:b testcase 46
pointers/error-dereference-expected-pointer.yo
:i returncode 1
:b stdout 0

:b stderr 85
pointers/error-dereference-expected-pointer.yo:2:6: ERROR: Expected pointer, got i64

:b testcase 38
pointers/error-reference-not-memory.yo
:i returncode 1
:b stdout 0

:b stderr 96
pointers/error-reference-not-memory.yo:2:6: ERROR: Cannot take reference of value not in memory

:b testcase 35
functions/no-arguments-no-return.yo
:i returncode 0
:b stdout 14
69
420
69
420

:b stderr 0

:b testcase 47
functions/no-arguments-no-return-first-class.yo
:i returncode 0
:b stdout 14
69
420
69
420

:b stderr 0

:b testcase 36
functions/yes-arguments-no-return.yo
:i returncode 0
:b stdout 7
69
420

:b stderr 0

:b testcase 48
functions/yes-arguments-no-return-first-class.yo
:i returncode 0
:b stdout 6
69
69

:b stderr 0

:b testcase 37
functions/yes-arguments-yes-return.yo
:i returncode 0
:b stdout 3
69

:b stderr 0

:b testcase 49
functions/yes-arguments-yes-return-first-class.yo
:i returncode 0
:b stdout 4
420

:b stderr 0

:b testcase 41
functions/arguments-as-local-variables.yo
:i returncode 0
:b stdout 14
69
420
69
420

:b stderr 0

:b testcase 30
functions/early-return-unit.yo
:i returncode 0
:b stdout 3
69

:b stderr 0

:b testcase 34
functions/early-return-not-unit.yo
:i returncode 0
:b stdout 3
69

:b stderr 0

:b testcase 22
functions/recursion.yo
:i returncode 0
:b stdout 43
1
2
6
24
120
720
5040
40320
362880
3628800

:b stderr 0

:b testcase 45
functions/recursion-of-entry-function-main.yo
:i returncode 0
:b stdout 23
10
9
8
7
6
5
4
3
2
1
0

:b stderr 0

:b testcase 33
functions/error-not-a-function.yo
:i returncode 1
:b stdout 0

:b stderr 73
functions/error-not-a-function.yo:2:5: ERROR: Expected function, got i64

:b testcase 43
functions/error-call-to-function-pointer.yo
:i returncode 1
:b stdout 0

:b stderr 110
functions/error-call-to-function-pointer.yo:6:5: ERROR: Cannot call pointer to function. Dereference it first

:b testcase 47
functions/error-direct-reference-to-function.yo
:i returncode 1
:b stdout 0

:b stderr 105
functions/error-direct-reference-to-function.yo:4:6: ERROR: Cannot take reference of value not in memory

:b testcase 42
functions/error-argument-count-mismatch.yo
:i returncode 1
:b stdout 0

:b stderr 83
functions/error-argument-count-mismatch.yo:4:8: ERROR: Expected 0 arguments, got 1

:b testcase 41
functions/error-argument-type-mismatch.yo
:i returncode 1
:b stdout 0

:b stderr 82
functions/error-argument-type-mismatch.yo:4:9: ERROR: Expected type i64, got bool

:b testcase 44
functions/error-return-type-expected-unit.yo
:i returncode 1
:b stdout 0

:b stderr 83
functions/error-return-type-expected-unit.yo:2:5: ERROR: Expected type (), got i64

:b testcase 48
functions/error-return-type-expected-not-unit.yo
:i returncode 1
:b stdout 0

:b stderr 87
functions/error-return-type-expected-not-unit.yo:2:5: ERROR: Expected type i64, got ()

:b testcase 39
functions/error-return-type-mismatch.yo
:i returncode 1
:b stdout 0

:b stderr 80
functions/error-return-type-mismatch.yo:2:5: ERROR: Expected type i64, got bool

:b testcase 26
type-cast/demonstration.yo
:i returncode 0
:b stdout 10
1
0
1
0
1

:b stderr 0

:b testcase 54
type-cast/error-cannot-cast-from-boolean-to-pointer.yo
:i returncode 1
:b stdout 0

:b stderr 98
type-cast/error-cannot-cast-from-boolean-to-pointer.yo:1:14: ERROR: Cannot cast from bool to &i64

:b testcase 54
type-cast/error-cannot-cast-from-pointer-to-boolean.yo
:i returncode 1
:b stdout 0

:b stderr 98
type-cast/error-cannot-cast-from-pointer-to-boolean.yo:2:11: ERROR: Cannot cast from &i64 to bool

:b testcase 56
type-cast/error-cannot-cast-from-function-to-anything.yo
:i returncode 1
:b stdout 0

:b stderr 100
type-cast/error-cannot-cast-from-function-to-anything.yo:3:13: ERROR: Cannot cast from fn () to i64

:b testcase 64
type-cast/error-cannot-cast-from-function-pointer-to-anything.yo
:i returncode 1
:b stdout 0

:b stderr 109
type-cast/error-cannot-cast-from-function-pointer-to-anything.yo:3:16: ERROR: Cannot cast from &fn () to i64

:b testcase 56
type-cast/error-cannot-cast-from-anything-to-function.yo
:i returncode 1
:b stdout 0

:b stderr 100
type-cast/error-cannot-cast-from-anything-to-function.yo:1:12: ERROR: Cannot cast from i64 to fn ()

:b testcase 64
type-cast/error-cannot-cast-from-anything-to-function-pointer.yo
:i returncode 1
:b stdout 0

:b stderr 109
type-cast/error-cannot-cast-from-anything-to-function-pointer.yo:1:12: ERROR: Cannot cast from i64 to &fn ()

