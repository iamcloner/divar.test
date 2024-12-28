package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type Price struct {
	Title      string `json:"title" form:"title" json:"title"`
	Price      int64  `json:"price" form:"price" json:"price"`
	Fixed      bool   `json:"fixed" form:"fixed" json:"fixed"`
	Negotiable bool   `json:"negotiable" form:"negotiable" json:"negotiable"`
}

type Posts struct {
	ID           primitive.ObjectID `json:"id" form:"-" bson:"_id,omitempty"`
	UserId       primitive.ObjectID `json:"userId" form:"userId" bson:"userId"`
	Title        string             `json:"title" form:"title" bson:"title"`
	Description  string             `json:"description" form:"description" bson:"description"`
	CategoryCode int                `json:"categoryCode" form:"categoryCode" bson:"categoryCode"`
	AreaCode     int                `json:"areaCode" form:"areaCode" bson:"areaCode"`
	Images       *[]string          `json:"images" form:"images" bson:"images"`
	Type         string             `json:"type" form:"type" bson:"type"`
	Prices       *[]Price           `json:"prices" form:"prices" bson:"prices"`
	////////////////////////////////////////////////////////////////////////////////
	Area             *int    `json:"Area" form:"Area" bson:"Area,omitempty"`
	Rooms            *int    `json:"rooms" form:"rooms" bson:"rooms,omitempty"`
	YearConstruction *int    `json:"YearConstruction" form:"YearConstruction" bson:"YearConstruction,omitempty"`
	Floors           *int    `json:"floors" form:"floors" bson:"floors,omitempty"`
	Floor            *int    `json:"floor" form:"floor" bson:"floor,omitempty"`
	Elevator         *bool   `json:"elevator" form:"elevator" bson:"elevator,omitempty"`
	Parking          *bool   `json:"parking" form:"parking" bson:"parking,omitempty"`
	Warehouse        *bool   `json:"warehouse" form:"warehouse" bson:"warehouse,omitempty"`
	Balcony          *bool   `json:"balcony" form:"balcony" bson:"balcony,omitempty"`
	ApplicationType  *string `json:"applicationType" form:"applicationType" bson:"applicationType,omitempty"`
	DocumentType     *string `json:"documentType" form:"documentType" json:"documentType"`
	//////////////////////////////////////////////////////////////////////////////
	Infrastructure *int  `json:"Infrastructure" form:"Infrastructure" bson:"Infrastructure,omitempty"`
	Yard           *bool `json:"yard" form:"yard" bson:"yard,omitempty"`
	/////////////////////////////////////////////////////////////////////////////
	Color     *string `json:"color" form:"color" bson:"color,omitempty"`
	Model     *int    `json:"model" form:"model" bson:"model,omitempty"`
	Mileage   *int    `json:"mileage" form:"mileage" bson:"mileage,omitempty"`
	Fuel      *string `json:"fuel" form:"fuel" bson:"fuel,omitempty"`
	Insurance *int    `json:"insurance" form:"insurance" bson:"insurance,omitempty"`
	Chassis   *string `json:"chassis" form:"chassis" bson:"chassis,omitempty"`
	Status    *string `json:"status" form:"status" bson:"status,omitempty"`
	Brand     *string `json:"brand" form:"brand" bson:"brand,omitempty"`
	Gearbox   *string `json:"gearbox" form:"gearbox" bson:"gearbox,omitempty"`
	/////////////////////////////////////////////////////////////////////////////
	Sims    *int    `json:"sims" form:"sims" bson:"sims,omitempty"`
	Storage *string `json:"storage" form:"storage" bson:"storage,omitempty"`
	Ram     *string `json:"ram" form:"ram" bson:"ram,omitempty"`
	Os      *string `json:"os" form:"os" bson:"os,omitempty"`
}

type ApartmentProperties struct {
	Area             int    `json:"Area" form:"Area" json:"Area"`
	Rooms            int    `json:"rooms" form:"rooms" json:"rooms"`
	YearConstruction int    `json:"YearConstruction" form:"YearConstruction" json:"YearConstruction"`
	Floors           int    `json:"floors" form:"floors" json:"floors"`
	Floor            int    `json:"floor" form:"floor" json:"floor"`
	Elevator         bool   `json:"elevator" form:"elevator" json:"elevator"`
	Parking          bool   `json:"parking" form:"parking" json:"parking"`
	Storage          bool   `json:"storage" form:"storage" json:"storage"`
	Balcony          bool   `json:"balcony" form:"balcony" json:"balcony"`
	DocumentType     string `json:"documentType" form:"documentType" json:"documentType"`
}
type VilaProperties struct {
	Infrastructure   int    `json:"Infrastructure" form:"Infrastructure" json:"Infrastructure"`
	Area             int    `json:"Area" form:"Area" json:"Area"`
	YearConstruction int    `json:"YearConstruction" form:"YearConstruction" json:"YearConstruction"`
	Rooms            int    `json:"rooms" form:"rooms" json:"rooms"`
	Parking          bool   `json:"Parking" form:"Parking" json:"Parking"`
	Storage          bool   `json:"Storage" form:"Storage" json:"Storage"`
	Yard             bool   `json:"yard" form:"yard" json:"yard"`
	DocumentType     string `json:"documentType" form:"documentType" json:"documentType"`
}
type LandProperties struct {
	Area         int    `json:"Area" form:"Area" json:"Area"`
	Type         string `json:"Type" form:"Type" json:"Type"`
	DocumentType string `json:"documentType" form:"documentType" json:"documentType"`
}
type CarProperties struct {
	Color     string `json:"color" form:"color" json:"color"`
	Model     int    `json:"model" form:"model" json:"model"`
	Mileage   int    `json:"mileage" form:"mileage" json:"mileage"`
	Fuel      string `json:"fuel" form:"fuel" json:"fuel"`
	Insurance int    `json:"insurance" form:"insurance" json:"insurance"`
	Chassis   string `json:"chassis" form:"chassis" json:"chassis"`
	Status    string `json:"status" form:"status" json:"status"`
	Brand     string `json:"brand" form:"brand" json:"brand"`
	Gearbox   string `json:"gearbox" form:"gearbox" json:"gearbox"`
}
type MotorcycleProperties struct {
	Color     string `json:"color" form:"color" json:"color"`
	Model     int    `json:"model" form:"model" json:"model"`
	Mileage   int    `json:"mileage" form:"mileage" json:"mileage"`
	Fuel      string `json:"fuel" form:"fuel" json:"fuel"`
	Insurance int    `json:"insurance" form:"insurance" json:"insurance"`
	Status    string `json:"status" form:"status" json:"status"`
	Brand     string `json:"brand" form:"brand" json:"brand"`
	Gearbox   string `json:"gearbox" form:"gearbox" json:"gearbox"`
}
type MobileProperties struct {
	Brand   string `json:"brand" form:"brand" json:"brand"`
	Status  string `json:"status" form:"status" json:"status"`
	Sims    int    `json:"sims" form:"sims" json:"sims"`
	Storage string `json:"storage" form:"storage" json:"storage"`
	Ram     string `json:"ram" form:"ram" json:"ram"`
	Color   string `json:"color" form:"color" json:"color"`
	Os      string `json:"os" form:"os" json:"os"`
}
