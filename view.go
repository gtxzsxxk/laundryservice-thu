package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", "")
}

var headers = map[string]string{
	"Host":            "api.cleverschool.cn",
	"User-Agent":      "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/Windows WindowsWechat",
	"content-type":    "application/json",
	"Referer":         "https://servicewechat.com/wx8034ff6b2ab33a9e/28/page-frame.html",
	"Accept-Encoding": "gbk",
}

func post_data(url, r_body string) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	request, _ := http.NewRequest("POST", url, strings.NewReader(r_body))
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	return body
}

func get_dorms(c *gin.Context) {
	response := post_data("https://api.cleverschool.cn/washapi4/device/tower", "{}")
	var json_src map[string]interface{}
	err := json.Unmarshal(response, &json_src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true})
		return
	}
	dorms := json_src["data"]
	c.JSON(http.StatusOK, dorms)
}

type Washer struct {
	Number  string `json:"number"`
	Name    string `json:"name"`
	Status  string `json:"status"` /*working idle illegally_operated error*/
	Id      string `json:"id"`
	Eta     string `json:"eta"`
	Updated string `json:"updated"`
}

func has_chinese_number(str string) bool {
	val := (strings.Contains(str, "负") || strings.Contains(str, "零") || strings.Contains(str, "一") || strings.Contains(str, "二") || strings.Contains(str, "三") || strings.Contains(str, "四") || strings.Contains(str, "五") || strings.Contains(str, "六") || strings.Contains(str, "七") || strings.Contains(str, "八") || strings.Contains(str, "九") || strings.Contains(str, "十"))
	return val
}

func get_number(str string) int {
	floor_tmp := str
	floor_tmp = strings.ReplaceAll(floor_tmp, "负", "-")
	floor_tmp = strings.ReplaceAll(floor_tmp, "三十一", "31")
	floor_tmp = strings.ReplaceAll(floor_tmp, "三十", "30")
	floor_tmp = strings.ReplaceAll(floor_tmp, "二十九", "29")
	floor_tmp = strings.ReplaceAll(floor_tmp, "二十八", "28")
	floor_tmp = strings.ReplaceAll(floor_tmp, "二十七", "27")
	floor_tmp = strings.ReplaceAll(floor_tmp, "二十六", "26")
	floor_tmp = strings.ReplaceAll(floor_tmp, "二十五", "25")
	floor_tmp = strings.ReplaceAll(floor_tmp, "二十四", "24")
	floor_tmp = strings.ReplaceAll(floor_tmp, "二十三", "23")
	floor_tmp = strings.ReplaceAll(floor_tmp, "二十二", "22")
	floor_tmp = strings.ReplaceAll(floor_tmp, "二十一", "21")
	floor_tmp = strings.ReplaceAll(floor_tmp, "二十", "20")
	floor_tmp = strings.ReplaceAll(floor_tmp, "十九", "19")
	floor_tmp = strings.ReplaceAll(floor_tmp, "十八", "18")
	floor_tmp = strings.ReplaceAll(floor_tmp, "十七", "17")
	floor_tmp = strings.ReplaceAll(floor_tmp, "十六", "16")
	floor_tmp = strings.ReplaceAll(floor_tmp, "十五", "15")
	floor_tmp = strings.ReplaceAll(floor_tmp, "十四", "14")
	floor_tmp = strings.ReplaceAll(floor_tmp, "十三", "13")
	floor_tmp = strings.ReplaceAll(floor_tmp, "十二", "12")
	floor_tmp = strings.ReplaceAll(floor_tmp, "十一", "11")
	floor_tmp = strings.ReplaceAll(floor_tmp, "十", "10")
	floor_tmp = strings.ReplaceAll(floor_tmp, "零", "0")
	floor_tmp = strings.ReplaceAll(floor_tmp, "一", "1")
	floor_tmp = strings.ReplaceAll(floor_tmp, "二", "2")
	floor_tmp = strings.ReplaceAll(floor_tmp, "三", "3")
	floor_tmp = strings.ReplaceAll(floor_tmp, "四", "4")
	floor_tmp = strings.ReplaceAll(floor_tmp, "五", "5")
	floor_tmp = strings.ReplaceAll(floor_tmp, "六", "6")
	floor_tmp = strings.ReplaceAll(floor_tmp, "七", "7")
	floor_tmp = strings.ReplaceAll(floor_tmp, "八", "8")
	floor_tmp = strings.ReplaceAll(floor_tmp, "九", "9")

	final_str := ""
	for i := 0; i < len(floor_tmp); i++ {
		what := floor_tmp[i]
		if what == '-' || (what >= byte('0') && what <= byte('9')) {
			final_str += string(floor_tmp[i])
		}
	}
	r, _ := strconv.Atoi(final_str)
	return r
}

