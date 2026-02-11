package dbase

// Query represents a generic database query with conditions, ordering, and pagination.
type Query struct {
	Conditions []Condition
	Limit      int
	Offset     int
	OrderBy    []Order
}

// Condition represents a single query condition.
type Condition struct {
	Field    string
	Operator Operator
	Value    any
	Or       bool // true = OR, false = AND
}

// Operator represents a comparison operator for query conditions.
type Operator string

const (
	OpEqual        Operator = "eq"
	OpNotEqual     Operator = "ne"
	OpGreater      Operator = "gt"
	OpGreaterEqual Operator = "ge"
	OpLess         Operator = "lt"
	OpLessEqual    Operator = "le"
	OpIn           Operator = "in"
	OpNotIn        Operator = "not_in"
	OpLike         Operator = "like"
	OpPrefix       Operator = "prefix"
	OpIsNull       Operator = "is_null"
	OpNotNull      Operator = "not_null"
)

// Order represents a sort order for query results.
type Order struct {
	Field      string
	Descending bool
}

// NewQuery creates a new empty query.
func NewQuery() *Query {
	return &Query{}
}

// Where adds an AND condition to the query.
func (q *Query) Where(field string, op Operator, value any) *Query {
	q.Conditions = append(q.Conditions, Condition{
		Field:    field,
		Operator: op,
		Value:    value,
	})
	return q
}

// Or adds an OR condition to the query.
func (q *Query) Or(field string, op Operator, value any) *Query {
	q.Conditions = append(q.Conditions, Condition{
		Field:    field,
		Operator: op,
		Value:    value,
		Or:       true,
	})
	return q
}

// OrderByAsc adds an ascending sort order.
func (q *Query) OrderByAsc(field string) *Query {
	q.OrderBy = append(q.OrderBy, Order{Field: field, Descending: false})
	return q
}

// OrderByDesc adds a descending sort order.
func (q *Query) OrderByDesc(field string) *Query {
	q.OrderBy = append(q.OrderBy, Order{Field: field, Descending: true})
	return q
}

// SetLimit sets the maximum number of results to return.
func (q *Query) SetLimit(limit int) *Query {
	q.Limit = limit
	return q
}

// SetOffset sets the number of results to skip.
func (q *Query) SetOffset(offset int) *Query {
	q.Offset = offset
	return q
}

// IsEmpty reports whether the query has no conditions.
func (q *Query) IsEmpty() bool {
	return q == nil || len(q.Conditions) == 0
}

// --- Shorthand constructors ---

// Eq creates a query with a single equality condition.
func Eq(field string, value any) *Query {
	return NewQuery().Where(field, OpEqual, value)
}

// Ne creates a query with a single not-equal condition.
func Ne(field string, value any) *Query {
	return NewQuery().Where(field, OpNotEqual, value)
}

// Gt creates a query with a single greater-than condition.
func Gt(field string, value any) *Query {
	return NewQuery().Where(field, OpGreater, value)
}

// Lt creates a query with a single less-than condition.
func Lt(field string, value any) *Query {
	return NewQuery().Where(field, OpLess, value)
}

// In creates a query with a single IN condition.
func In(field string, values ...any) *Query {
	return NewQuery().Where(field, OpIn, values)
}

// Like creates a query with a single LIKE condition.
func Like(field string, pattern string) *Query {
	return NewQuery().Where(field, OpLike, pattern)
}
