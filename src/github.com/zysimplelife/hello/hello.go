package main

import (
    "fmt"
    "math/rand"
	"github.com/zysimplelife/stringutil"
)

const Pi = 3.14 // const shouldn't use :=

func main() {
    fmt.Println("Generate radom munber ", rand.Intn(19))
    fmt.Println("add result = ", add(5,9))
    a, b := swap("abc","bcd")
    fmt.Println("swap result = ",a,b)
    a, b = namedReturn("abc","bcd")
    fmt.Println("default result ",a,b)
    varcall()
    forloop()
    forloopNoIndex()
    whileInGoSpelledFor()
	fmt.Println(stringutil.Reverse("!oG ,olleH"))
}

func add(x, y int) int {
    return x + y
}

func swap(x, y string) (string,string){
    return y, x
}

func namedReturn(x, y string) (r1 ,r2 string){
    r2 = "hello world!"
    return
}

var c, python, java bool

func varcall(){
    var i int
    fmt.Println(i, c, python, java)
}

func forloop(){
    sum := 0
    for i:=0; i < 10 ; i++ {
        sum +=i
    }
    fmt.Println(sum)
}


func forloopNoIndex(){
    sum := 1
    for ;sum < 1000 ; {
        sum += sum
    }
    fmt.Println(sum)
}

func whileInGoSpelledFor(){
    sum :=1
    for sum < 1000 {
        sum += sum
    }
    fmt.Println(sum)
}



