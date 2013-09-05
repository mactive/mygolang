package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

const baseurl = "http://jandan.net/%s/page-%d"
const column = "ooxx"
const start = 100
const end = 120

/*
 *	根据图片url保存到本地
 */
func urlretrieve(downurl string, file string) error {

	resp, err := http.Get(downurl)
	if err != nil {
		fmt.Println("err:", err)
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	fout, err := os.Create(file)
	if err != nil {
		fmt.Println("err:", err)
		return err
	}

	defer func() {
		resp.Body.Close()
		fout.Close()
	}()

	fout.Write(body)
	return nil
}

/*
 *	返回每一页的图片url列表
 */
func getImageOnepage(downurl string) []string {

	var img_re []string

	//读取页面消息主体
	resp, err := http.Get(downurl)
	if err != nil {
		fmt.Println(err)
		return img_re
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return img_re
	}

	defer resp.Body.Close()

	//去掉HTML消息体的空白字符,替换为空格。
	re_n := regexp.MustCompile(`\s`)
	bodystr := re_n.ReplaceAllString(string(body), " ")

	//查找<ol class="commentlist"><"/ol">间内容
	re1 := regexp.MustCompile(`<ol class="commentlist".*</ol>`)
	li_comment_m := re1.Find([]byte(bodystr))

	//查找<p><img src="xxx" /></p>间内容
	re2 := regexp.MustCompile(`<p><img src="(.+?)"`)
	img_urls := re2.FindAllSubmatch(li_comment_m, -1)

	for _, img_url := range img_urls {
		img_re = append(img_re, string(img_url[1]))
	}

	return img_re
}

/*
 *	返回所有的图片url列表
 */
func getImage() []string {

	var img_urls []string

	for i := start; i < end; i++ {
		url := fmt.Sprintf(baseurl, column, i)
		img_urls2 := getImageOnepage(url)
		for _, v := range img_urls2 {
			img_urls = append(img_urls, v)
		}

	}

	for _, img_url := range img_urls {
		fmt.Println(img_url)
	}
	return img_urls
}

/*
 *	保存列表中的图片
 */
func saveImage(img_urls []string) {
	for _, img_url := range img_urls {

		//解析图片的名称
		re3 := regexp.MustCompile(`.+?/`)
		img_name := re3.ReplaceAllString(img_url, "")

		urlretrieve(img_url, img_name)
	}

}

func main() {
	saveImage(getImage())
}