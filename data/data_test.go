package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalObj(t *testing.T) {
	expected := map[string]interface{}{
		"foo":  map[interface{}]interface{}{"bar": "baz"},
		"one":  1.0,
		"true": true,
	}

	test := func(actual map[string]interface{}) {
		assert.Equal(t, expected["foo"], actual["foo"])
		assert.Equal(t, expected["one"], actual["one"])
		assert.Equal(t, expected["true"], actual["true"])
	}
	test(JSON(`{"foo":{"bar":"baz"},"one":1.0,"true":true}`))
	test(YAML(`foo:
  bar: baz
one: 1.0
true: true
`))
}

func TestUnmarshalArray(t *testing.T) {

	expected := []string{"foo", "bar"}

	test := func(actual []interface{}) {
		assert.Equal(t, expected[0], actual[0])
		assert.Equal(t, expected[1], actual[1])
	}
	test(JSONArray(`["foo","bar"]`))
	test(YAMLArray(`
- foo
- bar
`))
}

func TestToJSON(t *testing.T) {
	expected := `{"down":{"the":{"rabbit":{"hole":true}}},"foo":"bar","one":1,"true":true}`
	in := map[string]interface{}{
		"foo":  "bar",
		"one":  1,
		"true": true,
		"down": map[interface{}]interface{}{
			"the": map[interface{}]interface{}{
				"rabbit": map[interface{}]interface{}{
					"hole": true,
				},
			},
		},
	}
	assert.Equal(t, expected, ToJSON(in))
}

func TestToJSONPretty(t *testing.T) {
	expected := `{
  "down": {
    "the": {
      "rabbit": {
        "hole": true
      }
    }
  },
  "foo": "bar",
  "one": 1,
  "true": true
}`
	in := map[string]interface{}{
		"foo":  "bar",
		"one":  1,
		"true": true,
		"down": map[string]interface{}{
			"the": map[string]interface{}{
				"rabbit": map[string]interface{}{
					"hole": true,
				},
			},
		},
	}
	assert.Equal(t, expected, ToJSONPretty("  ", in))
}

func TestToYAML(t *testing.T) {
	expected := `d: !!timestamp 2006-01-02T15:04:05.999999999-07:00
foo: bar
? |-
  multi
  line
  key
: hello: world
one: 1
"true": true
`
	mst, _ := time.LoadLocation("MST")
	in := map[string]interface{}{
		"foo":  "bar",
		"one":  1,
		"true": true,
		`multi
line
key`: map[string]interface{}{
			"hello": "world",
		},
		"d": time.Date(2006, time.January, 2, 15, 4, 5, 999999999, mst),
	}
	assert.Equal(t, expected, ToYAML(in))
}

func TestCSV(t *testing.T) {
	in := "first,second,third\n1,2,3\n4,5,6"
	expected := [][]string{
		{"first", "second", "third"},
		{"1", "2", "3"},
		{"4", "5", "6"},
	}
	assert.Equal(t, expected, CSV(in))

	in = "first;second;third\r\n1;2;3\r\n4;5;6\r\n"
	assert.Equal(t, expected, CSV(";", in))

	in = ""
	assert.Equal(t, [][]string{nil}, CSV(in))

	in = "\n"
	assert.Equal(t, [][]string{nil}, CSV(in))

	in = "foo"
	assert.Equal(t, [][]string{[]string{"foo"}}, CSV(in))
}

func TestCSVByRow(t *testing.T) {
	in := "first,second,third\n1,2,3\n4,5,6"
	expected := []map[string]string{
		{
			"first":  "1",
			"second": "2",
			"third":  "3",
		},
		{
			"first":  "4",
			"second": "5",
			"third":  "6",
		},
	}
	assert.Equal(t, expected, CSVByRow(in))

	in = "1,2,3\n4,5,6"
	assert.Equal(t, expected, CSVByRow("first,second,third", in))

	in = "1;2;3\n4;5;6"
	assert.Equal(t, expected, CSVByRow(";", "first;second;third", in))

	in = "first;second;third\r\n1;2;3\r\n4;5;6"
	assert.Equal(t, expected, CSVByRow(";", in))

	expected = []map[string]string{
		{"A": "1", "B": "2", "C": "3"},
		{"A": "4", "B": "5", "C": "6"},
	}

	in = "1,2,3\n4,5,6"
	assert.Equal(t, expected, CSVByRow("", in))

	expected = []map[string]string{
		{"A": "1", "B": "1", "C": "1", "D": "1", "E": "1", "F": "1", "G": "1", "H": "1", "I": "1", "J": "1", "K": "1", "L": "1", "M": "1", "N": "1", "O": "1", "P": "1", "Q": "1", "R": "1", "S": "1", "T": "1", "U": "1", "V": "1", "W": "1", "X": "1", "Y": "1", "Z": "1", "AA": "1", "BB": "1", "CC": "1", "DD": "1"},
	}

	in = "1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1"
	assert.Equal(t, expected, CSVByRow("", in))
}

