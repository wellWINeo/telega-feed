package repositories

import (
	"context"
	"errors"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"io"
)

func mapResultSet[T any](rows query.Result, ctx context.Context) ([]*T, error) {
	results := make([]*T, 0)

	for {
		resultSet, err := rows.NextResultSet(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		for row, err := range sugar.UnmarshalRows[T](resultSet.Rows(ctx)) {
			if err != nil {
				return nil, err
			}

			result := row
			results = append(results, &result)
		}
	}

	return results, nil
}
