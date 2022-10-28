package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

//var (
//	ftime string
//)

func main() {
	//fmt.Println("Please enter the expected time: ")
	//fmt.Scanln(&ftime)
	//time.Now().Format("2006-01-02")
	f, err := os.Open("test_out.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f_csv := csv.NewReader(f)

	f_out, err := os.Create("Distribution/" + time.Now().Format("2006-01-02") + "_Distribution.csv") //ftime
	if err != nil {
		panic(err)
	}
	defer f_out.Close()

	ss, err := f_csv.Read()
	//m := map[int]string{}
	//for i := 16; i < len(ss); i++ {
	//	m[i] = ss[i]
	//}

	for i := 0; i < len(ss); i++ {
		f_out.WriteString(ss[i])
		f_out.WriteString(",")
	}
	f_out.WriteString("\n")

	for {
		ss, err = f_csv.Read()
		if err != nil {
			break
		}
		nums := make([]int, 0, 10000)
		sum := 0
		quan := 0
		quan_sum := 0
		for i, i_quan := 16, 1; i < len(ss); i, i_quan = i+1, i_quan+1 {
			t, err := strconv.Atoi(ss[i])
			if err != nil {
				log.Println("字符串转数字失败")
				continue
			}
			nums = append(nums, t)
			sum += t
			quan += i_quan * t
			quan_sum += i_quan
		}

		avg := fmt.Sprintf("%.3f", float64(quan)/float64(sum))

		zhe := []float64{0.1, 0.25, 0.5, 0.75, 0.9}
		th := []int{}
		for i := 0; i < len(zhe); i++ {
			th = append(th, int(zhe[i]*float64(sum)))
		}

		th_results := make([]int, 5)
		th_cha := make([]int, 5)

		for i := 0; i < len(th_cha); i++ {
			th_cha[i] = sum
		}

		max := 0
		min := 0
		max_i := 0
		min_i := 0
		sum_2 := 0
		for i := 0; i < len(nums); i++ {
			//if max < nums[i] {
			//	max = nums[i]
			//	max_i = i
			//}
			if min == 0 {
				if nums[i] != 0 {
					min = nums[i]
					min_i = i
				}
			}

			sum_2 += nums[i]
			for j := 0; j < len(th); j++ {
				if sum_2 >= th[j] {
					th_results[j] = i + 1
					th[j] = math.MaxInt
				}
			}
		}
		for i := len(nums) - 1; i >= 0; i-- {
			if max == 0 {
				if nums[i] != 0 {
					max = nums[i]
					max_i = i
				}
			}

		}
		ss[7] = strconv.Itoa(min_i + 1)
		ss[8] = strconv.Itoa(th_results[0])
		ss[9] = strconv.Itoa(th_results[1])
		ss[10] = strconv.Itoa(th_results[2])
		ss[11] = strconv.Itoa(th_results[3])
		ss[12] = strconv.Itoa(th_results[4])
		ss[13] = strconv.Itoa(max_i + 1)
		ss[14] = strconv.Itoa(sum)
		ss[15] = avg
		for i := 0; i < len(ss); i++ {
			f_out.WriteString(ss[i])
			f_out.WriteString(",")
		}
		f_out.WriteString("\n")
	}

}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}