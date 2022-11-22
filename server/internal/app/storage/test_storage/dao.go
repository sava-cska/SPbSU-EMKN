package test_storage

import (
	"errors"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"time"
)

type TestDAO struct {
	LoginToUser           map[string]*models.User
	TokenToRegistration   map[string]*TestRegistrationData
	CurrentPeriodId       uint32
	Periods               map[uint32]*models.Period
	Courses               map[uint32]*TestCourseData
	TokenToChangePassword map[string]*TestChangePasswordData
	UserAvatars           map[uint32]*models.Profile
}

func (dao *TestDAO) AddNewUser(id uint32, login string, passwd string, email string, fname string, lname string) {
	dao.LoginToUser[login] = &models.User{
		Login:     login,
		Password:  passwd,
		Email:     email,
		ProfileId: id,
		FirstName: fname,
		LastName:  lname,
	}
}

func (dao *TestDAO) AddNewCourse(course *TestCourseData) {
	dao.Courses[course.Course.Id] = course
}

func (dao *TestDAO) AddRegDetails(details *TestRegistrationData) {
	dao.TokenToRegistration[details.Token] = details
}

func (dao *TestDAO) AddPeriod(id uint32, text string) {
	dao.Periods[id] = &models.Period{
		Id:   id,
		Text: text,
	}
}

func (dao *TestDAO) AddChangePwdData(data *TestChangePasswordData) {
	dao.TokenToChangePassword[data.Token] = data
}

func (dao *TestDAO) AddUserProfile(userId uint32, avatarUrl string, fname string, lname string) {
	dao.UserAvatars[userId] = &models.Profile{
		ProfileId: userId,
		AvatarUrl: avatarUrl,
		FirstName: fname,
		LastName:  lname,
	}
}

func (dao *TestDAO) ExistsLogin(login string) bool {
	_, ok := dao.LoginToUser[login]
	return ok
}

func (dao *TestDAO) ExistsEmail(email string) bool {
	for _, user := range dao.LoginToUser {
		if user.Email == email {
			return true
		}
	}
	return false
}

func (dao *TestDAO) AddUser(user *models.User) error {
	for _, other := range dao.LoginToUser {
		if other.Email == user.Email {
			return errors.New("user with this email exists")
		}
		if other.Login == user.Login {
			return errors.New("user with this login exists")
		}
	}

	dao.LoginToUser[user.Login] = user
	return nil
}

func (dao *TestDAO) FindUser(email string) (models.User, error) {
	for _, user := range dao.LoginToUser {
		if user.Email == email {
			return *user, nil
		}
	}
	return models.User{}, errors.New("not found")
}

func (dao *TestDAO) FindUserByLogin(login string) (models.User, error) {
	user, ok := dao.LoginToUser[login]
	if ok {
		return *user, nil
	} else {
		return models.User{}, errors.New("not found")
	}
}

func (dao *TestDAO) GetPassword(login string) (string, error) {
	user, ok := dao.LoginToUser[login]
	if ok {
		return user.Password, nil
	} else {
		return "", nil
	}
}

func (dao *TestDAO) UpdatePassword(login string, newPassword string) error {
	user, ok := dao.LoginToUser[login]
	if ok {
		user.Password = newPassword
		return nil
	} else {
		return errors.New("not found")
	}
}

func (dao *TestDAO) Upsert(
	token string,
	user *models.User,
	expireDate time.Time,
	verificationCode string,
) error {
	dao.TokenToRegistration[token] = &TestRegistrationData{
		Token:            token,
		Login:            user.Login,
		Password:         user.Password,
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		ExpireData:       &expireDate,
		VerificationCode: verificationCode,
	}
	return nil
}

func (dao *TestDAO) FindRegistration(token string) (models.User, time.Time, string, error) {
	reg, ok := dao.TokenToRegistration[token]
	if !ok {
		return models.User{}, time.Time{}, "", errors.New("not found")
	}
	return models.User{
		Login:     reg.Login,
		Password:  reg.Password,
		Email:     reg.Email,
		FirstName: reg.FirstName,
		LastName:  reg.LastName,
	}, *reg.ExpireData, reg.VerificationCode, nil
}

func (dao *TestDAO) FindRegistrationAndDelete(token string) (models.User, time.Time, string, error) {
	user, exp, code, err := dao.FindRegistration(token)
	if err != nil {
		return user, exp, code, err
	}
	delete(dao.TokenToRegistration, token)
	return user, exp, code, err
}

func (dao *TestDAO) GetInfo() (*models.GeneralInfo, error) {
	return &models.GeneralInfo{CurrentPeriodId: dao.CurrentPeriodId}, nil
}

func (dao *TestDAO) ExistPeriod(periodId uint32) (bool, error) {
	_, ok := dao.Periods[periodId]
	return ok, nil
}

