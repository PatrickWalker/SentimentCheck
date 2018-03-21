package main

import (
    "github.com/ChimeraCoder/anaconda"
    "github.com/spf13/viper"
    "github.com/SocialHarvest/sentiment"    
    
    "fmt"
    "strings"
)

var (
  ignoreList  map[string]struct{}
)
func initConfig() error {
    viper.SetConfigType("yaml")
    viper.SetConfigFile("./twitterConfig")
    
    if err := viper.ReadInConfig(); err != nil {
        fmt.Println(err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
            
		}
		return err
	}
    temp := strings.Split(viper.GetString("ignore-list"),",")
    ignoreList = make(map[string]struct{}, len(temp))
    for _, s := range temp {
        ignoreList[s] = struct{}{}
    }
    return nil
}

    

func main(){
    fmt.Println("Hello")
    err := initConfig()
    if err  != nil {
        panic("Couldn't get config so we can't do anything")
    }
    api := anaconda.NewTwitterApiWithCredentials(viper.GetString("access-token"), viper.GetString("access-token-secret"),viper.GetString("consumer-key"), viper.GetString("consumer-secret"))    
    res, _ := api.GetSearch("Cool Stuff", nil)
    //SearchResponse has a Statuses which is a []Tweet. We really need strings 
    //Also the resutls can be replenished by calling get Next so we almost need to take the res object and process it. Get each tweet. Convert to string and process
    TwitterSentimentAnalysis(res)
}


func TwitterSentimentAnalysis(sr anaconda.SearchResponse){
    neg :=   make([]string, 1)
    pos :=   make([]string, 1)
    an := sentiment.NewAnalyzer()   
    fmt.Println("Doing some result stuff")
    //Loop for more Tweets here
    fmt.Printf("Got %v results \n", len(sr.Statuses))
    for _ , tweet := range sr.Statuses {
        if !shouldIgnore(tweet){
            //Should check the tweet language and switch/set the sentiment language
            res := an.Classify(tweet.Text)
            if res < 0 {
                //Add to Negative
               neg = append(neg,tweet.Text)
            } else if res > 0 {
                //Add to Positive
               pos = append(pos,tweet.Text)
            }
        }

    }
    //Got all these results now so what do we do. Well we analyze them of course..
    analyzeSentimentResults(neg,pos)
}
func shouldIgnore(tweet anaconda.Tweet) bool{
    
    _, ok := ignoreList[tweet.User.ScreenName]
    
    notEng := tweet.Lang != "en"
    return ok || notEng

}


func analyzeSentimentResults(neg,pos []string){
    for _, neg_val := range neg {
        fmt.Println(neg_val)
    }
    fmt.Println("pOsitive")
        for _, pos_val := range pos {
        fmt.Println(pos_val)
    }
}
