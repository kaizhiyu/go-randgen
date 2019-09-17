package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"go-randgen/gendata"
	"go-randgen/grammar"
	"go-randgen/resource"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
)

var format bool
var breake bool
var zzPath string
var yyPath string
var outPath string
var queries int
var maxRecursive int
var root string
var dsn1 string
var dsn2 string
var debug bool

var rootCmd = &cobra.Command{
	Use:   "go port for randgen",
	Short: "random generate sql with yy and zz like mysql randgen",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if yyPath == "" {
			return errors.New("yy are required")
		}

		if (dsn1 == "" && dsn2 != "") || (dsn1 != "" && dsn2 == "") {
			return errors.New("dsn must have a pair")
		}

		return nil
	},
	Run: randgenAction,
}

// init command flag
func init() {
	rootCmd.Flags().BoolVarP(&format, "format", "F", true,
		"generate sql that is convenient for reading")
	rootCmd.Flags().BoolVarP(&breake, "break", "B", false,
		"break zz yy result to two resource")
	rootCmd.Flags().StringVarP(&zzPath, "zz", "Z","", "zz file path, go randgen have a default zz")
	rootCmd.Flags().StringVarP(&yyPath, "yy", "Y","", "yy file path, required")
	rootCmd.Flags().StringVarP(&outPath, "output", "o","output", "sql output file path")
	rootCmd.Flags().IntVarP(&queries, "queries", "Q", 100, "random sql num generated by zz")
	rootCmd.Flags().StringVarP(&root, "root", "R", "query", "root bnf expression to generate sqls")
	rootCmd.Flags().StringVar(&dsn1, "dsn1", "", "one of compare dsn")
	rootCmd.Flags().StringVar(&dsn2, "dsn2", "", "one of compare dsn")
	rootCmd.Flags().IntVar(&maxRecursive, "maxrecur", 5,
		"yy expression most recursive number, if you want recursive without limit ,set it <= 0")
	rootCmd.Flags().BoolVar(&debug, "debug", false,
		"print detail generate path")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
		os.Exit(1)
	}
}

// generate all sqls and write them into file
func randgenAction(cmd *cobra.Command, args []string) {
	var zzBs []byte
	var err error
	if zzPath == "" {
		log.Println("load default zz")
		zzBs, err = resource.Asset("resource/default.zz.lua")
	} else {
		zzBs, err = ioutil.ReadFile(zzPath)
	}

	if err != nil {
		log.Fatalf("load zz fail, %v\n", err)
	}

	zz := string(zzBs)

	yyBs, err := ioutil.ReadFile(yyPath)
	if err != nil {
		log.Fatalf("load yy from %s fail, %v\n", yyPath, err)
	}

	yy := string(yyBs)

	ddls, keyf, err := gendata.ByZz(zz)
	if err != nil {
		log.Fatalln(err)
	}

	if maxRecursive <= 0 {
		maxRecursive = math.MaxInt32
	}

	randomSqls, err := grammar.ByYy(yy, queries, root, maxRecursive, keyf, debug)
	if err != nil {
		log.Fatalln("Error: " + err.Error())
	}

	if breake {
		err := ioutil.WriteFile(outPath+".data.sql",
			[]byte(strings.Join(ddls, ";\n") + ";"), os.ModePerm)
		if err != nil {
			log.Printf("write ddl in dist fail, %v\n", err)
		}

		err = ioutil.WriteFile(outPath+".rand.sql",
			[]byte(strings.Join(randomSqls, ";\n") + ";"), os.ModePerm)
		if err != nil {
			log.Printf("write random sql in dist fail, %v\n", err)
		}
	} else {
		allSqls := make([]string, 0)
		allSqls = append(allSqls, ddls...)
		allSqls = append(allSqls, randomSqls...)

		err = ioutil.WriteFile(outPath + ".sql",
			[]byte(strings.Join(allSqls, ";\n") + ";"), os.ModePerm)
		if err != nil {
			log.Printf("sql output error, %v\n", err)
		}
	}

	if dsn1 != "" && dsn2 != "" {
		// compare two dsn
	}

}
