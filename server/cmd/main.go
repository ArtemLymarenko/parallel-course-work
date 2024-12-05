package main

import (
	"fmt"
	syncMap "parallel-course-work/pkg/sync_map"
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

	smap := syncMap.NewSyncHashMap[[]string](10, 0.75, 10)
	err := smap.Insert("123", []string{"value"})
	if err != nil {
		panic(err)
	}

	res, ok := smap.Get("123")
	if !ok {
		panic(err)
	}

	res.Value = append(res.Value, "value2")
	res2, ok := smap.Get("123")
	if !ok {
		panic(err)
	}
	fmt.Println(res2)
}