func get_dorms_devices(c *gin.Context) {
	id := c.Param("id")
	response := post_data("https://api.cleverschool.cn/washapi4/device/status", "{\"towerKey\":\""+id+"\",\"deviceType\":\"\"}")
	var json_src map[string]interface{}
	err := json.Unmarshal(response, &json_src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true})
		return
	}
	data := json_src["data"].([]interface{})
	floor_map := make(map[string][]Washer)
	for _, value := range data {
		v := value.(map[string]interface{})
		/* 状态：待机、工作（运转）、脱水开盖、机器出错 */
		machine_info_arr := strings.Split(v["macUnionCode"].(string), " ")
		status_shit := strings.Split(v["status"].(string), " ")
		status_str := v["status"].(string)
		var status string
		var eta string
		var updated string
		eta = "N/A"
		updated = "N/A"
		if strings.Contains(status_str, "待") {
			status = "idle"
		} else if strings.Contains(status_str, "工") || strings.Contains(status_str, "运") {
			status = "working"
			for _, value := range status_shit {
				if strings.Contains(value, "剩余") {
					eta = strings.Split(value, ":")[1]
				}
			}
		} else if strings.Contains(status_str, "脱水") || strings.Contains(status_str, "开盖") {
			status = "illegally_operated"
		} else {
			status = "error"
		}
		for _, value := range status_shit {
			if strings.Contains(value, "时间") {
				updated = strings.Split(value, ":")[1]
			}
		}
		floor := v["floorName"].(string)
		washer := Washer{
			Number:  fmt.Sprintf("#%d", (len(floor_map[floor]) + 1)),
			Name:    machine_info_arr[0],
			Id:      machine_info_arr[1],
			Status:  status,
			Eta:     eta,
			Updated: updated,
		}
		if floor_map[floor] == nil {
			floor_map[floor] = make([]Washer, 0)
		}
		floor_map[floor] = append(floor_map[floor], washer)
	}
	/* 极其低效 */
	keys := make([]string, 0)
	last_index := 0
	floor_sort := make(map[string][]Washer)
	for key := range floor_map {
		if has_chinese_number(key) {
			keys = append(keys, key)
			last_index += 1
		}
	}
	for key := range floor_map {
		if !has_chinese_number(key) {
			keys = append(keys, key)
		}
	}
	for i := 0; i < last_index; i++ {
		floor_sort[strconv.Itoa(get_number(keys[i]))+"层"] = floor_map[keys[i]]
	}
	for i := last_index; i < len(keys); i++ {
		floor_sort[keys[i]] = floor_map[keys[i]]
	}

	/*if need_sort {
		/*floor_sort := make(map[string][]Washer)
		for key := range floor_map {
			if has_chinese_number(key) {
				keys = append(keys, key)
				last_index += 1
			}
		}
		for key := range floor_map {
			if !has_chinese_number(key) {
				keys = append(keys, key)
			}
		}
		for i := 0; i < last_index; i++ {
			for j := i; j < last_index; j++ {
				a := get_number(keys[j])
				b := get_number(keys[i])
				if a < b {
					tmp := keys[j]
					keys[j] = keys[i]
					keys[i] = tmp
				}
			}
		}
		for i := 0; i < len(keys); i++ {
			floor_sort[strconv.Itoa(get_number(keys[i]))+"层"] = floor_map[keys[i]]
		}
		byt, _ := json.Marshal(floor_sort)
		//c.String(http.StatusOK, string(byt))
		return
	}*/
	c.JSON(http.StatusOK, floor_sort)
}
