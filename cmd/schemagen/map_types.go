package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/cassandra-gocql-driver/v2"
)

var types = map[string]string{
	"ascii":     "string",
	"bigint":    "int64",
	"blob":      "[]byte",
	"boolean":   "bool",
	"counter":   "int",
	"date":      "time.Time",
	"decimal":   "inf.Dec",
	"double":    "float64",
	"duration":  "gocql.Duration",
	"float":     "float32",
	"inet":      "string",
	"int":       "int32",
	"smallint":  "int16",
	"text":      "string",
	"time":      "time.Duration",
	"timestamp": "time.Time",
	"timeuuid":  "[16]byte",
	"tinyint":   "int8",
	"uuid":      "[16]byte",
	"varchar":   "string",
	"varint":    "int64",
}

func mapScyllaToGoType(s string) string {
	frozenRegex := regexp.MustCompile(`frozen<([a-z]*)>`)
	match := frozenRegex.FindAllStringSubmatch(s, -1)
	prefix := ""
	if match != nil {
		origin := s
		s = match[0][1]
		mapRegex := regexp.MustCompile(`map<([a-z]*), frozen<([a-z]*)>>`)
		match = mapRegex.FindAllStringSubmatch(origin, -1)
		if match != nil {
			key := match[0][1]
			prefix = fmt.Sprintf("map[%s]", types[key])
		}

		setRegex := regexp.MustCompile(`set<frozen<([a-z]*)>>`)
		match = setRegex.FindAllStringSubmatch(origin, -1)
		if match != nil {
			prefix = "[]"
		}

		listRegex := regexp.MustCompile(`list<frozen<([a-z]*)>>`)
		match = listRegex.FindAllStringSubmatch(origin, -1)
		if match != nil {
			prefix = "[]"
		}
	}

	mapRegex := regexp.MustCompile(`map<([a-z]*), ([a-z]*)>`)
	setRegex := regexp.MustCompile(`set<([a-z]*)>`)
	listRegex := regexp.MustCompile(`list<([a-z]*)>`)
	tupleRegex := regexp.MustCompile(`tuple<(?:([a-z]*),? ?)*>`)
	match = mapRegex.FindAllStringSubmatch(s, -1)
	if match != nil {
		key := match[0][1]
		value := match[0][2]

		return "map[" + types[key] + "]" + types[value]
	}

	match = setRegex.FindAllStringSubmatch(s, -1)
	if match != nil {
		key := match[0][1]

		return "[]" + types[key]
	}

	match = listRegex.FindAllStringSubmatch(s, -1)
	if match != nil {
		key := match[0][1]

		return "[]" + types[key]
	}

	match = tupleRegex.FindAllStringSubmatch(s, -1)
	if match != nil {
		tuple := match[0][0]
		subStr := tuple[6 : len(tuple)-1]
		types := strings.Split(subStr, ", ")

		typeStr := "struct {\n"
		for i, t := range types {
			typeStr = typeStr + "\t\tField" + strconv.Itoa(i+1) + " " + mapScyllaToGoType(t) + "\n"
		}
		typeStr = typeStr + "\t}"

		return typeStr
	}

	t, exists := types[s]
	if exists {
		return t
	}

	return prefix + camelize(s) + "UserType"
}

// cqlPrimitiveTypes maps gocql native type identifiers to CQL type names.
// The names must match the keys of the types map.
var cqlPrimitiveTypes = map[gocql.Type]string{
	gocql.TypeAscii:     "ascii",
	gocql.TypeBigInt:    "bigint",
	gocql.TypeBlob:      "blob",
	gocql.TypeBoolean:   "boolean",
	gocql.TypeCounter:   "counter",
	gocql.TypeDate:      "date",
	gocql.TypeDecimal:   "decimal",
	gocql.TypeDouble:    "double",
	gocql.TypeDuration:  "duration",
	gocql.TypeFloat:     "float",
	gocql.TypeInet:      "inet",
	gocql.TypeInt:       "int",
	gocql.TypeSmallInt:  "smallint",
	gocql.TypeText:      "text",
	gocql.TypeTime:      "time",
	gocql.TypeTimestamp: "timestamp",
	gocql.TypeTimeUUID:  "timeuuid",
	gocql.TypeTinyInt:   "tinyint",
	gocql.TypeUUID:      "uuid",
	gocql.TypeVarchar:   "varchar",
	gocql.TypeVarint:    "varint",
}

// typeToString renders a gocql TypeInfo as a CQL type string, e.g. "int",
// "set<text>" or "map<uuid, text>". Types the driver cannot resolve (e.g. user
// defined types) are reported with their raw validator string, such as
// "frozen<address>", so that mapScyllaToGoType can map them to the generated
// UserType structs.
func typeToString(t gocql.TypeInfo) string {
	switch info := t.(type) {
	case gocql.CollectionType:
		switch info.Type() {
		case gocql.TypeMap:
			return "map<" + typeToString(info.Key) + ", " + typeToString(info.Elem) + ">"
		case gocql.TypeList:
			return "list<" + typeToString(info.Elem) + ">"
		case gocql.TypeSet:
			return "set<" + typeToString(info.Elem) + ">"
		}
	case gocql.TupleTypeInfo:
		elems := make([]string, len(info.Elems))
		for i, elem := range info.Elems {
			elems[i] = typeToString(elem)
		}
		return "tuple<" + strings.Join(elems, ", ") + ">"
	case gocql.UDTTypeInfo:
		return info.Name
	}

	if name, ok := cqlPrimitiveTypes[t.Type()]; ok {
		return name
	}

	// Unknown types (e.g. user defined types) are held by the driver as a
	// string type carrying the raw validator, formatting it exposes that string.
	return fmt.Sprintf("%v", t)
}
