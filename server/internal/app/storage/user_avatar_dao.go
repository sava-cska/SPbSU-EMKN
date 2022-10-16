package storage

import (
	"fmt"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"strconv"
	"strings"
)

type UserAvatarDAO struct {
	Storage *Storage
}

func (dao *UserAvatarDAO) GetProfileById(profileIds []int32) ([]models.Profile, error) {
	var strProfileIds []string
	for _, val := range profileIds {
		strProfileIds = append(strProfileIds, strconv.Itoa(int(val)))
	}
	rows, err := dao.Storage.Db.Query(
		fmt.Sprintf(`SELECT profile_id, avatar_url, first_name, last_name
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
			&profile.AvatarUrl,
			&profile.FirstName,
			&profile.LastName); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, err
}
