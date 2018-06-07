# v2

## what

It's a JavaScript engine. Because the world doesn't have enough of those already.

This one is written entirely in Go, currently as two packages - a parser, and a
VM. The VM generates bytecode, and then executes it, rather than operating on
the AST directly.

## why?

I felt like learning more about what's involved in JavaScript. I work with the
details of JavaScript a lot already [on QML](https://doc.qt.io/qt-5/qtqml-index.html)'s
JavaScript engine, v4, but since I didn't write most of that, I didn't get the
chance to learn about it from the ground up. v2 is my attempt at learning some
of those things, and having some fun along the way.

The name (v2) was chosen because I'm unimaginative. There's v8, there's v4. And
since I'm writing something that is much less capable than either of them, that
only leaves me with a limited set of smaller numbers. :)

## status

This thing can run some limited amounts of code, but it will die if you look at
it wrong. There's a complete shortage of error handling, and generally many
things are not yet implemented {fully,properly,at all}.

A (very incomplete) list of missing things:
 
* exceptions
* arrays
* most built-in objects
* spec compliance
* implicit semicolon handling at parse time
* regular expressions

I expect this will improve, assuming the project holds enough interest for me
to keep working on it, but consider yourself warned.
