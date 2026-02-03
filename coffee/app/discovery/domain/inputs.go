package domain

type UpsertProfileUserDataInput struct {
	ProfileAdminDetails *ProfileAdminDetails `json:"profileAdminDetails"`
	ProfileUserDetails  ProfileUserDetails   `json:"profileUserDetails"`
	Socials             []SocialAccountInput `json:"socials"`
}

type UpsertProfileAdminDataInput struct {
	ProfileAdminDetails ProfileAdminDetails  `json:"profileAdminDetails"`
	Socials             []SocialAccountInput `json:"socials"`
}

type PatchProfileAdminDetailsInput struct {
	ProfileAdminDetails *ProfileAdminDetails `json:"profileAdminDetails"`
}

type SocialAccountInput struct {
	Platform string `json:"platform"`
	Handle   string `json:"handle"`
}
