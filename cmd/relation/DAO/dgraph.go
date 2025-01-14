package DAO

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"log"
	"sync"
)

type RespUser struct {
	Uid            string     `json:"uid,omitempty"`
	Name           string     `json:"name,omitempty"`
	Follows        []RespUser `json:"follows,omitempty"`
	Followers      []RespUser `json:"~follows,omitempty"`
	Friends        []RespUser `json:"friends,omitempty"`
	FollowsCount   int64      `json:"count(follows)"`
	FollowersCount int64      `json:"count(~follows)"`
}
type UserList struct {
	FollowList []RespUser `json:"followlist,omitempty"`
}

type DgraphUser struct {
	Uid     string       `json:"uid,omitempty"`
	Name    string       `json:"name,omitempty"`
	Follows []DgraphUser `json:"follows,omitempty"`
}

var DGO *dgo.Dgraph
var once sync.Once

func InitDGO() error {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	once.Do(func() {
		dialOpts := append([]grpc.DialOption{},
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
		d, err := grpc.Dial("localhost:9080", dialOpts...) //这里注意修改grpc端口，我的win10 9080端口是预留端口
		if err != nil {
			log.Fatalln(err)
		}
		DGO = dgo.NewDgraphClient(
			api.NewDgraphClient(d),
		)
	})

	return nil
}
func CreatSchema(ctx context.Context, schema string) error {
	if DGO == nil {
		log.Println("DGO 未实例化")
		return errors.New("DGO 未实例化")
	}
	op := &api.Operation{Schema: schema}
	if err := DGO.Alter(ctx, op); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 注意调用GetFollowList和GetFollowerList接口时，复用了User结构体，无关数据返回零值！！！
func GetFollowList(ctx context.Context, u int64) (RespUser, error) {
	q := `
	{
	followlist(func: uid(%d)){
		uid
		name
		follows{
			uid
			name
			count(follows)
			count(~follows)
		}
	count(follows)
	}
	}
`
	response, err := DGO.NewTxn().Query(ctx, fmt.Sprintf(q, u))
	if err != nil {
		log.Println(err)
		return RespUser{}, err
	}
	log.Println(response.String())
	var userlist UserList
	err = json.Unmarshal(response.Json, &userlist)
	if err != nil {
		log.Println(err)
		return RespUser{}, err
	}
	if len(userlist.FollowList) == 0 {
		log.Println("FollowList 没有数据")
		return RespUser{}, err
	}
	return userlist.FollowList[0], err

}
func GetFollowerList(ctx context.Context, u int64) (RespUser, error) {
	q := `
	{
	followlist(func: uid(%d)){
		uid
		name
		follows{
			uid
			name
			count(follows)
			count(~follows)
		}
		~follows{
			uid
			name
			count(follows)
			count(~follows)
		}
	count(~follows)
	}
	}
`
	response, err := DGO.NewTxn().Query(ctx, fmt.Sprintf(q, u))
	if err != nil {
		log.Println(err)
		return RespUser{}, err
	}
	log.Println(response.String())
	var userlist UserList
	err = json.Unmarshal(response.Json, &userlist)
	if err != nil {
		log.Println(err)
		return RespUser{}, err
	}
	if len(userlist.FollowList) == 0 {
		log.Println("FollowerList 没有数据")
		return RespUser{}, err
	}
	log.Println(userlist.FollowList[0])

	return userlist.FollowList[0], err

}
func GetFriendList(ctx context.Context, u int64) (RespUser, error) {
	follows, err := GetFollowList(ctx, u)
	if err != nil {
		return RespUser{}, err
	}
	followers, err := GetFollowerList(ctx, u)
	if err != nil {
		return RespUser{}, err
	}
	user := RespUser{Uid: follows.Uid, Name: follows.Name}
	m := make(map[string]struct{})

	for i := range follows.Follows {
		m[follows.Follows[i].Uid] = struct{}{}
	}
	for i := range followers.Followers {
		if _, ok := m[followers.Followers[i].Uid]; ok {
			user.Friends = append(user.Friends, followers.Followers[i])
		}
	}
	log.Println(user)
	return user, nil

}
func CreatDefaultSchema(ctx context.Context) error {
	schema := `
		name: string .
		follow: [uid] .
		type RespUser {
			name: string
			follow: [RespUser]
		}
	`
	log.Println("TODO: 反向边索引还未添加")
	return CreatSchema(ctx, schema)
}
func UpsertUser(ctx context.Context, u DgraphUser) error {
	pb, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
		return err
	}
	return Set(ctx, pb, nil)
}

func Set(ctx context.Context, json []byte, nquads []byte) error {
	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   json,
		SetNquads: nquads,
	}
	response, err := DGO.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(response.String())
	return nil
}
func Delete(ctx context.Context, json []byte, nquads []byte) error {
	mu := &api.Mutation{
		CommitNow:  true,
		DeleteJson: json,
		DelNquads:  nquads,
	}
	response, err := DGO.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(response.String())
	return nil
}
