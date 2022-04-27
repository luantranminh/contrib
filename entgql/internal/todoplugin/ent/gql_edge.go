// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func (c *Category) Todos(ctx context.Context) ([]*Todo, error) {
	result, err := c.Edges.TodosOrErr()
	if IsNotLoaded(err) {
		result, err = c.QueryTodos().All(ctx)
	}
	return result, err
}

func (t *Todo) Parent(ctx context.Context) (*Todo, error) {
	result, err := t.Edges.ParentOrErr()
	if IsNotLoaded(err) {
		result, err = t.QueryParent().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (t *Todo) Children(
	ctx context.Context, after *Cursor, first *int, before *Cursor, last *int, orderBy *TodoOrder, where *TodoWhereInput,
	opts ...TodoPaginateOption,
) (*TodoConnection, error) {
	totalCount := t.Edges.totalCount[1]
	if nodes, err := t.Edges.ChildrenOrErr(); err == nil {
		conn := &TodoConnection{Edges: []*TodoEdge{}}
		if totalCount != nil {
			conn.TotalCount = *totalCount
		}
		opts = append(opts, WithTodoOrder(orderBy))
		pager, err := newTodoPager(opts)
		if err != nil {
			return nil, err
		}
		conn.build(nodes, pager, first, last)
		return conn, nil
	}
	query := t.QueryChildren()
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newTodoPager(opts)
	if err != nil {
		return nil, err
	}
	if query, err = pager.applyFilter(query); err != nil {
		return nil, err
	}
	conn := &TodoConnection{Edges: []*TodoEdge{}}
	if !hasCollectedField(ctx, edgesField) || first != nil && *first == 0 || last != nil && *last == 0 {
		if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
			if totalCount != nil {
				conn.TotalCount = *totalCount
			} else if conn.TotalCount, err = query.Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
		return conn, nil
	}

	if (after != nil || first != nil || before != nil || last != nil) && hasCollectedField(ctx, totalCountField) {
		count, err := query.Clone().Count(ctx)
		if err != nil {
			return nil, err
		}
		conn.TotalCount = count
	}

	query = pager.applyCursors(query, after, before)
	query = pager.applyOrder(query, last != nil)
	if limit := paginateLimit(first, last); limit != 0 {
		query.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := query.collectField(ctx, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}

	nodes, err := query.All(ctx)
	if err != nil || len(nodes) == 0 {
		return conn, err
	}
	conn.build(nodes, pager, first, last)
	return conn, nil
}
