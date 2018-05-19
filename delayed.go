package pomfu

import "sync"

type Delay func(Response) Request

func merge(r1 Response, r2 Response) Response {
	for k, v := range r2 {
		r1[k] = v
	}
	return r1
}

func (r Request) AddDelayed(delay Delay) {
	r.delay = append(r.delay, delay)
}

func (r Request) processDelays(p *Pomf) (Response, error) {
	var nr Request
	for _, delay := range r.delay {
		req := delay(r.prev)
		nr.name = append(nr.name, req.name...)
		nr.out = append(nr.out, req.out...)
	}
	return nr.Upload(p)
}

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
				resp, err := sr.reqest.Upload(sr.pomf)
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
