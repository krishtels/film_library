package models

import (
	"time"

	"film-library/src/internal/tools"
)

func ToQueryableObject(a *Actor) *tools.QueryableObject {
	qo := tools.NewQueryableObject()

	if len(a.Name) != 0 {
		qo.Add("actor_name", a.Name)
	}

	if len(a.Sex) != 0 {
		qo.Add("sex", a.Sex)
	}

	if !a.Birthday.IsZero() {
		qo.Add("birthday", a.Birthday)
	}

	return qo
}

func ToActorResponse(a *Actor) *ActorResponse {
	return &ActorResponse{
		ID: int(a.ID),
		Info: ActorInfo{
			Name:     a.Name,
			Sex:      a.Sex,
			Birthday: a.Birthday.Format(time.DateOnly),
		},
		Films: a.Films,
	}
}

func ToActor(ai *ActorInfo) *Actor {
	birthday, _ := time.Parse(time.DateOnly, ai.Birthday)

	return &Actor{
		Name:     ai.Name,
		Sex:      ai.Sex,
		Birthday: birthday,
	}
}
