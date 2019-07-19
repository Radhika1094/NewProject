package model

type UserInformation struct {
	UserName       string  `json:"userName" bson:"userName"`
	Email          string  `json:"email" bson:"email"`
	Password       string  `json:"password" bson:"password"`
	ContactNumber  int     `json:"contactNumber" bson:"contactNumber"`
	DateOfBirth    string  `json:"dateOfBirth" bson:"dateOfBirth"`
	AddressDetails Address `json:"addressDetails" bson:"addressDetails"`
}

type UserCredential struct {
	UserName string `json:"userName" bson:"userName"`
	Password string `json:"password" bson:"password"`
}

type Address struct {
	State    State    `json:"state" bson:"state"`
	District District `json:"district" bson:"district"`
	Taluka   Taluka   `json:"taluka" bson:"taluka"`
	City     City     `json:"city" bson:"city"`
}
type State struct {
	StateName string `json:"stateName" bson:"stateName"`
}
type District struct {
	DistrictName string `json:"districtName" bson:"districtName"`
}
type Taluka struct {
	TalukaName string `json:"talukaName" bson:"talukaName"`
}
type City struct {
	CityName string `json:"cityName" bson:"cityName"`
}
