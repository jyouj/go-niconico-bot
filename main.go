package main

import (
  "fmt"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "os"
  "strconv"

  "github.com/bwmarrin/discordgo"
  "golang.org/x/text/encoding/japanese"
  "golang.org/x/text/transform"
)

type Niconico struct {
	Meta struct {
		Status     int    `json:"status"`
		TotalCount int    `json:"totalCount"`
		ID         string `json:"id"`
	} `json:"meta"`
	Data []struct {
		ContentID   string `json:"contentId"`
		Title       string `json:"title"`
		ViewCounter int    `json:"viewCounter"`
	} `json:"data"`
}

func main() {
  d, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
  if err != nil {
    fmt.Println("failed to create discord session", err)
  }

  // bot, err := d.User("@" + os.Getenv("DISCORD_CLIENT_ID"))
  if err != nil {
    fmt.Println("failed to access account", err)
  }

  d.AddHandler(handleCmd)
  err = d.Open()
  if err != nil {
    fmt.Println("unable to establish connection", err)
  }

  defer d.Close()

  <-make(chan struct{})
}


func handleCmd(s *discordgo.Session, msg *discordgo.MessageCreate) {
  user := msg.Author
  if user.ID == s.State.User.ID || user.Bot {
    return
  }

  content := msg.Content

  url := "https://api.search.nicovideo.jp/api/v2/snapshot/video/contents/search?q=" + content + "&targets=title&fields=contentId,title,viewCounter&filters[viewCounter][gte]=1000&_sort=-viewCounter&_offset=0&_limit=3&_context=apiguide"
  fmt.Println(url)
  resp, _ := http.Get(url)
  defer resp.Body.Close()
  byteArray, _ := ioutil.ReadAll(resp.Body)

  fmt.Println(string(byteArray))

  jsonBytes := ([]byte)(byteArray)
  data := new(Niconico)

  if err := json.Unmarshal(jsonBytes, data); err != nil {
    fmt.Println("JSON Unmarshal error:", err)
    return
  }

  if len(data.Data) == 0 {
    s.ChannelMessageSend(msg.ChannelID, content + "は探せませんでした。英語で入力してください")
  } else {
    var vc string
    vc = strconv.Itoa(data.Data[0].ViewCounter)

    s.ChannelMessageSend(msg.ChannelID, data.Data[0].Title + " 再生回数：" + vc)
    s.ChannelMessageSend(msg.ChannelID, "https://nico.ms/" + data.Data[0].ContentID)
  }

  fmt.Printf("Message: %+v\n", msg.Message)
}
