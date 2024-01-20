package db

import (
	"fmt"
	"os"

	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/nedpals/supabase-go"
)

type PremiumTable struct {
	Id       int64  `json:"id"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Domain   string `json:"domain"`
	Quota    int64  `json:"quota"`
	CC       string `json:"cc"`
	Adblock  bool   `json:"adblock"`
}

type SniTable struct {
	Id     int64  `json:"id"`
	Server string `json:"server"`
}
type DomainTable struct {
	Location string `json:"location"`
	Domain   string `json:"domain"`
	Populate int    `json:"populate"`
	Code     string `json:"code"`
}

func Connect() *supabase.Client {
	return supabase.CreateClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"))
}

func GetDomainList() []DomainTable {
	var (
		domains = []DomainTable{}
	)

	if err := Connect().DB.From("domains").Select("*").Execute(&domains); err != nil {
		fmt.Println(err)
	}

	return domains
}

func GetSniList() []string {
	var (
		sniList = []string{}
		rows    = []SniTable{}
	)

	if err := Connect().DB.From("sni").Select("*").Execute(&rows); err != nil {
		panic(err)
	}

	for _, sni := range rows {
		sniList = append(sniList, sni.Server)
	}

	return sniList
}

func GetPremiumList() map[string][]PremiumTable {
	var (
		premiumList = map[string][]PremiumTable{}
		rows        = []PremiumTable{}
	)

	if err := Connect().DB.From("premium").Select("*").Neq("type", "dummy").Execute(&rows); err != nil {
		panic(err)
	}

	for _, premium := range rows {
		if premium.Quota > 0 {
			premiumList[premium.Type] = append(premiumList[premium.Type], premium)
		}
	}

	return premiumList
}

func UpdatePremiumQuota(name string) bool {
	rows := []PremiumTable{}
	if err := Connect().DB.From("premium").Select("*").Eq("id", name).Execute(&rows); err != nil {
		fmt.Println(err)
		return true
	}

	row := rows[0]
	row.Quota = row.Quota - (helper.GetUserStats(name) / 1000000)
	if err := Connect().DB.From("premium").Update(row).Eq("id", name).Execute(&rows); err != nil {
		fmt.Println(err)
	}

	if row.Quota > 0 {
		return true
	}
	return false
}
