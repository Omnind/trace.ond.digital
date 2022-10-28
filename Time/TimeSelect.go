package main

import (
	"encoding/csv"
	//"fmt"
	"log"
	"os"
)

//var (
//	ftime string
//)

func main() {
	//fmt.Println("Please enter the expected time: ")
	//fmt.Scanln(&ftime)
	f, err := os.Open("demodata/large1/LDG_N199_Housing TI_fqc_2022-06-20_2022-10-20.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f_csv := csv.NewReader(f)

	f_out, err := os.Create("out.csv")
	if err != nil {
		panic(err)
	}
	defer f_out.Close()
	m := map[string]bool{}

	for {
		ss, err := f_csv.Read()
		if err != nil {
			break
		}
		if ss[1][:10] == "2022-10-20" {
			m[ss[0]] = true
		}
	}

	f_dir, err := os.Open("demodata/large1")
	if err != nil {
		panic(err)
	}

	fs, err := f_dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(fs); i++ {
		fs_r, err := os.Open("demodata/large1/" + fs[i].Name())
		if err != nil {
			log.Println("Open" + fs[i].Name() + "Failed")
			break
		}
		fs_out, err := os.Create("demodata/large/" + fs[i].Name())
		if err != nil {
			log.Println("Create" + fs[i].Name() + "_out.scv" + "Failed")
			break
		}
		fs_r_csv := csv.NewReader(fs_r)
		ss, err := fs_r_csv.Read()
		if err != nil {
			log.Println("Read Header of" + fs[i].Name() + "Failed")
			continue
		}
		for j := 0; j < len(ss); j++ {
			fs_out.WriteString(ss[j])
			fs_out.WriteString(",")
		}
		fs_out.WriteString("\n")
		for {
			ss, err = fs_r_csv.Read()
			if err != nil {
				break
			}
			if m[ss[0]] {
				for j := 0; j < len(ss); j++ {
					fs_out.WriteString(ss[j])
					fs_out.WriteString(",")
				}
				fs_out.WriteString("\n")
			}
		}
	}

}
