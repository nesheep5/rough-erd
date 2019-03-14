package rough_er

import "sort"

type Schema struct {
	Name   string
	Tables []*Table
	Memo   string
}

type Table struct {
	Name       string
	Type       string
	Comment    string
	Columns    []*Column
	Memo       string
	Definition string
}

type Column struct {
	Position uint
	Name     string
	Type     string
	Default  string
	Nullable string
	Comment  string
	Memo     string
}

func UpdateSchemas(old []*Schema, new []*Schema) []*Schema {
	var ss []*Schema
	for _, o := range old {
		for _, n := range new {
			if o.Name == n.Name {
				// merge
				o.Tables = updateTables(o.Tables, n.Tables)
				ss = append(ss, o)
			}
		}
	}
	for _, n := range new {
		exists := false
		for _, s := range ss {
			if n.Name == s.Name {
				exists = true
			}
		}
		if !exists {
			ss = append(ss, n)
		}
	}

	sort.SliceStable(ss, func(i, j int) bool { return ss[i].Name < ss[j].Name })
	return ss
}

func updateTables(old []*Table, new []*Table) []*Table {
	var ts []*Table
	for _, o := range old {
		for _, n := range new {
			if o.Name == n.Name {
				// merge
				o.Type = n.Type
				o.Comment = n.Comment
				o.Definition = n.Definition
				o.Columns = updateColumns(o.Columns, n.Columns)
				ts = append(ts, o)
			}
		}
	}
	for _, n := range new {
		exists := false
		for _, s := range ts {
			if n.Name == s.Name {
				exists = true
			}
		}
		if !exists {
			ts = append(ts, n)
		}
	}
	sort.SliceStable(ts, func(i, j int) bool { return ts[i].Name < ts[j].Name })
	return ts
}

func updateColumns(old []*Column, new []*Column) []*Column {
	var cs []*Column
	for _, o := range old {
		for _, n := range new {
			if o.Name == n.Name {
				// merge
				o.Position = n.Position
				o.Type = n.Type
				o.Default = n.Default
				o.Nullable = n.Nullable
				o.Comment = n.Comment
				cs = append(cs, o)
			}
		}
	}
	for _, n := range new {
		exists := false
		for _, s := range cs {
			if n.Name == s.Name {
				exists = true
			}
		}
		if !exists {
			cs = append(cs, n)
		}
	}
	sort.SliceStable(cs, func(i, j int) bool { return cs[i].Position < cs[j].Position })
	return cs
}
