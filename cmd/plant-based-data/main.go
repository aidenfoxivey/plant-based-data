package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/hello", helloHandler)
	log.Println("Listing for requests at http://localhost:8000/hello")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func levenshteinDistance(s string, t string) int {
	if s == t {
		return 0
	}

	m := len(s)
	n := len(t)
	d := make([][]int, m)
	distances := make([]int, n*m)

	// Loop over the rows, slicing each row from the front of the remaining pixels slice.
	for i := range d {
		d[i], distances = distances[:n], distances[n:]
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			d[i][j] = 0
		}
	}
	for i := 0; i < n; i++ {
		d[i][0] = i
	}
	for j := 0; j < n; j++ {
		d[0][j] = j
	}

	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			substitutionCost := 1
			if s[i] == t[j] {
				substitutionCost = 0
			}

			d[i][j] = min(min(d[i-1][j]+1,
				d[i][j-1]+1),
				d[i-1][j-1]+substitutionCost)
		}
	}
	return d[m][n]
}
