package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"html/template"
	"io"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//分页计算函数
func Pages(results_count int, page int, pagesize int) (pages int, pageout int, beginnum int, endnum int, offset int) {
	//取得记录总数，计算总页数用
	//results_count,总共有results_count条记录

	//设定每一页显示的记录数
	if pagesize < 0 || pagesize < 1 {
		pagesize = 10 //如无设置，则默认每页显示10条记录
	}

	//计算总页数
	pages = int(math.Ceil(float64(results_count) / float64(pagesize)))
	//返回pages

	//判断页数设置,否则，设置为第一页
	if page < 0 || page < 1 {
		page = 1
	}
	if page > pages {
		page = pages
	}
	//返回page

	beginnum = page - 4
	endnum = page + 5

	if page < 5 {
		beginnum = 1
		endnum = 10 //可用链接数，现在是当前页加前后两页共5页，if条件为可用链接数的一半
	}
	if page > pages-5 {
		beginnum = pages - 9
		endnum = pages
	}
	if beginnum < 1 {
		beginnum = 1
	}
	if endnum > pages {
		endnum = pages
	}
	//返回beginnum
	//返回endnum

	//计算记录偏移量
	offset = int((page - 1) * pagesize)
	return int(pages), int(page), int(beginnum), int(endnum), offset
}

func Pagesbar(keyword string, results_max int, pages int, page int, beginnum int, endnum int, style int) (output template.HTML) {
	var raw string
	switch {
	case style == 1:
		/*
				"<div class='pagination'>"
			        "<span class='page-numbers'>共"+strconv.Itoa(pages)+"页</span>"
			        "<span class='page-numbers current'>"+strconv.Itoa(page)+"</span>"
			        "<a class='page-numbers' href='?page="+strconv.Itoa(beginnum)+"'>"+strconv.Itoa(beginnum)+"</a>"
			        "<a class="next page-numbers" href="?page="+strconv.Itoa(endnum)+">Next</a>"
			    "</div>"
		*/

		raw = "<div class='pagination'>"
		if results_max > 0 {
			raw = raw + "<span class='page-numbers'>找到相关结果" + strconv.Itoa(int(results_max)) + "个，共" + strconv.Itoa(int(pages)) + "页</span>"
			count := int(pages + 1)
			for i := 1; i < count; i++ {
				if i == page {
					raw = raw + "<span class='page-numbers current'>" + strconv.Itoa(i) + "</span>"
				} else {
					raw = raw + "<a class='page-numbers' href='?" + keyword + "page=" + strconv.Itoa(i) + "'>" + strconv.Itoa(i) + "</a>"
				}
			}
			if (page != pages) && (page < pages) {
				raw = raw + "<a class='next page-numbers' href='?" + keyword + "page=" + strconv.Itoa(page+1) + "'>Next</a>"
			}

		} else {
			raw = raw + "<h2>Shit!No Data!</h2>"
			raw = raw + "<span class='page-numbers'>共0页</span>"
		}
		output = template.HTML(raw + "</div>")

	case style == 2:
		/*
			<div class="pagination"><ul>
					<li><a href="#">&laquo;</a></li>
					<li class="active"><a href="#">1</a></li>
					<li><a href="#">2</a></li>
					<li><a href="#">3</a></li>
					<li><a href="#">4</a></li>
					<li><a href="#">&raquo;</a></li>
			</ul></div>
		*/

		if results_max > 0 {
			raw = "<div class='pagination pagination-centered'><ul>"
			count := pages + 1
			//begin page
			if (page != beginnum) && (page > beginnum) {
				raw = raw + "<li><a href='?" + keyword + "page=" + strconv.Itoa(page-1) + "'>&laquo;</a></li>"
			}
			for i := 1; i < count; i++ {
				//current page and loop pages
				if i == page {
					raw = raw + "<li class='active'><a href='javascript:void();'>" + strconv.Itoa(i) + "</a></li>"
				} else {
					raw = raw + "<li><a href='?" + keyword + "page=" + strconv.Itoa(i) + "'>" + strconv.Itoa(i) + "</a></li>"
				}
				//next page
				if (page != endnum) && (page < endnum) && (i == pages) {
					raw = raw + "<li><a href='?" + keyword + "page=" + strconv.Itoa(page+1) + "'>&raquo;</a></li>"
				}
			}
			raw = raw + "</ul></div>"
		}

		output = template.HTML(raw)

	case style == 3:
		/*
			<div class="pagenav">
				<p>
					<a href="" class="on">1</a>
					<a href="">2</a>
					<a href="">3</a>
					<a href="">4</a>
				</p>
			</div>
		*/
		raw = "<div class=\"pagenav\">"
		if results_max > 0 {
			raw = raw + "<p>"
			count := int(pages + 1)
			for i := 1; i < count; i++ {
				if i == page { //当前页
					raw = raw + "<a onclick=\"javascript:void();\" class=\"on\">" + strconv.Itoa(i) + "</a>"
				} else { //普通页码链接
					raw = raw + "<a href='?" + keyword + "page=" + strconv.Itoa(i) + "'>" + strconv.Itoa(i) + "</a>"
				}
			}
			if (page != pages) && (page < pages) { //下一页
				raw = raw + "<a class='next' href='?" + keyword + "page=" + strconv.Itoa(page+1) + "'>下一页</a>"
			}

		} else {
			raw = raw + "<h2>No Data!</h2>"
			raw = raw + "<span class='page-numbers'>共0页</span>"
		}
		raw = raw + "</p>"
		output = template.HTML(raw + "</div>")

	}

	return output
}

