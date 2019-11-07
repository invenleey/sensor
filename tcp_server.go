package sensor

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func OpenKit() {

	http.HandleFunc("/", index)
	http.HandleFunc("/test/", test)

	if err := http.ListenAndServe("0.0.0.0:6666", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("Form: ", r.Form)
	fmt.Println("Path: ", r.URL.Path)
	fmt.Println(r.Form["a"])
	fmt.Println(r.Form["b"])
	for k, v := range r.Form {
		fmt.Println(k, "=>", v, strings.Join(v, "-"))
	}
	fmt.Fprint(w, "It works !")
}

type RequestBody struct {
	SData  []byte
	NodeIP string
}

func CreateRequest(r *http.Request) (RequestBody, error) {
	_ = r.ParseForm()
	if len(r.Form["addr"]) == 0 || len(r.Form["reg"]) == 0 || len(r.Form["nodeIP"]) == 0 {
		return RequestBody{}, errors.New("error request")
	}
	reg := r.Form["reg"][0]
	i, _ := strconv.Atoi(r.Form["addr"][0])
	var sr []byte
	sr = append(sr, byte(i))
	sr = append(sr, InfoMK["ReadFunc"]...)
	sr = append(sr, InfoMK[reg]...)
	sr = append(sr, CreateCRC(sr)...)
	var rq RequestBody
	rq.SData = sr
	rq.NodeIP = r.Form["nodeIP"][0]
	return rq, nil
}

func test(w http.ResponseWriter, r *http.Request) {
	if rq, err := CreateRequest(r); err != nil {
		var rd ReadResult
		rd.Status = 2
		if bs, err := json.Marshal(rd); err == nil {
			_, err := w.Write(bs)
			if err != nil {
				log.Println("发送操作失败: ", err)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		b, _ := GetDeviceSession(rq.NodeIP)
		p, err := b.MeasureRequest(rq.SData, []string{"测量值", "温度"})
		if err == nil {
			if bs, err := json.Marshal(p); err == nil {
				_, err := w.Write(bs)
				if err != nil {
					log.Println("发送操作失败: ", err)
				}
			} else {
				fmt.Println(err)
			}
		}
	}
}
