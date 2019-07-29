package grpccompanyserv

import (
	"log"
	"testing"

	api "companyserv/grpccompanyserv"
	"context"
	"google.golang.org/grpc"
	"time"

	//      "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	srv api.CompanyServiceClient
	ctx context.Context
}

const address = "localhost:7070"

func TestSuite(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	assert.Nil(t, err)
	defer conn.Close()
	c := api.NewCompanyServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	myTest := &ServerSuite{
		srv: c,
		ctx: ctx,
	}
	suite.Run(t, myTest)
}

func (ss *ServerSuite) xTestCreateCompany() {
	r, err := ss.srv.CreateCompany(ss.ctx, &api.CreateCompanyRequest{Name: "a3", Title: "aaМащА«2»"})
	require.Nil(ss.T(), err)

	// New record
	//	assert.Equal(ss.T(), int32(0), r.ErrorCode)
	//	assert.NotEqual(ss.T(), int64(0), r.Id)

	// Record exists
	assert.Equal(ss.T(), &api.CreateCompanyReply{ErrorCode: 1003}, r)

}

func (ss *ServerSuite) xTestGetCompanyIDs() {
	req := &api.GetCompanyIDsRequest{
		Page:          2, // начинаем с 0
		PerPage:       2,
		SortField:     1,
		SortDirection: 0,
		AndFilters: []*api.CompanyFilter{
			&api.CompanyFilter{Filter: &api.CompanyFilter_Name{Name: "a"}},
			&api.CompanyFilter{Filter: &api.CompanyFilter_Title{Title: "a"}},
			&api.CompanyFilter{Filter: &api.CompanyFilter_CreatedAtTsSince{CreatedAtTsSince: 1}},
			&api.CompanyFilter{Filter: &api.CompanyFilter_CreatedAtTsUntil{CreatedAtTsUntil: time.Now().Unix()}},
		},
	}
	//	req.AndFilters = append(req.AndFilters,&api.CompanyFilter_Name{Name: "a"})
	r, err := ss.srv.GetCompanyIDs(ss.ctx, req)
	require.Nil(ss.T(), err)

	assert.NotEqual(ss.T(), 0, len(r.CompanyIds))
	log.Printf("IDS: %+v", r.CompanyIds)
}

func (ss *ServerSuite) xTestGetCompaniesByIDs() {
	r, err := ss.srv.GetCompaniesByIDs(ss.ctx, &api.GetCompaniesByIDsRequest{CompanyIds: []int64{32, 7}})
	require.Nil(ss.T(), err)
	assert.NotEqual(ss.T(), 0, len(r.Companies))
	for k, v := range r.Companies {
		log.Printf("ID %d: %v", k, v)
	}
}

func (ss *ServerSuite) xTestAddUserToCompany() {
	r, err := ss.srv.AddUserToCompany(ss.ctx, &api.AddUserToCompanyRequest{UserId: -1, CompanyId: 32})
	require.Nil(ss.T(), err)
	assert.NotEqual(ss.T(), 0, r.CompanyUserId)

}

func (ss *ServerSuite) xTestSwitchCUserStatus() {

	r, err := ss.srv.SwitchCUserStatus(ss.ctx, &api.SwitchCUserStatusRequest{
		// Идентификатор пользователя компании.
		CompanyUserId: 13,
		// Статус, с которого нужно переключиться.
		FromStatus: 0,
		// Статус, на который нужно переключиться.
		ToStatus: 1,
	})
	require.Nil(ss.T(), err)
	assert.Equal(ss.T(), int32(0), r.ErrorCode)
}

func (ss *ServerSuite) xTestGetCUserIDs() {
	req := &api.GetCUserIDsRequest{
		Page:          2, // начинаем с 0
		PerPage:       2,
		SortField:     1,
		SortDirection: 0,
		AndFilters: []*api.CUserFilter{
			&api.CUserFilter{Filter: &api.CUserFilter_CompanyId{CompanyId: -1}},
			&api.CUserFilter{Filter: &api.CUserFilter_UserId{UserId: -1}},
			&api.CUserFilter{Filter: &api.CUserFilter_Status{Status: 0}},
			&api.CUserFilter{Filter: &api.CUserFilter_CreatedAtTsSince{CreatedAtTsSince: 1}},
			&api.CUserFilter{Filter: &api.CUserFilter_CreatedAtTsUntil{CreatedAtTsUntil: time.Now().Unix()}},
			&api.CUserFilter{Filter: &api.CUserFilter_UpdatedAtTsSince{UpdatedAtTsSince: 1}},
			&api.CUserFilter{Filter: &api.CUserFilter_UpdatedAtTsUntil{UpdatedAtTsUntil: time.Now().Unix()}},
		},
	}
	r, err := ss.srv.GetCUserIDs(ss.ctx, req)
	require.Nil(ss.T(), err)

	assert.NotEqual(ss.T(), 0, len(r.CompanyUserIds))
	log.Printf("IDS: %+v", r.CompanyUserIds)
}

func (ss *ServerSuite) xTestGetCUsersByIDs() {
	r, err := ss.srv.GetCUsersByIDs(ss.ctx, &api.GetCUsersByIDsRequest{CompanyUserIds: []int64{13, 7}})
	require.Nil(ss.T(), err)
	assert.NotEqual(ss.T(), 0, len(r.CompanyUsers))
	for k, v := range r.CompanyUsers {
		log.Printf("ID %d: %v", k, v)
	}
}

func (ss *ServerSuite) xTestSetSessionCUser() {
	r, err := ss.srv.SetSessionCUser(ss.ctx, &api.SetSessionCUserRequest{SessionId: 1, CompanyUserId: 13})
	require.Nil(ss.T(), err)
	assert.Equal(ss.T(), int32(0), r.ErrorCode)
}

func (ss *ServerSuite) TestGetSessionCUser() {
	r, err := ss.srv.GetSessionCUser(ss.ctx, &api.GetSessionCUserRequest{SessionId: 1})
	require.Nil(ss.T(), err)
	assert.NotEqual(ss.T(), int64(0), r.CompanyUserId)
}