/** 微博时间格式化显示
 * @param timestamp，标准时间戳
 */
func SmcTimeSince(timeAt time.Time) string {
	now := time.Now()
	since := math.Abs(float64(now.Unix() - timeAt.Unix()))

	output := ""
	switch {
	case since < 60:
		output = "刚刚"
	case since < 60*60:
		output = fmt.Sprintf("%v分钟前", math.Floor(since/60))
	case since < 60*60*24:
		output = fmt.Sprintf("%v小时前", math.Floor(since/3600))
	case since < 60*60*24*2:
		output = fmt.Sprintf("昨天%v", timeAt.Format("15:04"))
	case since < 60*60*24*3:
		output = fmt.Sprintf("前天%v", timeAt.Format("15:04"))
	case timeAt.Format("2006") == now.Format("2006"):
		output = timeAt.Format("1月2日 15:04")
	default:
		output = timeAt.Format("2006年1月2日 15:04")
	}
	// if math.Floor(since/3600) > 0 {
	//     if timeAt.Format("2006-01-02") == now.Format("2006-01-02") {
	//         output = "今天 "
	//         output += timeAt.Format("15:04")
	//     } else {
	//         if timeAt.Format("2006") == now.Format("2006") {
	//             output = timeAt.Format("1月2日 15:04")
	//         } else {
	//             output = timeAt.Format("2006年1月2日 15:04")
	//         }
	//     }
	// } else {
	//     m := math.Floor(since / 60)
	//     if m > 0 {
	//         output = fmt.Sprintf("%v分钟前", m)
	//     } else {
	//         output = "刚刚"
	//     }
	// }
	return output
}

//获取这个小时的开始点
func ThisHour() time.Time {
	t := time.Now()
	year, month, day := t.Date()
	hour, _, _ := t.Clock()

	return time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
}

