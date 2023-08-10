package repository

import (
	"awesomeProject/webook/internal/domain"
	"awesomeProject/webook/internal/repository/dao"
	"fmt"
	"golang.org/x/net/context"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {

	u, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	ctime := changeTime(u.Ctime)
	utime := changeTime(u.Utime)
	return domain.User{
		Id:    u.Id,
		Email: u.Email,
		//Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		Abstract: u.Abstract,
		Ctime:    ctime,
		Utime:    utime,
	}, nil
}

func changeTime(srcTime int64) string {
	millisecondTimestamp := int64(srcTime) // 假设这是一个毫秒时间戳示例

	// 将毫秒时间戳除以1000得到秒级时间戳
	secondTimestamp := millisecondTimestamp / 1000

	// 将秒级时间戳转换为time.Time类型
	timestampTime := time.Unix(secondTimestamp, 0)

	// 载入中国时区
	chinaTimezone, _ := time.LoadLocation("Asia/Shanghai")

	// 将时间调整为中国时区
	timestampTimeInChina := timestampTime.In(chinaTimezone)

	// 格式化为日期时间字符串
	formattedDate := timestampTimeInChina.Format("2006-01-02 15:04:05") // 格式可以根据需求进行调整

	fmt.Println("Formatted Date in China Timezone:", formattedDate)
	return formattedDate

}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	println(1111)
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

//func (r *UserRepository) FindById(int64) {
//	//先从cache里面找
//	//再从dao里面找
//	//找到了回写cache
//}

func (r *UserRepository) Edit(ctx context.Context, id int64, u domain.User) error {
	return r.dao.EditUserProfile(ctx, id, dao.User{
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		Abstract: u.Abstract,
	})
}
