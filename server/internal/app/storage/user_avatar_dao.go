package storage

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
)

type UserAvatarDAO interface {
	GetProfileById(profileIds []int32) ([]models.Profile, error)
	UpdateProfile(profile models.Profile) error
}

type userAvatarDAO struct {
	Storage *DbStorage
}

func (dao *userAvatarDAO) GetProfileById(profileIds []int32) ([]models.Profile, error) {
	var strProfileIds []string
	for _, val := range profileIds {
		strProfileIds = append(strProfileIds, strconv.Itoa(int(val)))
	}
	rows, err := dao.Storage.Db.Query(
		fmt.Sprintf(`SELECT profile_id, avatar_url
		FROM user_avatar_base
		WHERE profile_id IN (%s)`, strings.Join(strProfileIds, ", ")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []models.Profile

	for rows.Next() {
		var profile models.Profile
		if err := rows.Scan(&profile.ProfileId,
			&profile.AvatarUrl); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, err
}

func (dao *userAvatarDAO) UpdateProfile(profile models.Profile) error {
	_, err := dao.Storage.Db.Exec(`INSERT INTO
			   user_avatar_base (profile_id, avatar_url)
		       VALUES ($1, $2)
		       ON CONFLICT (profile_id) DO UPDATE SET avatar_url = EXCLUDED.avatar_url`,
		profile.ProfileId,
		profile.AvatarUrl,
	)
	return err
}
