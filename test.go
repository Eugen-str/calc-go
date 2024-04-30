package main

import (
    "fmt"
)

type EvalTest struct{
    expression  string
    expected    float64
}

func eval_test(){
    var var_list []Var
    tests := []EvalTest{
        {"3 + 4", 7},
        {"3 + 4 + 5", 12},
        {"12 * 4", 48},
        {"4 * 12", 48},
        {"-6 - 4", -10},
        {"-6 -4", -10},
        {"3 + 4 - 5 + 6 - 7", 1},
        {"3 * 4 - 20 / 5 * 7", -16},
        {"3 * ( 3 + 4 )", 21},
        {"3 ^ 2", 9},

        {"3^2", 9},
        {"((7))", 7},
        {"((7)*(5 + 4))", 63},
        {"( 3 + 1 ) * ( 4 + 5 * ( 3 + 6 ) / 3 )", 76},
        {"(3 + 1) * (4+5*(3+6 ) /       3 )", 76},
        {"h = 6.626", 6.626},
        {"h", 6.626},
        {"a = 4", 4},
        {"h = 12345 + a", 12349},
        {"h + 4", 12353},
    }

    for i, test := range(tests){
        tokens, err := tokenize(test.expression)
        if(err != nil){
            fmt.Print(err.Error())
            continue
        }
        res, new_list, err := eval(tokens, var_list)
        var_list = new_list
        if(err != nil){
            fmt.Print(err.Error())
            continue
        }
        if(res == test.expected){
            fmt.Printf("\x1b[32m")
            fmt.Printf("TEST: test #%d passed - %.8f = %.8f\n", i + 1, test.expected, res)
            fmt.Printf("\x1B[37m")
        } else {
            fmt.Printf("\x1B[31m")
            fmt.Printf("TEST FAILED: test #%d failed - %.8f != %.8f\n", i + 1, test.expected, res)
            fmt.Printf("\x1B[37m")
        }
    }
}
