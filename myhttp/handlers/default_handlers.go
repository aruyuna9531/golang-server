package handlers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile(materialDir + "mainpage.html")
	if err != nil {
		panic("http mainhandler error: " + err.Error())
	}
	fmt.Fprint(w, string(b))
}
func MainHandlerByGin(c *gin.Context) {
	b, err := os.ReadFile(materialDir + "mainpage.html")
	if err != nil {
		panic("http mainhandler error: " + err.Error())
	}
	fmt.Fprint(c.Writer, string(b))
}

func AMoney(c *gin.Context) {
	b, err := os.ReadFile(materialDir + "AATool.html")
	if err != nil {
		panic("http mainhandler error: " + err.Error())
	}
	fmt.Fprint(c.Writer, string(b))
}

func AACalcResult(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusBadRequest, "提交的表单不合法")
		log.Println("AACalcResult error: " + err.Error())
		return
	}
	fmt.Printf("post form values: %v\n", c.Request.Form)
	totalprice := AtoF(c.Request.Form.Get("totalprice"))
	if totalprice < 1e-2 {
		fmt.Fprint(c.Writer, "金额总和为0，免费场还要A钱（？）")
		return
	}
	totalWeight := 0
	sep := map[string]int{} // val-min
	originH := map[string]string{}
	for i := 1; ; i++ {
		nk := fmt.Sprintf("nameinput%d", i)
		if !c.Request.Form.Has(nk) {
			break
		}
		n := c.Request.Form.Get(nk)
		h := c.Request.Form.Get(fmt.Sprintf("timehour%d", i))
		t := TransHourToMin(h)
		totalWeight += t
		sep[n] = t
		originH[n] = h
	}
	if totalWeight == 0 {
		fmt.Fprint(c.Writer, "时间总和为0，没有人参与本次A钱（？）")
		return
	}
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("总额：¥%.2f\n", totalprice))
	for n, t := range sep {
		b.WriteString(fmt.Sprintf("%s（%s小时）：¥%.2f\n", n, originH[n], float64(t)*totalprice/float64(totalWeight)))
	}
	fmt.Fprint(c.Writer, b.String())
}

func AtoF(A string) (F float64) {
	pointAt := -1
	for i := 0; i < len(A); i++ {
		if A[i] >= '0' && A[i] <= '9' {
			if pointAt == -1 {
				F *= 10
				F += float64(A[i] - '0')
			} else {
				v := float64(A[i] - '0')
				for k := pointAt; k < i; k++ {
					v /= 10
				}
				F += v
			}
		} else if A[i] == '.' {
			if pointAt == -1 {
				pointAt = i
			} else {
				panic(fmt.Sprintf("TransHourToMin duplicated point ."))
			}
		} else {
			panic(fmt.Sprintf("TransHourToMin illegal character %c", A[i]))
		}
	}
	return
}

func TransHourToMin(hourStr string) (min int) {
	pointAt := -1
	for i := 0; i < len(hourStr); i++ {
		if hourStr[i] >= '0' && hourStr[i] <= '9' {
			if pointAt == -1 {
				min = min*10 + int(hourStr[i]-'0')*60
			} else {
				v := int(hourStr[i]-'0') * 60
				for k := pointAt; k < i; k++ {
					v /= 10
				}
				min += v
			}
		} else if hourStr[i] == '.' {
			if pointAt == -1 {
				pointAt = i
			} else {
				panic(fmt.Sprintf("TransHourToMin duplicated point ."))
			}
		} else {
			panic(fmt.Sprintf("TransHourToMin illegal character %c", hourStr[i]))
		}
	}
	return
}
