## Канальный логер

### Примеры

```go
...
func (c *Client) writePump() {
...
		case message, ok := <-c.send:
            log.Printf("msg:%s",message)
            // Сообщение в текстовом формате
            	c.hub.WebSocketOutput <- message
            // Отправляем сообщение другому обработчику
			...
        case <-ticker.C:
...