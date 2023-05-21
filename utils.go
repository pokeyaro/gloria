// Copyright (c) 2023 Pokeya Boa <pokeya.mystic@gmail.com>, All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package gloria

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// isEmptyString checks if a string is empty or equals to a "-" value.
// The 's' parameter is the string to be checked.
// It returns true if the string is empty or equals to a "-" value, false otherwise.
func isEmptyString(s string) bool {
	return s == "" || s == signHorizontal
}

// containsMethod checks if a string is present in the list of query methods.
// The 's' parameter is the string to be checked.
// It returns true if the string is found in the list, false otherwise.
func containsMethod(s string) bool {
	for _, str := range QueryMethods {
		if s == str {
			return true
		}
	}
	return false
}

// isValidMethod checks if a string is a valid query method.
// The 's' parameter is the string to be checked.
// If the string is not a valid query method, it throws a panic with an error message.
func isValidMethod(s string) {
	if !containsMethod(s) {
		panic(fmt.Errorf(`Must choose one of "%s"`, strings.Join(QueryMethods, ", ")))
	}
}

// isValidHost checks if a string is a valid host.
// The 'host' parameter is the string to be checked.
// It returns true if the host is valid, and false otherwise.
func isValidHost(host string) bool {
	// Used to match valid domain name patterns
	ipReStr := `^(localhost|([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)?(:\d{1,5})?)$`
	validHostRegex := regexp.MustCompile(ipReStr)
	return validHostRegex.MatchString(host)
}

// isValidIPAddrPort checks if a string is a valid IP address and port combination.
// The 'ipAddrPort' parameter is the string to be checked.
// It returns true if the IP address and port combination is valid, and false otherwise.
func isValidIPAddrPort(ipAddrPort string) bool {
	switch parts := strings.Split(ipAddrPort, ":"); len(parts) {
	case 2:
		if port, err := strconv.Atoi(parts[1]); err != nil || port < 0 || port > 65535 {
			return false
		}
		fallthrough
	case 1:
		return net.ParseIP(parts[0]) != nil
	default:
		return false
	}
}

// isEmpty checks if a value is empty.
// The 'value' parameter is the value to be checked.
// It returns true if the value is empty, and false otherwise.
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	switch v := value.(type) {
	case string:
		return len(v) == 0
	case map[string]string:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case error:
		return v == nil
	case []*http.Cookie:
		return len(v) == 0
	case *log.Logger:
		return v == nil
	case *Exception:
		if v == nil || (v.CodeLocation == "" && v.PanicError == nil && v.FailureReason == "" && v.OccurrenceTime == 0) {
			return true
		}
		return false
	default:
		return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
	}
}

// urlSegments parses the URL path and returns a rawUrl struct.
// The 'urlpath' parameter is the URL path to be parsed.
// It returns a pointer to the rawUrl struct.
func urlSegments(urlpath string) (seg *rawUrl) {
	parsedURL, err := url.Parse(urlpath)
	if err != nil {
		panic(fmt.Errorf("URL parsing error: %w", err))
	}

	scheme := parsedURL.Scheme
	host := parsedURL.Host
	path := parsedURL.Path

	// parse query parameters
	queryParams := parsedURL.Query()

	// Parse query parameters as map[string]any
	params := make(H)
	for key, values := range queryParams {
		params[key] = strings.Join(values, ",")
	}

	// split path
	pathSegments := strings.Split(path, signSlash)

	seg = &rawUrl{
		urls: urls{
			scheme:   scheme,
			host:     host,
			baseURI:  RootURL,
			endpoint: RootURL,
		},
		params: params,
	}

	if len(pathSegments) == 2 {
		seg.endpoint = path
	} else if len(pathSegments) > 2 {
		seg.baseURI += pathSegments[1]
		seg.endpoint += strings.Join(pathSegments[2:], signSlash)
	}
	return
}

// convertToSMap converts a map of values to a string map.
// The 'input' parameter is the input map to be converted.
// It returns the converted string map.
func convertToSMap(input H) SMap {
	output := make(SMap, len(input))
	for key, value := range input {
		switch v := value.(type) {
		case string:
			output[key] = v
		case int:
			output[key] = strconv.Itoa(v)
		case float64:
			output[key] = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			output[key] = strconv.FormatBool(v)
		case []string:
			output[key] = strings.Join(v, ",")
		case []int:
			vs := make([]string, 0, len(v))
			for _, n := range v {
				vs = append(vs, strconv.Itoa(n))
			}
			output[key] = strings.Join(vs, ",")
		default:
			panic(fmt.Sprintf("Unsupported value type for key '%s': %T", key, value))
		}
	}
	return output
}

// getUserAgent generates the User-Agent string for the HTTP request.
// It combines the application name, version, operating system, architecture, and Go version.
// It returns the User-Agent string.
func getUserAgent() string {
	appName := Title
	appVer := Version
	os := runtime.GOOS
	arch := runtime.GOARCH
	version := runtime.Version()[2:]
	parts := strings.Split(version, ".")
	major, minor, patch := parts[0], parts[1], parts[2]
	ua := fmt.Sprintf("%s/%s (%s %s) Go/%s.%s.%s", appName, appVer, os, arch, major, minor, patch)
	return ua
}

// getBearerAuth generates the Bearer authentication header value.
// The 'token' parameter is the token to be included in the header.
// It returns the Bearer authentication header value.
func getBearerAuth(token string) string {
	return fmt.Sprintf("%s %s", AuthTypeBearer, token)
}

// getBasicAuth generates the Basic authentication header value.
// The 'username' and 'password' parameters are used to construct the authentication value.
// It returns the Basic authentication header value.
func getBasicAuth(username, password string) string {
	auth := fmt.Sprintf("%s:%s", username, password)
	credentials := base64.StdEncoding.EncodeToString([]byte(auth))
	return fmt.Sprintf("%s %s", AuthTypeBasic, credentials)
}
