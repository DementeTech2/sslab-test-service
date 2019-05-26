package fetch

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

var GetCountryRegex *regexp.Regexp
var GetOrgNameRegex *regexp.Regexp

type WhoIs struct {
	Ip     string
	Result string
}

func (w *WhoIs) GetInfo() {
	cmd := exec.Command("whois", w.Ip)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Result = string(cmdOutput.Bytes())
}

func (w *WhoIs) GetCountry() string {
	res := GetCountryRegex.FindAllStringSubmatch(w.Result, 1)

	if len(res) == 0 {
		return "NOT FOUND"
	}

	return res[0][1]
}

func (w *WhoIs) GetOwner() string {
	res := GetOrgNameRegex.FindAllStringSubmatch(w.Result, 1)

	if len(res) == 0 {
		return "NOT FOUND"
	}

	return res[0][1]
}
