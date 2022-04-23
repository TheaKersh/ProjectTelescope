package main

type Session struct {
	X_vals Dataset `json:"X_Vals"`
	Y_vals Dataset `json:"Y_Vals"`
	//Specifies which dataset to fill
	FillY bool `json:"FillY"`
}

type Dataset struct {
	Name   string   `json:"Name"`
	Id     string   `json:"Id"`
	Params []string `json:"Params"`
}

func EmptySess() Session {
	return Session{X_vals: Dataset{Name: ""}, Y_vals: Dataset{Name: ""}}
}

func MakeSession(x_name string, y_name string) Session {
	x_set := Dataset{Name: x_name}
	y_set := Dataset{Name: y_name}
	sess := Session{X_vals: x_set, Y_vals: y_set, FillY: false}
	return sess
}

func (sess Session) SetXFeats(x_feats []string) {
	sess.X_vals.Params = x_feats
}

func (sess Session) SetYFeats(y_feats []string) {
	sess.Y_vals.Params = y_feats
}