func (dao *TestDAO) GetPeriods() ([]*models.Period, error) {
	res := make([]*models.Period, 0)
	for _, period := range dao.Periods {
		res = append(res, period)
	}
	return res, nil
}

func (dao *TestDAO) GetDescription(courseId uint32) (*string, error) {
	course, ok := dao.Courses[courseId]
	if !ok {
		return nil, nil
	} else {
		return &course.Course.ShortDescription, nil
	}
}

func (dao *TestDAO) GetDescriptionTimestamp(courseId uint32) (*time.Time, error) {
	course, ok := dao.Courses[courseId]
	if !ok {
		return nil, nil
	} else {
		return course.DescTimestamp, nil
	}
}

func (dao *TestDAO) GetCoursesByPeriod(periodId uint32) ([]*models.CourseInDB, error) {
	res := make([]*models.CourseInDB, 0)
	for _, c := range dao.Courses {
		if c.PeriodId == periodId {
			res = append(res, c.Course)
		}
	}
	return res, nil
}

func (dao *TestDAO) ExistCourse(courseId uint32) (bool, error) {
	_, ok := dao.Courses[courseId]
	return ok, nil
}

// GetVerificationCodeInfo returns (if verification code is valid, expiresAt, error). Returns empty string if not found
func (dao *TestDAO) GetVerificationCodeInfo(identificationToken string) (string, *time.Time, error) {
	data, ok := dao.TokenToChangePassword[identificationToken]
	if !ok {
		return "", nil, errors.New("not found")
	} else {
		return data.VerificationCode, data.ExpireDate, nil
	}
}

// SetChangePasswordToken remembers changePasswordToken for identificationToken issued in accounts/begin_change_password
func (dao *TestDAO) SetChangePasswordToken(identificationToken string, changeTime time.Time, changePasswordToken string) error {
	data, ok := dao.TokenToChangePassword[identificationToken]
	if !ok {
		return errors.New("not found")
	}
	data.ChangePasswordToken = changePasswordToken
	data.ChangePasswordExpireTime = &changeTime
	return nil
}

func (dao *TestDAO) UpsertChangePasswordData(token string, login string, expiredTime time.Time, verificationCode string) error {
	dao.TokenToChangePassword[token] = &TestChangePasswordData{
		Token:                    token,
		Login:                    login,
		ExpireDate:               &expiredTime,
		VerificationCode:         verificationCode,
		ChangePasswordToken:      token,
		ChangePasswordExpireTime: &time.Time{},
	}
	return nil
}

func (dao *TestDAO) FindTokenAndDelete(token string) (string, error) {
	data, ok := dao.TokenToChangePassword[token]
	if !ok {
		return "", errors.New("not found")
	}
	delete(dao.TokenToChangePassword, token)
	return data.Login, nil
}

func (dao *TestDAO) FindPwdToken(changePwdToken string) (string, time.Time, error) {
	for _, data := range dao.TokenToChangePassword {
		if data.ChangePasswordToken == changePwdToken {
			return data.Login, *data.ChangePasswordExpireTime, nil
		}
	}
	return "", time.Time{}, errors.New("not found")
}

func (dao *TestDAO) GetProfileById(profileIds []int32) ([]models.Profile, error) {
	res := make([]models.Profile, 0)
	for _, id := range profileIds {
		profile, ok := dao.UserAvatars[uint32(id)]
		if ok {
			res = append(res, *profile)
		}
	}
	return res, nil
}

func (dao *TestDAO) UpdateProfile(profile models.Profile) error {
	existedProfile, ok := dao.UserAvatars[profile.ProfileId]
	if ok {
		existedProfile.AvatarUrl = profile.AvatarUrl
	} else {
		dao.UserAvatars[profile.ProfileId] = &profile
	}
	return nil
}

func (dao *TestDAO) ExistRecord(profileId uint32, courseId uint32) (bool, error) {
	course, ok := dao.Courses[courseId]
	if !ok {
		return false, nil
	}
	for id := range course.Students {
		if id == profileId {
			return true, nil
		}
	}
	return false, nil
}

func (dao *TestDAO) AddRecord(profileId uint32, courseId uint32) error {
	course, ok := dao.Courses[courseId]
	if !ok {
		return errors.New("not found")
	}
	course.Students[profileId] = true
	return nil
}

func (dao *TestDAO) DeleteRecord(profileId uint32, courseId uint32) error {
	course, ok := dao.Courses[courseId]
	if ok {
		delete(course.Students, profileId)
	}
	return nil
}

func (dao *TestDAO) GetTeachersByCourse(courseId uint32) ([]uint32, error) {
	course, ok := dao.Courses[courseId]
	if !ok {
		return nil, errors.New("not found")
	}
	teachers := make([]uint32, 0)
	for tid := range course.Teachers {
		teachers = append(teachers, tid)
	}
	return teachers, nil
}
