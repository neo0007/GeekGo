package types

type Dao struct {
	i int
}

func NewDao(i int) *Dao {
	return &Dao{i: i}
}

func (d *Dao) AddDao(j int) int {
	println(d.i + j)
	return d.i + j
}

type Repo struct {
	dao *Dao
}

func NewRepo(dao *Dao) *Repo {
	return &Repo{dao}
}

func (r *Repo) AddRepo(j int) int {
	return r.dao.AddDao(j)
}

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo}
}

func (s *Service) AddService(j int) int {
	return s.repo.AddRepo(j)
}

type Web struct {
	svc *Service
}

func NewWeb(svc *Service) *Web {
	return &Web{svc: svc}
}

func (w *Web) AddWeb(j int) int {
	return w.svc.AddService(j)
}

func main() {
	d := NewDao(3)
	r := NewRepo(d)
	s := NewService(r)
	w := NewWeb(s)
	println(w.svc.repo.dao.i)
	w.AddWeb(200)
}
