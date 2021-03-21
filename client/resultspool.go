// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

var freeResults = make(resultsPool, maxResults)

const (
	maxResults     = 1000
	defResultItems = 1000
)

type resultsPool chan []*result

func (p *resultsPool) get() []*result {
	select {
	case results := <-*p:
		return results
	default:
		return make([]*result, 0, defResultItems)
	}
}

func (p *resultsPool) put(results []*result) {
	results = results[:0]
	select {
	case *p <- results: // added
	default: // full
	}
}
