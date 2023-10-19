## Канальный вебсокет логер Ws Logger Plugin. Плагин для Chlogger

### Описание
Возможность отправлять команды

### Примеры
Gin example router. 
Подключаться к консоли: http://localhost:8187/

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xela07ax/chLogger"
	"github.com/xela07ax/chLogger/inputRpc"
	"github.com/xela07ax/chLogger/wsLoggerPlugin"
	"time"
)

var Cxlogger chan [4]string

func main() {
	// Создаем ws плагин логер
	logErWs := wsLoggerPlugin.NewWsLogger()
	go logErWs.Run()

	// Создаем логер
	logEr := chLogger.NewChLoger(&chLogger.Config{
		Dir:       "x-loger",
		Broadcast: logErWs.Input,
	})
	logEr.RunLogerDaemon()
	logErWs.Loger = logEr.ChInLog
	Cxlogger = logEr.ChInLog
	Cxlogger <- [4]string{"Welcome", "nil", "Вас приветствует Контроллер v1.1"}

	services := make(map[string]func([]byte))
	rpc := inputRpc.NewRpc(Cxlogger, services)
	logErWs.Interpretator = rpc.InputMsg
	services["go.tracker.svc.repiter"] = func(bytes []byte) {
		Cxlogger <- [4]string{"go.tracker.svc.repiter", "nil", fmt.Sprintf("%s", bytes)}
	}

	// Полезная нагрузка
	go checkWsConnection(1)

	// Для коннекта нужно прокинуть наружу
	router := gin.New()

	router.GET("/", gin.WrapF(logErWs.HomePageWs))
	router.GET("/sentws", gin.WrapF(logErWs.SentWS))
	router.GET("/ws", gin.WrapF(logErWs.ServeWs))

	go checkWsConnection(1)
	port := "8187"
	// Listen and serve on defined port
	Cxlogger <- [4]string{"Welcome", "nil", fmt.Sprintf("Listening on port %s", port)}
	Cxlogger <- [4]string{"Welcome", "nil", fmt.Sprintf("Please see: http://localhost:%s/", port)}
	router.Run(":" + port)

}

func checkWsConnection(i int) {
	// Запускаем таймер для периодической проверки состояния соединения
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	Cxlogger <- [4]string{"checkWsConnection", "nil", fmt.Sprintf("Запускаем таймер для периодической проверки:%d", i)}
	for range ticker.C {
		Cxlogger <- [4]string{"checkWsConnection", "nil", fmt.Sprintf("Каая то ошибка проверки подключения:%d |ok", i), "ERROR"}
	}
	fmt.Println("good by")
}

```
Пример из консоли:
```
█║ (•̪●)=︻╦̵̵̿╤──:{ "service": "go.tracker.svc.repiter", "request": {"tom":"req","cid":123,"msg":"version"} }
(╯°o°）╯│▌ 2023-10-19 13:31:53 | FUNC:WS_Client | UNIT: readPump | TIP: |TEXT: 【{ "service": "go.tracker.svc.repiter", "request": {"tom":"req","cid":123,"msg":"version"} }】
(╯°o°）╯│▌ 2023-10-19 13:31:53 | FUNC:ⓇⓅⒸ | UNIT: ⚡𝓻𝓮𝓺𝓾𝓮𝓼𝓽⚡ | TIP:HTTP_READ |TEXT: 【🅱🅾🅳🆈【{ "service": "go.tracker.svc.repiter", "request": {"tom":"req","cid":123,"msg":"version"} }】】
(╯°o°）╯│▌ 2023-10-19 13:31:53 | FUNC:go.tracker.svc.repiter | UNIT: nil | TIP: |TEXT: 【{"tom":"req","cid":123,"msg":"version"}】
(╯°o°）╯│▌ 2023-10-19 13:31:54 | FUNC:checkWsConnection | UNIT: nil | TIP:ERROR |TEXT: 【Каая то ошибка проверки подключения:1 |ok】
(╯°o°）╯│▌ 2023-10-19 13:30:12 | FUNC:WS_Client | UNIT: readPump | TIP: |TEXT: 【init】
(╯°o°）╯│▌ 2023-10-19 13:30:12 | FUNC:WS_Client | UNIT: readPump | TIP: |TEXT: 【wait new message from ws client】
█║ (•̪●)=︻╦̵̵̿╤──:ваы
█║ ¯\_(ツ)_/¯ err:COM:Ошибка чтения RPC: invalid character 'Ð' looking for beginning of value | ERTX:can't read RPC
(╯°o°）╯│▌ 2023-10-19 13:30:15 | FUNC:WS_Client | UNIT: readPump | TIP: |TEXT: 【ваы】
(╯°o°）╯│▌ 2023-10-19 13:30:15 | FUNC:ⓇⓅⒸ | UNIT: ⚡𝓻𝓮𝓺𝓾𝓮𝓼𝓽⚡ | TIP:HTTP_READ |TEXT: 【🅱🅾🅳🆈【ваы】】
(╯°o°）╯│▌ 2023-10-19 13:30:15 | FUNC:ⓇⓅⒸ | UNIT: nil | TIP:1 |TEXT: 【COM:Ошибка чтения RPC: invalid character 'Ð' looking for beginning of value | ERTX:can't read RPC】
(╯°o°）╯│▌ 2023-10-19 13:30:24 | FUNC:checkWsConnection | UNIT: nil | TIP:ERROR |TEXT: 【Каая то ошибка проверки подключения:1 |ok】
(╯°o°）╯│▌ 2023-10-19 13:30:39 | FUNC:checkWsConnection | UNIT: nil | TIP:ERROR |TEXT: 【Каая то ошибка проверки подключения:1 |ok】
(╯°o°）╯│▌ 2023-10-19 13:30:54 | FUNC:checkWsConnection | UNIT: nil | TIP:ERROR |TEXT: 【Каая то ошибка проверки подключения:1 |ok】
```