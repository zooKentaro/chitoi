package row

func MainTableStructs() []interface{} {
	return []interface{}{
		&BroadcastMessage{},
		&Business{},
		&User{},
		&UserBusiness{},
		&UserRank{},
		&PersonalMessage{},
		&Room{},
	}
}