func TestCSVByColumn(t *testing.T) {
	in := "first,second,third\n1,2,3\n4,5,6"
	expected := map[string][]string{
		"first":  {"1", "4"},
		"second": {"2", "5"},
		"third":  {"3", "6"},
	}
	assert.Equal(t, expected, CSVByColumn(in))

	in = "1,2,3\n4,5,6"
	assert.Equal(t, expected, CSVByColumn("first,second,third", in))

	in = "1;2;3\n4;5;6"
	assert.Equal(t, expected, CSVByColumn(";", "first;second;third", in))

	in = "first;second;third\r\n1;2;3\r\n4;5;6"
	assert.Equal(t, expected, CSVByColumn(";", in))

	expected = map[string][]string{
		"A": {"1", "4"},
		"B": {"2", "5"},
		"C": {"3", "6"},
	}

	in = "1,2,3\n4,5,6"
	assert.Equal(t, expected, CSVByColumn("", in))
}

func TestAutoIndex(t *testing.T) {
	assert.Equal(t, "A", autoIndex(0))
	assert.Equal(t, "B", autoIndex(1))
	assert.Equal(t, "Z", autoIndex(25))
	assert.Equal(t, "AA", autoIndex(26))
	assert.Equal(t, "ZZ", autoIndex(51))
	assert.Equal(t, "AAA", autoIndex(52))
	assert.Equal(t, "YYYYY", autoIndex(128))
}

func TestToCSV(t *testing.T) {
	in := [][]string{
		{"first", "second", "third"},
		{"1", "2", "3"},
		{"4", "5", "6"},
	}
	expected := "first,second,third\r\n1,2,3\r\n4,5,6\r\n"

	assert.Equal(t, expected, ToCSV(in))

	expected = "first;second;third\r\n1;2;3\r\n4;5;6\r\n"

	assert.Equal(t, expected, ToCSV(";", in))
}

func TestTOML(t *testing.T) {
	in := `# This is a TOML document. Boom.

title = "TOML Example"

[owner]
name = "Tom Preston-Werner"
organization = "GitHub"
bio = "GitHub Cofounder & CEO\nLikes tater tots and beer."
dob = 1979-05-27T07:32:00Z # First class dates? Why not?

[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 5000
enabled = true

[servers]

  # You can indent as you please. Tabs or spaces. TOML don't care.
  [servers.alpha]
  ip = "10.0.0.1"
  dc = "eqdc10"

  [servers.beta]
  ip = "10.0.0.2"
  dc = "eqdc10"

[clients]
data = [ ["gamma", "delta"], [1, 2] ] # just an update to make sure parsers support it

# Line breaks are OK when inside arrays
hosts = [
  "alpha",
  "omega"
]
`
	expected := map[string]interface{}{
		"title": "TOML Example",
		"owner": map[string]interface{}{
			"name":         "Tom Preston-Werner",
			"organization": "GitHub",
			"bio":          "GitHub Cofounder & CEO\nLikes tater tots and beer.",
			"dob":          time.Date(1979, time.May, 27, 7, 32, 0, 0, time.UTC),
		},
		"database": map[string]interface{}{
			"server":         "192.168.1.1",
			"ports":          []interface{}{int64(8001), int64(8001), int64(8002)},
			"connection_max": int64(5000),
			"enabled":        true,
		},
		"servers": map[string]interface{}{
			"alpha": map[string]interface{}{
				"ip": "10.0.0.1",
				"dc": "eqdc10",
			},
			"beta": map[string]interface{}{
				"ip": "10.0.0.2",
				"dc": "eqdc10",
			},
		},
		"clients": map[string]interface{}{
			"data": []interface{}{
				[]interface{}{"gamma", "delta"},
				[]interface{}{int64(1), int64(2)},
			},
			"hosts": []interface{}{"alpha", "omega"},
		},
	}

	assert.Equal(t, expected, TOML(in))
}

func TestToTOML(t *testing.T) {
	expected := `foo = "bar"
one = 1
true = true

[down]
  [down.the]
    [down.the.rabbit]
      hole = true
`
	in := map[string]interface{}{
		"foo":  "bar",
		"one":  1,
		"true": true,
		"down": map[interface{}]interface{}{
			"the": map[interface{}]interface{}{
				"rabbit": map[interface{}]interface{}{
					"hole": true,
				},
			},
		},
	}
	assert.Equal(t, expected, ToTOML(in))
}
