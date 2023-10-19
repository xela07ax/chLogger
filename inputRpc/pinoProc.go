package inputRpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Rpc struct {
	Alias string
	//cfg *model.Config
	Services map[string]func([]byte)
	Loger    chan [4]string
}

// JSON MICRO EXAMPLE
//
//	{
//		"service": "go.tracker.svc.repiter",
//		"request": {"tom":"req","cid":123,"msg":"version"}
//	}
type Micro struct {
	Service string          `json:"service"`
	Request json.RawMessage `json:"request"`
}

func NewRpc(loger chan [4]string, services map[string]func([]byte)) *Rpc {
	return &Rpc{Alias: "ⓇⓅⒸ", Loger: loger, Services: services}
}

func (rpc *Rpc) InputRpc(w http.ResponseWriter, r *http.Request) {
	rpcMsg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ertx := fmt.Sprintf("COM:Ошибка чтения тела: %s | ERTX:can't read body", err)
		rpc.Loger <- [4]string{rpc.Alias, "nil", ertx, "1"}
		http.Error(w, ertx, http.StatusConflict) // 409
		return
	}
	err, resp := rpc.InputMsg(rpcMsg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		ertx := fmt.Sprintf("COM: не удалось отправить ответ | ERTX:%v", err)
		rpc.Loger <- [4]string{rpc.Alias, "nil", ertx, "1"}
		http.Error(w, ertx, http.StatusBadRequest) // 400
		return
	}
	rpc.Loger <- [4]string{rpc.Alias, "⚡𝓼𝓽𝓪𝓽𝓾𝓼 𝟮𝟬𝟬⚡", fmt.Sprintf("🅱🅾🅳🆈【%s】", string(resp)), "HTTP_WRITE"}
}
func (rpc *Rpc) InputMsg(rpcMsg []byte) (err error, resp []byte) {
	rpc.Loger <- [4]string{rpc.Alias, "⚡𝓻𝓮𝓺𝓾𝓮𝓼𝓽⚡", fmt.Sprintf("🅱🅾🅳🆈【%s】", rpcMsg), "HTTP_READ"}

	microRout := &Micro{}
	err = json.Unmarshal(rpcMsg, microRout)
	if err != nil {
		err = fmt.Errorf("COM:Ошибка чтения RPC: %s | ERTX:can't read RPC", err)
		rpc.Loger <- [4]string{rpc.Alias, "nil", err.Error(), "1"}
		return
	}
	svc, ok := rpc.Services[microRout.Service]
	if !ok {
		err = fmt.Errorf("COM:Ошибка поиска| service: %v | ERTX:can't find Service| use:%v", microRout.Service, rpc.Services)
		rpc.Loger <- [4]string{rpc.Alias, "nil", err.Error(), "1"}
		return
	}
	breakdat := make(chan struct{})
	go func() {
		svc(microRout.Request)
		breakdat <- struct{}{}
	}()

	select {
	case <-breakdat:
	case <-time.After(30 * time.Second):
		err = fmt.Errorf("COM: Сервис[%s] не принял данные |Err:_timeout_%d_", microRout.Service, 30)
		rpc.Loger <- [4]string{rpc.Alias, "nil", err.Error(), "1"}
		return
	}
	return
}
