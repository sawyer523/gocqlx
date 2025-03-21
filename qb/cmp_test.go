// Copyright (C) 2017 ScyllaDB
// Use of this source code is governed by a ALv2-style
// license that can be found in the LICENSE file.

package qb

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCmp(t *testing.T) {
	table := []struct {
		C Cmp
		S string
		N []string
	}{
		// Basic comparators
		{
			C: Eq("eq"),
			S: "eq=?",
			N: []string{"eq"},
		},
		{
			C: Eq("token"),
			S: "\"token\"=?",
			N: []string{"token"},
		},
		{
			C: Eq("Eq"),
			S: "\"Eq\"=?",
			N: []string{"eq"},
		},
		{
			C: Eq("\"Eq\""),
			S: "\"Eq\"=?",
			N: []string{"eq"},
		},
		{
			C: Eq("\"EQ\""),
			S: "\"EQ\"=?",
			N: []string{"eq"},
		},
		{
			C: Ne("ne"),
			S: "ne!=?",
			N: []string{"ne"},
		},
		{
			C: Ne("\"Ne\""),
			S: "\"Ne\"!=?",
			N: []string{"ne"},
		},
		{
			C: Ne("\"NE\""),
			S: "\"NE\"!=?",
			N: []string{"ne"},
		},
		{
			C: NeTuple("ne", 3),
			S: "ne!=(?,?,?)",
			N: []string{"ne[0]", "ne[1]", "ne[2]"},
		},
		{
			C: NeTuple("\"NE\"", 3),
			S: "\"NE\"!=(?,?,?)",
			N: []string{"ne[0]", "ne[1]", "ne[2]"},
		},
		{
			C: Lt("lt"),
			S: "lt<?",
			N: []string{"lt"},
		},
		{
			C: Lt("\"Lt\""),
			S: "\"Lt\"<?",
			N: []string{"lt"},
		}, {
			C: Lt("\"LT\""),
			S: "\"LT\"<?",
			N: []string{"lt"},
		},
		{
			C: LtTuple("lt", 2),
			S: "lt<(?,?)",
			N: []string{"lt[0]", "lt[1]"},
		},
		{
			C: LtTuple("\"Lt\"", 2),
			S: "\"Lt\"<(?,?)",
			N: []string{"lt[0]", "lt[1]"},
		},
		{
			C: LtTuple("\"LT\"", 2),
			S: "\"LT\"<(?,?)",
			N: []string{"lt[0]", "lt[1]"},
		},
		{
			C: LtOrEq("lt"),
			S: "lt<=?",
			N: []string{"lt"},
		},
		{
			C: LtOrEq("\"Lt\""),
			S: "\"Lt\"<=?",
			N: []string{"lt"},
		},
		{
			C: LtOrEq("\"LT\""),
			S: "\"LT\"<=?",
			N: []string{"lt"},
		},
		{
			C: LtOrEqTuple("lt", 2),
			S: "lt<=(?,?)",
			N: []string{"lt[0]", "lt[1]"},
		},
		{
			C: LtOrEqTuple("\"Lt\"", 2),
			S: "\"Lt\"<=(?,?)",
			N: []string{"lt[0]", "lt[1]"},
		},
		{
			C: LtOrEqTuple("\"LT\"", 2),
			S: "\"LT\"<=(?,?)",
			N: []string{"lt[0]", "lt[1]"},
		},
		{
			C: Gt("gt"),
			S: "gt>?",
			N: []string{"gt"},
		},
		{
			C: Gt("\"Gt\""),
			S: "\"Gt\">?",
			N: []string{"gt"},
		},
		{
			C: Gt("\"GT\""),
			S: "\"GT\">?",
			N: []string{"gt"},
		},
		{
			C: GtTuple("gt", 2),
			S: "gt>(?,?)",
			N: []string{"gt[0]", "gt[1]"},
		},
		{
			C: GtTuple("\"Gt\"", 2),
			S: "\"Gt\">(?,?)",
			N: []string{"gt[0]", "gt[1]"},
		},
		{
			C: GtTuple("\"GT\"", 2),
			S: "\"GT\">(?,?)",
			N: []string{"gt[0]", "gt[1]"},
		},
		{
			C: GtOrEq("gt"),
			S: "gt>=?",
			N: []string{"gt"},
		},
		{
			C: GtOrEq("\"Gt\""),
			S: "\"Gt\">=?",
			N: []string{"gt"},
		},
		{
			C: GtOrEq("\"GT\""),
			S: "\"GT\">=?",
			N: []string{"gt"},
		},
		{
			C: GtOrEqTuple("gt", 2),
			S: "gt>=(?,?)",
			N: []string{"gt[0]", "gt[1]"},
		},
		{
			C: GtOrEqTuple("\"Gt\"", 2),
			S: "\"Gt\">=(?,?)",
			N: []string{"gt[0]", "gt[1]"},
		},
		{
			C: GtOrEqTuple("\"GT\"", 2),
			S: "\"GT\">=(?,?)",
			N: []string{"gt[0]", "gt[1]"},
		},
		{
			C: In("in"),
			S: "\"in\" IN ?",
			N: []string{"in"},
		},
		{
			C: In("\"In\""),
			S: "\"In\" IN ?",
			N: []string{"in"},
		},
		{
			C: In("\"IN\""),
			S: "\"IN\" IN ?",
			N: []string{"in"},
		},
		{
			C: InTuple("in", 2),
			S: "\"in\" IN (?,?)",
			N: []string{"in[0]", "in[1]"},
		},
		{
			C: InTuple("\"In\"", 2),
			S: "\"In\" IN (?,?)",
			N: []string{"in[0]", "in[1]"},
		},
		{
			C: InTuple("\"IN\"", 2),
			S: "\"IN\" IN (?,?)",
			N: []string{"in[0]", "in[1]"},
		},
		{
			C: Contains("cnt"),
			S: "cnt CONTAINS ?",
			N: []string{"cnt"},
		},
		{
			C: Contains("\"Cnt\""),
			S: "\"Cnt\" CONTAINS ?",
			N: []string{"cnt"},
		},
		{
			C: Contains("\"CNT\""),
			S: "\"CNT\" CONTAINS ?",
			N: []string{"cnt"},
		},
		{
			C: ContainsTuple("cnt", 2),
			S: "cnt CONTAINS (?,?)",
			N: []string{"cnt[0]", "cnt[1]"},
		},
		{
			C: ContainsTuple("\"Cnt\"", 2),
			S: "\"Cnt\" CONTAINS (?,?)",
			N: []string{"cnt[0]", "cnt[1]"},
		},
		{
			C: ContainsTuple("\"CNT\"", 2),
			S: "\"CNT\" CONTAINS (?,?)",
			N: []string{"cnt[0]", "cnt[1]"},
		},
		{
			C: ContainsKey("cntKey"),
			S: "\"cntKey\" CONTAINS KEY ?",
			N: []string{"cnt_key"},
		},
		{
			C: ContainsKey("\"CntKey\""),
			S: "\"CntKey\" CONTAINS KEY ?",
			N: []string{"cnt_key"},
		},
		{
			C: ContainsKeyTuple("cntKey", 2),
			S: "\"cntKey\" CONTAINS KEY (?,?)",
			N: []string{"cnt_key[0]", "cnt_key[1]"},
		},
		{
			C: ContainsKeyTuple("\"CntKey\"", 2),
			S: "\"CntKey\" CONTAINS KEY (?,?)",
			N: []string{"cnt_key[0]", "cnt_key[1]"},
		},
		{
			C: Like("like"),
			S: "like LIKE ?",
			N: []string{"like"},
		},
		{
			C: Like("\"Like\""),
			S: "\"Like\" LIKE ?",
			N: []string{"like"},
		},
		{
			C: Like("\"LIKE\""),
			S: "\"LIKE\" LIKE ?",
			N: []string{"like"},
		},
		{
			C: LikeTuple("like", 2),
			S: "like LIKE (?,?)",
			N: []string{"like[0]", "like[1]"},
		},
		{
			C: LikeTuple("\"Like\"", 2),
			S: "\"Like\" LIKE (?,?)",
			N: []string{"like[0]", "like[1]"},
		},
		{
			C: LikeTuple("\"LIKE\"", 2),
			S: "\"LIKE\" LIKE (?,?)",
			N: []string{"like[0]", "like[1]"},
		},

		// Custom bind names
		{
			C: EqNamed("eq", "name"),
			S: "eq=?",
			N: []string{"name"},
		},
		{
			C: NeNamed("ne", "name"),
			S: "ne!=?",
			N: []string{"name"},
		},
		{
			C: LtNamed("lt", "name"),
			S: "lt<?",
			N: []string{"name"},
		},
		{
			C: LtOrEqNamed("lt", "name"),
			S: "lt<=?",
			N: []string{"name"},
		},
		{
			C: GtNamed("gt", "name"),
			S: "gt>?",
			N: []string{"name"},
		},
		{
			C: GtOrEqNamed("gt", "name"),
			S: "gt>=?",
			N: []string{"name"},
		},
		{
			C: InNamed("in", "name"),
			S: "\"in\" IN ?",
			N: []string{"name"},
		},
		{
			C: ContainsNamed("cnt", "name"),
			S: "cnt CONTAINS ?",
			N: []string{"name"},
		},
		{
			C: ContainsKeyNamed("cntKey", "name"),
			S: "\"cntKey\" CONTAINS KEY ?",
			N: []string{"name"},
		},
		{
			C: LikeTupleNamed("like", 2, "name"),
			S: "like LIKE (?,?)",
			N: []string{"name[0]", "name[1]"},
		},
		// Custom bind names on tuples
		{
			C: EqTupleNamed("eq", 2, "name"),
			S: "eq=(?,?)",
			N: []string{"name[0]", "name[1]"},
		},
		{
			C: NeTupleNamed("ne", 3, "name"),
			S: "ne!=(?,?,?)",
			N: []string{"name[0]", "name[1]", "name[2]"},
		},
		{
			C: LtTupleNamed("lt", 2, "name"),
			S: "lt<(?,?)",
			N: []string{"name[0]", "name[1]"},
		},
		{
			C: LtOrEqTupleNamed("lt", 2, "name"),
			S: "lt<=(?,?)",
			N: []string{"name[0]", "name[1]"},
		},
		{
			C: GtTupleNamed("gt", 2, "name"),
			S: "gt>(?,?)",
			N: []string{"name[0]", "name[1]"},
		},
		{
			C: GtOrEqTupleNamed("gt", 2, "name"),
			S: "gt>=(?,?)",
			N: []string{"name[0]", "name[1]"},
		},
		{
			C: InTupleNamed("in", 2, "name"),
			S: "\"in\" IN (?,?)",
			N: []string{"name[0]", "name[1]"},
		},
		{
			C: ContainsTupleNamed("cnt", 2, "name"),
			S: "cnt CONTAINS (?,?)",
			N: []string{"name[0]", "name[1]"},
		},
		{
			C: ContainsKeyTupleNamed("cntKey", 2, "name"),
			S: "\"cntKey\" CONTAINS KEY (?,?)",
			N: []string{"name[0]", "name[1]"},
		},

		// Literals
		{
			C: EqLit("eq", "litval"),
			S: "eq=litval",
		},
		{
			C: NeLit("ne", "litval"),
			S: "ne!=litval",
		},
		{
			C: LtLit("lt", "litval"),
			S: "lt<litval",
		},
		{
			C: LtOrEqLit("lt", "litval"),
			S: "lt<=litval",
		},
		{
			C: GtLit("gt", "litval"),
			S: "gt>litval",
		},
		{
			C: GtOrEqLit("gt", "litval"),
			S: "gt>=litval",
		},
		{
			C: InLit("in", "litval"),
			S: "\"in\" IN litval",
		},
		{
			C: ContainsLit("cnt", "litval"),
			S: "cnt CONTAINS litval",
		},

		// Functions
		{
			C: EqFunc("eq", Fn("fn", "arg0", "arg1")),
			S: "eq=fn(?,?)",
			N: []string{"arg0", "arg1"},
		},
		{
			C: EqFunc("eq", MaxTimeuuid("arg0")),
			S: "eq=maxTimeuuid(?)",
			N: []string{"arg0"},
		},
		{
			C: EqFunc("eq", MinTimeuuid("arg0")),
			S: "eq=minTimeuuid(?)",
			N: []string{"arg0"},
		},
		{
			C: EqFunc("eq", Now()),
			S: "eq=now()",
		},
		{
			C: NeFunc("ne", Fn("fn", "arg0", "arg1", "arg2")),
			S: "ne!=fn(?,?,?)",
			N: []string{"arg0", "arg1", "arg2"},
		},
		{
			C: LtFunc("eq", Now()),
			S: "eq<now()",
		},
		{
			C: LtOrEqFunc("eq", MaxTimeuuid("arg0")),
			S: "eq<=maxTimeuuid(?)",
			N: []string{"arg0"},
		},
		{
			C: GtFunc("eq", Now()),
			S: "eq>now()",
		},
		{
			C: GtOrEqFunc("eq", MaxTimeuuid("arg0")),
			S: "eq>=maxTimeuuid(?)",
			N: []string{"arg0"},
		},
	}

	buf := bytes.Buffer{}
	for _, test := range table {
		buf.Reset()
		name := test.C.writeCql(&buf)
		if diff := cmp.Diff(test.S, buf.String()); diff != "" {
			t.Error(diff)
		}
		if diff := cmp.Diff(test.N, name); diff != "" {
			t.Error(diff)
		}
	}
}
