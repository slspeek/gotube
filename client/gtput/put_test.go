package main

import (
  "strings"
  "testing"
)

func TestParseArgsNoneGiven(t *testing.T) {
  _, _,_, err := parseArgs([]string{"gtput"})
  if err.Error() != "usage: gtput filename [name] [description]" {
    t.Fatal(err)
  }
}

func TestParseArgsOneGiven(t *testing.T) {
  fn, name, desc, err := parseArgs([]string{"gtput", "foo.data"})
  if err != nil {
    t.Fatal(err)
  }
  if fn != "foo.data" {
    t.Fatal("Expected fn to be foo.data; was ", fn)
  }
  if name != "foo.data" {
    t.Fail()
  }
  if desc != "" {
    t.Fail()
  }
}
func TestParseArgsTwoGiven(t *testing.T) {
  fn, name, desc, err := parseArgs([]string{"gtput", "foo.data", "foo"})
  if err != nil {
    t.Fatal(err)
  }
  if fn != "foo.data" {
    t.Fatal("Expected fn to be foo.data; was ", fn)
  }
  if name != "foo" {
    t.Fail()
  }
  if desc != "" {
    t.Fail()
  }
}
func TestParseArgsThreeGiven(t *testing.T) {
  fn, name, desc, err := parseArgs([]string{"gtput", "foo.data", "foo", "fooish"})
  if err != nil {
    t.Fatal(err)
  }
  if fn != "foo.data" {
    t.Fatal("Expected fn to be foo.data; was ", fn)
  }
  if name != "foo" {
    t.Fail()
  }
  if desc != "fooish" {
    t.Fail()
  }
}


func TestParseConfig(t *testing.T) {
  conf := ` { "Username":"steven", "Password":"gnu",
  "EndPoint":"http://localhost:8080" }`
  cfg, err := parseConfigImpl(strings.NewReader(conf))
  if err != nil {
    t.Fatal(err)
  }
  
  if cfg.Username != "steven" {
    t.Fatal("Expected : steven");
  }

  if cfg.Password != "gnu" {
    t.Fatal("Expected : steven");
  }
}

