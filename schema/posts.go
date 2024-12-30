package schema

import (
	"divar.ir/api/repositories/categoryRepositories"
	"divar.ir/internal/listPicker"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
	"time"
)

type Price struct {
	Title      string `json:"title" form:"title" json:"title"`
	Price      int64  `json:"price" form:"price" json:"price"`
	Fixed      bool   `json:"fixed" form:"fixed" json:"fixed"`
	Negotiable bool   `json:"negotiable" form:"negotiable" json:"negotiable"`
}

type Posts struct {
	ID           primitive.ObjectID `json:"id" form:"-" bson:"_id,omitempty"`
	UserId       primitive.ObjectID `json:"userId" form:"userId" bson:"userId"`
	Title        string             `json:"title" form:"title" bson:"title" validate:"required"`
	Description  string             `json:"description" form:"description" bson:"description" validate:"required"`
	CategoryCode int                `json:"categoryCode" form:"categoryCode" bson:"categoryCode" validate:"required"`
	AreaCode     int                `json:"areaCode" form:"areaCode" bson:"areaCode" validate:"required"`
	Images       []string           `json:"images" form:"images" bson:"images" validate:"required"`
	Type         string             `json:"type" form:"type" bson:"type" validate:"required"`
	Prices       []Price            `json:"prices" form:"prices" bson:"prices" validate:"required"`
	Attributes   map[string]any     `json:"attributes" form:"attributes" bson:"attributes" validate:"required"`
}

func (p Posts) checkAttrType(key string, t string, min int, max int, pickListName string, pickLimit int) error {
	if p.Attributes[key] == nil {
		return errors.New(key + " is required.")
	}
	switch t {
	case "string":
		value, ok := p.Attributes[key].(string)
		if !ok {
			return errors.New(key + " is not string")
		}
		if len(value) < min || len(value) > max {
			return errors.New(key + " length must be between " + strconv.Itoa(min) + " and " + strconv.Itoa(max))
		}

	case "int":
		floatValue, ok := p.Attributes[key].(float64)
		if !ok {
			return errors.New(key + " is not number")
		}
		value := int(floatValue)
		if value < min || value > max {
			return errors.New(key + " must be between " + strconv.Itoa(min) + " and " + strconv.Itoa(max))
		}
	case "float":
		_, ok := p.Attributes[key].(float64)
		if !ok {
			return errors.New(key + " is not float")
		}
	case "bool":
		_, ok := p.Attributes[key].(bool)
		if !ok {
			return errors.New(key + " is not bool")
		}
	case "datetime":
		_, ok := p.Attributes[key].(time.Time)
		if !ok {
			return errors.New(key + " is not datetime")
		}
	case "pick":
		pickList := strings.Split(p.Attributes[key].(string), ",")
		if len(pickList) > pickLimit {
			return errors.New(key + " has too many values")
		}
		for i := range pickList {
			check, err := listPicker.CheckItem(pickListName, pickList[i])
			if err != nil {
				return err
			}
			if !check {
				return errors.New(key + " not valid value picked")
			}
		}

	default:
		return errors.New("invalid type")
	}
	return nil
}

func (p Posts) Validate() error {
	if p.Title == "" {
		return errors.New("title is required")
	}
	if check, _ := listPicker.CheckItem("postTypes", p.Type); !check {
		return errors.New("invalid type")
	}
	if p.Description == "" {
		return errors.New("description is required")
	}
	if p.Prices == nil || len(p.Prices) < 1 {
		return errors.New("prices is required")
	}
	if p.CategoryCode == 0 {
		return errors.New("category code is required")
	}
	result, err := categoryRepositories.GetCategories(p.CategoryCode)
	if err != nil {
		return errors.New("invalid Category code")
	}
	resultMap, ok := result.(primitive.M)
	if !ok {
		return errors.New("invalid Requirement")
	}
	requires, ok := resultMap["requires"].(primitive.A)
	if !ok {
		return errors.New("invalid Requirements")
	}
	for _, requirement := range requires {
		requirementStr, ok := requirement.(string)
		if !ok {
			return errors.New("invalid requirement format")
		}
		requirementValues := strings.Split(requirementStr, ",")
		if len(requirementValues) < 2 {
			return errors.New("invalid requirement format")
		}
		requirementName := requirementValues[0]
		requirementType := requirementValues[1]
		if requirementType == "string" || requirementType == "int" {
			if len(requirementValues) < 4 {
				return errors.New("invalid requirement format")
			}
			low, err := strconv.Atoi(requirementValues[2])
			if err != nil {
				return errors.New("invalid requirement format")
			}
			high, err := strconv.Atoi(requirementValues[3])
			if err != nil {
				return errors.New("invalid requirement format")
			}
			err = p.checkAttrType(requirementName, requirementType, low, high, "", 0)
			if err != nil {
				return err
			}
		} else if requirementType == "pick" {
			if len(requirementValues) < 4 {
				return errors.New("invalid requirement format")
			}
			pickList := requirementValues[2]
			pickLimit, err := strconv.Atoi(requirementValues[3])
			if err != nil {
				return errors.New("invalid requirement format")
			}
			err = p.checkAttrType(requirementName, requirementType, 0, 0, pickList, pickLimit)
			if err != nil {
				return err
			}
		} else {
			err = p.checkAttrType(requirementName, requirementType, 0, 0, "", 0)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