//获取今天的开始点
func ThisDate() time.Time {
	t := time.Now()
	year, month, day := t.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

//获取这周的开始点
func ThisWeek() time.Time {
	t := time.Now()
	year, month, day := t.AddDate(0, 0, -1*int(t.Weekday())).Date()

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

//获取这月的开始点
func ThisMonth() time.Time {
	t := time.Now()
	year, month, _ := t.Date()

	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

//获取今年的开始点
func ThisYear() time.Time {
	t := time.Now()
	year, _, _ := t.Date()

	return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
}

//生成规定范围内的整数
//设置起始数字范围，0开始,n截止
func RangeRand(n int) int {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(n)

}

//标准正态分布随机整数，n为随机个数,从0开始
func Nrand(n int64) float64 {
	//sample = NormFloat64() * desiredStdDev + desiredMean
	// 默认位置参数(期望desiredMean)为0,尺度参数(标准差desiredStdDev)为1.

	var i, sample int64 = 0, 0
	desiredMean := 0.0
	desiredStdDev := 100.0

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i < n {
		rn := int64(r.NormFloat64()*desiredStdDev + desiredMean)
		sample = rn % n
		i += 1
	}

	return math.Abs(float64(sample))
}

// 对字符串进行md5哈希,
// 返回32位小写md5结果
/*
func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
*/
func MD5(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	result := hex.EncodeToString(hash.Sum(nil))
	return result
}

// 对字符串进行md5哈希,
// 返回16位小写md5结果
func MD5_16(s string) string {
	return MD5(s)[8:24]
}

// 对字符串进行sha1哈希,
// 返回42位小写sha1结果
func SHA1(s string) string {

	hasher := sha1.New()
	hasher.Write([]byte(s))

	//result := fmt.Sprintf("%x", (hasher.Sum(nil)))
	result := hex.EncodeToString(hasher.Sum(nil))
	return result
}

// RSA加密
func RsaEncrypt(origData []byte, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(crand.Reader, pub, origData)
}

// RSA解密
func RsaDecrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(crand.Reader, priv, ciphertext)
}

func FixedpathByNumber(n int, layer int) string {

	hash := md5.New()
	o := ""
	for i := 1; i < layer+1; i++ {

		s := strconv.Itoa(int((n * i) + i))
		hash.Write([]byte(s))
		result := hex.EncodeToString(hash.Sum(nil))
		r := result[0:2]
		o += r + "/"
	}

	return o
}

func FixedpathByString(s string, layer int) string {

	hash := md5.New()
	output := ""
	for i := 1; i < layer+1; i++ {

		s += s + strconv.Itoa(i+i*i)
		hash.Write([]byte(s))
		result := hex.EncodeToString(hash.Sum(nil))
		r := result[0:2]
		output += r + "/"
	}

	return output
}

func Filehash(path string) (string, error) {

	if file, err := os.Open(path); err != nil {
		return "", err
	} else {

		h := sha1.New()

		if _, erro := io.Copy(h, file); erro != nil {
			return "", erro
		} else {

			//return fmt.Srintf("%x", h.Sum(nil))
			result := hex.EncodeToString(h.Sum(nil))
			//result := fmt.Sprintf("%d", h.Sum(nil))
			//result, _ := fmt.Printf("%d", h.Sum(nil))
			return result, nil
		}
	}

}

func Filehash_number(path string) (int, error) {

	if file, err := os.Open(path); err != nil {
		return 0, err
	} else {

		h := sha1.New()

		if _, erro := io.Copy(h, file); erro != nil {
			return 0, erro
		} else {

			//dst, _ := strconv.Atoi(fmt.Sprintf("%d", h.Sum(nil)))
			//return fmt.Srintf("%x", h.Sum(nil))
			//result := fmt.Sprintf("%d", h.Sum(nil))
			result, _ := fmt.Printf("%d", h.Sum(nil))
			return result, nil
		}
	}

}

func Filehash_block(path string, block int64) string {
	file, err := os.Open(path)
	defer file.Close()
	hash := ""

	if err != nil {
		return ""
	}

	data := make([]byte, block)
	for {
		n, err := file.Read(data)

		if n != 0 {
			//hash = MD5(string(data))
			hash = SHA1(string(data))
		} else {
			break
		}

		if err != nil && err != io.EOF {
			//panic(err)
			return ""
		}
	}

	return hash
}

/**
* user : example@example.com login smtp server user
* password: xxxxx login smtp server password
* host: smtp.example.com:port   smtp.163.com:25
* to: example@example.com;example1@163.com;example2@sina.com.cn;...
* subject:The subject of mail
* body: The content of mail
* mailtype: mail type html or text
 */
func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func GetSensitiveInfoRemovedEmail(email string) string {
	const (
		mail_separator_sign = "@"
		min_mail_id_length  = 2
	)

	emailSepPos := strings.Index(email, mail_separator_sign)

	if emailSepPos < 0 {
		return email
	}

	mailId, mailDomain := email[:emailSepPos], email[emailSepPos+1:]

	if mailIdLength := len(mailId); mailIdLength > min_mail_id_length {
		firstChar, lastChar := string(mailId[0]), string(mailId[mailIdLength-1])
		stars := "***"
		switch mailIdLength - min_mail_id_length {
		case 1:
			stars = "*"
		case 2:
			stars = "**"
		}
		mailId = firstChar + stars + lastChar
	}

	result := mailId + mail_separator_sign + mailDomain
	return result
}

func Html2str(html string) string {
	src := string(html)

	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}

//截取字符
func Substr(str string, start, length int, symbol string) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end]) + symbol
}

