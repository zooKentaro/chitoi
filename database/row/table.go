package row

func MainTableStructs() []interface{} {
	return []interface{}{
		&BroadcastMessage{},
		&Business{},
		&Contact{},
		&User{},
		&UserBusiness{},
		&UserRank{},
		&PersonalMessage{},
		&Room{},
	}
}
