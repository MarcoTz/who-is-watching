package anilist_api

import (
  "fmt"
  "bytes"
  "io"
  "strings"
  "net/http"
  "encoding/json"
)

func QueryApi() string{
  json_str,err := json.Marshal(`{"query":
  "query ($id: Int) { 
    Media (id: $id, type: ANIME) { 
      id
      title {
        romaji
        english
        native
      }
    }
  };", "variables":{"id":21} }`)
  if err != nil { return fmt.Sprintf("%s",err) }
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
