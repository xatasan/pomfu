package pomfu

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

// Add a request to be processed after the "regular" requests. This
// "delayed request" uses a Delayed type (ie. a function from Response
// -> Request) that "calculates" the "actual" request based on the
// response of all the files uploaded until now.
func (r Request) AddDelayed(delay func(Response) Request) {
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
	return p.upload(nr)
}
