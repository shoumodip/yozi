:i count 4
:b shell 25
../yozi -r arithmetics.yo
:i returncode 0
:b stdout 17
69
420
69
420
69

:b stderr 0

:b shell 19
../yozi -r block.yo
:i returncode 0
:b stdout 8
3
2
4
3

:b stderr 0

:b shell 23
../yozi -r condition.yo
:i returncode 0
:b stdout 2
2

:b stderr 0

:b shell 18
../yozi -r loop.yo
:i returncode 0
:b stdout 0

:b stderr 0

