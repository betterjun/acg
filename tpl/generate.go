package tpl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

func GenerateTemplate(tmplFile, outFile string, data interface{}) error {
	// 文件带_auto.的，为需要自动生成，且覆盖的
	if strings.Index(outFile, "_auto.") == -1 {
		existed, err := PathExists(outFile)
		if err != nil {
			return err
		}
		if existed {
			return fmt.Errorf("file existed")
		}
	}

	templ, err := ioutil.ReadFile(tmplFile)
	if err != nil {
		return err
	}
	t := template.Must(template.New("escape").Funcs(template.FuncMap{
		"formatName": FormatFieldName,
	}).Parse(string(templ)))

	//t := template.Must(template.New("escape").Parse(string(templ)))
	nb := bytes.NewBuffer(nil)
	if err := t.Execute(nb, data); err != nil {
		return err
	}

	err = ioutil.WriteFile(outFile, nb.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	ExecuteShell("gofmt", []string{"-w", outFile})

	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 执行命令
func ExecuteShell(cmdName string, cmdArgs []string) (ret string, err error) {
	fmt.Printf("ExecuteShell begin:%v %v\n", cmdName, cmdArgs)
	cmd := exec.Command(cmdName, cmdArgs...)
	out_bytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("ExecuteShell error:cmd=%v args=%v err=%v\n", cmdName, cmdArgs, err)
	}
	ret = string(out_bytes)
	fmt.Printf("ExecuteShell result:cmd=%v args=%v ret=%v\n", cmdName, cmdArgs, ret)

	return ret, err
}
