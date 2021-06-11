package pipeline

import (
	"fmt"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/domain/status"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/usecases/link"
	"net/http"
	"time"
)

const workerCount = 4

func LinkStatusUpdater(luc *link.LinkUseCases) {
	go func() {
		links := make(chan link.Link, workerCount)

		go func() {
			for {
				userLinks, err := luc.LinkStorage.GetAllUserLinks()
				if err != nil {
					fmt.Println("Bad `GetAllUserLinks` request to link database")
				}
				for _, lnk := range userLinks {
					links <- link.Link{LinkId: lnk.LinkId, Link: lnk.Link, LinkStatus: lnk.LinkStatus}
				}
				time.Sleep(5 * time.Second)
			}
		}()

		for i := 0; i < workerCount; i++ {
			go func() {
				for {
					lnk, ok := <-links
					if !ok {
						fmt.Println("links channel has been closed and drained")
						return
					}
					res, _ := http.Get(lnk.Link)
					s := status.OK
					if res.StatusCode == http.StatusNotFound {
						s = status.Failed
					}
					err := luc.LinkStorage.UpdateLinkStatusByLinkId(lnk.LinkId, s)
					if err != nil {
						fmt.Println(err)
					}
					if lnk.LinkStatus != s {
						fmt.Printf("%v status changed from %v to %v\n",
							lnk.LinkId, lnk.LinkStatus, s)
					}
				}
			}()
		}
	}()
}
