package fetcher

import (
	"bufio"
	"crawler/distributed/config"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//var cookie = "FSSBBIl1UgzbN7NO=5fHorTb66sf5JfnG1FCKu1fSQmSg1S2eQL2WiCfcH4oSwyuFuLPEm2TGn9SSpg37z2e4LKxTJelLUX3gZVw.jjG; FSSBBIl1UgzbN7NP=53.Jiabtus.LqqqDrODCFcqxDgiFSCV_N2cMSunsfOkX2Fvk2hapijIcrNoYLhZoLJa0wF0BB4fWkc6.n.Y8w75mlChAZUJNz1.N18sD9u03qq3U66UYV7f7rF4N7HvttltUOiZNm3DTMP6i5X3ERolAVWwVnYd4lXc.Fk8NtiqjJDHnges3d7fn3KP0RV11P54B.kL8uZzY4iTWIQrRRE.2nmAZRmlT1e9YEIfoFJ6a5PqGz_pooSTOgRl6ricpu4BxhJjuis1g2glZ_d_iULOjpISnXHOo5Zn4EpTFQGHq_qsx5Hgal.mluVK4l_SGm0; sid=1f657ef5-a307-48ba-af57-a83290bf3b03; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1655916224; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1655918267; _exid=vNulRNaNn88Mtb63xyqFvDQiiUP%2B8cIgjBUvCFEGJiUidRDy0qG%2BR56vBDKT2pYpLjbHLe0xwzb2yPX1sE8E4Q%3D%3D; ec=oDVvmNMI-1655916225536-e1b6fcb633c89365862116; _efmdata=NHzEJ0JS%2FeiANl7aLUwTQZI6lH8BVaisxFD0ri6VN1wAROtDFJKCFXZO7cC7%2FzBihJbsw5Yk0ylWx%2FtPImgGzZtY7PKHpSvAtqfsFUAGOSA%3D"
var rateLimiter = time.Tick(time.Second / config.Qps)

// Fetch 拉取网页url内容
// 返回[]byte数据
func Fetch(url string) ([]byte, error) {
	<-rateLimiter //限制请求太快
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	//client := &http.Client{}
	//req, err := http.NewRequest("GET", url, nil)
	//if err != nil {
	//	return nil, err
	//}

	//req.Header.Set("Cookie", "FSSBBIl1UgzbN7NO=5fHorTb66sf5JfnG1FCKu1fSQmSg1S2eQL2WiCfcH4oSwyuFuLPEm2TGn9SSpg37z2e4LKxTJelLUX3gZVw.jjG; FSSBBIl1UgzbN7NP=53.JHAbtusuLqqqDrODZ2ha80dOTppPW8BpAilanme9YHGKw5rDQgAIf1HlqnAUaY.lFnpZJZiAOT0L9BnEPy63hrHgTBOZFWsPqqu3f0f8Fl4BDFk5jpXEmgr4VI.IlcTi6d90ZfLadJaY_JEHz9eONZrEeKoxpF2fSNbSK8ItMo2LzSMuGNyQmtB171xyIBRSagzl_eJY4JC30ALmmVR42v_.9WBU1HkVTvOXiZHpkjfWcfnjM.dxh9WWmqL3wQGA.FZun4GhRIcsajz6mjmJ; sid=1f657ef5-a307-48ba-af57-a83290bf3b03; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1655916224; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1655918395; _exid=vNulRNaNn88Mtb63xyqFvDQiiUP%2B8cIgjBUvCFEGJiUidRDy0qG%2BR56vBDKT2pYpLjbHLe0xwzb2yPX1sE8E4Q%3D%3D; ec=oDVvmNMI-1655916225536-e1b6fcb633c89365862116; _efmdata=NHzEJ0JS%2FeiANl7aLUwTQZI6lH8BVaisxFD0ri6VN1wAROtDFJKCFXZO7cC7%2FzBihJbsw5Yk0ylWx%2FtPImgGzZtY7PKHpSvAtqfsFUAGOSA%3D")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:101.0) Gecko/20100101 Firefox/101.0")
	//resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Fetch Error: status code %d", resp.StatusCode)
	}

	//检测编码并替换为utf-8
	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

// determineEncoding 猜测r的编码类型(确认编码)
// 返回r的字符编码， 解析失败返回utf8
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bts, err := r.Peek(1024)
	if err != nil { //错误默认返回utf8
		log.Printf("determineEncoding Error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bts, "")
	return e
}

//func Fetch(url string) ([]byte, error) {
//
//	client := &http.Client{}
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		log.Fatalln("NewRequest is err ", err)
//		return nil, fmt.Errorf("NewRequest is err %v\n", err)
//	}
//
//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36")
//
//	//返送请求获取返回结果
//	resp, err := client.Do(req)
//
//	//直接用http.Get(url)进行获取信息，爬取时可能返回403，禁止访问
//	//resp, err := http.Get(url)
//
//	if err != nil {
//		return nil, fmt.Errorf("Error: http Get, err is %v\n", err)
//	}
//
//	//关闭response body
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("Error: StatusCode is %d\n", resp.StatusCode)
//	}
//
//	//utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
//	bodyReader := bufio.NewReader(resp.Body)
//	utf8Reader := transform.NewReader(bodyReader, determineEncoding(bodyReader).NewDecoder())
//
//	return ioutil.ReadAll(utf8Reader)
//}
//
///**
//确认编码格式
//*/
//func determineEncoding(r *bufio.Reader) encoding.Encoding {
//
//	//这里的r读取完得保证resp.Body还可读
//	body, err := r.Peek(1024)
//
//	//如果解析编码类型时遇到错误,返回UTF-8
//	if err != nil {
//		log.Printf("determineEncoding error is %v", err)
//		return unicode.UTF8
//	}
//
//	//这里简化,不取是否确认
//	e, _, _ := charset.DetermineEncoding(body, "")
//	return e
//}
