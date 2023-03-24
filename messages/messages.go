package messages

const (
	// KindConnected is sent when user connects
	KindConnected = iota + 1
	// KindUserJoined is sent when someone else joins
	KindUserJoined
	// KindUserLeft is sent when someone leaves
	KindUserLeft
	// KindStroke message specifies a drawn stroke by a user
	KindStroke
	// KindClear message is send when a user clears the screen
	KindClear
)

type User struct {
	ID    string `json:"id"`
	Color string `json:"color"`
}

func NewUserJoined(userID int, color string) *UserJoined {
	return &UserJoined{
		Kind: KindUserJoined,
		User: User{ID: string(rune(userID)), Color: color},
	}
}

type UserJoined struct {
	Kind int  `json:"kind"`
	User User `json:"user"`
}

type UserLeft struct {
	Kind   int `json:"kind"`
	UserID int `json:"userId"`
}

func NewUserLeft(userID int) *UserLeft {
	return &UserLeft{
		Kind:   KindUserLeft,
		UserID: userID,
	}
}

type Connected struct {
	Kind  int    `json:"kind"`
	Color string `json:"color"`
	Users []User `json:"users"`
}

func NewConnected(color string, users []User) *Connected {
	return &Connected{
		Kind:  KindConnected,
		Color: color,
		Users: users,
	}
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Stroke struct {
	Kind   int     `json:"kind"`
	UserID int     `json:"userId"`
	Points []Point `json:"points"`
	Finish bool    `json:"finish"`
}

type Clear struct {
	Kind   int `json:"kind"`
	UserID int `json:"userId"`
}