func GetFile(file_url string, file_path string, useragent string, referer string) error {
	f, err := os.OpenFile(file_path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	stat, err := f.Stat() //获取文件状态
	if err != nil {
		return err
	}

	ss, _ := strconv.Atoi(fmt.Sprintf("%v", stat.Size))
	f.Seek(int64(ss), 0) //把文件指针指到文件末

	req, err := http.NewRequest("GET", file_url, nil)
	if err != nil {
		return err
	}

	if useragent == "default" {
		useragent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.64 Safari/537.31"
	}

	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	req.Header.Set("User-Agent", useragent)
	req.Header.Set("Range", "bytes="+fmt.Sprintf("%s", stat.Size)+"-")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer f.Close()
	defer resp.Body.Close()

	if written, err := io.Copy(f, resp.Body); err != nil {
		return err
	} else {

		if fs, e := os.Stat(file_path); e != nil {
			if ferr := os.Remove(file_path); ferr != nil {
				fmt.Println("Remove file error:", ferr)
			}
			return err
		} else {

			if rh, e := strconv.Atoi(resp.Header.Get("Content-Length")); e != nil || fs.Size() != int64(rh) {

				if fs.Size() != int64(rh) {

					er := errors.New(file_url + " save failed")
					fmt.Println(er)

					if ferr := os.Remove(file_path); ferr != nil {
						fmt.Println("Remove file error:", ferr)
					}
					return er

				}
				return e
			} else {

				fmt.Println(file_url + " download success!")
				fmt.Println("written: ", written)
			}
		}
	}
	return err
}

func PostFile(filepath string, actionurl string, fieldname string) (*http.Response, error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	// use the body_writer to write the Part headers to the buffer
	_, err := body_writer.CreateFormFile(fieldname, filepath)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	// the file data will be the second part of the body
	fh, err := os.Open(filepath)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}
	defer fh.Close()
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()
	close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	close_buf := bytes.NewBufferString(close_string)
	// use multi-reader to defer the reading of the file data until writing to the socket buffer.
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", filepath)
		return nil, err
	}

	if req, err := http.NewRequest("POST", actionurl, request_reader); err != nil {
		return nil, err
	} else {

		// Set headers for multipart, and Content Length
		req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
		req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())

		return http.DefaultClient.Do(req)
	}

}

func WriteFile(path string, filename string, content string) error {
	//path = path[0 : len(path)-len(filename)]
	filename = path + filename
	os.MkdirAll(path, 0644)
	file, err := os.Create(filename)
	if err != nil {
		return err
	} else {
		writer := bufio.NewWriter(file)
		writer.WriteString(content)
		writer.Flush()
	}
	defer file.Close()
	return nil
}

func MoveFile(frompath string, topath string) error {

	if fromfile, err := os.Open(frompath); err != nil {
		return err
	} else {

		if tofile, err := os.OpenFile(topath, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
			return err
		} else {
			io.Copy(tofile, fromfile)
			os.Remove(frompath)
			/*
				io.Copy 在一般情况下拷贝不会出错，多个携程访问的时候可能会出现“read ./data/*.png: Access is denied.”的错误，
				造成这个错误的原因很可能是由于多个协程争抢打开文件导致，然而对实际情况是没有任何坏影响的。
				如果我们根据这个错误作出判断的话就会错上加错，所以在这里不做任何判断，完全由上帝决定好了。
			*/
			return nil

		}
	}

}

func Htmlquote(text string) string {
	//HTML编码为实体符号
	/*
	   Encodes `text` for raw use in HTML.
	       >>> htmlquote("<'&\\">")
	       '&lt;&#39;&amp;&quot;&gt;'
	*/

	text = strings.Replace(text, "&", "&amp;", -1) // Must be done first!
	text = strings.Replace(text, "<", "&lt;", -1)
	text = strings.Replace(text, ">", "&gt;", -1)
	text = strings.Replace(text, "'", "&#39;", -1)
	text = strings.Replace(text, "\"", "&quot;", -1)
	text = strings.Replace(text, "“", "&ldquo;", -1)
	text = strings.Replace(text, "”", "&rdquo;", -1)
	text = strings.Replace(text, " ", "&nbsp;", -1)
	return text
}

