package comm

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"lottery/conf"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// NowUnix 当前时间的时间戳
func NowUnix() int {
	return int(time.Now().In(conf.SysTimeLocation).Unix())
}

// FormatFormUnixTime 将unix时间戳格式化为yyyymmdd H:i:s格式
func FormatFormUnixTime(t int64) string {
	if t > 0 {
		return time.Unix(t, 0).Format(conf.SysTimeForm)
	} else {
		return time.Now().Format(conf.SysTimeForm)
	}
}

// FormatFormUnixTimeShort 将unix时间戳格式化为yyyymmdd
func FormatFormUnixTimeShort(t int64) string {
	if t > 0 {
		return time.Unix(t, 0).Format(conf.SysTimeFormShort)
	} else {
		return time.Now().Format(conf.SysTimeFormShort)
	}
}

// ParseTime 将字符串转成时间
func ParseTime(str string) (time.Time, error) {
	return time.ParseInLocation(conf.SysTimeForm, str, conf.SysTimeLocation)
}

// Random 得到一个随机数
func Random(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if max < 1 {
		return r.Int()
	} else {
		return r.Intn(max)
	}
}

// 对一个字符串进行加密 fixme
func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	//if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	//	return nil, err
	//}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// 对一个字符串进行解密 fixme
func decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Addslashes addslashes() 函数返回在预定义字符之前添加反斜杠的字符串。 fixme
// 预定义字符是：
// 单引号（'）
// 双引号（"）
// 反斜杠（\）
func Addslashes(str string) string {
	tmpRune := []rune{}
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

// Stripslashes 删除由Addslashes()函数添加的反斜杠
func Stripslashes(str string) string {
	dstRune := []rune{}
	strRune := []rune(str)
	strLenth := len(strRune)
	for i := 0; i < strLenth; i++ {
		if strRune[i] == []rune{'\\'}[0] {
			i++
		}
		dstRune = append(dstRune, strRune[i])
	}
	return string(dstRune)
}

// Ip4toInt 将字符串的IP转化为数字
func Ip4toInt(ip string) int64 {
	bits := strings.Split(ip, ".")
	if len(bits) == 4 {
		b0, _ := strconv.Atoi(bits[0])
		b1, _ := strconv.Atoi(bits[1])
		b2, _ := strconv.Atoi(bits[2])
		b3, _ := strconv.Atoi(bits[3])
		var sum int64
		sum += int64(b0) << 24
		sum += int64(b1) << 16
		sum += int64(b2) << 8
		sum += int64(b3)
		return sum
	} else {
		return 0
	}
}

// NextDayDuration 得到当前时间到下一个零点的延时
func NextDayDuration() time.Duration {
	year, month, day := time.Now().Add(time.Hour * 24).Date() //Add给当前时间增加一定的时间间隔，得到新的时间
	next := time.Date(year, month, day, 0, 0, 0, 0, conf.SysTimeLocation)
	return next.Sub(time.Now()) //sub返回next-time.Now(), next是下一个零点，如果当前是8月4号， 那么next是2022-08-05 00:00:00 +0800 CST
}

// GetInt64 从接口类型安全获取到int64
func GetInt64(i interface{}, d int64) int64 {
	if i == nil {
		return d
	}
	switch i.(type) {
	case string:
		num, err := strconv.Atoi(i.(string))
		if err != nil {
			return d
		} else {
			return int64(num)
		}
	case []byte:
		bits := i.([]byte)
		if len(bits) == 8 {
			return int64(binary.LittleEndian.Uint64(bits))
		} else if len(bits) <= 4 {
			num, err := strconv.Atoi(string(bits))
			if err != nil {
				return d
			} else {
				return int64(num)
			}
		}
	case uint:
		return int64(i.(uint))
	case uint8:
		return int64(i.(uint8))
	case uint16:
		return int64(i.(uint16))
	case uint32:
		return int64(i.(uint32))
	case uint64:
		return int64(i.(uint64))
	case int:
		return int64(i.(int))
	case int8:
		return int64(i.(int8))
	case int16:
		return int64(i.(int16))
	case int32:
		return int64(i.(int32))
	case int64:
		return int64(i.(int64))
	case float32:
		return int64(i.(float32))
	case float64:
		return int64(i.(float64))
	}
	return d
}

// GetString 从接口类型安全获取到字符串类型
func GetString(str interface{}, d string) string {
	if str == nil {
		return d
	}
	switch str.(type) {
	case string:
		return str.(string) // 类型断言就是将接口类型的值(x)，转换成类型(T)。格式为：x.(T)
	case []byte:
		return string(str.([]byte))
	}
	return fmt.Sprintf("%s", str)
}

// GetInt64FromMap 从map中得到指定的key
func GetInt64FromMap(dm map[string]interface{}, key string, dft int64) int64 {
	data, ok := dm[key]
	if !ok {
		return dft
	}
	return GetInt64(data, dft)
}

// GetStringFromMap 从map中得到指定的key
func GetStringFromMap(dm map[string]string, key string, dft string) string {
	data, ok := dm[key]
	if !ok {
		return dft
	}
	return data
}
