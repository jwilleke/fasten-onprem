package models

import (
	"fmt"
	"regexp"
	"strings"
)

// allowedAggregationFunctions is the allowlist of SQL aggregation functions that
// may be interpolated into a query. The aggregation `fn` is concatenated directly
// into the SQL string (e.g. `max(...)`), so it MUST be constrained to this set to
// prevent SQL injection. An empty function (no wrapper) is also permitted.
var allowedAggregationFunctions = map[string]bool{
	"":      true,
	"count": true,
	"sum":   true,
	"avg":   true,
	"min":   true,
	"max":   true,
}

// aggregationFieldRegex constrains an aggregation `field` value. The field name is
// already validated against the per-resource search-parameter allowlist, but the
// optional `:modifier` portion is interpolated raw into a SQL JSON path (`'$.<mod>'`),
// so the whole token is restricted to a safe charset. The literal `*` (count-all) is
// handled separately by the caller.
var aggregationFieldRegex = regexp.MustCompile(`^[A-Za-z0-9_-]+(:[A-Za-z0-9_-]+)?$`)

// validateAggregation enforces the injection-safety constraints on a single
// aggregation clause (`count_by` / `group_by` / `order_by`): a populated field drawn
// from a safe charset and a function drawn from the allowlist.
func validateAggregation(name string, agg *QueryResourceAggregation) error {
	if len(agg.Field) == 0 {
		return fmt.Errorf("if '%s' is present, field must be populated", name)
	}
	if agg.Field != "*" && !aggregationFieldRegex.MatchString(agg.Field) {
		return fmt.Errorf("%s field contains invalid characters: %q", name, agg.Field)
	}
	if !allowedAggregationFunctions[strings.ToLower(agg.Function)] {
		return fmt.Errorf("%s has unsupported aggregation function: %q", name, agg.Function)
	}
	return nil
}

// maps to frontend/src/app/models/widget/dashboard-widget-query.ts
type QueryResource struct {
	Use    string                 `json:"use"`
	Select []string               `json:"select"`
	From   string                 `json:"from"`
	Where  map[string]interface{} `json:"where"`
	Limit  *int                   `json:"limit,omitempty"`
	Offset *int                   `json:"offset,omitempty"`

	//aggregation fields
	Aggregations *QueryResourceAggregations `json:"aggregations"`
}

type QueryResourceAggregations struct {
	CountBy *QueryResourceAggregation `json:"count_by,omitempty"` //alias for both groupby and orderby, cannot be used together

	GroupBy *QueryResourceAggregation `json:"group_by,omitempty"`
	OrderBy *QueryResourceAggregation `json:"order_by,omitempty"`
}

type QueryResourceAggregation struct {
	Field    string `json:"field"`
	Function string `json:"fn"` //built-in SQL aggregation functions (eg. Count, min, max, etc).
}

func (q *QueryResource) Validate() error {
	if len(q.Use) > 0 {
		return fmt.Errorf("'use' is not supported yet")
	}

	if len(q.From) == 0 {
		return fmt.Errorf("'from' is required")
	}

	if q.Aggregations != nil {
		if len(q.Select) > 0 {
			return fmt.Errorf("cannot use 'select' and 'aggregations' together")
		}

		if q.Aggregations.CountBy != nil {
			if err := validateAggregation("count_by", q.Aggregations.CountBy); err != nil {
				return err
			}
		}
		if q.Aggregations.GroupBy != nil {
			if err := validateAggregation("group_by", q.Aggregations.GroupBy); err != nil {
				return err
			}
		}
		if q.Aggregations.OrderBy != nil {
			if err := validateAggregation("order_by", q.Aggregations.OrderBy); err != nil {
				return err
			}
		}

		if q.Aggregations.CountBy != nil {
			if q.Aggregations.GroupBy != nil {
				return fmt.Errorf("cannot use 'count_by' and 'group_by' together")
			}
			if q.Aggregations.OrderBy != nil {
				return fmt.Errorf("cannot use 'count_by' and 'order_by' together")
			}
		}
		if q.Aggregations.CountBy == nil && q.Aggregations.OrderBy == nil && q.Aggregations.GroupBy == nil {
			return fmt.Errorf("aggregations must have at least one of 'count_by', 'group_by', or 'order_by'")
		}

	}

	if q.Limit != nil && *q.Limit < 0 {
		return fmt.Errorf("'limit' must be greater than or equal to zero")
	}
	if q.Offset != nil && *q.Offset < 0 {
		return fmt.Errorf("'offset' must be greater than or equal to zero")
	}

	return nil
}
