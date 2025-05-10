# Black Box Behavior Tests

## Run tests
```console
$ ./rere.py replay test.list
```

## How to add a test?
1. Make sure tests are currently passing

```console
$ ./rere.py replay test.list
```

1. Create the test file
1. Add the following to `test.list`

```
../yozi -r <file>
```

1. Record the tests

```
../rere.py record test.list
```

1. Run the tests again to ensure reproducibilty

```
../rere.py replay test.list
```
