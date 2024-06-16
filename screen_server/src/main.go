/*
 * @Author: gongluck
 * @Date: 2024-04-01 01:25:52
 * @Last Modified by: gongluck
 * @Last Modified time: 2024-05-19 01:52:39
 */

package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	fmt.Println("start server")

	//curl http://www.gongluck.fun/static/screen.txt
	//curl https://www.gongluck.fun/static/screen.txt
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := http.ListenAndServe(":8001", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := http.ListenAndServeTLS(":8002", "../../cert/gongluck.fun.crt", "../../cert/gongluck.fun.key", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()

	wg.Wait()
}
