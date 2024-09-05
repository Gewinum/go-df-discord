package server

type NewUserHandler func(user *User)

type Service struct {
	repo     Repository
	codeStr  CodeStore
	handlers []NewUserHandler
}

func NewService(repo Repository, codeStr CodeStore) *Service {
	return &Service{repo: repo, codeStr: codeStr}
}

func (s *Service) AddHandler(handler NewUserHandler) {
	s.handlers = append(s.handlers, handler)
}

func (s *Service) IssueCode(xuid string) (*CodeInformation, error) {
	existing, _ := s.repo.GetUserByXUID(xuid)
	if existing != nil {
		return nil, NewApplicationError(40000, "Minecraft account is already bound to ID "+existing.Discord)
	}
	return s.codeStr.Issue(xuid)
}

func (s *Service) CheckCode(code string) (*CodeInformation, error) {
	info, err := s.codeStr.GetInformation(code)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *Service) RevokeCode(code string) error {
	return s.codeStr.Revoke(code)
}

func (s *Service) GetUserByXUID(xuid string) (*User, error) {
	return s.repo.GetUserByXUID(xuid)
}

func (s *Service) GetUserByDiscord(discord string) (*User, error) {
	return s.repo.GetUserByDiscord(discord)
}

func (s *Service) CreateUser(discord, xuid string) (*User, error) {
	user, err := s.repo.CreateUser(discord, xuid)
	if err != nil {
		return nil, err
	}
	for _, handler := range s.handlers {
		handler(user)
	}
	return user, nil
}

func (s *Service) DeleteUserByDiscord(discord string) error {
	return s.repo.DeleteUserByDiscord(discord)
}

func (s *Service) DeleteUserByXUID(xuid string) error {
	return s.repo.DeleteUserByXUID(xuid)
}
