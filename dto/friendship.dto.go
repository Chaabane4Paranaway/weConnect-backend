package dto

type AddFriendshipDTO struct {
	Pseudo string `json:"pseudo" binding:"required"`
	// User2ID string `json:"user2ID"`
}

// Préparer la liste des amis
type FriendInfo struct {
	ID     string `json:"id"`
	Pseudo string `json:"pseudo"`
	Online bool   `json:"online"`
}
