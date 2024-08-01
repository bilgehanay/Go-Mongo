package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bilgehanay/Go-Mongo/ResponseHandler"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

var (
	configFile string
	configType string
	config     ConfigModel
	e          errgroup.Group
	userdb     *mongo.Collection
	orderdb    *mongo.Collection
	cl         = make(map[string]*mongo.Collection)
	ctx        = context.Background()
)

func init() {
	flag.StringVar(&configFile, "c", "conf_dev", "Config File Name")
	flag.StringVar(&configType, "t", "json", "Config File Type")
	flag.Parse()

	viper.SetConfigName(configFile)
	viper.SetConfigType(configType)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	for _, mcnf := range config.Mongo {
		mongoConn := options.Client().ApplyURI(mcnf.ConnectionString)
		mongoConn.SetAppName(mcnf.ConnectionName)
		mc, err := mongo.Connect(ctx, mongoConn)
		fmt.Println("mc", mcnf.ConnectionString)
		if err != nil {
			panic(err)
		}
		for _, dc := range mcnf.Collection {
			cl[dc.N] = mc.Database(dc.D).Collection(dc.C)
		}
	}
	userdb = cl[config.Mongo["example"].Collection["example"].N]
	orderdb = cl[config.Mongo["example"].Collection["order"].N]
	if userdb == nil || orderdb == nil {
		fmt.Println("Db can not initilazied")
	}
	fmt.Println("Db initilazied")
	if err := ResponseHandler.LoadMessages(); err != nil {
		panic(err)
	}
}
