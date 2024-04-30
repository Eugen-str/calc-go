# calc-go

It's a calculator in golang for simple mathematical expressions

idk how go modules work so it's all just in a main file, maybe i'll fix that later

Help menu from the program:

```console
> help
| calc-go
| help menu
|
| it's just a calculator
| supported operations: + - * / ^
| brackets ( ) as well i guess
| for example: (3 + 1) * (4 + 5 * (3 + 6) / 3)
|
| defining variables
| for example: hi = 5
|
| You can also define functions like in math (NOT YET IMPLEMENTED)
| for example: f(x) = x^2 + 3*x + 5
|
| type "quit" when you're done (or just CTRL-C idc)
```

For now the biggest problem is the input method, it does not allow going left/right with arrow keys, making any kind of bigger equation pretty hard to write. Also error reporting is probably not very good, i haven't spent a lot of time catching the edge-cases, so some 'wrong' expressions might just crash or something.
