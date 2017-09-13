package utils

import(
    "fmt"
    //"fmt"
    //"html"
    //"net/mail"
    "os"
    "time"
    "bytes"
    "regexp"
    "os/exec"
    //"reflect"
    "strings"
    "io/ioutil"
    "encoding/base64"
    "golang.org/x/net/html/charset"
    //"errors"
    //"strings"
    //"encoding/json"
    //"gopkg.in/mgo.v2"
    //"gopkg.in/mgo.v2/bson"
    //"git.gojus.com.br/eugeniucarvalho/gojus/shared"
    //"git.gojus.com.br/eugeniucarvalho/gojus/config"
    //"git.gojus.com.br/eugeniucarvalho/gojus/services"
)

func Now() int64{
    return time.Now().Unix()
}

func ToUTF8(str, origEncoding string) string {
    byteReader  := bytes.NewReader([]byte(str))
    reader, _   := charset.NewReaderLabel(origEncoding, byteReader)
    strBytes, _ := ioutil.ReadAll(reader)
    return string(strBytes)
}

func FileGetContents(filename string) (string, error){
    b, err := ioutil.ReadFile(filename)
    if err != nil { return "", err }
    return string(b), nil
}

func FilePutContents(filename string, content string, perm os.FileMode) error {
    cont := []byte(content)
    if err := ioutil.WriteFile(filename, cont, perm) ; err != nil { return err }
    f, err1 := os.Create(filename)
    if err1 != nil { return err1 }
    defer f.Close()
    _, err2 := f.Write(cont)
    //if err2 != nil { return err2 }
    //f.Sync()
    return err2
}

func RemoveFile(file string) error {
    err := os.Remove(file)
    if err != nil {
        return err
    }
    return nil
}

func Encode(delimiter string, values ...string) string {
    return ToBase64(strings.Join(values, delimiter))
}

func ToBase64(s string) string {
    return  base64.StdEncoding.EncodeToString([]byte(s))
}

func FromBase64(value string) (string, error) {
    dec, err := base64.StdEncoding.DecodeString(value)
    return string(dec), err
}

func Decode(delimiter, value string) ([]string, error) {
    decoded, err := FromBase64(value)
    if err != nil { return []string{}, err }
    return strings.Split(string(decoded), delimiter), nil
}

type Thumbnail struct {
    Resolution string

// Adiciona o arquivo de entrada
// path e o caminho do arquivo e pages sao as paginas a serem convertidas
// pages = "[0]" - representa a pagina 1
// pages = "" - representa todas as paginas
    Input      string

    Output     string
    InExt      string
    Page       string
}

func NewThumbnail() *Thumbnail {
    t := &Thumbnail{}
    t.Resolution = "x300"
    return t;
}
var (
    ThumbSoffice = regexp.MustCompile("(doc|ppt|pps|pot|xls|xlt|xlw|dot|csv|txt|rtf)(x|m)?")
)
func(t *Thumbnail) Gen() error {
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

        index     := strings.LastIndex(t.Output,"/")
        outputdir := t.Output[0: index]

        if err = Exec(fmt.Sprintf(cmd, outputdir, t.Input)); err != nil {
            return err
        }

        t.Input = t.Output
        // strings.Replace(t.Input, t.InExt, ".png", -1)
        fmt.Println("new input -", t.InExt,"-", t.Input)
        remove = true
    }
    // Cria o thumbnail
    cmd = "convert -thumbnail %s -background white -alpha remove %s %s"
    Exec(fmt.Sprintf(cmd, t.Resolution, t.Input + t.Page, t.Output))
    // Apaga o arquivo temporario gerado pela conversao intermediaria
    if remove {
    //    RemoveFile(t.Input)
    }
    return nil;
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