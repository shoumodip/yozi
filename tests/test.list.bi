:i count 22
:b testcase 14
arithmetics.yo
:i returncode 0
:b stdout 70
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
69
420
69
420
69
420

:b stderr 0

:b testcase 11
booleans.yo
:i returncode 0
:b stdout 8
1
0
0
1

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
global-variables/error-assignment-does-not-match-type-in-definition.yo:2:18: ERROR: Expected type bool, got i64

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

