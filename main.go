package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const(
    // TOKENS
    TYPE_PLUS  = iota
    TYPE_MINUS = iota
    TYPE_MULT  = iota
    TYPE_DIV   = iota
    TYPE_EXP   = iota

    TYPE_NUM   = iota

    TYPE_OPEN_PAREN   = iota
    TYPE_CLOSED_PAREN = iota

    TYPE_VAR    = iota
    TYPE_ASSIGN = iota

    // PRIORITIES
    ASSIGN_PRI = 1
    PLUS_PRI   = 2
    MINUS_PRI  = 2
    MULT_PRI   = 3
    DIV_PRI    = 3
    EXP_PRI    = 4
    PAREN_PRI  = 5
    MAX_PRI = 5
)

type Token struct{
    ttype int
    priority int
    value string
}

type Var struct{
    name string
    value string
}

func isNumeric(c byte) bool{
    return c >= '0' && c <= '9'
}

func isAlph(c byte) bool {
    return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func parse(input string) []string {
    var isOperation = func(c byte) bool{
        switch c{
        case '*', '+', '-', '/', '^':
            return true
        default:
            return false
        }
    }
    var symbols []string

    input = strings.Trim(input, "\n")
    input = strings.Trim(input, " ")

    p := 0
    curr := ""
    var prev byte
    for p < len(input){
        if(curr != "" && input[p] == '.' && isNumeric(input[p+1])){
            curr += string(input[p])
            curr += string(input[p+1])
            prev = input[p+1]
            p += 2
            continue
        }
        if(curr == "" && input[p] == '-' && isNumeric(input[p+1])){
            curr += string(input[p])
            curr += string(input[p+1])
            prev = input[p+1]
            p += 2
            continue
        }
        if(isNumeric(prev) && !isNumeric(input[p]) || isOperation(prev) && !isOperation(input[p]) && curr != ""){
            symbols = append(symbols, curr)
            curr = ""
        }
        switch input[p]{
        case ' ':
            if(curr != ""){
                symbols = append(symbols, curr)
            }
            curr = ""
        case '(', ')':
            symbols = append(symbols, string(input[p]))
            curr = ""
        default:
            curr += string(input[p])
        }

        prev = input[p]
        p += 1
    }
    if(curr != ""){
        symbols = append(symbols, curr)
    }

    /*
    fmt.Printf("SYMBOLS:\n")
    for _, s := range symbols{
        fmt.Printf(".%s.", s)
    }
    */
    return symbols
}

func isValidName(name string) bool{
    if(len(name) == 0 || !isAlph(name[0])){
        return false
    }
    for i := 1; i < len(name); i++{
        if(!isAlph(name[i]) || !isNumeric(name[i])){
            return false
        }
    }
    return true
}

func getToken(str string) (Token, error) {
    _, notNum := strconv.ParseFloat(str, 64)

    if(notNum != nil){
        switch str{
        case "+":
            return Token{ttype: TYPE_PLUS, priority: PLUS_PRI}, nil
        case "-":
            return Token{ttype: TYPE_MINUS, priority: MINUS_PRI}, nil
        case "*":
            return Token{ttype: TYPE_MULT, priority: MULT_PRI}, nil
        case "/":
            return Token{ttype: TYPE_DIV, priority: DIV_PRI}, nil
        case "^", "**":
            return Token{ttype: TYPE_EXP, priority: EXP_PRI}, nil
        case "(":
            return Token{ttype: TYPE_OPEN_PAREN, priority: PAREN_PRI}, nil
        case ")":
            return Token{ttype: TYPE_CLOSED_PAREN, priority: PAREN_PRI}, nil
        case "=":
            return Token{ttype: TYPE_ASSIGN, priority: ASSIGN_PRI}, nil
        default:
            if(isValidName(str)){
                return Token{ttype: TYPE_VAR, value: str}, nil
            }
            return Token{}, fmt.Errorf("* ERROR: unknown token \"%s\"\n", str)
        }
    }else{
        return Token{ttype: TYPE_NUM, value: str}, nil
    }
}

func tokenize(input string) ([]Token, error) {
    var tokens []Token

    symbols := parse(input)
    //fmt.Printf("TOKENS:")
    for _, str := range(symbols){
        token, err := getToken(str)
        if(err != nil){
            return nil, fmt.Errorf("%s", err)
        }
        
        tokens = append(tokens, token)
        //fmt.Printf("%v  ", token)
    }
    if(len(tokens) == 0){
        return []Token{{ttype: TYPE_NUM, value: "0"}}, nil
    }
    return tokens, nil
}

func removeElement(list []Token, idx int) []Token{
    return append(list[:idx], list[idx+1:]...)
}

func removeElements(list []Token, start, end int) []Token{
    return append(list[:start], list[end:]...)
}

func insertAt(list []Token, el Token, idx int) []Token{
    list = append(list[:idx+1], list[idx:]...)
    list[idx] = el
    return list
}

func eval(tokens []Token, var_list []Var) (float64, []Var, error) {
    var remove_operands = func(tokens []Token, c int) []Token {
        tokens = removeElement(tokens, c+2)
        tokens = removeElement(tokens, c+1)
        return tokens
    }
    var get_max_priority = func(tokens []Token) int {
        max := 1
        for _, token := range(tokens){
            if(token.priority > max){
                max = token.priority
            }
        }
        return max
    }
    var get_value_from_var_list = func(var_list []Var, name string) (string, bool) {
        for _, va := range var_list{
            if(va.name == name){
                return va.value, true
            }
        }
        return "", false
    }
    var get_operator = func(t Token, var_list []Var, isAssign bool) (float64, bool, error) {
        if(isAssign){ return -1, true, nil }
        if(t.ttype == TYPE_VAR){
            val, exists := get_value_from_var_list(var_list, t.value)
            if(!exists){
                return -1, true, fmt.Errorf("* ERROR: variable %s not defined\n", t.value)
            }
            num, err := strconv.ParseFloat(val, 64)
            if(err != nil){
                return -1, true, fmt.Errorf("* ERROR1: error converting \"%s\" to number\n", t.value)
            }
            return num, true, nil
        } else {
            num, err := strconv.ParseFloat(t.value, 64)
            if(err != nil){
                return -1, false, fmt.Errorf("* ERROR1: error converting \"%s\" to number\n", t.value)
            }
            return num, true, nil
        }
    }
    var assign_var = func(va Var, var_list []Var) []Var {
        for i, _ := range var_list{
            if(var_list[i].name == va.name){
                var_list[i].value = va.value
                return var_list
            }
        }
        return append(var_list, va)
    }

    c := 0
    priority := get_max_priority(tokens)
    for len(tokens) != 1 {
        //if(len(tokens) == 2){
        //    return -1, var_list, fmt.Errorf("* ERROR: invalid expression\n")
        //}
        if(c == len(tokens) - 1){
            c = 0
            priority = get_max_priority(tokens)
        }
        if(priority <= 0) { 
            return -1, var_list, fmt.Errorf("* ERROR: cannot evaluate expression\n")
        }

        t1 := tokens[c]
        op := tokens[c+1]

        if(op.ttype == TYPE_NUM && op.value[0] == '-'){
            tokens = insertAt(tokens, Token{ttype: TYPE_MINUS, priority: MINUS_PRI}, c+1)
            t1 = tokens[c]
            op = tokens[c+1]
            tokens[c+2].value = tokens[c+2].value[1:]
            priority = get_max_priority(tokens)
        }
        t2 := tokens[c+2]

        k := 0
        if(t2.ttype == TYPE_OPEN_PAREN){
            k = 2
        }

        if(t1.ttype == TYPE_OPEN_PAREN || t2.ttype == TYPE_OPEN_PAREN){
            var block []Token

            s := 1
            block_len := 0
            for{
                if(tokens[c + 1 + k + block_len].ttype == TYPE_CLOSED_PAREN) {s--}
                if(s == 0) { break }
                if(tokens[c + 1 + k + block_len].ttype == TYPE_OPEN_PAREN) {s++}
                block = append(block, tokens[c + 1 + k + block_len])
                block_len++
            }

            block_val, _, err := eval(block, var_list)
            if(err != nil){
                return -1, var_list, fmt.Errorf("* ERROR: error evaluating block\n")
            }

            tokens = removeElements(tokens, c + k, c + block_len + 2 + k)
            tokens = insertAt(tokens, Token{ttype: TYPE_NUM, value: strconv.FormatFloat(block_val, 'f', -1, 64)}, c + k)
            continue
        }

        if(op.priority == priority){
            var result float64

            num1, isVar1, err := get_operator(t1, var_list, op.ttype == TYPE_ASSIGN)
            if(isVar1){
                tokens[c].ttype = TYPE_NUM
            }
            if(err != nil){
                return -1, var_list, err
            }
            num2, isVar2, err := get_operator(t2, var_list, false)
            if(isVar2){
                tokens[c].ttype = TYPE_NUM
            }
            if(err != nil){
                return -1, var_list, err
            }

            switch op.ttype{
            case TYPE_MULT:
                result = num1 * num2
            case TYPE_DIV:
                if(num2 == 0) { return -1, nil, fmt.Errorf("* MATH ERROR: division by zero\n")}
                result = num1 / num2
            case TYPE_PLUS:
                result = num1 + num2
            case TYPE_MINUS:
                result = num1 - num2
            case TYPE_EXP:
                result = math.Pow(num1, num2)
            case TYPE_ASSIGN:
                if(t1.ttype != TYPE_VAR && !isVar1){
                    return -1, var_list, fmt.Errorf("* ERROR: wrong use of assignment operator\n")
                }
                var_list = assign_var(Var{name: t1.value, value: t2.value}, var_list)
                result = num2
                tokens[c].ttype = TYPE_NUM
            default:
                return -1, var_list, fmt.Errorf("* ERROR: use of unimplemented operation\n")
            }
            tokens = remove_operands(tokens, c)
            tokens[c].value = strconv.FormatFloat(result, 'f', -1, 64)
            continue
        }
        c += 2
    }

    final, _, err := get_operator(tokens[0], var_list, false)
    if(err != nil){
        return -1, var_list, err
    }

    return final, var_list, nil
}

func print_help(){
    fmt.Println("| calc-go")
    fmt.Println("| help menu")
    fmt.Println("| ")
    fmt.Println("| it's just a calculator")
    fmt.Println("| supported operations: + - * / ^")
    fmt.Println("| brackets ( ) as well i guess")
    fmt.Println("| for example: (3 + 1) * (4 + 5 * (3 + 6) / 3)")
    fmt.Println("| ")
    fmt.Println("| defining variables")
    fmt.Println("| for example: hi = 5")
    fmt.Println("| ")
    fmt.Println("| You can also define functions like in math (NOT YET IMPLEMENTED)")
    fmt.Println("| for example: f(x) = x^2 + 3*x + 5")
    fmt.Println("| ")
    fmt.Println("| type \"quit\" when you're done (or just CTRL-C idc)")
}

func main2(){
    eval_test()
}

func main(){
    var var_list []Var
    fmt.Println("* calc-go")
    fmt.Println("| type \"help\" to see the help menu!")
    in := bufio.NewReader(os.Stdin)
    for{
        fmt.Print("> ")
        input, err := in.ReadString('\n')

        if(err != nil){
            fmt.Println("* ERROR: input error")
            break
        }

        if(input == "quit\n"){
            break
        } else if(input == "vars\n"){
            if(len(var_list) == 0){
                fmt.Println("| no defined variables")
                continue
            }

            fmt.Println("| variables:")
            for _, v := range var_list{
                fmt.Printf("| %s = %s\n", v.name, v.value)
            }
            continue
        } else if(input == "help\n"){
            print_help()
            continue
        }

        tokens, err := tokenize(input)

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

        fmt.Printf("| %.8v\n", res)
    }
    fmt.Println("* bye!")
}
