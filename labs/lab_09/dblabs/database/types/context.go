package types

import (
	"fmt"
	"github.com/Inspirate789/go-randomdata"
	"math/rand"
	"strings"
)

type Context struct {
	ID, TermsCount, WordsCount int
	RegDate, Text              string
}

func RandContext(terms []Term) Context {
	var context Context

	context.RegDate = randomdata.FullDateInRange("2022-09-01", "2022-09-17")
	context.TermsCount = rand.Intn(10) + 1

	builder := strings.Builder{}
	for j := 0; j < context.TermsCount; j++ {
		curTerm := terms[rand.Intn(len(terms))]
		builder.WriteString(curTerm.Text)
		builder.WriteRune(' ')
		context.WordsCount += curTerm.WordsCount
	}
	context.Text = builder.String()

	return context
}

func ContextToSlice(context Context) []string {
	return []string{
		fmt.Sprint(context.TermsCount),
		fmt.Sprint(context.WordsCount),
		context.RegDate,
		context.Text,
	}
}
