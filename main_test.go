package main

import (
	"strconv"
	"testing"

	"gorm.io/gorm/clause"
)

// GORM_REPO: https://github.com/go-gorm/gorm.git
// GORM_BRANCH: master
// TEST_DRIVERS: sqlite, mysql, postgres, sqlserver

func TestGORM(t *testing.T) {
	const N = 10
	var foos []*Foo
	for i := 0; i < N; i++ {
		isrt := strconv.Itoa(i)
		foos = append(foos, NewFoo(isrt))
	}

	if err := DB.Create(foos).Error; err != nil {
		t.Errorf("Create %v", err)
	}

	for i, foo := range foos {
		isrt := strconv.Itoa(i)
		foo.Name = "updated-" + isrt
	}

	copyFoos := make(map[string]*Foo)
	for _, foo := range foos {
		copyFoos[foo.ID] = &Foo{
			ID:   foo.ID,
			Name: foo.Name,
		}
	}

	res := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(foos)

	if res.Error != nil {
		t.Errorf("create on conflict update %v", res.Error)
	}
	if res.RowsAffected != int64(len(foos)) {
		t.Errorf("expect %v updated %v", len(foos), res.RowsAffected)
	}

	var fail bool
	for _, act := range foos {
		exp, ok := copyFoos[act.ID]
		if !ok {
			t.Errorf("unexpected ID %v", act.ID)
		}
		if act.Name != exp.Name {
			fail = true
			t.Logf("invalid name for id %v expect %v got %v", act.ID, exp.Name, act.Name)
		}
	}
	if fail {
		t.Error("incorrect foos after upsert")
	}
}
