package utils

import (
    "fmt"
    "bytes"
    "errors"
    "io/ioutil"
    "encoding/json"

    log "github.com/sirupsen/logrus"

    "github.com/PSauerborn/lifelink/pkg/utils"
    "github.com/PSauerborn/lifelink/pkg/users"
)

var (
    // define custom errors
    ErrUserDoesNotExist  = errors.New("User does not exist")
    ErrUserAlreadyExists = errors.New("User already exists")
)

type UsersAPIAccessor struct {
    *utils.BaseAPIAccessor
}

// function to generate new API accessor for Texas Real Foods API
func NewUsersApiAccessor(host, protocol string, port *int) *UsersAPIAccessor {
    baseAccessor := utils.BaseAPIAccessor{
        Host: host,
        Port: port,
        Protocol: protocol,
    }
    return &UsersAPIAccessor{
        &baseAccessor,
    }
}

// function to generate new API accessor for Texas Real Foods API
func NewUsersApiAccessorFromConfig(config utils.APIDependencyConfig) *UsersAPIAccessor {
    baseAccessor := utils.NewAPIAccessorFromConfig(config)
    return &UsersAPIAccessor{
        baseAccessor,
    }
}

type UserDetailsResponse struct {
    HttpCode int        `json:"http_code"`
    Success  bool       `json:"success"`
    User     users.User `json:"user"`
}

// API function used to retrieve user details for a
// given user
func(accessor *UsersAPIAccessor) GetUserDetails(uid, targetUser string) (UserDetailsResponse, error) {
    log.Debug(fmt.Sprintf("retrieving user details for %s", targetUser))
    var response UserDetailsResponse
    url := accessor.FormatURL(fmt.Sprintf("/users/details/%s", targetUser))

    headers := map[string]string{"X-Authenticated-Userid": uid}
    req, err := accessor.NewJSONRequest("GET", url, nil, headers)
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
    case 404:
        log.Error("cannot retrieve user details from API: user does not exist")
        return response, ErrUserDoesNotExist
    default:
        // parse response body and log
        responseBody, _ := ioutil.ReadAll(resp.Body)
        log.Error(fmt.Errorf("received invalid response from API with status code %d: %+v", 
            resp.StatusCode, responseBody))
        return response, utils.ErrInvalidAPIResponse
    }
}

// API function used to retrieve user details for a
// given user
func(accessor *UsersAPIAccessor) CreateUser(uid string, user interface{}) (bool, error) {
    log.Debug("creating new user")
    url := accessor.FormatURL("/users/new")

    // convert request body to JSON
    body, err := json.Marshal(user)
    if err != nil {
        log.Error(fmt.Errorf("unable to serialise data to JSON: %+v", err))
        return false, err
    }
    // generate request headers and request instance
    headers := map[string]string{"X-Authenticated-Userid": uid}
    req, err := accessor.NewJSONRequest("POST", url, bytes.NewBuffer(body), headers)
    if err != nil {
        log.Error(fmt.Errorf("unable to generate new HTTP request: %+v", err))
        return false, err
    }
    // execute HTTP request
    resp, err := accessor.ExecuteRequest(req)
    if err != nil {
        log.Error(fmt.Errorf("unable to execute API request: %+v", err))
        return false, err
    }
    defer resp.Body.Close()

    switch resp.StatusCode {
    case 200:
        return true, nil
    case 400:
        log.Error(fmt.Errorf("cannot create user: invalid request body"))
        return false, utils.ErrInvalidRequestBodyJSON
    case 401:
        log.Error(fmt.Errorf("cannot create user: unauthorized"))
        return false, utils.ErrUnauthorized
    case 409:
        log.Error("unable to create new user: user already exists")
        return false, ErrUserAlreadyExists
    default:
        // parse response body and log
        responseBody, _ := ioutil.ReadAll(resp.Body)
        log.Error(fmt.Errorf("received invalid response from API with status code %d: %+v", 
            resp.StatusCode, responseBody))
        return false, utils.ErrInvalidAPIResponse
    }
}
