package bramble

import "context"

var service *Service = &Service{}

type Query struct {
}

func (q *Query) Service(ctx context.Context) (*Service, error) {
	return service, nil
}
