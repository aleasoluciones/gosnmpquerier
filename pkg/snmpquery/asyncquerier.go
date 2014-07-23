package snmpquery

type AsyncQuerier struct {
	Input      chan Query
	Output     chan Query
	Contention int
}

func NewAsyncQuerier(contention int) *AsyncQuerier {
	querier := AsyncQuerier{
		Input:      make(chan Query, 10),
		Output:     make(chan Query, 10),
		Contention: contention,
	}
	go querier.process()
	return &querier
}

func (querier *AsyncQuerier) process() {
	m := make(map[string]chan Query)

	for query := range querier.Input {
		_, exists := m[query.Destination]
		if exists == false {
			channel_tmp := make(chan Query, 10)
			m[query.Destination] = channel_tmp
			for i := 0; i < querier.Contention; i++ {
				go processQueriesFromChannel(channel_tmp, querier.Output)
			}
		}
		m[query.Destination] <- query
	}
}

func handleQuery(query *Query) {
	switch query.Cmd {
	case WALK:
		query.Response, query.Error = walk(query.Destination, query.Community, query.Oid, query.Timeout, query.Retries)
	case GET:
		query.Response, query.Error = get(query.Destination, query.Community, query.Oid, query.Timeout, query.Retries)
	}
}

func processQueriesFromChannel(input chan Query, processed chan Query) {
	for query := range input {
		handleQuery(&query)
		processed <- query
	}
}
