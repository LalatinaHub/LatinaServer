package helper

import (
	"context"
	"fmt"

	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/v2fly/v2ray-core/v5/app/stats/command"
	"google.golang.org/grpc"
)

func connect() command.StatsServiceClient {
	conn, err := grpc.Dial(CS.V2rayAPIAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	return command.NewStatsServiceClient(conn)
}

func GetUserStats(name string) int64 {
	resp, err := connect().GetStats(context.Background(), &command.GetStatsRequest{
		Name:   fmt.Sprintf("user>>>%s>>>traffic>>>downlink", name),
		Reset_: true,
	})
	if err != nil {
		fmt.Println(err)
	}

	if resp == nil {
		return 0
	}

	fmt.Printf("[V2Ray] %s has spent %d MB of data", name, resp.Stat.Value/1000000)
	return resp.Stat.Value
}
