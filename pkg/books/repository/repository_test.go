package repository_test

import "testing"

func TestFetchAll(t *testing.T) {
	testCases := []struct {
		name   string
		limit  int
		offset int
	}{
		{
			name:   "should return books with given offset and limit",
			limit:  5,
			offset: 1,
		},
		{
			name:   "should return books with given offset and limit",
			limit:  10,
			offset: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//	todo start in memory db
			//	todo inject db to repository
			//	todo fetch based on limit and offset
		})
	}
}
