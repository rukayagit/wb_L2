```Golang
package main
 
func main() {
    ch := make(chan int)
    go func() {
        for i := 0; i < 10; i++ {
            ch <- i
        }
    }()
 
    for n := range ch {
        println(n)
    }
}

```

Ответ: 
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
Deadlock

Читаем из канала даже после того, как в него ничего не записывается, отсюда получаем deadlock.
Если закрыть канал после цикла в горутине, программа отработает без deadlock.