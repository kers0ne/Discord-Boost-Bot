THIS SHOULD BE A LINTER ERRORpackage main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
    "sync"
)

var (
    client = &http.Client{}
    wg     sync.WaitGroup
)

func setTitle(title string) {
    // Function to set the console title
    fmt.Printf("\033]0;%s\007", title)
}

func clear() {
    // Function to clear the console screen
    fmt.Print("\033[H\033[2J")
}

func joinServer(token, inv string) {
    defer wg.Done()

    req, _ := http.NewRequest("POST", "https://discord.com/api/v9/invites/"+inv, nil)
    req.Header.Set("Authorization", token)
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9001 Chrome/83.0.4103.122 Electron/9.3.5 Safari/537.36")
    req.Header.Set("accept", "*/*")
    req.Header.Set("accept-language", "en-US")
    req.Header.Set("cookie", "__cfduid="+randomString(43)+"; __dcfduid="+randomString(32)+"; locale=en-US")
    req.Header.Set("DNT", "1")
    req.Header.Set("origin", "https://discord.com")
    req.Header.Set("sec-fetch-dest", "empty")
    req.Header.Set("sec-fetch-mode", "cors")
    req.Header.Set("sec-fetch-site", "same-origin")
    req.Header.Set("referer", "https://discord.com/channels/@me")
    req.Header.Set("TE", "Trailers")
    req.Header.Set("X-Super-Properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC45MDAxIiwib3NfdmVyc2lvbiI6IjEwLjAuMTkwNDIiLCJvc19hcmNoIjoieDY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6ODMwNDAsImNsaWVudF9ldmVudF9zb3VyY2UiOm51bGx9")

    resp, err := client.Do(req)
    if err != nil {
        fmt.Printf("Failed To Join Server: %v\n", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusCreated {
        fmt.Println("Successfully Joined Server")
    } else {
        fmt.Printf("Failed To Join Server, Status Code: %d\n", resp.StatusCode)
    }
}

func boostServer(guildID, token string) {
    defer wg.Done()

    req, _ := http.NewRequest("PUT", "https://discord.com/api/v9/guilds/"+guildID+"/premium/subscriptions", strings.NewReader(`{"user_premium_guild_subscription_slot_ids":["`+guildID+`"]}`))
    req.Header.Set("Authorization", token)
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9001 Chrome/83.0.4103.122 Electron/9.3.5 Safari/537.36")
    req.Header.Set("accept", "*/*")
    req.Header.Set("accept-language", "en-US")
    req.Header.Set("cookie", "__cfduid="+randomString(43)+"; __dcfduid="+randomString(32)+"; locale=en-US")
    req.Header.Set("DNT", "1")
    req.Header.Set("origin", "https://discord.com")
    req.Header.Set("sec-fetch-dest", "empty")
    req.Header.Set("sec-fetch-mode", "cors")
    req.Header.Set("sec-fetch-site", "same-origin")
    req.Header.Set("referer", "https://discord.com/channels/@me")
    req.Header.Set("TE", "Trailers")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Super-Properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC45MDAxIiwib3NfdmVyc2lvbiI6IjEwLjAuMTkwNDIiLCJvc19hcmNoIjoieDY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6ODMwNDAsImNsaWVudF9ldmVudF9zb3VyY2UiOm51bGx9")

    resp, err := client.Do(req)
    if err != nil {
        fmt.Printf("Failed To Boost Server: %v\n", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        fmt.Println("Successfully Boosted Server")
    } else {
        fmt.Printf("Failed To Boost Server, Status Code: %d\n", resp.StatusCode)
    }
}


func startJoin(inv string) {
    setTitle("Server Joiner - [Blust#9380]")
    fmt.Printf("Joining servers for invite code: %s\n", inv)

    file, err := os.Open("tokens.txt")
    if err != nil {
        fmt.Printf("Error opening tokens file: %v\n", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        token := scanner.Text()
        wg.Add(1)
        go joinServer(token, inv)
    }

    wg.Wait()
}

func startBoost(id string) {
    setTitle("Server Booster - [Blust#9380]")
    fmt.Printf("Boosting server with ID: %s\n", id)

    file, err := os.Open("tokens.txt")
    if err != nil {
        fmt.Printf("Error opening tokens file: %v\n", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        token := scanner.Text()
        wg.Add(1)
        go boostServer(id, token)
    }

    wg.Wait()
}

func randomString(length int) string {
    const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    result := make([]byte, length)
    for i := range result {
        result[i] = chars[rand.Intn(len(chars))]
    }
    return string(result)
}

func main() {
    clear()
    fmt.Println("Server Tool")
    fmt.Println("--------------------------------")

    fmt.Println("[1] Server Joiner")
    fmt.Println("[2] Server Booster")

    var choice string
    fmt.Print("Choice: ")
    fmt.Scanln(&choice)

    switch choice {
    case "1":
        var inviteCode string
        fmt.Print("Enter Invite Code: discord.gg/")
        fmt.Scanln(&inviteCode)
        startJoin(inviteCode)
    case "2":
        var guildID string
        fmt.Print("Enter Guild ID: ")
        fmt.Scanln(&guildID)
        startBoost(guildID)
    default:
        fmt.Println("Invalid Choice")
    }
}
