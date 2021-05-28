package pipeline

import (
	"fmt"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/domain/status"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/usecases/link"
	"net/http"
	"sync"
	"time"
)

func LinkStatusUpdater(luc *link.LinkUseCases) {
	go func() {
		links := make(chan link.Link)

		go func() {
			for {
				userLinks, err := luc.LinkStorage.GetAllUserLinks()
				if err != nil {
					panic(err)
				}
				for _, lnk := range userLinks {
					links <- link.Link{LinkId: lnk.LinkId, Link: lnk.Link, LinkStatus: lnk.LinkStatus}
				}
				time.Sleep(5 * time.Second)
			}
		}()

		var wg sync.WaitGroup
		for i := 0; i < 4; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					lnk, ok := <-links
					if !ok {
						return
					}
					res, _ := http.Get(lnk.Link)
					s := status.OK
					if res.StatusCode == http.StatusNotFound {
						s = status.Failed
					}
					_ = luc.LinkStorage.UpdateLinkStatusByLinkId(lnk.LinkId, s)
					if lnk.LinkStatus != s {
						fmt.Printf("%v status changed from %v to %v\n",
							lnk.LinkId, lnk.LinkStatus, s)
					}
				}
			}()
		}

		wg.Wait()
	}()
}
