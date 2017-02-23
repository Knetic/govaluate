package main

/*
	Courtesy of abrander
	ref: https://gist.github.com/abrander/fa05ae9b181b48ffe7afb12c961b6e90
*/
import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	// this is horrible and the author knows it. Fortunately, it's just for testing!
	".."
)

var (
	hello  = "hello"
	empty  struct{}
	empty2 *string
	empty3 *int

	values = []interface{}{
		-1,
		0,
		12,
		13,
		"",
		"hello",
		&hello,
		nil,
		"nil",
		empty,
		empty2,
		true,
		false,
		time.Now(),
		rune('r'),
		int64(34),
		time.Duration(0),
		"true",
		"false",
		"\ntrue\n",
		"\nfalse\n",
		"12",
		"nil",
		"arg1",
		"arg2",
		int(12),
		int32(12),
		int64(12),
		complex(1.0, 1.0),
		[]byte{0, 0, 0},
		[]int{0, 0, 0},
		[]string{},
		"[]",
		"{}",
		"\"\"",
		"\"12\"",
		"\"hello\"",
		".*",
		"==",
		"!=",
		">",
		">=",
		"<",
		"<=",
		"=~",
		"!~",
		"in",
		"&&",
		"||",
		"^",
		"&",
		"|",
		">>",
		"<<",
		"+",
		"-",
		"*",
		"/",
		"%",
		"**",
		"-",
		"!",
		"~",
		"?",
		":",
		"??",
		"+",
		"-",
		"*",
		"/",
		"%",
		"**",
		"&",
		"|",
		"^",
		">>",
		"<<",
		",",
		"(",
		")",
		"[",
		"]",
		"\n",
		"\000",
	}

	panics = 0
)


const (
	ITERATIONS = 100000000
)

func main() {

	seed := seedRandom();
	fmt.Printf("Beginning torture test, seed: %d, iterations: %d\n", seed, ITERATIONS);

	for i := 0; i < ITERATIONS; i++ {

		num := rand.Intn(3) + 2
		expression := ""

		for n := 0; n < num; n++ {
			expression += fmt.Sprintf(" %s", getRandom(values))
		}

		tryit(expression)
	}

	fmt.Printf("Done. %d/%d panics.\n", panics, ITERATIONS)
	if panics > 0 {
		os.Exit(1)
	}
}

func tryit(expression string) {
	parameters := make(map[string]interface{})
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Panic: \"%s\". Expression: \"%s\". Parameters: %+v\n", r, expression, parameters)
			panics++
		}
	}()

	eval, _ := govaluate.NewEvaluableExpression(expression)
	if eval == nil {
		return
	}

	vars := eval.Vars()
	for _, v := range vars {
		parameters[v] = getRandom(values)
	}

	eval.Evaluate(parameters)
}

func seedRandom() int64 {

	var seed int64

	if len(os.Args) > 1 {
		seed, _ = strconv.ParseInt(os.Args[1], 10, 64)
	} else {
		seed = time.Now().UnixNano()
	}

	rand.Seed(seed)
	return seed
}

func getRandom(haystack []interface{}) interface{} {
	i := rand.Intn(len(haystack))

	return haystack[i]
}
