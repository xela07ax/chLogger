## Канальный логер

### Описание
Вывод в STDOUT и сразу вывод в файл, группируются файлы по первым ячейкам. Возможность транслировать к подписчикам, к примеру на ws logger. Не надо тянуть тип/пакет, так как всё на стандартных типах, только канал с 4-я текстовыми ячейками. Первая ячейка - это название паета/объекта, вторая - наименование функции/"nil", третья - содержание, четвертая - error("1")/norm(""/"0"). Имеет фильтры по первым двум ячейкам, можно отключать логирование блоков мешающих аналитике исполнения ПО.

### Применимость
Очень удобно при дебаге проекта который в продакшене. Не задерживает выполнение программы, можно добавлять обработчиков. 
Не очень хорошо ведет себя при неожиданном завершении программы, можно потерять много логов, поскольку они в очереди! Имейте в виду, объект логера имеет функции корректного завершения.

### Дополнительно
Имеет дополнительный запас полезных функций из пакета toolPack/tp)

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
	logEr.ChInLog <- [4]string{"Welcome", "nil", "Вас приветствует Silika-FileКонтроллер v1.1"}
	fmt.Println("-main->wait")
	logEr.ChInLog <- [4]string{"Welcome", "nil", "Передаем ошибку", "1"}
	time.Sleep(1 * time.Second)

}
```


```go
package main

import (
	"fmt"
	"github.com/xela07ax/chLogger"
	"time"
)

var Cxlogger chan [4]string

func main() {
	fmt.Println("Testing ch logger")
	logEr := chLogger.NewChLoger(&chLogger.Config{
		ConsolFilterFn: map[string]int{"Front Http Server": 0},
		ConsolFilterUn: map[string]int{"Pooling": 1},
		Mode:           0,
		Dir:            "x-loger",
	})
	logEr.RunLogerDaemon()
	Cxlogger = logEr.ChInLog
	Cxlogger <- [4]string{"Welcome", "nil", "Вас приветствует Silika-FileКонтроллер v1.1"}
	go checkWsConnection(1)
	time.Sleep(1 * time.Minute)
}

func checkWsConnection(i int) {
	// Запускаем таймер для периодической проверки состояния соединения
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	Cxlogger <- [4]string{"checkWsConnection", "nil", fmt.Sprintf("Запускаем таймер для периодической проверки:%d", i)}
	for range ticker.C {
		Cxlogger <- [4]string{"checkWsConnection", "nil", fmt.Sprintf("проверка подключения:%d. Ошибка |ok", i), "ERROR"}
	}
	fmt.Println("good by")
}
```

