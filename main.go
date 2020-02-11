package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "strings"
    "time"
    "strconv"
)

/** ツイートJSONデコード用に構造体定義 */
type Post struct {
    Tweet struct {
        Id        string `json:"id"`
        Text      string `json:"full_text"`
        CreatedAt string `json:"created_at"`
        Favos     string `json:"favorite_count"`
    } `json:"tweet"`
}

/** 設定JSONデコード用に構造体定義 */
type Settings struct {
    IsTweet       bool `json:"is_tweet"`
    UserId        string `json:"user_id"`
    MinFavos      int `json:"min_favos"`
    Message       string `json:"message"`
}

func main() {
    // JSONファイル読み込み
    st := loadSettings("settings.json")
    posts := loadPosts("tweet.js")
    // 一番良いツイートを頼む
    msg := getMessage(st.Message, st.UserId, st.MinFavos, posts)
    if msg != "" {
        if st.IsTweet {
            Tweet(msg)
        } else {
            fmt.Println(msg)
        }
    } else {
        fmt.Println("該当なし")
    }
}

func loadPosts(path string) []Post {
    bytes, err := ioutil.ReadFile(path)
    // 余計なヘッダを消す
    source := strings.TrimLeft(string(bytes), "window.YTD.tweet.part0 = ")
    if err != nil {
        log.Fatal(err)
    }
    // JSONデコード
    var posts []Post
    if err := json.Unmarshal([]byte(source), &posts); err != nil {
        log.Fatal(err)
    }
    return posts
}

func getMessage(rawMsg string, user string, minFavos int, posts []Post) (string) {
    layout := "Mon Jan 02 15:04:05 -0700 2006"
    now := time.Now()
    jst := time.FixedZone("Asia/Tokyo", 9*60*60)
    // 集計
    dat := time.Now()
    maxFavos := -1
    pId := ""
    for _, p := range posts {
        t := p.Tweet
        tu, _ := time.Parse(layout, t.CreatedAt)
        tj := tu.In(jst)
        favos, _ := strconv.Atoi(t.Favos)
        if tj.Month() == now.Month() && tj.Day() == now.Day() && favos >= minFavos && favos > maxFavos {
            dat = tj
            maxFavos = favos
            pId = t.Id
        }
    }
    msg := ""
    if maxFavos >= 0 {
        msg = strings.Replace(rawMsg, "{{year}}", strconv.Itoa(dat.Year()), -1)
        msg = strings.Replace(msg, "{{month}}", strconv.Itoa(int(dat.Month())), -1)
        msg = strings.Replace(msg, "{{day}}", strconv.Itoa(dat.Day()), -1)
        msg = strings.Replace(msg, "{{favos}}", strconv.Itoa(maxFavos), -1)
        url := fmt.Sprintf("https://twitter.com/%s/status/%s", user, pId)
        msg = strings.Replace(msg, "{{url}}", url, -1)
    }
    return msg
}


func loadSettings(path string) Settings {
    bytes, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatal(err)
    }
    // JSONデコード
    var settings Settings
    if err := json.Unmarshal(bytes, &settings); err != nil {
        log.Fatal(err)
    }
    return settings
}
