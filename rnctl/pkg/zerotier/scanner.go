// scanner.go
// worker pool端口扫描

package zerotier

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Scanner struct {
	TargetPort int
	Timeout    time.Duration
	Workers    int
}

func (st *Scanner) worker(id int, jobs <-chan string, results chan<- string, workboard *sync.WaitGroup) {
	for nwip := range jobs {
		addr := fmt.Sprintf("%s:%d", nwip, st.TargetPort)
		conn, err := net.DialTimeout("tcp", addr, st.Timeout)
		if err != nil {
			workboard.Done()
			continue
		}
		conn.Close()
		results <- nwip
		workboard.Done()
		// dont seem like the optimal solution. 255 workboard.Done()-s? (如果扫描255个子网段的话)
		// 对吗，好像对的，好像不对
		// 能跑通再说吧
	}
}

// 好像这个(worker)id放这里没吊用。先放着先好了
func (st *Scanner) Scan(nwips []string) []string {
	nwipsChan := make(chan string, len(nwips))
	resultChan := make(chan string, len(nwips))
	var workboard sync.WaitGroup

	for i := 0; i < st.Workers; i++ { // 启动 s.workers个goroutine
		go st.worker(i, nwipsChan, resultChan, &workboard)
	}

	// 把nwips所有切片丢进去
	for _, singleIP := range nwips {
		workboard.Add(1)
		nwipsChan <- singleIP
	}
	close(nwipsChan)
	// 收尾防卡死
	workboard.Wait()
	close(resultChan)

	var finalRes []string
	for ips := range resultChan {
		finalRes = append(finalRes, ips)
	}
	return finalRes
}
