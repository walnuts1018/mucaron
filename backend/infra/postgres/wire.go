package postgres

import (
	"github.com/google/wire"
	"github.com/walnuts1018/mucaron/usecase/subjects"
)

var Set = wire.NewSet(
	NewPostgres,
	wire.Bind(new(subjects.SubjectRepository), new(*PostgresClient)),
)
