package main

type Wallet struct{
	username string
	password string
}

func New(username, password string) *Wallet  {
	return &Wallet{
		username: username,
		password: password,
	}
}

func (w *Wallet) ToString()(string, string){
	return w.username, w.password
}

func (w *Wallet)Save()(int64, error){
	id, err:=InsertRow(*w)
	    if err != nil {
		return 0, err
    }
	return id, nil
}

