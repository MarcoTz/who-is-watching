package anilist_api

import (
  "fmt"
  "bytes"
  "io"
  "strings"
  "net/http"
)

func QueryApi() string{
  json_str := []byte(`{"query": "query ($name: String!) { Page { media(search:$name, type: ANIME)  { id title { romaji english native } } } }" , "variables":{"name":"ONE PIECE"} } `)

  req,err := http.NewRequest(http.MethodPost,"https://graphql.anilist.co/",bytes.NewBuffer(json_str))
  if err != nil { return fmt.Sprintf("%s",err) }

  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Accept", "application/json")
  res,err := http.DefaultClient.Do(req)
  if err != nil { return fmt.Sprintf("%s",err) }

  buf := new(strings.Builder)
  _,err = io.Copy(buf,res.Body)
  if err != nil { return fmt.Sprintf("%s",err) }

  return fmt.Sprintf(`Response Code: %d<br/>Content: %s`,res.StatusCode, buf.String())
}
