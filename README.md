## Канальный логер

### Примеры

```go
package main

import (
	"fmt"
	"github.com/xela07ax/chLogger"
	"time"
)

func main() {
	fmt.Println("Testing ch logger")
	logEr := chLogger.NewChLoger(&chLogger.Config{
		ConsolFilterFn: map[string]int{"Front Http Server": 0},
		ConsolFilterUn: map[string]int{"Pooling": 1},
		Mode:           0,
		Dir:            "x-loger",
	})
	logEr.RunLogerDaemon()
	logEr.ChInLog <- [4]string{"Welcome", "nil", fmt.Sprintf("Вас приветствует \"Silika-FileManager Контроллер\" v1.1 (11112020) \n")}
	fmt.Println("-main->wait")
	time.Sleep(1 * time.Second)

}

...