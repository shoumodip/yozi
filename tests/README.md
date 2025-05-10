# Black Box Behavior Tests

## Run tests
```console
$ ./rere.py replay test.list
```

## How to add a test?
- Make sure tests are currently passing

```console
$ ./rere.py replay test.list
```

- Create the test file
- Add the following to `test.list`

```
../yozi -r <file>
```

- Record the tests

```
../rere.py record test.list
```

- Run the tests again to ensure reproducibilty

```
../rere.py replay test.list
```
