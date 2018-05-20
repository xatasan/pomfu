package pomfu

import "sync"

// The Delay type enables the user to generate requests based on
// previous responses. This can be helpful when uploading a HTML page,
// while embedding images or any equivalent interrelation of uploads
// where one file is dependant on another
type Delay func(Response) Request

// merge merges two responses, by overwriting the first (r1) with the
// keys from the second (r2). Hence if r1 and r2 contain the same key,
// the resulting response will use r2's keys instead of r1.
//
// Nevertheless, this behaviour should be avoided wherever possible by
// the user, to avoid confusion.
func merge(r1 Response, r2 Response) Response {
	for k, v := range r2 {
		r1[k] = v
	}
	return r1
}

func (r *Request) GenSubrequest(p *Pomf) Request {
	var req Request
	r.subrq = append(r.subrq, struct {
		pomf   *Pomf
		reqest Request
	}{p, req})
	return req
}

// Add a request to be processed after the "regular" requests. This
// "delayed request" uses a Delayed type (ie. a function from Response
// -> Request) that "calculates" the "actual" request based on the
// response of all the files uploaded until now.
func (r Request) AddDelayed(delay Delay) {
	r.delay = append(r.delay, delay)
}

// an internal method called by Upload to create and upload a the
// delayed request
func (r Request) processDelays(p *Pomf) (Response, error) {
	var nr Request
	for _, delay := range r.delay {
		req := delay(r.prev)
		nr.name = append(nr.name, req.name...)
		nr.out = append(nr.out, req.out...)
	}
	return nr.Upload(p)
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