func Htmlunquote(text string) string {
	//实体符号解释为HTML
	/*
	   Decodes `text` that's HTML quoted.
	       >>> htmlunquote('&lt;&#39;&amp;&quot;&gt;')
	       '<\\'&">'
	*/

	// strings.Replace(s, old, new, n)
	// 在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换

	text = strings.Replace(text, "&nbsp;", " ", -1)
	text = strings.Replace(text, "&rdquo;", "”", -1)
	text = strings.Replace(text, "&ldquo;", "“", -1)
	text = strings.Replace(text, "&quot;", "\"", -1)
	text = strings.Replace(text, "&#39;", "'", -1)
	text = strings.Replace(text, "&gt;", ">", -1)
	text = strings.Replace(text, "&lt;", "<", -1)
	text = strings.Replace(text, "&amp;", "&", -1) // Must be done last!
	return text
}

func CheckPassword(password string) (b bool) {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", password); !ok {
		return false
	}
	return true
}

func CheckUsername(username string) (b bool) {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", username); !ok {
		return false
	}
	return true
}

func CheckEmail(email string) (b bool) {
	if ok, _ := regexp.MatchString(`^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+((\.[a-zA-Z0-9_-]{2,3}){1,2})$`, email); !ok {
		return false
	}
	return true
}

func Thumbnail(input_file string, output_file string, output_size string, output_align string, background string) error {
	cmd := exec.Command("convert", "-thumbnail", output_size, "-background", background, "-gravity", output_align, "-extent", output_size, input_file, output_file)
	err := cmd.Run()

	if err != nil {
		//fmt.Println("convert failed.")
		return err

	} else {
		//fmt.Println("convert okay!")
		return nil
	}

}

func Watermark(watermark_file string, input_file string, output_file string, output_align string) error {
	//composite -gravity southeast -dissolve 30 -geometry +15%+15%  lhslogo.png input_file.jpg output_file.jpg
	cmd := exec.Command("composite", "-gravity", output_align, "-dissolve", "100", watermark_file, input_file, output_file)

	err := cmd.Run()

	if err != nil {
		return err
	} else {
		return nil
	}

}

func Rex(text string, iregexp string) (b bool) {
	if ok, _ := regexp.MatchString(iregexp, text); !ok {
		return false
	}
	return true
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//发送报文 是否加密 HTTP状态 动作URL 数据内容 RSA公匙
func SendingPackets(encrypt bool, status string, actionurl string, content []byte, publicKey []byte) (*http.Response, error) {

	//只有公钥则只能加密  公钥私钥都有才能解密 所以私匙不能在客户端公开  客户端获取的内容由服务端的权限控制决定
	var body_buf io.Reader
	if encrypt {
		// RSA加密
		if rsa_encrypt_content, err := RsaEncrypt(content, publicKey); err != nil {
			return nil, err
		} else {
			body_buf = bytes.NewBufferString(string(rsa_encrypt_content))
		}
	} else {
		//无需加密
		body_buf = bytes.NewBufferString(string(content))
	}

	//hash就是各种内容的混合体加key的hash值,验证这个hash是否一致来保证内容不被非法更改
	createtime := strconv.Itoa(int(time.Now().UnixNano()))
	actionurl = actionurl + "?hash=" + Encrypt_hash(status+createtime+string(content)+string(publicKey), nil) + "-" + createtime

	if req, err := http.NewRequest(status, actionurl, body_buf); err != nil {
		return nil, err

	} else {
		hd, err := http.DefaultClient.Do(req)
		return hd, err
	}
}

func ReceivingPackets(decrypt bool, hash string, status string, content []byte, publicKey []byte, privateKey []byte) ([]byte, error) {
	//若 decrypt为真则进行RSA解密
	if decrypt {

		if rsa_decrypt_content, err := RsaDecrypt(content, privateKey); err != nil {
			return nil, err
		} else {
			content = rsa_decrypt_content
		}
	}

	//防擅改校验数据
	if hash != "" {
		s := strings.Split(hash, "-")
		hash = s[0]
		createtime := s[1]
		if (hash != "") && (createtime != "") {

			if Validate_hash(hash, status+createtime+string(content)+string(publicKey)) {
				//返回数据明文
				return content, nil
			} else {
				return nil, errors.New("报文无法通过数据校验")
			}
		}
	}
	return nil, errors.New("数据校验HASH值为空")
}
