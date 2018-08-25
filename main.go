package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type printAncestorFunc func(string, string, string)

func main() {

	dataFileName := os.Args[1]
	rootPersonId := os.Args[2]
	absDataFileName, err := filepath.Abs(dataFileName)
	if err != nil {
		fmt.Printf("ERROR:%s", err.Error())
		os.Exit(1)
	}
	dataDir := filepath.Dir(absDataFileName)
	splitDataFile(dataDir, dataFileName, []string{"place", "person", "family", "child"})

	//places, err := loadCsv(dataDir, "place")
	persons, err := loadCsv(dataDir, "person")
	family, err := loadCsv(dataDir, "family")
	child, err := loadCsv(dataDir, "child")

	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return
	}

	//placesMaps := csv2map(places)
	personsMaps := csv2map(persons)
	familyMaps := csv2map(family)
	childMaps := csv2map(child)

	//placesById := mapByKey(placesMaps, "Place")
	personsById := mapByKey(personsMaps, "Person")
	familyById := mapByKey(familyMaps, "Marriage")
	familyByChild := mapByKey(childMaps, "Child")

	// var printAncestor printAncestorFunc
	// printAncestor = func(personId string, prefix string, indent string) {

	// 	personList := personsById[personId]
	// 	if personList != nil {
	// 		fmt.Printf("%s%s%s %s nar. %s, %s\n", indent, prefix, personList[0]["Given"], personList[0]["Surname"], personList[0]["Birth date"], personList[0]["Suffix"])
	// 	} else {
	// 		fmt.Printf("%s%sNení k dispozici %s\n", indent, prefix, personId)
	// 	}

	// 	families := familyByChild[personId]
	// 	if families == nil {
	// 		return
	// 	}
	// 	familyId := familyByChild[personId][0]["Family"]
	// 	family := familyById[familyId][0]
	// 	if family != nil {
	// 		printAncestor(family["Husband"], "👨: ", strings.Join([]string{indent, "|--"}, ""))
	// 		printAncestor(family["Wife"], "👩: ", strings.Join([]string{indent, "|--"}, ""))
	// 	}
	// 	return
	// }
	// printAncestor(rootPersonId, "", "")

	var graphAncestor printAncestorFunc
	graphAncestor = func(personId string, prefix string, indent string) {

		personList := personsById[personId]
		if personList != nil {
			fmt.Printf("%s [label=\"%s %s %s\\nnar. %s\\n%s\"]\n", raw(personId), prefix, personList[0]["Given"], personList[0]["Surname"], personList[0]["Birth date"], personList[0]["Suffix"])
		}

		families := familyByChild[personId]
		if families == nil {
			return
		}
		familyId := familyByChild[personId][0]["Family"]
		family := familyById[familyId][0]
		if family != nil {
			fatherId := family["Husband"]
			if fatherId != "" {
				fmt.Printf("%s -> %s [label=\"otec\"]\n", raw(personId), raw(fatherId))
				graphAncestor(fatherId, "👨", strings.Join([]string{indent, "|--"}, ""))
			}
			motherId := family["Wife"]
			if motherId != "" {
				fmt.Printf("%s -> %s [label=\"matka\"]\n", raw(personId), raw(motherId))
				graphAncestor(motherId, "👩", strings.Join([]string{indent, "|--"}, ""))
			}
		}
	}

	fmt.Printf("digraph ancestors {\n graph[rankdir=LR]")
	graphAncestor(rootPersonId, "🌳", "")
	fmt.Printf("}\n")

}

func raw(id string) string {
	return strings.Replace(strings.Replace(id, "[", "", -1), "]", "", -1)
}

func loadCsv(dir string, name string) ([][]string, error) {
	file, _ := os.Open(filepath.Join(dir, fmt.Sprintf("%s.csv", name)))
	defer file.Close()
	return csv.NewReader(file).ReadAll()
}

func csv2map(csv [][]string) []map[string]string {
	columns := make([]string, 0)
	ret := make([]map[string]string, 0)
	for idx, row := range csv {
		if idx == 0 {
			// fmt.Printf("processing header row row %d: %v\n", idx, row)
			for _, col := range row {
				columns = append(columns, col)
			}
		} else {
			rowMap := make(map[string]string)
			for colIdx, colName := range columns {
				// fmt.Printf("processing column %s idx %d\n", colName, colIdx)
				rowMap[colName] = row[colIdx]
			}
			ret = append(ret, rowMap)
		}
	}
	return ret
}

func mapByKey(maps []map[string]string, key string) map[string][]map[string]string {
	ret := make(map[string][]map[string]string)
	for _, m := range maps {
		if ret[m[key]] == nil {
			ret[m[key]] = make([]map[string]string, 0)
		}
		ret[m[key]] = append(ret[m[key]], m)
	}
	return ret
}

func splitDataFile(dir string, dataFileName string, names []string) error {
	file, err := os.Open(dataFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for _, name := range names {
		outFileName := filepath.Join(dir, fmt.Sprintf("%s.csv", name))
		outFile, err := os.Create(outFileName)
		if err != nil {
			return err
		}
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				_, err := outFile.WriteString(strings.Join([]string{line, "\n"}, ""))
				if err != nil {
					return err
				}
			} else {
				break
			}
		}
		outFile.Close()
	}
	return nil
}
