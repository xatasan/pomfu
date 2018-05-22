package pomfu

import "sync"

// return a new sub request, which will be processed after the current request is finished
func (r Request) NewSubrequest(p *Pomf) Request {
	var req Request
	r.subrq = append(r.subrq, struct {
		pomf   *Pomf
		reqest Request
	}{p, req})
	return req
}

// internal method called by Upload to start subrequests in in parallel
func (r Request) processSubreq(res Response) ([]Response, error) {
	var (
		wg, fwg sync.WaitGroup
		donech  = make(chan bool, 1)
		errch   = make(chan error, 1)
		resch   = make(chan Response)
		resp    []Response
	)

	wg.Add(len(r.subrq))
	go func() {
		for _, subreq := range r.subrq {
			sr := subreq
			go func() {
				sr.reqest.prev = res
				resp, err := sr.pomf.upload(sr.reqest)
				if err != nil {
					errch <- err
				}
				resch <- resp
				wg.Done()
			}()
		}

		wg.Wait()
		donech <- true
		close(resch)
	}()

	fwg.Add(1)
	go func() {
		for r := range resch {
			resp = append(resp, r)
		}
		fwg.Done()
	}()

	select {
	case <-donech:
		fwg.Wait()
		return resp, nil
	case err := <-errch:
		return nil, err
	}
}
