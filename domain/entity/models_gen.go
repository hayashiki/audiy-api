// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package entity

import (
	"fmt"
	"io"
	"strconv"
)

type Connection interface {
	IsConnection()
}

type Edge interface {
	IsEdge()
}

type Node interface {
	IsNode()
}

type AudioConnection struct {
	PageInfo *PageInfo    `json:"pageInfo"`
	Edges    []*AudioEdge `json:"edges"`
}

func (AudioConnection) IsConnection() {}

type AudioEdge struct {
	Cursor string `json:"cursor"`
	Node   *Audio `json:"node"`
}

func (AudioEdge) IsEdge() {}

type AudioFilter struct {
	Played *bool `json:"played"`
	Stared *bool `json:"stared"`
	Liked  *bool `json:"liked"`
}

type AudiosInput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CommentConnection struct {
	PageInfo *PageInfo      `json:"pageInfo"`
	Edges    []*CommentEdge `json:"edges"`
}

func (CommentConnection) IsConnection() {}

type CommentEdge struct {
	Cursor string   `json:"cursor"`
	Node   *Comment `json:"node"`
}

func (CommentEdge) IsEdge() {}

type CommentOrder struct {
	Field     *CommentOrderField `json:"field"`
	Direction *SortDirection     `json:"direction"`
}

type CreateCommentInput struct {
	AudioID string `json:"audioID"`
	Body    string `json:"body"`
}

type CreateFeedInput struct {
	AudioID string `json:"audioID"`
}

type CreatePlayPayload struct {
	Result bool  `json:"result"`
	Play   *Play `json:"play"`
}

type CreateUserInput struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type DeleteCommentResult struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type DeleteFeedResult struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type DeleteLikeInput struct {
	ID string `json:"id"`
}

type DeleteStarInput struct {
	ID string `json:"id"`
}

type FeedConnection struct {
	PageInfo *PageInfo   `json:"pageInfo"`
	Edges    []*FeedEdge `json:"edges"`
}

func (FeedConnection) IsConnection() {}

type FeedEdge struct {
	Cursor string `json:"cursor"`
	Node   *Feed  `json:"node"`
}

func (FeedEdge) IsEdge() {}

type FeedFilter struct {
	State *FeedEvent `json:"state"`
}

type PageInfo struct {
	Cursor    string `json:"cursor"`
	TotalPage int    `json:"totalPage"`
	HasMore   bool   `json:"hasMore"`
}

type QuerySpec struct {
	Order  []AudioOrder `json:"order"`
	Cursor string       `json:"cursor"`
	Limit  *int         `json:"limit"`
}

type UpdateAudioInput struct {
	AudioID string `json:"audioID"`
}

type UpdateCommentInput struct {
	ID string `json:"id"`
}

type UpdateFeedInput struct {
	ID    string    `json:"id"`
	Event FeedEvent `json:"event"`
}

type Version struct {
	Hash    string `json:"hash"`
	Version string `json:"version"`
}

type AudioOrder string

const (
	AudioOrderPublishedAtAsc  AudioOrder = "PUBLISHED_AT_ASC"
	AudioOrderPublishedAtDesc AudioOrder = "PUBLISHED_AT_DESC"
)

var AllAudioOrder = []AudioOrder{
	AudioOrderPublishedAtAsc,
	AudioOrderPublishedAtDesc,
}

func (e AudioOrder) IsValid() bool {
	switch e {
	case AudioOrderPublishedAtAsc, AudioOrderPublishedAtDesc:
		return true
	}
	return false
}

func (e AudioOrder) String() string {
	return string(e)
}

func (e *AudioOrder) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AudioOrder(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AudioOrder", str)
	}
	return nil
}

func (e AudioOrder) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type CommentOrderField string

const (
	CommentOrderFieldID CommentOrderField = "ID"
)

var AllCommentOrderField = []CommentOrderField{
	CommentOrderFieldID,
}

func (e CommentOrderField) IsValid() bool {
	switch e {
	case CommentOrderFieldID:
		return true
	}
	return false
}

func (e CommentOrderField) String() string {
	return string(e)
}

func (e *CommentOrderField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CommentOrderField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CommentOrderField", str)
	}
	return nil
}

func (e CommentOrderField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FeedEvent string

const (
	FeedEventPlayed   FeedEvent = "PLAYED"
	FeedEventUnplayed FeedEvent = "UNPLAYED"
	FeedEventStared   FeedEvent = "STARED"
	FeedEventUnstared FeedEvent = "UNSTARED"
	FeedEventLiked    FeedEvent = "LIKED"
	FeedEventUnliked  FeedEvent = "UNLIKED"
)

var AllFeedEvent = []FeedEvent{
	FeedEventPlayed,
	FeedEventUnplayed,
	FeedEventStared,
	FeedEventUnstared,
	FeedEventLiked,
	FeedEventUnliked,
}

func (e FeedEvent) IsValid() bool {
	switch e {
	case FeedEventPlayed, FeedEventUnplayed, FeedEventStared, FeedEventUnstared, FeedEventLiked, FeedEventUnliked:
		return true
	}
	return false
}

func (e FeedEvent) String() string {
	return string(e)
}

func (e *FeedEvent) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FeedEvent(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FeedEvent", str)
	}
	return nil
}

func (e FeedEvent) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "ASC"
	SortDirectionDesc SortDirection = "DESC"
)

var AllSortDirection = []SortDirection{
	SortDirectionAsc,
	SortDirectionDesc,
}

func (e SortDirection) IsValid() bool {
	switch e {
	case SortDirectionAsc, SortDirectionDesc:
		return true
	}
	return false
}

func (e SortDirection) String() string {
	return string(e)
}

func (e *SortDirection) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortDirection(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortDirection", str)
	}
	return nil
}

func (e SortDirection) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type State string

const (
	StatePlayed   State = "Played"
	StateUnplayed State = "Unplayed"
)

var AllState = []State{
	StatePlayed,
	StateUnplayed,
}

func (e State) IsValid() bool {
	switch e {
	case StatePlayed, StateUnplayed:
		return true
	}
	return false
}

func (e State) String() string {
	return string(e)
}

func (e *State) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = State(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid State", str)
	}
	return nil
}

func (e State) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
