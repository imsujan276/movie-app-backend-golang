package handlers

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
)

var movies []*models.Movie

// graphql schema definition
var fields = graphql.Fields{
	"movie": &graphql.Field{
		Type:        movieType,
		Description: "Get movie by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, ok := p.Args["id"].(int)
			if ok {
				for _, movie := range movies {
					if movie.ID == id {
						return movie, nil
					}
				}
			}
			return nil, nil
		},
	},
	"movies": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "get all movies",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return movies, nil
		},
	},
	"search": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "search movies by title",
		Args: graphql.FieldConfigArgument{
			"titleContains": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var theList []*models.Movie
			search, ok := p.Args["titleContains"].(string)
			if ok {
				for _, movie := range movies {
					if strings.Contains(strings.ToLower(movie.Title), strings.ToLower(search)) {
						theList = append(theList, movie)
					}
				}
			}
			return theList, nil
		},
	},
}

var movieType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Movie",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"year": &graphql.Field{
				Type: graphql.Int,
			},
			"release_date": &graphql.Field{
				Type: graphql.DateTime,
			},
			"runtime": &graphql.Field{
				Type: graphql.Int,
			},
			"rating": &graphql.Field{
				Type: graphql.Int,
			},
			"mpaa_rating": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"poster": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

func (m *Repository) GetAllMoviesGraphQL(w http.ResponseWriter, r *http.Request) {
	movies, _ = m.App.Models.DB.All()

	q, err := io.ReadAll(r.Body)
	if err != nil {
		m.errorJSON(w, err)
		return
	}

	query := string(q)

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		m.errorJSON(w, errors.New("failed to create Schema"))
		return
	}

	params := graphql.Params{Schema: schema, RequestString: query}
	response := graphql.Do(params)
	if len(response.Errors) > 0 {
		m.errorJSON(w, errors.New(fmt.Sprintf("failed: %+v", response.Errors)))
		return
	}

	j, _ := json.MarshalIndent(response, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
