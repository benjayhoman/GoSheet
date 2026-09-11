package evaluatefuncs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

const (
	HTTP_TIMEOUT      = 4 * time.Second
	JSON_DEFAULT_PATH = ""
)

func HttpFunctions() map[string]ExprFunction {
	return map[string]ExprFunction{
		"HTTPGET": func(args ...any) (any, error) {
			if err := validateArgs("HTTPGET", args, 1, 2); err != nil {
				return nil, err
			}

			url := args[0].(string)
			path := getPathArgument(args)

			body, err := httpGet(url)
			if err != nil {
				return nil, err
			}

			jsonString := string(body)

			if path == JSON_DEFAULT_PATH {
				return jsonString, nil
			}

			result := gjson.Get(jsonString, path)
			return result.String(), nil
		},
	}
}

func getPathArgument(args []any) string {
	if len(args) > 1 && strings.TrimSpace(args[1].(string)) != "" {
		return args[1].(string)
	}
	return JSON_DEFAULT_PATH
}

func httpGet(url string) (string, error) {
	client := &http.Client{
		Timeout: HTTP_TIMEOUT,
	}

	resp, err := client.Get(url)
	if err != nil {
		if os.IsTimeout(err) {
			return "", fmt.Errorf("[Timeout to endpoint %s]", url)
		}
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
