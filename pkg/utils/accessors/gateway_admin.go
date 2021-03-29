package utils

import (
    "fmt"
    "bytes"
    "io/ioutil"
    "encoding/json"

    log "github.com/sirupsen/logrus"

    "github.com/PSauerborn/lifelink/pkg/utils"
)

type GatewayAdminAPIAccessor struct {
    *utils.BaseAPIAccessor
}

// function to generate new API accessor for Texas Real Foods API
func NewGatewayAdminApiAccessor(host, protocol string, port *int) *GatewayAdminAPIAccessor {
    baseAccessor := utils.BaseAPIAccessor{
        Host: host,
        Port: port,
        Protocol: protocol,
    }
    return &GatewayAdminAPIAccessor{
        &baseAccessor,
    }
}

// function to generate new API accessor for Texas Real Foods API
func NewGatewayAdminApiAccessorFromConfig(config utils.APIDependencyConfig) *GatewayAdminAPIAccessor {
    baseAccessor := utils.NewAPIAccessorFromConfig(config)
    return &GatewayAdminAPIAccessor{
        baseAccessor,
    }
}

type TokenResponse struct {
    HttpCode int    `json:"http_code"`
    Success  bool   `json:"success"`
    Token    string `json:"token"`
}

// API function used to retrieve user details for a
// given user
func(accessor *GatewayAdminAPIAccessor) GetAccessToken(uid string, admin bool) (TokenResponse, error) {
    log.Debug("creating new user")
    var response TokenResponse
    url := accessor.FormatURL("admin/token")

    // convert request body to JSON
    body, err := json.Marshal(map[string]interface{}{"uid": uid, "admin": admin})
    if err != nil {
        log.Error(fmt.Errorf("unable to serialise data to JSON: %+v", err))
        return response, err
    }
    // generate request headers and request instance
    headers := map[string]string{"X-Authenticated-Userid": uid}
    req, err := accessor.NewJSONRequest("POST", url, bytes.NewBuffer(body), headers)
    if err != nil {
        log.Error(fmt.Errorf("unable to generate new HTTP request: %+v", err))
        return response, err
    }
    // execute HTTP request
    resp, err := accessor.ExecuteRequest(req)
    if err != nil {
        log.Error(fmt.Errorf("unable to execute API request: %+v", err))
        return response, err
    }
    defer resp.Body.Close()

    switch resp.StatusCode {
    case 200:
        // decode JSON response and return
        if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
            log.Error(fmt.Errorf("unable to decode JSON response: %+v", err))
            return response, err
        }
        return response, nil
    default:
        // parse response body and log
        responseBody, _ := ioutil.ReadAll(resp.Body)
        log.Error(fmt.Errorf("received invalid response from API with status code %d: %+v",
            resp.StatusCode, responseBody))
        return response, utils.ErrInvalidAPIResponse
    }
}
