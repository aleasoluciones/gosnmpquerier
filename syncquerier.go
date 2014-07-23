package gosnmpquerier

type SyncQuerier struct {
	Input        chan QueryWithOutputChannel
	asyncQuerier *AsyncQuerier
}

func NewSyncQuerier(contention int) *SyncQuerier {
	querier := SyncQuerier{
		Input:        make(chan QueryWithOutputChannel),
		asyncQuerier: NewAsyncQuerier(contention),
	}
	go querier.processAndDispatchQueries()
	return &querier
}

func (querier *SyncQuerier) ExecuteQuery(query Query) Query {
	output := make(chan Query)
	querier.Input <- QueryWithOutputChannel{query, output}
	processedQuery := <-output
	return processedQuery
}

func (querier *SyncQuerier) processAndDispatchQueries() {

	m := make(map[int]chan Query)
	i := 0
	for {
		select {
		case queryWithOutputChannel := <-querier.Input:
			queryWithOutputChannel.query.Id = i
			i += 1
			m[queryWithOutputChannel.query.Id] = queryWithOutputChannel.responseChannel
			querier.asyncQuerier.Input <- queryWithOutputChannel.query
		case processedQuery := <-querier.asyncQuerier.Output:
			m[processedQuery.Id] <- processedQuery
			delete(m, processedQuery.Id)
		}
	}
}
