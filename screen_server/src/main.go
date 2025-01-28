/*
 * @Author: gongluck
 * @Date: 2024-04-01 01:25:52
 * @Last Modified by: gongluck
 * @Last Modified time: 2025-01-28 16:49:54
 */

package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	fmt.Println("start server")

	//curl http://xxx/screen.txt
	//curl https://xxx/screen.txt
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

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
		err := http.ListenAndServeTLS(":8002", "./tls.crt", "./tls.key", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()

	wg.Wait()
}
