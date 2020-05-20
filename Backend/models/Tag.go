package models

import "github.com/jinzhu/gorm"

type Tag struct {
	gorm.Model
	Name string
}

//func (t *Tags) GetTags() []string {
//	return t.Tags
//}
//
//func (t *Tags) HasTag(tagName string) (bool, int) {
//	for index, tag := range t.Tags {
//		if tag == tagName {
//			return true, index
//		}
//	}
//	return false, -1
//}
//
//func (t *Tags) AddTag(tagName string) {
//	hasTag, _ := t.HasTag(tagName)
//	if !hasTag {
//		t.Tags = append(t.Tags, tagName)
//	}
//}
//
//func (t *Tags) RemoveTag(tagName string) bool {
//	lenTags := len(t.GetTags())
//	hasTag, tagIndex := t.HasTag(tagName)
//	if hasTag {
//		t.Tags[tagIndex] = t.Tags[lenTags-1]
//		t.Tags[lenTags-1] = ""
//		t.Tags = t.Tags[:lenTags-1]
//	}
//	return false
//}
