package main

import (
	"fmt"
	syncMap "parallel-course-work/pkg/sync_map"
	"time"
)

func main() {
	//threadPool := threadpool.New(4, 1)
	//router := tcpRouter.New()
	//router.AddRoute(tcpRouter.POST, "status", func(ctx *tcpRouter.RequestContext) error {
	//	return ctx.ResponseJSON(tcpRouter.OK, "hello world!")
	//})
	//server := tcpServer.New(8080, threadPool, router)
	//err := server.Start()
	//if err != nil {
	//	log.Fatal(err)
	//}

	smap := syncMap.NewSyncHashMap[[]string](10, 0.75)
	var err error
	go func() {
		err := smap.Put("123", []string{"value", "value2"})

		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := smap.Put("123", []string{"value"})

		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(2 * time.Second)
	res, ok := smap.Get("123")
	if !ok {
		panic(err)
	}
	fmt.Println(res)
	res.Value = append(res.Value, "value3")
	res2, ok := smap.Get("123")
	if !ok {
		panic(err)
	}
	fmt.Println(res2)
}
