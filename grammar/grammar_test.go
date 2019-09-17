package grammar

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const yyWithOutKeyword = `
query:
    select

select:
    SELECT
        fields
    FROM (
        select
    ) as t
    WHERE t.val > 10
    | SELECT *
	FROM table1
	ORDER BY LOWER(fieldA), LOWER(fieldB)

fields:
    fieldA
    | fieldB
    | fieldA, fieldB
`

const yy = `
query:
    {if(a==nil) then a = 1 end} select

select:
    SELECT
           fieldA,
           fieldB,
           {print(string.format("field%d", a)); a = a + 1}
    FROM (
	SELECT _field AS fieldA, _field AS fieldB
	FROM _table
	ORDER BY LOWER(fieldA), LOWER(fieldB)
    ) as t
    WINDOW w AS (ORDER BY fieldA);
`

const luaYy = `
{
f={a=1, b=3}
arr={0,2,3,4}
}

query:
  {print(arr[f.a])} | {print(arr[f.b])}
`

func TestByYyWithoutKeyword(t *testing.T) {
	t.SkipNow()
	num := 10
	sqls, err := ByYy(yyWithOutKeyword, num, "query", 5, nil, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, num, len(sqls))

	for _, sql := range sqls {
		fmt.Println(sql)
	}
}


func TestByYy(t *testing.T) {
	t.SkipNow()
	sqls, err := ByYy(yy, 10, "query", 5, map[string]func() string{
		"_table": func() string {
			return "aaa_tabl"
		},
		"_field": func() string {
			return "ffff"
		},
	}, false)
	assert.Equal(t, nil, err)

	for _, sql := range sqls {
		fmt.Println(sql)
	}
}

func TestLuaYy(t *testing.T) {
	t.SkipNow()
	sqls, err := ByYy(luaYy, 10, "query", 5,nil, false)
	assert.Equal(t, nil, err)

	for _, sql := range sqls {
		fmt.Println(sql)
	}
}

const testFirstWriteYy = `
query:
	{n=1} select

select:
    sub_select | haha_select

sub_select:
    {n = 2} SELECT

haha_select:
    {m = 4} SELECT

`

func TestFirstWriteYy(t *testing.T) {
	sqls, err := ByYy(testFirstWriteYy, 50, "query", 5,nil, false)
	assert.Equal(t, nil, err)

	for _, sql := range sqls {
		assert.Equal(t, "SELECT", sql)
	}
}

const testSemiColon = `
query:
	select ; create

select:
    SET @stmt = "FFF";
	PREPARE stmt FROM @stmt_create ; EXECUTE stmt ;
	EXECUTE stmt;

create:
	CREATE OOO; CCC
`

func TestSemiColon(t *testing.T) {
	sqls, err := ByYy(testSemiColon, 6, "query", 5, nil, false)
	assert.Equal(t, nil, err)

	expected := []string{
		`SET @stmt = "FFF"`,
		`PREPARE stmt FROM @stmt_create`,
		`EXECUTE stmt`,
		`EXECUTE stmt`,
		`CREATE OOO`,
		`CCC`,
	}

	assert.Equal(t, len(expected), len(sqls))

	for i, sql := range sqls {
		//fmt.Println(sql)
		assert.Equal(t, expected[i], sql)
	}
}



