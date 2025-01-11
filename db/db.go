package db

import (
	"strconv"
	"time"

	database "github.com/FoolVPN-ID/megalodon-api/modules/db"
	"github.com/FoolVPN-ID/megalodon-api/modules/db/servers"
	"github.com/FoolVPN-ID/megalodon-api/modules/db/users"
	"github.com/LalatinaHub/LatinaServer/helper"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func GetKVList() map[string]any {
	var (
		kvList = map[string]any{}
		client = database.MakeDatabase().GetClient()
	)

	rows, err := client.Query("SELECT * FROM kv;")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var (
			id    int
			key   string
			value any
		)

		err := rows.Scan(&id, &key, &value)
		if err != nil {
			panic(err)
		}

		kvList[key] = value
	}

	return kvList
}

func GetServerList() []servers.ServerStruct {
	var (
		serverList = []servers.ServerStruct{}
		client     = database.MakeDatabase().GetClient()
	)

	rows, err := client.Query("SELECT * FROM servers;")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		server := servers.ServerStruct{}

		err := rows.Scan(
			&server.ID,
			&server.Code,
			&server.Domain,
			&server.IP,
			&server.Country,
			&server.UsersCount,
			&server.UsersMax,
		)

		if err != nil {
			panic(err)
		}

		serverList = append(serverList, server)
	}

	return serverList
}

func GetPremiumList() map[string][]users.UserStruct {
	var (
		userList = map[string][]users.UserStruct{}
		now      = time.Now()
		client   = database.MakeDatabase().GetClient()
	)

	rows, err := client.Query("SELECT * FROM users;")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		user := users.UserStruct{}

		err := rows.Scan(
			&user.ID,
			&user.Token,
			&user.Password,
			&user.Expired,
			&user.ServerCode,
			&user.Quota,
			&user.Relay,
			&user.Adblock,
			&user.VPN,
		)

		if err != nil {
			panic(err)
		}

		// Filter users
		isExpired, _ := time.Parse("2006-01-02", user.Expired)
		if user.Quota > 0 && isExpired.Compare(now) >= 0 && user.ServerCode != "" && user.VPN != "" {
			userList[user.VPN] = append(userList[user.VPN], user)
		}
	}

	return userList
}

func UpdatePremiumQuota(name string) bool {
	var (
		user   = users.UserStruct{}
		client = database.MakeDatabase().GetClient()
	)

	id, err := strconv.Atoi(name)
	if err != nil {
		panic(err)
	}

	row := client.QueryRow("SELECT * FROM users WHERE id = ?;", id)
	err = row.Scan(
		&user.ID,
		&user.Token,
		&user.Password,
		&user.Expired,
		&user.ServerCode,
		&user.Quota,
		&user.Relay,
		&user.Adblock,
		&user.VPN,
	)

	if err != nil {
		return true
	}

	user.Quota = user.Quota - int((helper.GetUserStats(name) / 1000000))
	_, err = client.Exec("UPDATE users SET quota = ? WHERE id = ?;", user.Quota, id)
	if err != nil {
		panic(err)
	}

	if user.Quota > 0 {
		return false
	}

	return true
}
