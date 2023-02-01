package handlers

import (
	"dblabs/database"
	"dblabs/database/types"
	"errors"
	"fmt"
	"github.com/vdobler/chart"
	"image/color"
	"strconv"
	"strings"
	"time"
)

type handler90 struct{}

func (h handler90) Title() string {
	return "Select parts of speech from Redis (cached once)"
}

func (h handler90) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	res, err := db.SelectPosRedisCachedOnce()
	if err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", errors.New("empty result")
	}

	builder := strings.Builder{}
	for _, pos := range res {
		builder.WriteString(pos.Name + "\n")
	}

	return builder.String(), nil
}

type handler91 struct{}

func (h handler91) Title() string {
	return "Select parts of speech from Postgres (cache time: 5s)"
}

func (h handler91) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	res, err := db.SelectPosPostgresCached5s()
	if err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", errors.New("empty result")
	}

	builder := strings.Builder{}
	for _, pos := range res {
		builder.WriteString(pos.Name + "\n")
	}

	return builder.String(), nil
}

type handler92 struct{}

func (h handler92) Title() string {
	return "Select parts of speech from Redis (cache time: 5s)"
}

func (h handler92) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 0 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 0, len(args)))
	}

	res, err := db.SelectPosRedisCached5s()
	if err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", errors.New("empty result")
	}

	builder := strings.Builder{}
	for _, pos := range res {
		builder.WriteString(pos.Name + "\n")
	}

	return builder.String(), nil
}

type handler93 struct{}

func (h handler93) Title() string {
	return "Time comparison with given iterations count (Postgres vs Redis)"
}

func (h handler93) Execute(db *database.Database, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(fmt.Sprintf("incorrect arguments count (expected %d, got %d)", 1, len(args)))
	}

	itersCount, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return "", errors.New("incorrect args")
	}

	var i, start, end, totalTime int64
	var postgresTimes, redisTimes []float64
	builder := strings.Builder{}

	for i = 0; i < itersCount; i++ {
		start = time.Now().UnixNano()
		err = db.InsertPosPostgres(types.RandPartOfSpeech(int(i)))
		end = time.Now().UnixNano()
		totalTime += end - start
	}
	builder.WriteString("Postgres insert: " + strconv.FormatInt(totalTime/itersCount, 10) + " ns/op\n")
	postgresTimes = append(postgresTimes, float64(totalTime/itersCount))

	totalTime = 0
	for i = 0; i < itersCount; i++ {
		start = time.Now().UnixNano()
		_, err = db.SelectPos(int(i))
		end = time.Now().UnixNano()
		totalTime += end - start
	}
	builder.WriteString("Postgres select: " + strconv.FormatInt(totalTime/itersCount, 10) + " ns/op\n")
	postgresTimes = append(postgresTimes, float64(totalTime/itersCount))

	totalTime = 0
	for i = 0; i < itersCount; i++ {
		start = time.Now().UnixNano()
		err = db.UpdatePosPostgres(int(i))
		end = time.Now().UnixNano()
		totalTime += end - start
	}
	builder.WriteString("Postgres update: " + strconv.FormatInt(totalTime/itersCount, 10) + " ns/op\n")
	postgresTimes = append(postgresTimes, float64(totalTime/itersCount))

	totalTime = 0
	for i = 0; i < itersCount; i++ {
		start = time.Now().UnixNano()
		err = db.DeletePosPostgres(int(i))
		end = time.Now().UnixNano()
		totalTime += end - start
	}
	builder.WriteString("Postgres delete: " + strconv.FormatInt(totalTime/itersCount, 10) + " ns/op\n")
	postgresTimes = append(postgresTimes, float64(totalTime/itersCount))

	totalTime = 0
	for i = 0; i < itersCount; i++ {
		start = time.Now().UnixNano()
		err = db.InsertPosRedis(strconv.Itoa(int(i)), types.RandPartOfSpeech(int(i)))
		end = time.Now().UnixNano()
		totalTime += end - start
	}
	builder.WriteString("Redis insert:    " + strconv.FormatInt(totalTime/itersCount, 10) + " ns/op\n")
	redisTimes = append(redisTimes, float64(totalTime/itersCount))

	totalTime = 0
	for i = 0; i < itersCount; i++ {
		start = time.Now().UnixNano()
		_, err = db.SelectPosFromRedis(strconv.Itoa(int(i)))
		end = time.Now().UnixNano()
		totalTime += end - start
	}
	builder.WriteString("Redis select:    " + strconv.FormatInt(totalTime/itersCount, 10) + " ns/op\n")
	redisTimes = append(redisTimes, float64(totalTime/itersCount))

	totalTime = 0
	for i = 0; i < itersCount; i++ {
		start = time.Now().UnixNano()
		err = db.UpdatePosRedis(strconv.Itoa(int(i)))
		end = time.Now().UnixNano()
		totalTime += end - start
	}
	builder.WriteString("Redis update:    " + strconv.FormatInt(totalTime/itersCount, 10) + " ns/op\n")
	redisTimes = append(redisTimes, float64(totalTime/itersCount))

	totalTime = 0
	for i = 0; i < itersCount; i++ {
		start = time.Now().UnixNano()
		err = db.DeletePosRedis(strconv.Itoa(int(i)))
		end = time.Now().UnixNano()
		totalTime += end - start
	}
	builder.WriteString("Redis delete:    " + strconv.FormatInt(totalTime/itersCount, 10) + " ns/op\n")
	redisTimes = append(redisTimes, float64(totalTime/itersCount))

	// Bar chart
	plot := chart.BarChart{Title: "Postgres vs Redis"}
	plot.XRange.Category = []string{"Insert", "Select", "Update", "Delete"}
	plot.XRange.Label, plot.YRange.Label = "Query type", "Time, ns"
	plot.Key.Pos, plot.Key.Cols, plot.Key.Border = "otc", 2, -1
	plot.YRange.ShowZero = true
	plot.ShowVal = 0
	plot.AddDataPair("Postgres", []float64{0, 1, 2, 3}, postgresTimes,
		chart.Style{Symbol: '#', LineColor: color.NRGBA{R: 0x30, G: 0x30, B: 0xff, A: 0xff}, LineWidth: 2, FillColor: color.NRGBA{R: 0xcb, G: 0xcb, B: 0xff, A: 0xff}})
	plot.AddDataPair("Redis", []float64{0, 1, 2, 3}, redisTimes,
		chart.Style{Symbol: 'O', LineColor: color.NRGBA{R: 0xe0, G: 0x44, B: 0x44, A: 0xff}, LineWidth: 2, FillColor: color.NRGBA{R: 0xf6, G: 0xb5, B: 0xcc, A: 0xff}})

	dumper := NewDumper("plot", 3, 3, 700, 400)
	defer dumper.Close()
	dumper.Plot(&plot)

	return builder.String(), nil
}
