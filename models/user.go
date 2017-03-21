package models

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"cloud-platform-ua/models/mongo"

	"golang.org/x/crypto/scrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User model definiton.
type User struct {
	ID       bson.ObjectId    `bson:"_id"      json:"_id,omitempty"`
	Phone    string    `bson:"phone"    json:"phone,omitempty"`
	Name     string    `bson:"name"     json:"name,omitempty"`
	Password string    `bson:"password" json:"password,omitempty"`
	Salt     string    `bson:"salt"     json:"salt,omitempty"`
	RegDate  time.Time `bson:"reg_date" json:"reg_date,omitempty"`
	NoEncPwd string    `bson:"no_enc_pwd" json:"no_enc_pwd,omitempty"`
	Email    string    `bson:"email" json:"email,omitempty"`
	IsAdmin  int       `bson:"is_admin" json:"is_admin,omitempty"`
}

const pwHashBytes = 64

func generateSalt() (salt string, err error) {
	buf := make([]byte, pwHashBytes)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", buf), nil
}

func generatePassHash(password string, salt string) (hash string, err error) {
	h, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, pwHashBytes)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h), nil
}

// NewUser alloc and initialize a user.
func NewUser(r *RegisterForm, t time.Time) (u *User, err error) {
	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}
	hash, err := generatePassHash(r.Password, salt)
	if err != nil {
		return nil, err
	}

	user := User{
		ID:          bson.NewObjectId(),
		Phone:    r.Phone,
		Name:     r.Name,
		Email:          r.Email,
		Password: hash,
		Salt:     salt,
		NoEncPwd: r.Password,
		IsAdmin: 0,
		RegDate:  t}

	return &user, nil
}

// Insert insert a document to collection.
func (u *User) Insert() (code int, err error) {
	mConn := mongo.Conn()
	defer mConn.Close()

	c := mConn.DB("cloud-platform").C("users")
	err = c.Insert(u)

	if err != nil {
		if mgo.IsDup(err) {
			code = ErrDupRows
		} else {
			code = ErrDatabase
		}
	} else {
		code = 0
	}
	return
}

// FindByID query a document according to input id.
func (u *User) FindByID(id string) (code int, err error) {
	mConn := mongo.Conn()
	defer mConn.Close()

	c := mConn.DB("cloud-platform").C("users")
	err = c.FindId(id).One(u)

	if err != nil {
		if err == mgo.ErrNotFound {
			code = ErrNotFound
		} else {
			code = ErrDatabase
		}
	} else {
		code = 0
	}
	return
}
// query user info by input "name"
func (u *User) FindByName(name string) (code int, err error) {
	mConn := mongo.Conn()
	defer mConn.Clone()

	c := mConn.DB("cloud-platform").C("users")
	err = c.Find(bson.M{"name":name}).One(u)

	if err != nil {
		if err == mgo.ErrNotFound {
			code = ErrNotFound
		} else {
			code = ErrDatabase
		}
	} else {
		code = 0
	}
	return
}
// update user information
func (u *User) UpdateUser() (err error) {
	mConn := mongo.Conn()
	defer mConn.Clone()
	c := mConn.DB("cloud-platform").C("users")
	selector := bson.M{"_id":u.ID}
	data := bson.M{"$set":bson.M{"name":u.Name, "email":u.Email}}
	_, err = c.Upsert(selector, data)
	return
}

// CheckPass compare input password.
func (u *User) CheckPass(pass string) (ok bool, err error) {
	hash, err := generatePassHash(pass, u.Salt)
	if err != nil {
		return false, err
	}

	return u.Password == hash, nil
}

// ClearPass clear password information.
func (u *User) ClearPass() {
	u.Password = ""
	u.Salt = ""
}

// ChangePass update password and salt information according to input id.
func ChangePass(id, oldPass, newPass string) (code int, err error) {
	mConn := mongo.Conn()
	defer mConn.Close()

	c := mConn.DB("").C("users")
	u := User{}
	err = c.FindId(id).One(&u)
	if err != nil {
		if err == mgo.ErrNotFound {
			return ErrNotFound, err
		}

		return ErrDatabase, err
	}

	oldHash, err := generatePassHash(oldPass, u.Salt)
	if err != nil {
		return ErrSystem, err
	}
	newSalt, err := generateSalt()
	if err != nil {
		return ErrSystem, err
	}
	newHash, err := generatePassHash(newPass, newSalt)
	if err != nil {
		return ErrSystem, err
	}

	err = c.Update(bson.M{"_id": id, "password": oldHash}, bson.M{"$set": bson.M{"password": newHash, "salt": newSalt}})
	if err != nil {
		if err == mgo.ErrNotFound {
			return ErrNotFound, err
		}

		return ErrDatabase, err
	}

	return 0, nil
}
