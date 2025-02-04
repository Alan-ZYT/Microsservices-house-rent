package handler

import (
	"context"

	"github.com/micro/go-log"

	example "go-1/GetHouses/proto/example"
	"go-1/homeweb/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"go-1/homeweb/utils"
	"strconv"
	"encoding/json"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetHouses(ctx context.Context, req *example.Request, rsp *example.Response) error {

	beego.Info("==============api/v1.0/houses  GETHouseData GET succ!!=============")


	//创建返回空间
	rsp.Errno  =  utils.RECODE_OK
	rsp.Errmsg  = utils.RecodeText(rsp.Errno)


	//获取url上的参数信息
	///api/v1.0/houses?aid=1&sd=&ed=&sk=price-inc&p=1
	var aid int //地区
	aid, _ = strconv.Atoi(req.Aid)
	var sd string //起始时间
	sd  = req.Sd
	var ed string //结束时间
	ed = req.Ed
	var sk string //第三栏的信息
	sk = req.Sk
	var page int // 页
	page ,_ = strconv.Atoi(req.P)
	beego.Info(aid, sd, ed, sk, page)


	/*返回json*/
	houses := []models.House{}
	//创建orm句柄
	o := orm.NewOrm()
	//设置查找的表
	qs := o.QueryTable("house")
	//根据查询条件来查找内容
	//查找传入地区的所有房屋
	num, err := qs.Filter("area_id", aid).All(&houses)
	if err != nil {

		rsp.Errno  =  utils.RECODE_PARAMERR
		rsp.Errmsg  = utils.RecodeText(rsp.Errno)
		return nil
	}
					//计算以下所有房屋/一页现实的数量
	total_page := int(num)/models.HOUSE_LIST_PAGE_CAPACITY + 1
	house_page := 1

	house_list := []interface{}{}
	for _, house := range houses {
		o.LoadRelated(&house, "Area")
		o.LoadRelated(&house, "User")
		o.LoadRelated(&house, "Images")
		o.LoadRelated(&house, "Facilities")
		house_list = append(house_list, house.To_house_info())
	}

	rsp.TotalPage = int64(total_page)
	rsp.CurrentPage = int64(house_page)
	rsp.Houses,_ = json.Marshal(house_list)

	return nil

}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Example) Stream(ctx context.Context, req *example.StreamingRequest, stream example.Example_StreamStream) error {
	log.Logf("Received Example.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&example.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Example) PingPong(ctx context.Context, stream example.Example_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&example.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
