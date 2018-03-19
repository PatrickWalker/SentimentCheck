package main

import (
    "github.com/ChimeraCoder/anaconda"
    "github.com/spf13/viper"
    "github.com/cdipaolo/sentiment"
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

func bootstrapSentimentModel() sentiment.Models {
    model, err := sentiment.Restore()
    if err != nil {
        panic(fmt.Sprintf("Could not restore model!\n\t%v\n", err))
    }
    return model
}

func TwitterSentimentAnalysis(sr anaconda.SearchResponse){
    var res []*sentiment.Analysis
    model :=  bootstrapSentimentModel()
    fmt.Println("Doing some result stuff")
    //Loop for more Tweets here
    for _ , tweet := range sr.Statuses {
        if !shouldIgnore(tweet.User.ScreenName){
            //Should check the tweet language and switch/set the sentiment language
            analysis := model.SentimentAnalysis(tweet.Text, sentiment.English)
            res = append(res,analysis)
        }

    }
    //Got all these results now so what do we do. Well we analyze them of course..
    analyzeSentimentResults(res)
}
func shouldIgnore(user string) bool{
    _, ok := ignoreList[user]
    return ok

}


func analyzeSentimentResults(res []*sentiment.Analysis){
    for _, sent := range res {
         if sent.Score  == 0 {
             fmt.Printf("Low Score For Statement %v \n ", sent.Sentences)
         } 

    }
}
