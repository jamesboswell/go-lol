package lol_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	lol "github.com/kdy1997/go-lol"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

var testKey = os.Getenv("RIOT_API_KEY")

func TestClient(t *testing.T) {
	if testKey == "" {
		t.Skipf("$RIOT_API_KEY is required")
		return
	}

	client := lol.New(nil, testKey)

	Convey("Client", t, func() {
		// credit: https://github.com/kevinohashi/php-riot-api/blob/master/testing.php
		const (
			testID     = 585897
			testName   = "RiotSchmick"
			testRegion = lol.NA
		)

		check := func(id int64, name string, region lol.Region, summoner *lol.Summoner) {
			So(summoner.ID, ShouldEqual, id)
			So(summoner.Name, ShouldEqual, name)
		}

		Convey(".Summoners() works", func() {
			summoners, err := client.Summoners(context.TODO(), testRegion, []int64{testID}).Do()
			So(err, ShouldBeNil)
			So(summoners, ShouldHaveLength, 1)
			check(testID, testName, testRegion, summoners[testID])
		})

		Convey(".SummonersByName() works", func() {
			summoners, err := client.SummonersByName(context.TODO(), testRegion, []string{testName}).Do()
			So(err, ShouldBeNil)
			So(summoners, ShouldHaveLength, 1)
			check(testID, testName, testRegion, summoners[strings.ToLower(testName)])
		})

		Convey(".SummonerNames() works", func() {
			names, err := client.SummonerNames(context.TODO(), testRegion, []int64{testID}).Do()
			So(err, ShouldBeNil)
			So(names, ShouldHaveLength, 1)
			So(names[testID], ShouldEqual, testName)
		})

	})
}

func TestUtil(t *testing.T) {
	Convey(".Normalize()", t, func() {
		for _, td := range []struct {
			Orig, Norm string
		}{
			{Orig: "RiotSchmick", Norm: "riotschmick"},
			{Orig: "SKT T1 Faker", Norm: "sktt1faker"},
		} {
			Convey(fmt.Sprintf("Returns (%q) for (%q)", td.Norm, td.Orig), func() {
				So(lol.Normalize(td.Orig), ShouldEqual, td.Norm)
			})
		}
	})
}
