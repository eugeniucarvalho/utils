package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"sync" //"fmt"
	"time" //"reflect"

	"golang.org/x/net/html/charset"
)

type MutexCounter struct {
	Value int64
	m     sync.Mutex
}

func (t *MutexCounter) Inc(inc int64) *MutexCounter {
	t.m.Lock()
	t.Value += inc
	t.m.Unlock()
	return t
}

func (t *MutexCounter) Dec(inc int64) *MutexCounter {
	t.m.Lock()
	t.Value -= inc
	t.m.Unlock()
	return t
}

func (t *MutexCounter) Eq(val int64) bool {
	return t.Value == val
}

func Now() int64 {
	return time.Now().Unix()
}

func ToUTF8(str, origEncoding string) string {
	byteReader := bytes.NewReader([]byte(str))
	reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
	strBytes, _ := ioutil.ReadAll(reader)
	return string(strBytes)
}

func GetFiles(directory string) ([]os.FileInfo, error) {

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return []os.FileInfo{}, err
	}
	return files, nil
}

func FileGetContents(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CopyDir(source, dest string) error {
	var (
		err        error
		file       *os.File
		sourceinfo os.FileInfo
		objects    []os.FileInfo
	)
	if sourceinfo, err = os.Stat(source); err != nil {
		return err
	}

	if err = os.MkdirAll(dest, sourceinfo.Mode()); err != nil {
		return err
	}
	file, _ = os.Open(source)
	if objects, err = file.Readdir(-1); err != nil {
		return err
	}

	for _, obj := range objects {
		srcPointer := source + "/" + obj.Name()
		dstPointer := dest + "/" + obj.Name()

		if obj.IsDir() {

			if err = CopyDir(srcPointer, dstPointer); err != nil {
				return err
			}
		} else if err = CopyFile(srcPointer, dstPointer); err != nil {
			return err
		}

	}
	return nil
}

func CopyFile(source, dest string) error {

	var (
		err        error
		sourcefile *os.File
		destfile   *os.File
	)

	if sourcefile, err = os.Open(source); err != nil {
		return err
	}

	defer sourcefile.Close()

	if destfile, err = os.Create(dest); err != nil {
		return err
	}

	defer destfile.Close()

	if _, err = io.Copy(destfile, sourcefile); err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}

	return nil
}

func FilePutContents(filename string, content string, perm os.FileMode) error {
	return FilePutContentsBytes(filename, []byte(content), perm)
}

func Mkdir(perm os.FileMode, paths ...string) (err error) {

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {

			if err = os.MkdirAll(path, perm); err != nil {
				return err
			}
		}
	}
	return
}
func FilePutContentsBytes(filename string, content []byte, perm os.FileMode) error {

	var (
		// err error
		// f   *os.File
		dir = path.Dir(filename)
	)

	// fmt.Println(filename, dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {

		if err = os.MkdirAll(dir, perm); err != nil {
			return err
		}
	}

	// if _, err := os.Stat(filename); os.IsNotExist(err) {
	// 	// if f, err = os.Create(filename); err != nil {
	// 	// 	return err
	// 	// }
	//         err = ioutil.WriteFile(fi.Name(), []byte(Value), 0644)
	//     } else {
	//         fmt.Println(" the word is not in the file")
	//     }

	// } else if f, err = os.OpenFile(filename, os.O_WRONLY, os.ModeAppend); err != nil {
	// 	return err
	// }

	// defer f.Close()
	// _, err = f.Write(content)
	// return err
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func RemoveFile(file string) error {

	if err := os.Remove(file); err != nil {
		return err
	}
	return nil
}

func Encode(delimiter string, values ...string) string {
	return ToBase64(strings.Join(values, delimiter))
}

func ToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func FromBase64(value string) (string, error) {
	dec, err := base64.StdEncoding.DecodeString(value)
	return string(dec), err
}

func Decode(delimiter, value string) ([]string, error) {
	decoded, err := FromBase64(value)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(decoded), delimiter), nil
}

type Thumbnail struct {
	Resolution string

	// Adiciona o arquivo de entrada
	// path e o caminho do arquivo e pages sao as paginas a serem convertidas
	// pages = "[0]" - representa a pagina 1
	// pages = "" - representa todas as paginas
	Input string

	Output string
	InExt  string
	Page   string
}

func NewThumbnail() *Thumbnail {
	t := &Thumbnail{}
	t.Resolution = "x300"
	return t
}

var (
	ThumbSoffice = regexp.MustCompile("(doc|ppt|pps|pot|xls|xlt|xlw|dot|csv|txt|rtf)(x|m)?")
)

func (t *Thumbnail) Gen() error {
	var err error

	if t.Input == "" {
		return fmt.Errorf("Input is empty")
	}
	if t.InExt == "" {
		// todo -- pegar a extensao do arquivo de entrada
	}
	cmd := ""
	remove := false
	// Converte alguns tipos de arquivo para a imagem para depois converter para thumbnail
	if ThumbSoffice.Match([]byte(strings.ToLower(t.InExt))) {
		//cmd = "convert -thumbnail %s -background white -alpha remove %s %s"
		cmd = "soffice --headless --convert-to png --outdir %s %s"

		index := strings.LastIndex(t.Output, "/")
		outputdir := t.Output[0:index]

		if err = Exec(fmt.Sprintf(cmd, outputdir, t.Input)); err != nil {
			return err
		}

		t.Input = t.Output
		// strings.Replace(t.Input, t.InExt, ".png", -1)
		fmt.Println("new input -", t.InExt, "-", t.Input)
		remove = true
	}
	// Cria o thumbnail
	cmd = "convert -thumbnail %s -background white -alpha remove %s %s"
	Exec(fmt.Sprintf(cmd, t.Resolution, t.Input+t.Page, t.Output))
	// Apaga o arquivo temporario gerado pela conversao intermediaria
	if remove {
		//    RemoveFile(t.Input)
	}
	return nil
}

func Exec(cmd string) error {
	fmt.Println("Executando .. ", cmd)
	out, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		fmt.Println("error occured ", err.Error())
		return err
	}
	fmt.Println("out >> ", string(out))

	return err

}

/*
func Push(slice interface{}, el interface{}, call func(interface{}) bool) {

    switch reflect.TypeOf(slice).Kind() {
    case reflect.Slice:
        s    := reflect.ValueOf(t)
        push := true
        item := reflect.ValueOf(el).Elem()
        for i := 0; i < s.Len(); i++ {
            push = call(s. ,item)
        }

        if push {
            s = append(s, item)
        }
    }

}
*/
