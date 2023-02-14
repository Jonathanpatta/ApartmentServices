package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CmdRun(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		errString := fmt.Sprintf(fmt.Sprint(err) + ": " + stderr.String())
		return "", errors.New(errString)
	}
	fmt.Println("Result: " + out.String())
	return "", nil
}

func main() {

	os.Setenv("GOOS", "linux")

	buildCmd := exec.Command("go", "build", "-o", "lambda", "main.go")
	_, err := CmdRun(buildCmd)
	if err != nil {
		panic(err)
	}

	zip := exec.Command("tar", "-a", "-c", "-f", "lambda.zip", "lambda")
	_, err = CmdRun(zip)
	if err != nil {
		panic(err)
	}

	b, err := os.ReadFile(".env") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	vars := string(b)
	vars = strings.Replace(vars, "\r\n", ",", -1)
	vars = strings.Replace(vars, "\n", ",", -1)
	vars = strings.TrimRight(vars, ",")
	vars = fmt.Sprintf("Variables={%v}", vars)
	fmt.Println(vars)

	setLambdaEnvVars := exec.Command("aws", "lambda", "update-function-configuration", "--function-name", "apartment-services", "--environment", vars)
	//fmt.Println(setLambdaEnvVars.String())
	_, err = CmdRun(setLambdaEnvVars)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	upload := exec.Command("aws", "lambda", "update-function-code", "--function-name", "apartment-services", "--zip-file", "fileb://lambda.zip")
	_, err = CmdRun(upload)
	if err != nil {
		panic(err)
	}

}
