```Golang
package main

import (
  "fmt"
)

func main() {
  var s = []string{"1", "2", "3"}
  modifySlice(s)
  fmt.Println(s)
}

func modifySlice(i []string) {
  i[0] = "3"
  i = append(i, "4")
  i[1] = "5"
  i = append(i, "6")
}
```

Ответ: [3 2 3]

Слайс - это структура из 3 полей: указатель на базовый массив, длина и capacity. 
При передаче в функцию структуры, слайс внутри функции будет работать с тем же массивом, что и внешний слайс.
Тем самым 0 элемент массива будет изменен. Но затем после первого append создается новая структура слайса, 
который ссылается на новый массива, и слайс s не ссылается на него.
