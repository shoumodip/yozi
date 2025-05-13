# Yozi
Compiler in Go

**NOTE: Due to technical difficulties with the C GO ABI, I am unable to maintain this project anymore. Therefore I am archiving this repo. I will rewrite this "language" in C, and maybe rebrand it under https://github.com/shoumodip/glos**

## Quick Start
Install [`clang`](https://clang.llvm.org/)

```console
$ go build
```

## Behaviour Tests
```console
$ cd tests
$ ./rere.py replay test.list
```

## Demonstration
**NOTE: This compiler is currently incomplete and under heavy development, thus
things can and will change at any moment.**

### Integers and Booleans
```rust
// This is a comment!
fn main() {
    #print 69
    #print 420

    #print true
    #print false
}
```

Integer types supported:
- `i8`
- `i16`
- `i32`
- `i64`
- `u8`
- `u16`
- `u32`
- `u64`

```rust
fn main() {
    #print 69u32 // Integer literals can be suffixed with their type
}
```

Typical arithmetic, bitwise, and logical operators work as expected.

### If Statements
```rust
fn main() {
    if true {
        #print 69
        #print 420
    }

    if false {
        #print 69
    } else {
        #print 1337
    }

    if false {
        #print 69
    } else if true {
        #print 80085
    } else {
        #print 420
    }
}
```

### Variables
```rust
let globalVar1 = 69
let globalVar2 bool = true
let globalVar3 i64

fn main() {
    let localVar1 = 69
    let localVar2 bool = true
    let localVar3 i64

    globalVar3 = 420
    localVar3 = 1337

    #print globalVar1
    #print globalVar2
    #print globalVar3

    #print localVar1
    #print localVar2
    #print localVar3
}
```

#### Local Variable Scoping
```rust
fn main() {
    {
        let x = 69
        #print x
    }

    #print x // ERROR: Undefined identifier 'x'
}
```

#### Local Variable Shadowing
```rust
fn main() {
    let x = 69
    #print x

    let x = x == 69
    #print x
}
```

### While Loops
```rust
fn main() {
    let i = 0
    while i < 10 {
        #print i
        i = i + 1
    }
}
```

### Functions
```rust
fn add(x i64, y i64) i64 {
    return x + y
}

fn double(x i64) i64 {
    return x * 2
}

fn printNum(x i64) {
    #print x
}

fn main() {
    printNum(double(add(100, 110)))
}
```

#### Recursion
```rust
fn factorial(n i64) i64 {
    if n < 2 {
        return 1
    }

    return n * factorial(n - 1)
}

fn main() {
    let i = 1
    while i <= 10 {
        #print factorial(i)
        i = i + 1
    }
}
```

#### First Class Functions
```rust
fn apply(x i64, f fn (i64) i64) i64 {
    return f(x)
}

fn double(x i64) i64 {
    return x * 2
}

fn main() {
    #print apply(210, double)
}
```

### Pointers
```rust
fn inc(x &i64) {
    *x = *x + 1
}

fn main() {
    let x = 68
    inc(&x)

    #print x
}
```

### Type Cast
```rust
fn main() {
    let x = 69
    if x as bool {
        #print true as i64
    } else {
        #print 1337 as &bool
    }
}
```
