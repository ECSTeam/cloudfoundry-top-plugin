package main

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
)

// [0] "/webappa/subapp1"
// [1] "/webappa"
// [2] ""
//
// findPath = "/webappabc"    => ""
// findPath = "/webappa"	  => "/webappa"
// findPath = "/webappa/"	  => "/webappa"
// findPath = "/webappa/doc"  => "/webappa"

func main() {
	fmt.Printf("start\n")
	test1()
	test2()
}

func test1() {
	hs := eventdata.NewHostStats("examplea")
	find(hs, "")
}

func test2() {
	hs := eventdata.NewHostStats("examplea")
	hs.AddPath("", "base")
	hs.AddPath("/webappa", "a level 1")
	hs.AddPath("/webappa/subapp1", "a level 2")
	hs.AddPath("/webappb", "b level 1")

	find(hs, "/webappabc")
	find(hs, "/webappa")
	find(hs, "/webappa/")
	find(hs, "/webappa/doc")
	find(hs, "/webappa/subapp")
	find(hs, "/webappa/subapp1")
	find(hs, "/webappa/subapp1/")
	find(hs, "/webappa/subapp1/xx")
	find(hs, "/")
	find(hs, "")
}
func find(hs *eventdata.HostStats, testPath string) {
	rs := hs.FindRouteStats(testPath)
	if rs != nil {
		fmt.Printf("findPath: %v  useRoute: %v\n\n", testPath, rs.Id())
	} else {
		fmt.Printf("No match\n")
	}
}
