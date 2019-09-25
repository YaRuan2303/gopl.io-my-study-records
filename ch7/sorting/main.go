// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 187.

// Sorting sorts a music playlist into a variety of orders.
package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

//!+main
type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{  //key,这是什么意思？？*Track指针类型的切片，切片内的元素为*Track指针类型，对*Track类型的切片值进行初始化操作；
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

//!-main

//!+printTracks
func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
}

//!-printTracks

//!+artistcode
type byArtist []*Track

func (x byArtist) Len() int           { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist } //表示从小到大的顺序排序？
func (x byArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] } //x[i]表示一个*Track类型的指针变量；

//!-artistcode

//!+yearcode
type byYear []*Track

func (x byYear) Len() int           { return len(x) }
func (x byYear) Less(i, j int) bool { return x[i].Year < x[j].Year }
func (x byYear) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

//!-yearcode

func main() {
	fmt.Println("byArtist:")
	sort.Sort(byArtist(tracks))
	printTracks(tracks)

	fmt.Println("\nReverse(byArtist):")
	sort.Sort(sort.Reverse(byArtist(tracks)))
	printTracks(tracks)

	fmt.Println("\nbyYear:")
	sort.Sort(byYear(tracks))
	printTracks(tracks)

	fmt.Println("\nCustom:")
	//!+customcall
	sort.Sort(customSort{tracks, func(x, y *Track) bool {
		if x.Title != y.Title {
			return x.Title < y.Title
		}
		if x.Year != y.Year {
			return x.Year < y.Year
		}
		if x.Length != y.Length {
			return x.Length < y.Length
		}
		return false
	}})
	//!-customcall
	printTracks(tracks)
}

/*
//!+artistoutput
Title       Artist          Album              Year  Length
-----       ------          -----              ----  ------
Go Ahead    Alicia Keys     As I Am            2007  4m36s
Go          Delilah         From the Roots Up  2012  3m38s
Ready 2 Go  Martin Solveig  Smash              2011  4m24s
Go          Moby            Moby               1992  3m37s
//!-artistoutput

//!+artistrevoutput
Title       Artist          Album              Year  Length
-----       ------          -----              ----  ------
Go          Moby            Moby               1992  3m37s
Ready 2 Go  Martin Solveig  Smash              2011  4m24s
Go          Delilah         From the Roots Up  2012  3m38s
Go Ahead    Alicia Keys     As I Am            2007  4m36s
//!-artistrevoutput

//!+yearoutput
Title       Artist          Album              Year  Length
-----       ------          -----              ----  ------
Go          Moby            Moby               1992  3m37s
Go Ahead    Alicia Keys     As I Am            2007  4m36s
Ready 2 Go  Martin Solveig  Smash              2011  4m24s
Go          Delilah         From the Roots Up  2012  3m38s
//!-yearoutput

//!+customout
Title       Artist          Album              Year  Length
-----       ------          -----              ----  ------
Go          Moby            Moby               1992  3m37s
Go          Delilah         From the Roots Up  2012  3m38s
Go Ahead    Alicia Keys     As I Am            2007  4m36s
Ready 2 Go  Martin Solveig  Smash              2011  4m24s
//!-customout
*/

//!+customcode
type customSort struct {
	t    []*Track
	less func(x, y *Track) bool  //一个func(x, y *Track) bool函数类型的函数变量；
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }

//!-customcode

func init() {
	//!+ints
	values := []int{3, 1, 4, 1} //初始化一个int类型的切片；
	fmt.Println(sort.IntsAreSorted(values)) // "false"
	sort.Ints(values)  //Sort排序是从小到大的？
	fmt.Println(values)                     // "[1 1 3 4]"
	fmt.Println(sort.IntsAreSorted(values)) // "true"
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	fmt.Println(values)                     // "[4 3 1 1]"  //why?? why not "[1 1 3 4]"
	fmt.Println(sort.IntsAreSorted(values)) // "false"
	//!-ints
}
