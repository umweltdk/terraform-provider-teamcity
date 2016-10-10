package main

import (
  "errors"
  "github.com/umweltdk/teamcity/teamcity"
  "os"
)

type Config struct {
  User string
  Password string
  URL string
  Insecure bool
  SkipCredsValidation bool
}

func (c *Config) Client() (interface{}, error) {
  if c.User == "" {
    c.User = os.Getenv("TEAMCITY_USER")
  }
  if c.User == "" {
    return nil, errors.New("Missing TeamCity user and TEAMCITY_USER not defined")
  }

  if c.Password == "" {
    c.Password = os.Getenv("TEAMCITY_PASSWORD")
  }
  if c.Password == "" {
    return nil, errors.New("Missing TeamCity password and TEAMCITY_PASSWORD not defined")
  }

  if c.URL == "" {
    c.URL = os.Getenv("TEAMCITY_URL")
  }
  if c.URL == "" {
    return nil, errors.New("Missing TeamCity URL and TEAMCITY_URL not defined")
  }

  client := teamcity.New(c.URL, c.User, c.Password)

  if !c.SkipCredsValidation {
    err := c.ValidateCredentials(client)
    if err != nil {
      return nil, err
    }
  }

  return client, nil
}

// Validate credentials early and fail before we do any graph walking.
func (c *Config) ValidateCredentials(client *teamcity.Client) error {
  server, err := client.Server()
  if err != nil {
    return err
  }
  if server == nil {
    return errors.New("Received no reply from server")
  }
  return nil
}
