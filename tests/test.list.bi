:i count 10
:b shell 44
../yozi -r -o arithmetics.exe arithmetics.yo
:i returncode 0
:b stdout 49
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

:b stderr 0

:b shell 32
../yozi -r -o block.exe block.yo
:i returncode 0
:b stdout 8
3
2
4
3

:b stderr 0

:b shell 40
../yozi -r -o condition.exe condition.yo
:i returncode 0
:b stdout 8
1
2
3
4

:b stderr 0

:b shell 30
../yozi -r -o loop.exe loop.yo
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

:b shell 76
../yozi -r -o global-variables/definition.exe global-variables/definition.yo
:i returncode 0
:b stdout 12
69
420
1337

:b stderr 0

:b shell 76
../yozi -r -o global-variables/assignment.exe global-variables/assignment.yo
:i returncode 0
:b stdout 7
69
420

:b stderr 0

:b shell 86
../yozi -r -o global-variables/error-undefined.exe global-variables/error-undefined.yo
:i returncode 1
:b stdout 0

:b stderr 75
global-variables/error-undefined.yo:1:7: ERROR: Undefined identifier 'foo'

:b shell 92
../yozi -r -o global-variables/error-redefinition.exe global-variables/error-redefinition.yo
:i returncode 1
:b stdout 0

:b stderr 152
global-variables/error-redefinition.yo:2:5: ERROR: Redefinition of global identifier 'x'
global-variables/error-redefinition.yo:1:5: NOTE: Defined here

:b shell 108
../yozi -r -o global-variables/error-assignment-undefined.exe global-variables/error-assignment-undefined.yo
:i returncode 1
:b stdout 0

:b stderr 86
global-variables/error-assignment-undefined.yo:1:1: ERROR: Undefined identifier 'foo'

:b shell 110
../yozi -r -o global-variables/error-assignment-not-memory.exe global-variables/error-assignment-not-memory.yo
:i returncode 1
:b stdout 0

:b stderr 97
global-variables/error-assignment-not-memory.yo:1:1: ERROR: Cannot assign to value not in memory

