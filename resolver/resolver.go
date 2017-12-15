package resolver

import (
	"context"
	"sync"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/util"

	"github.com/shurcooL/graphql"
)

type name struct {
	Value string `json:"value"`
}

type batch []name

func Run(n []string, qs chan *Query, gnf *gnfinder.GnFinder) {
	var client = graphql.NewClient(gnf.URL, nil)
	jobs := make(chan batch)
	var wgJobs sync.WaitGroup

	wgJobs.Add(gnf.Workers)
	for i := 1; i <= gnf.Workers; i++ {
		go resolverWorker(i, jobs, qs, &wgJobs, client)
	}

	go prepareJobs(n, jobs, gnf)
	wgJobs.Wait()
	close(qs)
}

func resolverWorker(i int, jobs <-chan batch, qs chan<- *Query,
	wg *sync.WaitGroup, client *graphql.Client) {
	defer wg.Done()
	for b := range jobs {
		var q Query
		variables := map[string]interface{}{"names": b}
		err := client.Query(context.Background(), &q, variables)
		util.Check(err)
		qs <- &q
	}
}

func prepareJobs(n []string, jobs chan<- batch, gnf *gnfinder.GnFinder) {
	l := len(n)

	for i := 0; i < l; i += gnf.BatchSize {
		end := i + gnf.BatchSize
		if end > l {
			end = l
		}
		jobs <- newBatch(n[i:end])
	}

	close(jobs)
}

func newBatch(n []string) (b batch) {
	b = make(batch, len(n))
	for i, v := range n {
		b[i] = name{Value: v}
	}
	return
}
