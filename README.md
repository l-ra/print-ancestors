# Export ancestors graph
The simmple utility exports the ancestors graph from CSV exported by Gramps.
I use it to export ancestors graph based on data downloaded from My Herritage.
The workflow:
1. export GEDCOM data from myherritage
2. import GEDCOM data into Gramps
3. export CSV from Gramps
4. print-ancestors data.csv rootPersonId

Where data.csv are the data exported from Gramps and rootPersonId is the name 
of the person who is in the root of graph.

Command I usually use:
`go run main.go ../../data/all/all.csv [IXXXXXX] >ancestors.dot && dot -O -Tpdf ancestors.dot  && evince ancestors.dot.pdf`

To print graph you can use PosteRazor of pdfposter.